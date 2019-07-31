package model

import (
	"bytes"
	"sort"
	"strings"
)

// Service describes an k8s service (e.g., catalog.mystore.com:8080)
// Each service has a fully qualified domain name (FQDN) and one or more
// ports where the service is listening for connections. *Optionally*, a
// service can have a single load balancer/virtual IP address associated
// with it, such that the DNS queries for the FQDN resolves to the virtual
// IP address (a load balancer IP).
//
// E.g., in kubernetes, a service foo is associated with
// foo.default.svc.cluster.local hostname, has a virtual IP of 10.0.1.1 and
// listens on ports 80, 8080
type Service struct {
	// Hostname of the service, e.g. "catalog.mystore.com"
	Hostname string `json:"hostname"`

	// Address specifies the service IPv4 address of the load balancer
	Address string `json:"address,omitempty"`

	// Ports is the set of network ports where the service is listening for
	// connections
	Ports PortList `json:"ports"`

	// ExternalName is only set for external services and holds the external
	// service DNS name.  External services are name-based solution to represent
	// external service instances as a service inside the cluster.
	ExternalName string `json:"external,omitempty"`

	// LoadBalancingDisabled indicates that no load balancing should be done for this service.
	LoadBalancingDisabled bool `json:"-"`
}

// External predicate checks whether the service is external
func (s *Service) External() bool {
	return s.ExternalName != ""
}

// Key generates a unique string referencing service instances for a given port.
// The separator character must be exclusive to the regular expressions allowed in the
// service declaration.
func (s *Service) Key(port *Port) string {
	// TODO: check port is non nil and membership of port in service
	return serviceKey(s.Hostname, PortList{port})
}

// Port represents a network port where a service is listening for
// connections. The port should be annotated with the type of protocol
// used by the port.
type Port struct {
	// Name ascribes a human readable name for the port object. When a
	// service has multiple ports, the name field is mandatory
	Name string `json:"name,omitempty"`

	// Port number where the service can be reached. Does not necessarily
	// map to the corresponding port numbers for the instances behind the
	// service. See NetworkEndpoint definition below.
	Port int `json:"port"`

	// Protocol to be used for the port.
	Protocol Protocol `json:"protocol"`
}

// PortList is a set of ports
type PortList []*Port

// Get retrieves a port declaration by name
func (ports PortList) Get(name string) (*Port, bool) {
	for _, port := range ports {
		if port.Name == name {
			return port, true
		}
	}
	return nil, false
}

// Protocol defines network protocols for ports
type Protocol string

const (
	// ProtocolGRPC declares that the port carries gRPC traffic
	ProtocolGRPC Protocol = "GRPC"
	// ProtocolHTTPS declares that the port carries HTTPS traffic
	ProtocolHTTPS Protocol = "HTTPS"
	// ProtocolHTTP2 declares that the port carries HTTP/2 traffic
	ProtocolHTTP2 Protocol = "HTTP2"
	// ProtocolHTTP declares that the port carries HTTP/1.1 traffic.
	// Note that HTTP/1.0 or earlier may not be supported by the proxy.
	ProtocolHTTP Protocol = "HTTP"
	// ProtocolTCP declares the the port uses TCP.
	// This is the default protocol for a service port.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP declares that the port uses UDP.
	// Note that UDP protocol is not currently supported by the proxy.
	ProtocolUDP Protocol = "UDP"
	// ProtocolMongo declares that the port carries mongoDB traffic
	ProtocolMongo Protocol = "Mongo"
	// ProtocolRedis declares that the port carries redis traffic
	ProtocolRedis Protocol = "Redis"
	// ProtocolUnsupported - value to signify that the protocol is unsupported
	ProtocolUnsupported Protocol = "UnsupportedProtocol"
)

