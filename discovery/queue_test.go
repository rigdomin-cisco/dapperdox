package discovery

import (
	"errors"
	"testing"
	"time"

	"github.com/kenjones-cisco/dapperdox/discovery/model"
)

func Test_queue(t *testing.T) {
	q := newQueue(1 * time.Microsecond)
	stop := make(chan struct{})
	out := 0
	err := true
	add := func(obj interface{}, event model.Event) error {
		t.Logf("adding %d, error: %t", obj.(int), err)
		out += obj.(int)
		if !err {
			return nil
		}
		err = false
		return errors.New("intentional error")
	}
	check := func(obj interface{}, event model.Event) error {
		if out != 4 {
			t.Errorf("queue => %d, want %d", out, 4)
		}
		close(stop)
		return nil
	}
	go q.Run(stop)

	q.Push(task{handler: add, obj: 1})
	q.Push(task{handler: add, obj: 2})
	q.Push(task{handler: check, obj: 0})
	<-stop
}

func Test_chainedHandler(t *testing.T) {
	q := newQueue(1 * time.Microsecond)
	stop := make(chan struct{})
	out := 0
	f := func(i int) handler {
		return func(obj interface{}, event model.Event) error {
			out += i
			return nil
		}
	}
	h := chainHandler{
		funcs: []handler{f(1), f(2)},
	}
	go q.Run(stop)

	q.Push(task{handler: h.Apply, obj: 0})
	q.Push(task{handler: func(obj interface{}, event model.Event) error {
		if out != 3 {
			t.Errorf("chainedHandler => %d, want %d", out, 3)
		}
		close(stop)
		return nil
	}, obj: 0})
}
