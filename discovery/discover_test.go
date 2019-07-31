package discovery

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/kenjones-cisco/dapperdox/config"
	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

func TestNewDiscoverer(t *testing.T) {
	viper.Set(config.DiscoverInterval, "5s")

	if _, err := NewDiscoverer(); err == nil {
		t.Errorf("NewDiscoverer() error = %v, wantErr %v", err, true)
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

	d, err := NewDiscoverer()
	if err != nil {
		t.Errorf("NewDiscoverer() error = %v, wantErr %v", err, true)
	}
	if d == nil {
		t.Errorf("NewDiscoverer() = %v, want not nil", d)
	}
}

func TestDiscoverer_run_nil_service(t *testing.T) {
	stop := make(chan struct{})

	d := &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: nil}
	go d.Run(stop)

	var once sync.Once
	once.Do(func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
	})
}

func TestDiscoverer_run_fake_service(t *testing.T) {
	stop := make(chan struct{})

	d := &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{}}
	go d.Run(stop)

	var once sync.Once
	once.Do(func() {
		time.Sleep(100 * time.Millisecond)
		close(stop)
	})
}

func TestDiscoverer_updateServices(t *testing.T) {
	// handle the initial run and run where data does not change
	d := &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{}, client: defaultTestClient}

	type args struct {
		s *model.Service
		e model.Event
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "initial run",
			args: args{s: testServices[0], e: model.EventAdd},
			want: []string{"http://abc.default.svc.cluster.local:80/swagger.json"},
		},
		{
			name: "next run with no data change",
			args: args{s: testServices[0], e: model.EventUpdate},
			want: []string{"http://abc.default.svc.cluster.local:80/swagger.json"},
		},
		{
			name: "next run with delete data",
			args: args{s: testServices[0], e: model.EventDelete},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.updateServices(tt.args.s, tt.args.e)
			got := d.APIList()
			t.Logf("APIs: %v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(updateServices) Discoverer.APIList() = %v, want %v", got, tt.want)
			}
		})
	}

	// trigger failure within discover()
	d = &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{}, client: defaultTestClient}
	d.updateServices(testServices[0], model.EventAdd)

	// trigger failure from services call
	d = &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{wantErr: true, wantNil: false}, client: defaultTestClient}
	d.updateServices(testServices[0], model.EventAdd)

	// trigger no data from services call
	d = &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{wantErr: false, wantNil: true}, client: defaultTestClient}
	d.updateServices(testServices[0], model.EventAdd)
}

func TestDiscoverer_updateInstances(t *testing.T) {
	// handle the initial run and run where data does not change
	d := &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{}, client: defaultTestClient}

	type args struct {
		s *model.ServiceInstance
		e model.Event
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "initial run",
			args: args{s: testInstances[0], e: model.EventAdd},
			want: []string{"http://192.168.1.1:80/swagger.json"},
		},
		{
			name: "next run with no data change",
			args: args{s: testInstances[0], e: model.EventUpdate},
			want: []string{"http://192.168.1.1:80/swagger.json"},
		},
		{
			name: "next run with delete data",
			args: args{s: testInstances[0], e: model.EventDelete},
			want: []string{"http://abc.default.svc.cluster.local:80/swagger.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.updateInstances(tt.args.s, tt.args.e)
			got := d.APIList()
			t.Logf("APIs: %v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("(updateInstances) Discoverer.APIList() = %v, want %v", got, tt.want)
			}
		})
	}

	// trigger failure within discover()
	d = &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{}, client: defaultTestClient}
	d.updateInstances(testInstances[0], model.EventAdd)

	// trigger failure from services call
	d = &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{wantErr: true, wantNil: false}, client: defaultTestClient}
	d.updateInstances(testInstances[0], model.EventAdd)

	// trigger no data from services call
	d = &Discoverer{data: state{services: emptyServiceMap, instances: emptyInstanceMap}, services: &fakeCatalog{wantErr: false, wantNil: true}, client: defaultTestClient}
	d.updateInstances(testInstances[0], model.EventAdd)
}
