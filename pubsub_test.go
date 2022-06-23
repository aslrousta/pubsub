package pubsub_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aslrousta/pubsub"
)

func TestPubSub(t *testing.T) {
	ps := pubsub.New[string, string]()
	wg := sync.WaitGroup{}

	totalCount := int32(0)
	for i := 0; i < 100; i++ {
		ch := make(chan string)
		ps.Subscribe("test", ch)

		wg.Add(1)
		go func() {
			count := int32(0)
			for range ch {
				if count++; count == 10 {
					break
				}
			}
			atomic.AddInt32(&totalCount, count)
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		ps.Publish("test", "message")
	}
	wg.Wait()

	if totalCount != 1000 {
		t.FailNow()
	}
}
