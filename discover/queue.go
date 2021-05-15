package discover

import (
	"sync"
	"time"

	"k8s.io/client-go/util/flowcontrol"

	"github.com/kenjones-cisco/dapperdox/discover/models"
)

const (
	qpsRate  float32 = 10
	qpsBurst int     = 100
)

// Queue of work tickets processed using a rate-limiting loop.
type Queue interface {
	// Push a ticket
	Push(Task)
	// Run the loop until a signal on the channel
	Run(<-chan struct{})
}

// Handler specifies a function to apply on an object for a given event type.
type Handler func(obj interface{}, event models.Event) error

// Task object for the event watchers; processes until handler succeeds.
type Task struct {
	handler Handler
	obj     interface{}
	event   models.Event
}

// NewTask creates a task from a work item.
func NewTask(handler Handler, obj interface{}, event models.Event) Task {
	return Task{handler: handler, obj: obj, event: event}
}

type queueImpl struct {
	delay   time.Duration
	queue   []Task
	lock    sync.Mutex
	closing bool
}

// NewQueue instantiates a queue with a processing function.
func NewQueue(errorDelay time.Duration) Queue {
	return &queueImpl{
		delay:   errorDelay,
		queue:   make([]Task, 0),
		closing: false,
		lock:    sync.Mutex{},
	}
}

func (q *queueImpl) Push(item Task) {
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
	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(qpsRate, qpsBurst)

	var item Task

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
				if err := item.handler(item.obj, item.event); err != nil {
					log().Errorf("Work item failed (%v), repeating after delay %v", err, q.delay)

					time.Sleep(q.delay)
				} else {
					break
				}
			}
		}
	}
}

// ChainHandler applies handlers in a sequence.
type ChainHandler struct {
	funcs []Handler
}

// Apply is the handler function.
func (ch *ChainHandler) Apply(obj interface{}, event models.Event) error {
	for _, f := range ch.funcs {
		if err := f(obj, event); err != nil {
			return err
		}
	}

	return nil
}

// Append a handler as the last handler in the chain.
func (ch *ChainHandler) Append(h Handler) {
	ch.funcs = append(ch.funcs, h)
}
