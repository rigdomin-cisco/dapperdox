package discovery

import (
	"errors"
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

// watcher defines an event loop.
//
// The watcher guarantees the following consistency requirement: registry
// view in the watcher is as AT LEAST as fresh as the moment notification
// arrives, but MAY BE more fresh (e.g. "delete" cancels an "add" event).  For
// example, an event for a service creation will see a service registry without
// the service if the event is immediately followed by the service deletion
// event.
//
// Handlers execute on the single worker queue in the order they are appended.
// Handlers receive the notification event and the associated object.  Note
// that all handlers must be appended before starting the watcher.
type watcher interface {
	// AppendServiceHandler notifies about changes to the service catalog.
	AppendServiceHandler(f func(*model.Service, model.Event))

	// AppendInstanceHandler notifies about changes to the service instances
	// for a service.
	AppendInstanceHandler(f func(*model.ServiceInstance, model.Event))

	// Run until a signal is received
	Run(stop <-chan struct{})
}

// itemizer enumerates service instances.
type itemizer interface {
	// ManagementPorts lists set of management ports associated with an IPv4 address.
	// These management ports are typically used by the platform for out of band management
	// tasks such as health checks, etc. In a scenario where the proxy functions in the
	// transparent mode (traps all traffic to and from the service instance IP address),
	// the configuration generated for the proxy will not manipulate traffic destined for
	// the management ports
	ManagementPorts(addr string) model.PortList
}

// svcItemWatcher provides service instance discovery within an event loop.
type svcItemWatcher interface {
	watcher
	itemizer
}

// catalogOptions stores the configurable attributes of a catalog.
type catalogOptions struct {
	// Namespace the watcher watches. If set to meta_v1.NamespaceAll (""), watcher watches all namespaces
	WatchedNamespace string
	ResyncPeriod     time.Duration
	DomainSuffix     string
}

// catalog is a collection of synchronized resource watchers
// caches are thread-safe
type catalog struct {
	domainSuffix string

	client    kubernetes.Interface
	queue     queue
	services  cacheHandler
	endpoints cacheHandler

	pods *podCache
}

type cacheHandler struct {
	informer cache.SharedIndexInformer
	handler  *chainHandler
}

// newClient creates an in cluster connection to Kubernetes API
func newClient() (kubernetes.Interface, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	return kubernetes.NewForConfig(config)
}

// newCatalog creates a new Kubernetes watcher
func newCatalog(client kubernetes.Interface, options catalogOptions) svcItemWatcher {
	log().Infof("Service watcher watching namespace %q", options.WatchedNamespace)

	// Queue requires a time duration for a retry delay after a handler error
	out := &catalog{
		domainSuffix: options.DomainSuffix,
		client:       client,
		queue:        newQueue(1 * time.Second),
	}

	out.services = out.createInformer(&v1.Service{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Services(options.WatchedNamespace).List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Services(options.WatchedNamespace).Watch(opts)
		})

	out.endpoints = out.createInformer(&v1.Endpoints{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).Watch(opts)
		})

	out.pods = newPodCache(out.createInformer(&v1.Pod{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).Watch(opts)
		}))

	return out
}

// AppendServiceHandler implements a service catalog operation
func (c *catalog) AppendServiceHandler(f func(*model.Service, model.Event)) {
	c.services.handler.Append(func(obj interface{}, event model.Event) error {
		svc := *obj.(*v1.Service)
		log().Debugf("(Handler) Service Details: %v", svc)

		// Do not handle "kube-system" services
		if svc.Namespace == meta_v1.NamespaceSystem {
			return nil
		}

		log().Infof("Handle service %s in namespace %s", svc.Name, svc.Namespace)

		if svcConv := convertService(svc, c.domainSuffix); svcConv != nil {
			f(svcConv, event)
		}
		return nil
	})
}

