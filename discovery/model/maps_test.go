package model

import (
	"reflect"
	"testing"
)

var testServices = []*Service{
	{
		Hostname: "abc.default.svc.cluster.local",
		Ports:    []*Port{{Name: "http", Port: 80, Protocol: ProtocolHTTP}},
	},
	{
		Hostname: "bob.default.svc.cluster.local",
		Address:  "10.0.1.1",
		Ports:    []*Port{{Name: "https", Port: 443, Protocol: ProtocolHTTPS}},
	},
	{
		Hostname: "cat.default.svc.cluster.local",
		Ports:    []*Port{{Name: "http", Port: 80, Protocol: ProtocolHTTP}},
	},
	{
		Hostname: "dog.default.svc.cluster.local",
		Address:  "10.0.1.2",
		Ports:    []*Port{{Name: "http", Port: 80, Protocol: ProtocolHTTP}, {Name: "grpc", Port: 90, Protocol: ProtocolGRPC}},
	},
	{
		Hostname: "foo.default.svc.cluster.local",
		Ports:    []*Port{{Name: "https", Port: 443, Protocol: ProtocolHTTPS}},
	},
}
var testServiceMap = NewServiceMap(testServices...)
var emptyServiceMap = NewServiceMap()

var testInstances = []*ServiceInstance{
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.1",
			Port:        80,
			ServicePort: testServices[0].Ports[0],
		},
		Service: testServices[0],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.2",
			Port:        8080,
			ServicePort: testServices[0].Ports[0],
		},
		Service: testServices[0],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.3",
			Port:        443,
			ServicePort: testServices[1].Ports[0],
		},
		Service: testServices[1],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.4",
			Port:        443,
			ServicePort: testServices[1].Ports[0],
		},
		Service: testServices[1],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.5",
			Port:        80,
			ServicePort: testServices[2].Ports[0],
		},
		Service: testServices[2],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.6",
			Port:        80,
			ServicePort: testServices[2].Ports[0],
		},
		Service: testServices[2],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.7",
			Port:        80,
			ServicePort: testServices[3].Ports[0],
		},
		Service: testServices[3],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.8",
			Port:        90,
			ServicePort: testServices[3].Ports[1],
		},
		Service: testServices[3],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.9",
			Port:        443,
			ServicePort: testServices[4].Ports[0],
		},
		Service: testServices[4],
	},
	{
		Endpoint: NetworkEndpoint{
			Address:     "192.168.1.10",
			Port:        443,
			ServicePort: testServices[4].Ports[0],
		},
		Service: testServices[4],
	},
}
var testInstanceMap = NewInstanceMap(testInstances...)
var emptyInstanceMap = NewInstanceMap()

