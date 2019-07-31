package discovery

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

var defaultTestClient = newTestClient(func(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: 200,
		// Send response to be tested
		Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
})

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func copyFile(src, dst string) {
	// Read all content of src to data
	data, _ := ioutil.ReadFile(src)
	// Write data to dst
	_ = ioutil.WriteFile(dst, data, 0644)
}

func copyKubeToken() {
	_ = os.MkdirAll("/var/run/secrets/kubernetes.io/serviceaccount", 0755)
	copyFile("testdata/token", "/var/run/secrets/kubernetes.io/serviceaccount/token")
}

func copyKubeCert() {
	_ = os.MkdirAll("/var/run/secrets/kubernetes.io/serviceaccount", 0755)
	copyFile("testdata/ca.crt", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
}

func genServerContent(f func() string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reply := f()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, reply)
	}))

	return srv
}

func genServerStatus(f func() int) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		status := f()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = fmt.Fprintln(w, "")
	}))

	return srv
}

func url2host(uri string) string {
	u, _ := url.Parse(uri)
	return u.Host
}

var testServices = []*model.Service{
	{
		Hostname: "abc.default.svc.cluster.local",
		Ports:    []*model.Port{{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}},
	},
	{
		Hostname: "bob.default.svc.cluster.local",
		Address:  "10.0.1.1",
		Ports:    []*model.Port{{Name: "https", Port: 443, Protocol: model.ProtocolHTTPS}},
	},
	{
		Hostname: "cat.default.svc.cluster.local",
		Ports:    []*model.Port{{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}},
	},
	{
		Hostname: "dog.default.svc.cluster.local",
		Address:  "10.0.1.2",
		Ports:    []*model.Port{{Name: "http", Port: 80, Protocol: model.ProtocolHTTP}, {Name: "grpc", Port: 90, Protocol: model.ProtocolGRPC}},
	},
	{
		Hostname: "foo.default.svc.cluster.local",
		Ports:    []*model.Port{{Name: "https", Port: 443, Protocol: model.ProtocolHTTPS}},
	},
}

/*
var sm = model.NewServiceMap(testServices...)
var testServiceMap = &sm
*/
var sm2 = model.NewServiceMap()
var emptyServiceMap = &sm2

var testInstances = []*model.ServiceInstance{
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.1",
			Port:        80,
			ServicePort: testServices[0].Ports[0],
		},
		Service: testServices[0],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.2",
			Port:        80,
			ServicePort: testServices[0].Ports[0],
		},
		Service: testServices[0],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.3",
			Port:        443,
			ServicePort: testServices[1].Ports[0],
		},
		Service: testServices[1],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.4",
			Port:        443,
			ServicePort: testServices[1].Ports[0],
		},
		Service: testServices[1],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.5",
			Port:        80,
			ServicePort: testServices[2].Ports[0],
		},
		Service: testServices[2],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.6",
			Port:        80,
			ServicePort: testServices[2].Ports[0],
		},
		Service: testServices[2],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.7",
			Port:        80,
			ServicePort: testServices[3].Ports[0],
		},
		Service: testServices[3],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.8",
			Port:        90,
			ServicePort: testServices[3].Ports[1],
		},
		Service: testServices[3],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.9",
			Port:        443,
			ServicePort: testServices[4].Ports[0],
		},
		Service: testServices[4],
	},
	{
		Endpoint: model.NetworkEndpoint{
			Address:     "192.168.1.10",
			Port:        443,
			ServicePort: testServices[4].Ports[0],
		},
		Service: testServices[4],
	},
}

/*
var im = model.NewInstanceMap(testInstances...)
var testInstanceMap = &im
*/
var im2 = model.NewInstanceMap()
var emptyInstanceMap = &im2

type fakeCatalog struct {
	wantErr bool
	wantNil bool
}

func (c *fakeCatalog) AppendServiceHandler(f func(*model.Service, model.Event)) {
}

func (c *fakeCatalog) AppendInstanceHandler(f func(*model.ServiceInstance, model.Event)) {
}

func (c *fakeCatalog) Run(stop <-chan struct{}) {
	<-stop
}

func (c *fakeCatalog) ManagementPorts(addr string) model.PortList {
	for _, inst := range testInstances {
		if inst.Endpoint.Address == addr {
			return inst.Service.Ports
		}
	}

	return nil
}
