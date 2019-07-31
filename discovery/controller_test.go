package discovery

import (
	"os"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

/*
func eventually(f func() bool, t *testing.T) {
	interval := 64 * time.Millisecond
	for i := 0; i < 10; i++ {
		if f() {
			return
		}
		t.Log("Sleeping ", interval)
		time.Sleep(interval)
		interval = 2 * interval
	}
	t.Errorf("Failed to satisfy function")
}
*/

const (
	// testService = "test"
	resync = 1 * time.Second
)

func Test_newClient(t *testing.T) {
	if _, err := newClient(); err == nil {
		t.Errorf("newClient() error = %v, wantErr %v", err, true)
	}

	envVars := []string{"KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT"}
	existingEnv := make(map[string]string)
	for _, k := range envVars {
		existingEnv[k] = os.Getenv(k)
	}
	defer func() {
		for k, v := range existingEnv {
			_ = os.Setenv(k, v)
		}
	}()

	_ = os.Setenv("KUBERNETES_SERVICE_HOST", "localhost")
	_ = os.Setenv("KUBERNETES_SERVICE_PORT", "8080")

	copyKubeToken()
	copyKubeCert()

	d, err := newClient()
	if err != nil {
		t.Errorf("newClient() error = %v, wantErr %v", err, true)
	}
	if d == nil {
		t.Errorf("newClient() = %v, want not nil", d)
	}
}

/*
func TestServices(t *testing.T) {
	cl := fake.NewSimpleClientset()
	t.Parallel()
	ns, err := createNamespace(t, cl)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer deleteNamespace(t, cl, ns)

	stop := make(chan struct{})
	defer close(stop)

	ctl := newCatalog(cl, catalogOptions{
		WatchedNamespace: ns,
		ResyncPeriod:     resync,
		DomainSuffix:     domainSuffix,
	})
	go ctl.Run(stop)

	hostname := serviceHostname(testService, ns, domainSuffix)

	var sds itemizer = ctl
	makeService(testService, ns, cl, t)
	eventually(func() bool {
		out, clientErr := sds.Services()
		if clientErr != nil {
			return false
		}
		t.Log("Services:", spew.Sdump(out))

		for _, item := range out {
			if item.Hostname == hostname &&
				len(item.Ports) == 1 &&
				item.Ports[0].Protocol == model.ProtocolHTTP {
				return true
			}
		}
		return false
	}, t)

	svc, err := sds.GetService(hostname)
	if err != nil {
		t.Errorf("GetService(%q) encountered unexpected error: %v", hostname, err)
	}
	if svc == nil {
		t.Errorf("GetService(%q) => should exists", hostname)
	}
	if svc.Hostname != hostname {
		t.Errorf("GetService(%q) => %q", hostname, svc.Hostname)
	}

	missing := serviceHostname("does-not-exist", ns, domainSuffix)
	svc, err = sds.GetService(missing)
	if err != nil {
		t.Errorf("GetService(%q) encountered unexpected error: %v", missing, err)
	}
	if svc != nil {
		t.Errorf("GetService(%q) => %s, should not exist", missing, svc.Hostname)
	}

	_, err = sds.GetService("fake")
	if err == nil {
		t.Errorf("GetService(%q) expected error: but no error returned", "fake")
	}

}

func makeService(n, ns string, cl kubernetes.Interface, t *testing.T) {
	_, err := cl.CoreV1().Services(ns).Create(&v1.Service{
		ObjectMeta: meta_v1.ObjectMeta{Name: n},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Port:     80,
					Name:     "http-example",
					Protocol: "TCP",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("Created service", n)
}
*/

/*
func Test_catalog_getPodAZ(t *testing.T) {

	pod1 := generatePod("pod1", "nsA", "", "node1", map[string]string{"app": "prod-app"})
	pod2 := generatePod("pod2", "nsB", "", "node2", map[string]string{"app": "prod-app"})
	testCases := []struct {
		name   string
		pods   []*v1.Pod
		wantAZ map[*v1.Pod]string
	}{
		{
			name: "should return correct az for given address",
			pods: []*v1.Pod{pod1, pod2},
			wantAZ: map[*v1.Pod]string{
				pod1: "region1/zone1",
				pod2: "region2/zone2",
			},
		},
		{
			name: "should return false if pod isnt in the cache",
			wantAZ: map[*v1.Pod]string{
				pod1: "",
				pod2: "",
			},
		},
		{
			name: "should return false if node isnt in the cache",
			pods: []*v1.Pod{pod1, pod2},
			wantAZ: map[*v1.Pod]string{
				pod1: "",
				pod2: "",
			},
		},
		{
			name: "should return false and empty string if node doesnt have zone label",
			pods: []*v1.Pod{pod1, pod2},
			wantAZ: map[*v1.Pod]string{
				pod1: "",
				pod2: "",
			},
		},
		{
			name: "should return false and empty string if node doesnt have region label",
			pods: []*v1.Pod{pod1, pod2},
			wantAZ: map[*v1.Pod]string{
				pod1: "",
				pod2: "",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			// Setup kube caches
			w := makeFakeKubeAPIcatalog("")
			addPods(t, w, c.pods...)
			for i, pod := range c.pods {
				ip := fmt.Sprintf("128.0.0.%v", i+1)
				id := fmt.Sprintf("%v/%v", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
				w.pods.keys[ip] = id
			}

			// Verify expected existing pod AZs
			for pod, wantAZ := range c.wantAZ {
				az, found := w.GetPodAZ(pod)
				if wantAZ != "" {
					if !reflect.DeepEqual(az, wantAZ) {
						t.Errorf("Wanted az: %s, got: %s", wantAZ, az)
					}
				} else {
					if found {
						t.Errorf("Unexpectedly found az: %s for pod: %s", az, pod.ObjectMeta.Name)
					}
				}
			}
		})
	}

}
*/

func setupInstances(w *catalog, namespace string, t *testing.T) {
	pods := []*v1.Pod{
		generatePod("pod1", namespace, "", "node1", map[string]string{"app": "test-app"}),
		generatePod("pod2", namespace, "", "node2", map[string]string{"app": "prod-app"}),
		generatePod("pod3", namespace, "", "node1", map[string]string{"app": "prod-app"}),
	}
	pods[1].Status = v1.PodStatus{PodIP: "128.0.0.2"}
	addPods(t, w, pods...)

	// Populate pod cache.
	w.pods.keys["128.0.0.1"] = namespace + "/pod1"
	w.pods.keys["128.0.0.2"] = namespace + "/pod2"
	w.pods.keys["128.0.0.3"] = namespace + "/pod3"

	createService(w, "svc1", namespace, nil, []int32{8080}, map[string]string{"app": "prod-app"}, t)
	createService(w, "svc2", namespace, nil, []int32{8081}, map[string]string{"app": "staging-app"}, t)

	// Endpoints are generated by Kubernetes from pod labels and service selectors.
	// Here we manually create them for mocking purpose.
	svc1Ips := []string{"128.0.0.2"}
	svc2Ips := []string{}
	portNames := []string{"test-port"}
	createEndpoints(w, "svc1", namespace, portNames, svc1Ips, t)
	createEndpoints(w, "svc2", namespace, portNames, svc2Ips, t)
}

func Test_catalog_ManagementPorts(t *testing.T) {
	stop := make(chan struct{})
	defer close(stop)

	namespace := "nsA"
	w := makeFakeKubeAPIcatalog(namespace)
	go w.Run(stop)

	setupInstances(w, namespace, t)

	ports := w.ManagementPorts("128.0.0.1")
	if len(ports) != 0 {
		t.Error("ManagementPorts Expected to retrieve no ports, but got: ", ports)
	}

	ports = w.ManagementPorts("128.0.0.9")
	if ports != nil {
		t.Error("ManagementPorts Expected to pod not to exist, but got: ", ports)
	}
}

/*
func Testcatalog_InstanceMethods(t *testing.T) {
	stop := make(chan struct{})
	defer close(stop)

	namespace := "nsA"
	w := makeFakeKubeAPIcatalog(namespace)
	go w.Run(stop)

	setupInstances(w, namespace, t)

	instances, err := w.Instances("svc1", []string{"test-port"})
	if err == nil && instances != nil {
		t.Error("Instances() Eexpected error but got nil")
	}

	hostname := serviceHostname("svc1", namespace, domainSuffix)
	instances, err = w.Instances(hostname, []string{"test-port"})
	if err != nil {
		t.Errorf("Instances() Unexpected error: %v", err)
	}
	if len(instances) == 0 {
		t.Error("Instances() Expected to retrieve >0 instances, but got: ", instances)
	}

	hostname = serviceHostname("svc2", namespace, domainSuffix)
	instances, err = w.Instances(hostname, []string{})
	if err != nil {
		t.Errorf("Instances() Unexpected error: %v", err)
	}
	if len(instances) != 0 {
		t.Error("Instances() Expected to retrieve 0 instances, but got: ", instances)
	}

	hostname = serviceHostname("svc2", "nsC", domainSuffix)
	instances, err = w.Instances(hostname, []string{})
	if err != nil {
		t.Errorf("Instances() Unexpected error: %v", err)
	}
	if len(instances) != 0 {
		t.Error("Instances() Expected to retrieve 0 instances, but got: ", instances)
	}

	instances, err = w.HostInstances(map[string]bool{"128.0.0.2": true})
	if err != nil {
		t.Errorf("HostInstances() Unexpected error: %v", err)
	}
	if len(instances) == 0 {
		t.Error("HostInstances() Expected to retrieve >0 instances, but got: ", instances)
	}

}
*/

func makeFakeKubeAPIcatalog(namespace string) *catalog {
	if namespace == "" {
		namespace = "default"
	}
	clientSet := fake.NewSimpleClientset()
	ctl := newCatalog(clientSet, catalogOptions{
		WatchedNamespace: namespace,
		ResyncPeriod:     resync,
		DomainSuffix:     domainSuffix,
	})
	w := ctl.(*catalog)

	w.AppendServiceHandler(updateServicesHandler)
	w.AppendInstanceHandler(updateInstances)

	return w
}

func updateServicesHandler(*model.Service, model.Event) {}

func updateInstances(*model.ServiceInstance, model.Event) {}

func createEndpoints(w *catalog, name, namespace string, portNames, ips []string, t *testing.T) {
	eas := []v1.EndpointAddress{}
	for _, ip := range ips {
		eas = append(eas, v1.EndpointAddress{IP: ip})
	}

	eps := []v1.EndpointPort{}
	for _, name := range portNames {
		eps = append(eps, v1.EndpointPort{Name: name})
	}

	endpoint := &v1.Endpoints{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Subsets: []v1.EndpointSubset{{
			Addresses: eas,
			Ports:     eps,
		}},
	}
	if err := w.endpoints.informer.GetStore().Add(endpoint); err != nil {
		t.Errorf("failed to create endpoints %s in namespace %s (error %v)", name, namespace, err)
	}
}

func createService(w *catalog, name, namespace string, annotations map[string]string,
	ports []int32, selector map[string]string, t *testing.T) {
	svcPorts := []v1.ServicePort{}
	for _, p := range ports {
		svcPorts = append(svcPorts, v1.ServicePort{
			Name:     "test-port",
			Port:     p,
			Protocol: "http",
		})
	}
	service := &v1.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "10.0.0.1",
			Ports:     svcPorts,
			Selector:  selector,
			Type:      v1.ServiceTypeClusterIP,
		},
	}
	if err := w.services.informer.GetStore().Add(service); err != nil {
		t.Errorf("Cannot create service %s in namespace %s (error: %v)", name, namespace, err)
	}
}

