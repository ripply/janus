package notifier

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/conversion"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func NewAgent(qtum *qtum.Qtum) *Agent {
	agent := &Agent{
		qtum:          qtum,
		mutex:         sync.RWMutex{},
		running:       false,
		stop:          make(chan interface{}),
		newHeads:      newSubscriptionRegistry(),
		logs:          newSubscriptionRegistry(),
		newPendingTxs: newSubscriptionRegistry(),
		syncing:       newSubscriptionRegistry(),
	}

	go agent.run()
	return agent
}

type subscriptionRegistry struct {
	mutex             sync.RWMutex
	subscriptionCount int
	subscriptions     map[string]*subscriptionInformation
}

func newSubscriptionRegistry() *subscriptionRegistry {
	return &subscriptionRegistry{
		mutex:             sync.RWMutex{},
		subscriptionCount: 0,
		subscriptions:     make(map[string]*subscriptionInformation),
	}
}

type subscriptionInformation struct {
	*Subscription
	params     *eth.EthSubscriptionRequest
	mutex      sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	running    bool
	qtum       *qtum.Qtum
}

func (s *subscriptionInformation) run() {
	if s.params == nil {
		return
	}

	if strings.ToLower(s.params.Method) != "logs" {
		return
	}

	s.mutex.Lock()
	if s.running {
		s.mutex.Unlock()
		return
	}
	s.running = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.running = false
	}()

	nextBlock := 0
	qtumTopics, err := eth.TranslateTopics(s.params.Params.Topics)
	if err != nil {
		s.qtum.GetDebugLogger().Log("msg", "Error translating logs topics", "error", err)
		return
	}
	req := &qtum.WaitForLogsRequest{
		FromBlock: nextBlock,
		ToBlock:   "latest",
		Filter: qtum.WaitForLogsFilter{
			Topics: &qtumTopics,
		},
	}

	failures := 0
	for {
		req.FromBlock = nextBlock
		resp, err := s.qtum.WaitForLogsWithContext(s.ctx, req)
		if err == nil {
			for _, qtumLog := range resp.Entries {
				ethLogs := conversion.ExtractETHLogsFromTransactionReceipt(&qtumLog)
				for _, ethLog := range ethLogs {
					s.notifier.Send(&eth.EthSubscription{
						SubscriptionID: s.Subscription.id,
						Result:         ethLog,
					})
				}
			}
			failures = 0
		} else {
			// error occurred
			s.qtum.GetDebugLogger().Log("subscriptionId", s.id, "err", err)
			failures = failures + 1
		}

		done := s.ctx.Done()

		select {
		case <-done:
			// err is wrapped so we can't detect (err == context.Cancelled)
			s.qtum.GetDebugLogger().Log("subscriptionId", s.id, "msg", "context closed, dropping subscription")
			return
		default:
		}

		backoffTime := getBackoff(failures, 0, 15*time.Second)

		if backoffTime > 0 {
			s.qtum.GetDebugLogger().Log("subscriptionId", s.id, "msg", fmt.Sprintf("backing off for %d miliseconds", backoffTime/time.Millisecond))
		}

		select {
		case <-done:
			return
		case <-time.After(backoffTime):
			// ok, try again
		}
	}
}

func getBackoff(count int, min time.Duration, max time.Duration) time.Duration {
	maxFailures := 10
	if count == 0 {
		return min
	}

	if count > maxFailures {
		return max
	}

	return ((max - min) / time.Duration(maxFailures)) * time.Duration(count)
}

type Agent struct {
	qtum          *qtum.Qtum
	mutex         sync.RWMutex
	running       bool
	stop          chan interface{}
	newHeads      *subscriptionRegistry
	logs          *subscriptionRegistry
	newPendingTxs *subscriptionRegistry
	syncing       *subscriptionRegistry
}

func (a *Agent) Stop() {
	a.mutex.Lock()
	a.lockAllRegistries(false)
	defer a.unlockAllRegistries(false)
	defer a.mutex.Unlock()

	closeSubscriptionRegistry(a.newHeads)
	closeSubscriptionRegistry(a.logs)
	closeSubscriptionRegistry(a.newPendingTxs)
	closeSubscriptionRegistry(a.syncing)
}

func closeSubscriptionRegistry(registry *subscriptionRegistry) {
	for _, sub := range registry.subscriptions {
		sub.cancelFunc()
	}
}

