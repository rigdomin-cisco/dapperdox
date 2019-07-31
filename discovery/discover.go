package discovery

import (
	"net/http"
	"sync"

	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

// Discoverer represents the state of the discovery mechanism
type Discoverer struct {
	mu       sync.RWMutex
	services svcItemWatcher
	data     state
	client   *http.Client
}

type state struct {
	services  *model.ServiceMap
	instances *model.InstanceMap

	apis []string
}

// NewDiscoverer configures a new instance of a Discoverer using Kubernetes client
func NewDiscoverer() (*Discoverer, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	ctl := newCatalog(client, catalogOptions{
		WatchedNamespace: viper.GetString(config.DiscoverNamespace),
		ResyncPeriod:     viper.GetDuration(config.DiscoverInterval),
		DomainSuffix:     "cluster.local",
	})

	sm := model.NewServiceMap()
	im := model.NewInstanceMap()

	d := &Discoverer{
		services: ctl,
		data: state{
			services:  &sm,
			instances: &im,
			apis:      make([]string, 0),
		},
		client: &http.Client{Timeout: viper.GetDuration(config.DiscoverTimeout)},
	}

	// register handlers
	ctl.AppendServiceHandler(d.updateServices)
	ctl.AppendInstanceHandler(d.updateInstances)

	return d, nil
}

// Run starts the discovery process
func (d *Discoverer) Run(stop <-chan struct{}) {
	d.findAPIPaths()
	if d.services != nil {
		go d.services.Run(stop)
	}
	<-stop
}

// APIList returns a list of the discovered APIs
func (d *Discoverer) APIList() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.data.apis
}

func (d *Discoverer) updateServices(s *model.Service, e model.Event) {
	log().Infof("(Discover Handler) Service: %v Event: %v", s, e)

	switch e {
	case model.EventAdd, model.EventUpdate:
		d.data.services.Insert(s)
	case model.EventDelete:
		d.data.services.Delete(s)
	}

	d.findAPIPaths()
}

func (d *Discoverer) updateInstances(s *model.ServiceInstance, e model.Event) {
	log().Infof("(Discover Handler) Instance: %v Event: %v", s, e)

	switch e {
	case model.EventAdd, model.EventUpdate:
		d.data.instances.Insert(s)
	case model.EventDelete:
		d.data.instances.Delete(s)
	}

	d.findAPIPaths()
}
