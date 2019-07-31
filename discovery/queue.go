package discovery

import (
	"sync"
	"time"

	"k8s.io/client-go/util/flowcontrol"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

// queue of work tickets processed using a rate-limiting loop
type queue interface {
	// Push a ticket
	Push(task)
	// Run the loop until a signal on the channel
	Run(<-chan struct{})
}

// handler specifies a function to apply on an object for a given event type
type handler func(obj interface{}, event model.Event) error

// task object for the event watchers; processes until handler succeeds
type task struct {
	handler handler
	obj     interface{}
	event   model.Event
}

type queueImpl struct {
	delay   time.Duration
	queue   []task
	lock    sync.Mutex
	closing bool
}

// newQueue instantiates a queue with a processing function
func newQueue(errorDelay time.Duration) queue {
	return &queueImpl{
		delay:   errorDelay,
		queue:   make([]task, 0),
		closing: false,
		lock:    sync.Mutex{},
	}
}

func (q *queueImpl) Push(item task) {
	q.lock.Lock()
	if !q.closing {
		q.queue = append(q.queue, item)
	}
	q.lock.Unlock()
}

func (q *queueImpl) Run(stop <-chan struct{}) {
	go func() {
		<-stop
		q.lock.Lock()
		q.closing = true
		q.lock.Unlock()
	}()

	// Throttle processing up to smoothed 10 qps with bursts up to 100 qps
	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(float32(10), 100)
	var item task
	for {
		rateLimiter.Accept()

		q.lock.Lock()
		if q.closing {
			q.lock.Unlock()
			return
		} else if len(q.queue) == 0 {
			q.lock.Unlock()
		} else {
			item, q.queue = q.queue[0], q.queue[1:]
			q.lock.Unlock()

			for {
				err := item.handler(item.obj, item.event)
				if err != nil {
					log().Infof("Work item failed (%v), repeating after delay %v", err, q.delay)
					time.Sleep(q.delay)
				} else {
					break
				}
			}
		}
	}
}

// chainHandler applies handlers in a sequence
type chainHandler struct {
	funcs []handler
}

// Apply is the handler function
func (ch *chainHandler) Apply(obj interface{}, event model.Event) error {
	for _, f := range ch.funcs {
		if err := f(obj, event); err != nil {
			return err
		}
	}
	return nil
}

// Append a handler as the last handler in the chain
func (ch *chainHandler) Append(h handler) {
	ch.funcs = append(ch.funcs, h)
}
