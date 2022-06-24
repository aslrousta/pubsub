package pubsub

import "sync"

// PubSub implements pub-sub pattern using channels.
type PubSub[T comparable, M any] struct {
	mutex         sync.RWMutex
	subscriptions map[T][]chan<- M
}

// New instantiates a new pub-sub.
func New[T comparable, M any]() *PubSub[T, M] {
	return &PubSub[T, M]{
		subscriptions: make(map[T][]chan<- M),
	}
}

// Publish sends msg to the topic.
func (ps *PubSub[T, M]) Publish(topic T, msg M) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if _, ok := ps.subscriptions[topic]; !ok {
		return
	}

	for _, ch := range ps.subscriptions[topic] {
		go func(ch chan<- M) {
			ch <- msg
		}(ch)
	}
}

// Subscribe subscribes to receive messages in the topic on ch.
func (ps *PubSub[T, M]) Subscribe(topic T, ch chan<- M) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if subs, ok := ps.subscriptions[topic]; ok {
		ps.subscriptions[topic] = append(subs, ch)
	} else {
		ps.subscriptions[topic] = []chan<- M{ch}
	}
}

// Unsubscribe removes the subscription on the topic.
func (ps *PubSub[T, M]) Unubscribe(topic T, ch chan<- M) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if _, ok := ps.subscriptions[topic]; !ok {
		return
	}

	index := -1
	for i, c := range ps.subscriptions[topic] {
		if c == ch {
			index = i
			break
		}
	}
	if index >= 0 {
		if subs := ps.subscriptions[topic]; len(subs) == 1 {
			delete(ps.subscriptions, topic)
		} else if index == 0 {
			ps.subscriptions[topic] = subs[1:]
		} else if index == len(subs)-1 {
			ps.subscriptions[topic] = subs[:len(subs)-1]
		} else {
			ps.subscriptions[topic] = append(subs[:index], subs[index+1:]...)
		}
	}
}

// Close terminates all topics and closes the subscribed channels.
func (ps *PubSub[T, M]) Close() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	for _, subs := range ps.subscriptions {
		for _, ch := range subs {
			close(ch)
		}
	}
	ps.subscriptions = nil
}