func (a *Agent) lockAllRegistries(readOnly bool) {
	if readOnly {
		a.newHeads.mutex.RLock()
		a.logs.mutex.RLock()
		a.newPendingTxs.mutex.RLock()
		a.syncing.mutex.RLock()
	} else {
		a.newHeads.mutex.Lock()
		a.logs.mutex.Lock()
		a.newPendingTxs.mutex.Lock()
		a.syncing.mutex.Lock()
	}
}

func (a *Agent) unlockAllRegistries(readOnly bool) {
	if readOnly {
		a.newHeads.mutex.RUnlock()
		a.logs.mutex.RUnlock()
		a.newPendingTxs.mutex.RUnlock()
		a.syncing.mutex.RUnlock()
	} else {
		a.newHeads.mutex.Unlock()
		a.logs.mutex.Unlock()
		a.newPendingTxs.mutex.Unlock()
		a.syncing.mutex.Unlock()
	}
}

func (a *Agent) subscriptionCount(acquireLocks bool) int {
	if acquireLocks {
		a.lockAllRegistries(true)
		defer a.unlockAllRegistries(true)
	}

	return a.newHeads.subscriptionCount +
		a.logs.subscriptionCount +
		a.newPendingTxs.subscriptionCount +
		a.syncing.subscriptionCount
}

func (a *Agent) unsubscribe(id string) {
	removeSubscription(id, a.newHeads)
	removeSubscription(id, a.logs)
	removeSubscription(id, a.newPendingTxs)
	removeSubscription(id, a.syncing)
}

func addSubscription(subscription *subscriptionInformation, registry *subscriptionRegistry) {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	_, collision := registry.subscriptions[subscription.id]
	registry.subscriptions[subscription.id] = subscription
	if !collision {
		registry.subscriptionCount = registry.subscriptionCount + 1
	}

	go subscription.run()
}

func removeSubscription(id string, registry *subscriptionRegistry) {
	registry.mutex.RLock()
	sub, exists := registry.subscriptions[id]
	registry.mutex.RUnlock()
	if exists {
		registry.mutex.Lock()
		_, exists = registry.subscriptions[id]
		if exists {
			delete(registry.subscriptions, id)
			registry.subscriptionCount = registry.subscriptionCount - 1
		}
		registry.mutex.Unlock()
	}

	if sub != nil {
		sub.cancelFunc()
	}
}

func (a *Agent) NewSubscription(notifier *Notifier, params *eth.EthSubscriptionRequest) (string, error) {
	subscription, err := notifier.Subscribe(a.unsubscribe)
	if err != nil {
		return "", err
	}

	wrappedContext, cancel := context.WithCancel(notifier.Context())

	wrappedSubscription := &subscriptionInformation{
		subscription,
		params,
		sync.RWMutex{},
		wrappedContext,
		cancel,
		false,
		a.qtum,
	}

	switch strings.ToLower(params.Method) {
	case "logs":
		addSubscription(wrappedSubscription, a.logs)
	case "newheads":
		addSubscription(wrappedSubscription, a.newHeads)
	case "newpendingtransactions":
		addSubscription(wrappedSubscription, a.newPendingTxs)
	case "syncing":
		addSubscription(wrappedSubscription, a.syncing)
	default:
		return "", errors.New(fmt.Sprintf("Unknown subscription type %s", params.Method))
	}

	a.mutex.RLock()
	if !a.running {
		// start processing subscriptions if nothing is running
		// only one routine will run at once so if multiple startup they will exit so only one runs
		go a.run()
	}
	a.mutex.RUnlock()

	return subscription.id, nil
}

func (a *Agent) run() {
	a.mutex.Lock()
	if a.running {
		a.mutex.Unlock()
		return
	}
	a.running = true
	a.mutex.Unlock()

	defer func() {
		a.mutex.Lock()
		defer a.mutex.Unlock()

		a.running = false
	}()

	lastBlock := int64(0)

	draining := true
	for draining {
		select {
		case <-a.stop:
			// drain
		default:
			draining = false
		}
	}

	return

	for {
		// infinite loop while we have subscriptions
		subscriptionCount := a.subscriptionCount(true)
		if subscriptionCount == 0 {
			return
		}

		blockchainInfo, err := a.qtum.GetBlockChainInfo()
		if err != nil {
			latestBlock := blockchainInfo.Blocks
			if latestBlock > lastBlock {

			}
		}

		select {
		case <-time.After(10 * time.Second):
			// continue
		case <-a.stop:
			return
		}
	}
}
