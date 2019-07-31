package discovery

import (
	"sync"

	v1 "k8s.io/api/core/v1"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

// podCache is an eventually consistent pod cache
type podCache struct {
	rwMu sync.RWMutex
	cacheHandler

	// keys maintains stable pod IP to name key mapping
	// this allows us to retrieve the latest status by pod IP
	keys map[string]string
}

func newPodCache(ch cacheHandler) *podCache {
	out := &podCache{
		cacheHandler: ch,
		keys:         make(map[string]string),
	}

	ch.handler.Append(func(obj interface{}, ev model.Event) error {
		out.rwMu.Lock()
		defer out.rwMu.Unlock()

		pod := *obj.(*v1.Pod)
		log().Debugf("(Handler) Pod Details: %v", pod)
		ip := pod.Status.PodIP
		if ip == "" {
			key := keyFunc(pod.Name, pod.Namespace)
			for k, v := range out.keys {
				if v == key {
					ip = k
					break
				}
			}
		}

		log().Infof("(Handler) Pod IP: %v", ip)
		if len(ip) > 0 {
			switch ev {
			case model.EventAdd, model.EventUpdate:
				out.keys[ip] = keyFunc(pod.Name, pod.Namespace)
			case model.EventDelete:
				delete(out.keys, ip)
			}
		}
		return nil
	})
	return out
}

// getPodByIp returns the pod or nil if pod not found or an error occurred
func (pc *podCache) getPodByIP(addr string) (*v1.Pod, bool) {
	pc.rwMu.RLock()
	defer pc.rwMu.RUnlock()

	key, exists := pc.keys[addr]
	if !exists {
		return nil, false
	}
	item, exists, err := pc.informer.GetStore().GetByKey(key)
	if !exists || err != nil {
		return nil, false
	}
	return item.(*v1.Pod), true
}
