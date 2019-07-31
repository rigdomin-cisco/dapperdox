package model

import (
	"strings"
	"testing"
)

var validServiceKeys = map[string]struct {
	service Service
}{
	"example-service1.default|grpc,http": {
		service: Service{
			Hostname: "example-service1.default",
			Ports:    []*Port{{Name: "http", Port: 80}, {Name: "grpc", Port: 90}}},
	},
	"my-service": {
		service: Service{
			Hostname: "my-service",
			Ports:    []*Port{{Name: "", Port: 80}}},
	},
	"svc.ns": {
		service: Service{
			Hostname: "svc.ns",
			Ports:    []*Port{{Name: "", Port: 80}}},
	},
	"svc": {
		service: Service{
			Hostname: "svc",
			Ports:    []*Port{{Name: "", Port: 80}}},
	},
	"svc|test": {
		service: Service{
			Hostname: "svc",
			Ports:    []*Port{{Name: "test", Port: 80}}},
	},
	"svc.default.svc.cluster.local|http-test": {
		service: Service{
			Hostname: "svc.default.svc.cluster.local",
			Ports:    []*Port{{Name: "http-test", Port: 80}}},
	},
}

// parseServiceKey is the inverse of the Service.String() method
func parseServiceKey(s string) (hostname string, ports PortList) {
	parts := strings.Split(s, "|")
	hostname = parts[0]

	var names []string
	if len(parts) > 1 {
		names = strings.Split(parts[1], ",")
	} else {
		names = []string{""}
	}

	for _, name := range names {
		ports = append(ports, &Port{Name: name})
	}
	return
}

func TestServiceString(t *testing.T) {
	for s, svc := range validServiceKeys {
		s1 := serviceKey(svc.service.Hostname, svc.service.Ports)
		if s1 != s {
			t.Errorf("ServiceKey => Got %s, expected %s", s1, s)
		}
		hostname, ports := parseServiceKey(s)
		if hostname != svc.service.Hostname {
			t.Errorf("ParseServiceKey => Got %s, expected %s for %s", hostname, svc.service.Hostname, s)
		}
		if len(ports) != len(svc.service.Ports) {
			t.Errorf("ParseServiceKey => Got %#v, expected %#v for %s", ports, svc.service.Ports, s)
		}
	}
}

func TestHTTPProtocol(t *testing.T) {
	if ProtocolUDP.IsHTTP() {
		t.Errorf("UDP is not HTTP protocol")
	}
	if !ProtocolGRPC.IsHTTP() {
		t.Errorf("gRPC is HTTP protocol")
	}
}

func TestConvertToProtocol(t *testing.T) {
	var testPairs = []struct {
		name string
		out  Protocol
	}{
		{"tcp", ProtocolTCP},
		{"http", ProtocolHTTP},
		{"HTTP", ProtocolHTTP},
		{"Http", ProtocolHTTP},
		{"https", ProtocolHTTPS},
		{"http2", ProtocolHTTP2},
		{"grpc", ProtocolGRPC},
		{"udp", ProtocolUDP},
		{"Mongo", ProtocolMongo},
		{"mongo", ProtocolMongo},
		{"MONGO", ProtocolMongo},
		{"Redis", ProtocolRedis},
		{"redis", ProtocolRedis},
		{"REDIS", ProtocolRedis},
		{"", ProtocolUnsupported},
		{"SMTP", ProtocolUnsupported},
	}

	for _, testPair := range testPairs {
		out := ConvertToProtocol(testPair.name)
		if out != testPair.out {
			t.Errorf("ConvertToProtocol(%q) => %q, want %q", testPair.name, out, testPair.out)
		}
	}
}
