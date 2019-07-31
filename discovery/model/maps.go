package model

import (
	"reflect"
	"sort"
)

// ServiceMap is a set of Service using Hostname as the unique key
type ServiceMap map[string]*Service

// NewServiceMap creates a ServiceMap from a list of Service.
func NewServiceMap(items ...*Service) ServiceMap {
	m := ServiceMap{}
	m.Insert(items...)
	return m
}

// Insert adds items to the set.
func (m ServiceMap) Insert(items ...*Service) {
	for _, item := range items {
		// if the item does not exist already or the value is not the same
		// then add/update the entry. Otherwise do nothing.
		if !m.Has(item) || !reflect.DeepEqual(m[item.Hostname], item) {
			m[item.Hostname] = item
		}
	}
}

// Delete removes all items from the set.
func (m ServiceMap) Delete(items ...*Service) {
	for _, item := range items {
		delete(m, item.Hostname)
	}
}

// Has returns true if and only if item is contained in the set.
func (m ServiceMap) Has(item *Service) bool {
	_, exists := m[item.Hostname]
	return exists
}

// HasAll returns true if and only if all items are contained in the set.
func (m ServiceMap) HasAll(items ...*Service) bool {
	for _, item := range items {
		if !m.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (m ServiceMap) HasAny(items ...*Service) bool {
	for _, item := range items {
		if m.Has(item) {
			return true
		}
	}
	return false
}

// List returns the contents as a sorted Service slice.
// sorted by service.Hostname
func (m ServiceMap) List() []*Service {
	services := make([]*Service, 0, len(m))
	for _, v := range m {
		services = append(services, v)
	}

	// sort by hostnames
	sort.SliceStable(services, func(i, j int) bool { return services[i].Hostname < services[j].Hostname })

	return services
}

// Len returns the size of the set.
func (m ServiceMap) Len() int {
	return len(m)
}

// InstanceMap is a set of ServiceInstance using Hostname as the unique key
type InstanceMap map[string]*ServiceInstance

// NewInstanceMap creates a InstanceMap from a list of ServiceInstance.
func NewInstanceMap(items ...*ServiceInstance) InstanceMap {
	m := InstanceMap{}
	m.Insert(items...)
	return m
}

// Insert adds items to the set.
func (m InstanceMap) Insert(items ...*ServiceInstance) {
	for _, item := range items {
		// if the item does not exist already or the value is not the same
		// then add/update the entry. Otherwise do nothing.
		if !m.Has(item) || !reflect.DeepEqual(m[item.Key()], item) {
			m[item.Key()] = item
		}
	}
}

// Delete removes all items from the set.
func (m InstanceMap) Delete(items ...*ServiceInstance) {
	for _, item := range items {
		delete(m, item.Key())
	}
}

// Has returns true if and only if item is contained in the set.
func (m InstanceMap) Has(item *ServiceInstance) bool {
	_, exists := m[item.Key()]
	return exists
}

// HasAll returns true if and only if all items are contained in the set.
func (m InstanceMap) HasAll(items ...*ServiceInstance) bool {
	for _, item := range items {
		if !m.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (m InstanceMap) HasAny(items ...*ServiceInstance) bool {
	for _, item := range items {
		if m.Has(item) {
			return true
		}
	}
	return false
}

// List returns the contents as a sorted ServiceInstance slice.
// sorted by service.Hostname
func (m InstanceMap) List() []*ServiceInstance {
	instances := make([]*ServiceInstance, 0, len(m))
	for _, v := range m {
		instances = append(instances, v)
	}

	// sort by hostname/ip/port
	sort.SliceStable(instances, func(i, j int) bool {
		return instances[i].Service.Hostname < instances[j].Service.Hostname ||
			(instances[i].Service.Hostname == instances[j].Service.Hostname &&
				instances[i].Endpoint.Port < instances[j].Endpoint.Port)
	})

	return instances
}

// ListByService returns the a slice of ServiceInstance where the associated Service
// matches the specified Service and ServicePort matches the specified Port.
func (m InstanceMap) ListByService(s *Service, p *Port) []*ServiceInstance {
	instances := make([]*ServiceInstance, 0)
	for _, v := range m {
		if v.Service.Hostname == s.Hostname && v.Endpoint.ServicePort.Port == p.Port {
			instances = append(instances, v)
		}
	}

	// sort by hostname/ip/port
	sort.SliceStable(instances, func(i, j int) bool {
		return instances[i].Service.Hostname < instances[j].Service.Hostname ||
			(instances[i].Service.Hostname == instances[j].Service.Hostname &&
				instances[i].Endpoint.Port < instances[j].Endpoint.Port)
	})

	return instances
}

// Len returns the size of the set.
func (m InstanceMap) Len() int {
	return len(m)
}
