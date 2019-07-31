package discovery

/*
import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

func Test_podCache(t *testing.T) {

	testCases := []struct {
		name         string
		pods         []*v1.Pod
		keys         map[string]string
		wantLabels   map[string]model.Labels
		wantNotFound bool
	}{
		{
			name: "Should find all addresses in the map",
			pods: []*v1.Pod{
				generatePod("pod1", "nsA", "", "", map[string]string{"app": "test-app"}),
				generatePod("pod2", "nsA", "", "", map[string]string{"app": "prod-app-1"}),
				generatePod("pod3", "nsB", "", "", map[string]string{"app": "prod-app-2"}),
			},
			keys: map[string]string{
				"128.0.0.1": "nsA/pod1",
				"128.0.0.2": "nsA/pod2",
				"128.0.0.3": "nsB/pod3",
			},
			wantLabels: map[string]model.Labels{
				"128.0.0.1": {"app": "test-app"},
				"128.0.0.2": {"app": "prod-app-1"},
				"128.0.0.3": {"app": "prod-app-2"},
			},
		},
		{
			name:         "Should fail if addr not in keys",
			wantLabels:   map[string]model.Labels{"128.0.0.1": nil},
			wantNotFound: true,
		},
		{
			name:         "Should fail if addr in keys but pod not in cache",
			wantLabels:   map[string]model.Labels{"128.0.0.1": nil},
			keys:         map[string]string{"128.0.0.1": "nsA/pod1"},
			wantNotFound: true,
		},
	}
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			clientSet := fake.NewSimpleClientset()
			ctl := newCatalog(clientSet, catalogOptions{
				WatchedNamespace: "default",
				ResyncPeriod:     resync,
				DomainSuffix:     domainSuffix,
			})
			w := ctl.(*catalog)

			// Populate podCache
			for _, pod := range c.pods {
				if err := w.pods.informer.GetStore().Add(pod); err != nil {
					t.Errorf("Cannot create %s in namespace %s (error: %v)", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, err)
				}
			}

			// Populate key
			w.pods.keys = c.keys

			// Verify podCache
			for addr, wantTag := range c.wantLabels {
				tag, found := w.pods.labelsByIP(addr)
				if !reflect.DeepEqual(wantTag, tag) {
					t.Errorf("Expected %v got %v", wantTag, tag)
				}
				if c.wantNotFound {
					if found {
						t.Error("Expected not found but was found")
					}
				}
			}

		})
	}

}
*/