// ConvertToProtocol converts a case-insensitive protocol to Protocol
func ConvertToProtocol(protocolAsString string) Protocol {
	switch strings.ToLower(protocolAsString) {
	case "tcp":
		return ProtocolTCP
	case "udp":
		return ProtocolUDP
	case "grpc":
		return ProtocolGRPC
	case "http":
		return ProtocolHTTP
	case "http2":
		return ProtocolHTTP2
	case "https":
		return ProtocolHTTPS
	case "mongo":
		return ProtocolMongo
	case "redis":
		return ProtocolRedis
	}

	return ProtocolUnsupported
}

// IsHTTP is true for protocols that use HTTP as transport protocol
func (p Protocol) IsHTTP() bool {
	switch p {
	case ProtocolHTTP, ProtocolHTTP2, ProtocolGRPC:
		return true
	default:
		return false
	}
}

// NetworkEndpoint defines a network address (IP:port) associated with an instance of the
// service. A service has one or more instances each running in a
// container/VM/pod. If a service has multiple ports, then the same
// instance IP is expected to be listening on multiple ports (one per each
// service port). Note that the port associated with an instance does not
// have to be the same as the port associated with the service. Depending
// on the network setup (NAT, overlays), this could vary.
//
// For e.g., if catalog.mystore.com is accessible through port 80 and 8080,
// and it maps to an instance with IP 172.16.0.1, such that connections to
// port 80 are forwarded to port 55446, and connections to port 8080 are
// forwarded to port 33333,
//
// then internally, we have two two endpoint structs for the
// service catalog.mystore.com
//  --> 172.16.0.1:54546 (with ServicePort pointing to 80) and
//  --> 172.16.0.1:33333 (with ServicePort pointing to 8080)
type NetworkEndpoint struct {
	// Address of the network endpoint, typically an IPv4 address
	Address string `json:"ip_address,omitempty"`

	// Port number where this instance is listening for connections This
	// need not be the same as the port where the service is accessed.
	// e.g., catalog.mystore.com:8080 -> 172.16.0.1:55446
	Port int `json:"port"`

	// Port declaration from the service declaration This is the port for
	// the service associated with this instance (e.g.,
	// catalog.mystore.com)
	ServicePort *Port `json:"service_port"`
}

// ServiceInstance represents an individual instance of a specific version
// of a service. It binds a network endpoint (ip:port), the service
// description (which is oblivious to various versions) with this instance.
//
// For example, the set of service instances associated with catalog.mystore.com
// are modeled like this
//      --> NetworkEndpoint(172.16.0.1:8888), Service(catalog.myservice.com)
//      --> NetworkEndpoint(172.16.0.2:8888), Service(catalog.myservice.com)
//      --> NetworkEndpoint(172.16.0.3:8888), Service(catalog.myservice.com)
//      --> NetworkEndpoint(172.16.0.4:8888), Service(catalog.myservice.com)
type ServiceInstance struct {
	Endpoint NetworkEndpoint `json:"endpoint,omitempty"`
	Service  *Service        `json:"service,omitempty"`
}

// Key generates a unique string referencing the instance details
func (s *ServiceInstance) Key() string {
	port := &Port{
		Name:     s.Endpoint.ServicePort.Name,
		Port:     s.Endpoint.Port,
		Protocol: s.Endpoint.ServicePort.Protocol,
	}
	return serviceKey(s.Endpoint.Address, PortList{port})
}

// serviceKey generates a service key for a collection of ports
func serviceKey(hostname string, servicePorts PortList) string {
	// example: name.namespace|http
	var buffer bytes.Buffer
	_, _ = buffer.WriteString(hostname)

	if len(servicePorts) > 0 {
		ports := make([]string, 0)
		for i := range servicePorts {
			if servicePorts[i].Name == "" {
				continue
			}
			ports = append(ports, servicePorts[i].Name)
		}

		if len(ports) > 0 {
			_, _ = buffer.WriteString("|")
		}

		sort.Strings(ports)
		for i := range ports {
			if i > 0 {
				_, _ = buffer.WriteString(",")
			}
			_, _ = buffer.WriteString(ports[i])
		}
	}
	return buffer.String()
}
