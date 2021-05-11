package notifier

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

var UnsubSignal = new(struct{})

type UnsubscribeCallback func(string)

type Subscription struct {
	id          string
	once        sync.Once
	unsubscribe UnsubscribeCallback
	notifier    *Notifier
}

func NewSubscription(notifier *Notifier, callback UnsubscribeCallback) (*Subscription, error) {
	id, err := getRandomSubscriptionId()
	if err != nil {
		return nil, err
	}
	return &Subscription{
		id:          id,
		once:        sync.Once{},
		unsubscribe: callback,
		notifier:    notifier,
	}, nil
}

func getRandomSubscriptionId() (string, error) {
	var subid [16]byte
	n, _ := rand.Read(subid[:])
	if n != 16 {
		return "", errors.New("Unable to generate subscription id")
	}
	return "0x" + hex.EncodeToString(subid[:]), nil
}

func (s *Subscription) Unsubscribe() {
	s.once.Do(func() {
		s.unsubscribe(s.id)
	})
}

type Notifier struct {
	runMutex              sync.Mutex
	mutex                 sync.Mutex
	ctx                   context.Context
	close                 func()
	send                  func(interface{}) error
	logger                log.Logger
	queue                 chan interface{}
	subscriptionIdPending *chan interface{}
	subscriptionsFlushed  *chan interface{}
	subscriptions         map[string]*Subscription
}

func NewNotifier(ctx context.Context, close func(), send func(interface{}) error, logger log.Logger) *Notifier {
	pending := make(chan interface{})
	flushed := make(chan interface{})
	return &Notifier{
		runMutex:              sync.Mutex{},
		mutex:                 sync.Mutex{},
		ctx:                   ctx,
		close:                 close,
		send:                  send,
		logger:                log.WithPrefix(logger, "component", "notifier"),
		queue:                 make(chan interface{}),
		subscriptionIdPending: &pending,
		subscriptionsFlushed:  &flushed,
		subscriptions:         make(map[string]*Subscription),
	}
}

func (n *Notifier) Context() context.Context {
	return n.ctx
}

func (n *Notifier) Subscribe(unsubscribeCallback UnsubscribeCallback) (*Subscription, error) {
	sub, err := NewSubscription(n, unsubscribeCallback)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (n *Notifier) ResponseSent() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if n.subscriptionIdPending != nil {
		close(*n.subscriptionIdPending)
		n.subscriptionIdPending = nil
	}
}

func (n *Notifier) ResponseRequired() {
	pending := make(chan interface{})
	n.subscriptionIdPending = &pending
}

func (n *Notifier) Send(event interface{}) {
	n.queue <- event
}

func (n *Notifier) closeSubscriptionsFlushed() {
	if n.subscriptionsFlushed != nil {
		close(*n.subscriptionsFlushed)
		n.subscriptionsFlushed = nil
	}
}

func (n *Notifier) Run() {
	n.runMutex.Lock()
	defer n.runMutex.Unlock()

	log.With(level.Debug(n.logger)).Log("msg", "Entering notifier loop")

	defer func() {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		defer func() {
			log.With(level.Debug(n.logger)).Log("msg", "Notifier loop exited")
		}()

		n.close()
		close(n.queue)
		n.closeSubscriptionsFlushed()
		for _, sub := range n.subscriptions {
			sub.Unsubscribe()
		}
	}()

	for {
		select {
		case <-n.ctx.Done():
			return
		case event := <-n.queue:
			log.With(level.Debug(n.logger)).Log("msg", event)
			<-*n.subscriptionIdPending
			if event == UnsubSignal {
				n.mutex.Lock()
				n.closeSubscriptionsFlushed()
				n.mutex.Unlock()
			} else {
				bytes, err := json.Marshal(event)
				if err != nil {
					panic(err)
				}
				err = n.send(bytes)
				if err != nil {
					// write failure, close connection and unsubscribe
					log.With(level.Debug(n.logger)).Log("msg", "Error writing response to websocket, closing it", "err", err)
					n.close()
					return
				}
			}
		}
	}
}