func TestNewServiceMap(t *testing.T) {
	type args struct {
		items []*Service
	}
	tests := []struct {
		name string
		args args
		want ServiceMap
	}{
		{
			name: "empty map",
			args: args{items: make([]*Service, 0)},
			want: emptyServiceMap,
		},
		{
			name: "map with items",
			args: args{items: testServices},
			want: testServiceMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServiceMap(tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServiceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_Insert(t *testing.T) {
	type args struct {
		items []*Service
	}
	tests := []struct {
		name string
		m    ServiceMap
		args args
		want int
	}{
		{
			name: "empty",
			m:    NewServiceMap(),
			args: args{items: make([]*Service, 0)},
			want: 0,
		},
		{
			name: "insert 1 into empty map",
			m:    NewServiceMap(),
			args: args{items: []*Service{testServices[0]}},
			want: 1,
		},
		{
			name: "insert 1 into non-empty map where value already exists",
			m:    NewServiceMap(testServices[0]),
			args: args{items: []*Service{testServices[0]}},
			want: 1,
		},
		{
			name: "insert 1 into non-empty map",
			m:    NewServiceMap(testServices[0]),
			args: args{items: []*Service{testServices[1]}},
			want: 2,
		},
		{
			name: "insert multiple into non-empty map where 1 value already exists",
			m:    NewServiceMap(testServices[0]),
			args: args{items: []*Service{testServices[0], testServices[1]}},
			want: 2,
		},
		{
			name: "insert multiple into non-empty map",
			m:    NewServiceMap(testServices[0]),
			args: args{items: []*Service{testServices[1], testServices[2]}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Insert(tt.args.items...)
			if got := tt.m.Len(); got != tt.want {
				t.Errorf("ServiceMap.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_Delete(t *testing.T) {
	type args struct {
		items []*Service
	}
	tests := []struct {
		name string
		m    ServiceMap
		args args
		want int
	}{
		{
			name: "delete 1 from empty map",
			m:    NewServiceMap(),
			args: args{items: []*Service{testServices[0]}},
			want: 0,
		},
		{
			name: "delete 1 from non-empty map",
			m:    NewServiceMap(testServices[0]),
			args: args{items: []*Service{testServices[0]}},
			want: 0,
		},
		{
			name: "delete multiple from non-empty map where 1 value does not already exist",
			m:    NewServiceMap(testServices[0]),
			args: args{items: []*Service{testServices[0], testServices[1]}},
			want: 0,
		},
		{
			name: "delete multiple from non-empty map",
			m:    NewServiceMap(testServices...),
			args: args{items: []*Service{testServices[1], testServices[2]}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Delete(tt.args.items...)
			if got := tt.m.Len(); got != tt.want {
				t.Errorf("ServiceMap.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_Has(t *testing.T) {
	type args struct {
		item *Service
	}
	tests := []struct {
		name string
		m    ServiceMap
		args args
		want bool
	}{
		{
			name: "empty map",
			m:    emptyServiceMap,
			args: args{item: testServices[0]},
			want: false,
		},
		{
			name: "non-empty map without value",
			m:    NewServiceMap(testServices[1]),
			args: args{item: testServices[0]},
			want: false,
		},
		{
			name: "non-empty map with value",
			m:    testServiceMap,
			args: args{item: testServices[3]},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Has(tt.args.item); got != tt.want {
				t.Errorf("ServiceMap.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_HasAll(t *testing.T) {
	type args struct {
		items []*Service
	}
	tests := []struct {
		name string
		m    ServiceMap
		args args
		want bool
	}{
		{
			name: "empty map",
			m:    emptyServiceMap,
			args: args{items: []*Service{testServices[0]}},
			want: false,
		},
		{
			name: "non-empty map without value",
			m:    NewServiceMap(testServices[1]),
			args: args{items: []*Service{testServices[0]}},
			want: false,
		},
		{
			name: "non-empty map with all values",
			m:    testServiceMap,
			args: args{items: []*Service{testServices[0], testServices[1]}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.HasAll(tt.args.items...); got != tt.want {
				t.Errorf("ServiceMap.HasAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_HasAny(t *testing.T) {
	type args struct {
		items []*Service
	}
	tests := []struct {
		name string
		m    ServiceMap
		args args
		want bool
	}{
		{
			name: "empty map",
			m:    emptyServiceMap,
			args: args{items: []*Service{testServices[0]}},
			want: false,
		},
		{
			name: "non-empty map without value",
			m:    NewServiceMap(testServices[1]),
			args: args{items: []*Service{testServices[0]}},
			want: false,
		},
		{
			name: "non-empty map with one value",
			m:    NewServiceMap(testServices[1]),
			args: args{items: []*Service{testServices[0], testServices[1]}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.HasAny(tt.args.items...); got != tt.want {
				t.Errorf("ServiceMap.HasAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_List(t *testing.T) {
	tests := []struct {
		name string
		m    ServiceMap
		want []*Service
	}{
		{
			name: "empty map",
			m:    emptyServiceMap,
			want: make([]*Service, 0),
		},
		{
			name: "non-empty map with a value",
			m:    NewServiceMap(testServices[1]),
			want: []*Service{testServices[1]},
		},
		{
			name: "non-empty map with many values",
			m:    testServiceMap,
			want: testServices,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceMap.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceMap_Len(t *testing.T) {
	tests := []struct {
		name string
		m    ServiceMap
		want int
	}{
		{
			name: "empty map",
			m:    emptyServiceMap,
			want: 0,
		},
		{
			name: "non-empty map with a value",
			m:    NewServiceMap(testServices[1]),
			want: 1,
		},
		{
			name: "non-empty map with many values",
			m:    testServiceMap,
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Len(); got != tt.want {
				t.Errorf("ServiceMap.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInstanceMap(t *testing.T) {
	type args struct {
		items []*ServiceInstance
	}
	tests := []struct {
		name string
		args args
		want InstanceMap
	}{
		{
			name: "empty map",
			args: args{items: make([]*ServiceInstance, 0)},
			want: emptyInstanceMap,
		},
		{
			name: "map with items",
			args: args{items: testInstances},
			want: testInstanceMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInstanceMap(tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInstanceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_Insert(t *testing.T) {
	type args struct {
		items []*ServiceInstance
	}
	tests := []struct {
		name string
		m    InstanceMap
		args args
		want int
	}{
		{
			name: "empty",
			m:    NewInstanceMap(),
			args: args{items: make([]*ServiceInstance, 0)},
			want: 0,
		},
		{
			name: "insert 1 into empty map",
			m:    NewInstanceMap(),
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: 1,
		},
		{
			name: "insert 1 into non-empty map where value already exists",
			m:    NewInstanceMap(testInstances[0]),
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: 1,
		},
		{
			name: "insert 1 into non-empty map",
			m:    NewInstanceMap(testInstances[0]),
			args: args{items: []*ServiceInstance{testInstances[1]}},
			want: 2,
		},
		{
			name: "insert multiple into non-empty map where 1 value already exists",
			m:    NewInstanceMap(testInstances[0]),
			args: args{items: []*ServiceInstance{testInstances[0], testInstances[1]}},
			want: 2,
		},
		{
			name: "insert multiple into non-empty map",
			m:    NewInstanceMap(testInstances[0]),
			args: args{items: []*ServiceInstance{testInstances[1], testInstances[2]}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Insert(tt.args.items...)
			if got := tt.m.Len(); got != tt.want {
				t.Errorf("InstanceMap.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_Delete(t *testing.T) {
	type args struct {
		items []*ServiceInstance
	}
	tests := []struct {
		name string
		m    InstanceMap
		args args
		want int
	}{
		{
			name: "delete 1 from empty map",
			m:    NewInstanceMap(),
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: 0,
		},
		{
			name: "delete 1 from non-empty map",
			m:    NewInstanceMap(testInstances[0]),
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: 0,
		},
		{
			name: "delete multiple from non-empty map where 1 value does not already exist",
			m:    NewInstanceMap(testInstances[0]),
			args: args{items: []*ServiceInstance{testInstances[0], testInstances[1]}},
			want: 0,
		},
		{
			name: "delete multiple from non-empty map",
			m:    NewInstanceMap(testInstances...),
			args: args{items: []*ServiceInstance{testInstances[1], testInstances[2]}},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Delete(tt.args.items...)
			if got := tt.m.Len(); got != tt.want {
				t.Errorf("InstanceMap.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_Has(t *testing.T) {
	type args struct {
		item *ServiceInstance
	}
	tests := []struct {
		name string
		m    InstanceMap
		args args
		want bool
	}{
		{
			name: "empty map",
			m:    emptyInstanceMap,
			args: args{item: testInstances[0]},
			want: false,
		},
		{
			name: "non-empty map without value",
			m:    NewInstanceMap(testInstances[1]),
			args: args{item: testInstances[0]},
			want: false,
		},
		{
			name: "non-empty map with value",
			m:    testInstanceMap,
			args: args{item: testInstances[3]},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Has(tt.args.item); got != tt.want {
				t.Errorf("InstanceMap.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_HasAll(t *testing.T) {
	type args struct {
		items []*ServiceInstance
	}
	tests := []struct {
		name string
		m    InstanceMap
		args args
		want bool
	}{
		{
			name: "empty map",
			m:    emptyInstanceMap,
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: false,
		},
		{
			name: "non-empty map without value",
			m:    NewInstanceMap(testInstances[1]),
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: false,
		},
		{
			name: "non-empty map with all values",
			m:    testInstanceMap,
			args: args{items: []*ServiceInstance{testInstances[0], testInstances[1]}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.HasAll(tt.args.items...); got != tt.want {
				t.Errorf("InstanceMap.HasAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_HasAny(t *testing.T) {
	type args struct {
		items []*ServiceInstance
	}
	tests := []struct {
		name string
		m    InstanceMap
		args args
		want bool
	}{
		{
			name: "empty map",
			m:    emptyInstanceMap,
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: false,
		},
		{
			name: "non-empty map without value",
			m:    NewInstanceMap(testInstances[1]),
			args: args{items: []*ServiceInstance{testInstances[0]}},
			want: false,
		},
		{
			name: "non-empty map with one value",
			m:    NewInstanceMap(testInstances[1]),
			args: args{items: []*ServiceInstance{testInstances[0], testInstances[1]}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.HasAny(tt.args.items...); got != tt.want {
				t.Errorf("InstanceMap.HasAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_List(t *testing.T) {
	tests := []struct {
		name string
		m    InstanceMap
		want []*ServiceInstance
	}{
		{
			name: "empty map",
			m:    emptyInstanceMap,
			want: make([]*ServiceInstance, 0),
		},
		{
			name: "non-empty map with a value",
			m:    NewInstanceMap(testInstances[0], testInstances[1]),
			want: []*ServiceInstance{testInstances[0], testInstances[1]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceMap.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_ListByService(t *testing.T) {
	type args struct {
		s *Service
		p *Port
	}
	tests := []struct {
		name string
		m    InstanceMap
		args args
		want []*ServiceInstance
	}{
		{
			name: "empty map",
			m:    emptyInstanceMap,
			args: args{s: &Service{}, p: &Port{}},
			want: make([]*ServiceInstance, 0),
		},
		{
			name: "non-empty map with a value",
			m:    testInstanceMap,
			args: args{s: testInstances[1].Service, p: testInstances[1].Endpoint.ServicePort},
			want: []*ServiceInstance{testInstances[0], testInstances[1]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ListByService(tt.args.s, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceMap.ListByService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceMap_Len(t *testing.T) {
	tests := []struct {
		name string
		m    InstanceMap
		want int
	}{
		{
			name: "empty map",
			m:    emptyInstanceMap,
			want: 0,
		},
		{
			name: "non-empty map with a value",
			m:    NewInstanceMap(testInstances[1]),
			want: 1,
		},
		{
			name: "non-empty map with many values",
			m:    testInstanceMap,
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Len(); got != tt.want {
				t.Errorf("InstanceMap.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