// AppendInstanceHandler implements a service catalog operation
func (c *catalog) AppendInstanceHandler(f func(*model.ServiceInstance, model.Event)) {
	c.endpoints.handler.Append(func(obj interface{}, event model.Event) error {
		ep := *obj.(*v1.Endpoints)
		log().Debugf("(Handler) Endpoint Details: %v", ep)

		// Do not handle "kube-system" endpoints
		if ep.Namespace == meta_v1.NamespaceSystem {
			return nil
		}

		log().Infof("Handle endpoint %s in namespace %s", ep.Name, ep.Namespace)
		addrs := make(map[string]bool)
		for _, ss := range ep.Subsets {
			for _, ea := range ss.Addresses {
				addrs[ea.IP] = true
			}
		}

		for _, instance := range c.hostInstances(addrs) {
			f(instance, event)
		}

		return nil
	})
}

// Run watcher until a signal is received
func (c *catalog) Run(stop <-chan struct{}) {
	go c.queue.Run(stop)
	go c.services.informer.Run(stop)
	go c.endpoints.informer.Run(stop)
	go c.pods.informer.Run(stop)

	<-stop
	log().Info("catalog terminated")
}

// ManagementPorts implements a service catalog operation
func (c *catalog) ManagementPorts(addr string) model.PortList {
	pod, exists := c.pods.getPodByIP(addr)
	if !exists {
		return nil
	}

	managementPorts, err := convertProbesToPorts(&pod.Spec)
	if err != nil {
		log().Infof("Error while parsing liveliness and readiness probe ports for %s => %v", addr, err)
	}

	// We continue despite the error because healthCheckPorts could return a partial
	// list of management ports
	return managementPorts
}

// notify is the first handler in the handler chain.
// Returning an error causes repeated execution of the entire chain.
func (c *catalog) notify(obj interface{}, event model.Event) error {
	if !c.services.informer.HasSynced() || !c.endpoints.informer.HasSynced() || !c.pods.informer.HasSynced() {
		return errors.New("waiting till full synchronization")
	}
	k, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log().Infof("Error retrieving key: %v", err)
	} else {
		log().Infof("Event %s: key %#v", event, k)
	}
	return nil
}

func (c *catalog) createInformer(o runtime.Object, resyncPeriod time.Duration, lf cache.ListFunc, wf cache.WatchFunc) cacheHandler {
	handler := &chainHandler{funcs: []handler{c.notify}}

	// TODO: finer-grained index (perf)
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{ListFunc: lf, WatchFunc: wf}, o,
		resyncPeriod, cache.Indexers{})

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			// TODO: filtering functions to skip over un-referenced resources (perf)
			AddFunc: func(obj interface{}) {
				c.queue.Push(task{handler: handler.Apply, obj: obj, event: model.EventAdd})
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					c.queue.Push(task{handler: handler.Apply, obj: cur, event: model.EventUpdate})
				}
			},
			DeleteFunc: func(obj interface{}) {
				c.queue.Push(task{handler: handler.Apply, obj: obj, event: model.EventDelete})
			},
		})

	return cacheHandler{informer: informer, handler: handler}
}

// serviceByKey retrieves a service by name and namespace
func (c *catalog) serviceByKey(name, namespace string) (*v1.Service, bool) {
	item, exists, err := c.services.informer.GetStore().GetByKey(keyFunc(name, namespace))
	if err != nil {
		log().Infof("serviceByKey(%s, %s) => error %v", name, namespace, err)
		return nil, false
	}
	if !exists {
		return nil, false
	}
	return item.(*v1.Service), true
}

// hostInstances implements a service catalog operation
func (c *catalog) hostInstances(addrs map[string]bool) []*model.ServiceInstance {
	var out []*model.ServiceInstance
	for _, item := range c.endpoints.informer.GetStore().List() {
		ep := *item.(*v1.Endpoints)
		for _, ss := range ep.Subsets {
			for _, ea := range ss.Addresses {
				if !addrs[ea.IP] {
					continue
				}

				item, exists := c.serviceByKey(ep.Name, ep.Namespace)
				if !exists {
					continue
				}
				svc := convertService(*item, c.domainSuffix)
				if svc == nil {
					continue
				}
				for _, port := range ss.Ports {
					svcPort, exists := svc.Ports.Get(port.Name)
					if !exists {
						continue
					}
					out = append(out, &model.ServiceInstance{
						Endpoint: model.NetworkEndpoint{
							Address:     ea.IP,
							Port:        int(port.Port),
							ServicePort: svcPort,
						},
						Service: svc,
					})
				}
			}
		}
	}
	return out
}