func addPods(t *testing.T, w *catalog, pods ...*v1.Pod) {
	for _, pod := range pods {
		if err := w.pods.informer.GetStore().Add(pod); err != nil {
			t.Errorf("Cannot create pod in namespace %s (error: %v)", pod.ObjectMeta.Namespace, err)
		}
	}
}

func generatePod(name, namespace, saName, node string, labels map[string]string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Labels:    labels,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			ServiceAccountName: saName,
			NodeName:           node,
		},
	}
}

/*
// source https://stackoverflow.com/a/31832326
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func createNamespace(t *testing.T, cl kubernetes.Interface) (string, error) {
	name := "my-test-" + randString(10)
	ns, err := cl.CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
	})
	if err != nil {
		return "", err
	}
	if ns.Name == "" {
		return "", errors.New("Namespace name not generated")
	}
	t.Log("Created namespace", ns.Name)
	return ns.Name, nil
}

func deleteNamespace(t *testing.T, cl kubernetes.Interface, ns string) {
	if ns != "" && ns != "default" {
		if err := cl.CoreV1().Namespaces().Delete(ns, &meta_v1.DeleteOptions{}); err != nil {
			t.Logf("Error deleting namespace: %v", err)
		}
		t.Log("Deleted namespace", ns)
	}
}
*/
