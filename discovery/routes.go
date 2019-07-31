package discovery

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

var (
	ignoredServices = []string{"apigw", "discovery", "auth-local"}
)

func (d *Discoverer) findAPIPaths() {
	d.mu.Lock()
	defer d.mu.Unlock()

	paths := make([]string, 0)

	for _, service := range d.data.services.List() {
		for _, port := range service.Ports {
			if port.Protocol.IsHTTP() && !isBlackListed(service.Hostname) {
				log().Infof("(findAPIPaths) name: %s port: %d", service.Hostname, port.Port)
				if api := d.getServiceAPI(service, port); api != "" {
					paths = append(paths, api)
				}
			}
		}
	}

	d.data.apis = paths
}

func (d *Discoverer) getServiceAPI(service *model.Service, port *model.Port) string {
	locations := make([]string, 0)

	for _, instance := range d.data.instances.ListByService(service, port) {
		for _, p := range d.services.ManagementPorts(instance.Endpoint.Address) {
			log().Infof("Management Port: %v", p)
			if p.Protocol.IsHTTP() {
				locations = append(locations, fmt.Sprintf("%s:%d", instance.Endpoint.Address, p.Port))
			}
		}
	}

	// add the service location as the last resort
	locations = append(locations, fmt.Sprintf("%s:%d", service.Hostname, port.Port))

	for _, loc := range locations {
		log().Infof("loading API paths from service location: %s", loc)
		if path := checkPath(d.client, loc); path != "" {
			return path
		}
	}

	return ""
}

func checkPath(client *http.Client, location string) string {
	u := &url.URL{
		Host:   location,
		Scheme: "http",
		Path:   "swagger.json",
	}
	log().Infof("apiLoader location: %s", u.String())

	resp, err := client.Get(u.String())
	if err != nil {
		log().Errorf("request failed: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log().Errorf("response error with code: %v", resp.StatusCode)
		return ""
	}

	return u.String()
}

func isBlackListed(name string) bool {
	for _, bl := range ignoredServices {
		if strings.HasPrefix(name, bl) {
			return true
		}
	}
	return false
}
