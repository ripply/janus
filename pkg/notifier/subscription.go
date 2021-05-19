package notifier

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/qtumproject/janus/pkg/conversion"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

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

	// this throttles QTUM api calls if waitforlogs is returning very quickly a lot
	limitToXApiCalls := 5
	inYSeconds := 10 * time.Second
	// if a QTUM API call returns quicker than this, we will wait until this time is reached
	// this prevents spamming the QTUM node too much
	minimumTimeBetweenCalls := 100 * time.Millisecond

	rolling := newRollingLimit(limitToXApiCalls)

	failures := 0
	for {
		req.FromBlock = nextBlock
		timeBeforeCall := time.Now()
		rolling.Push(timeBeforeCall)
		resp, err := s.qtum.WaitForLogsWithContext(s.ctx, req)
		timeAfterCall := time.Now()
		if err == nil {
			for _, qtumLog := range resp.Entries {
				ethLogs := conversion.ExtractETHLogsFromTransactionReceipt(&qtumLog)
				for _, ethLog := range ethLogs {
					s.qtum.GetDebugLogger().Log("subscriptionId", s.id, "msg", "notifying of logs")
					s.notifier.Send(&eth.EthSubscription{
						SubscriptionID: s.Subscription.id,
						Result:         ethLog,
					})
				}
			}
			oldest := rolling.Oldest()
			a := time.Now()
			if oldest != nil && a.Sub(*oldest) < inYSeconds {
				// too many request returning successfully too quickly, slow them down
				failures = failures + 1
			} else {
				failures = 0
			}
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

		timeCallTook := timeAfterCall.Sub(timeAfterCall)
		if timeCallTook < minimumTimeBetweenCalls {
			timeLeftUntilMinimumTimeBetweenCallsReached := minimumTimeBetweenCalls - timeCallTook
			backoffTime = time.Duration(math.Max(float64(backoffTime), float64(timeLeftUntilMinimumTimeBetweenCallsReached)))
		}

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

// implementes an array with a rolling index that returns the oldest inserted element
type rollingLimit struct {
	index int
	limit int
	times []*time.Time
}

func newRollingLimit(limit int) *rollingLimit {
	roll := &rollingLimit{
		index: 0,
		limit: limit,
		times: []*time.Time{},
	}

	for i := 0; i < limit; i++ {
		roll.times = append(roll.times, nil)
	}

	return roll
}

func (r *rollingLimit) oldest() int {
	return (r.index + 1) % r.limit
}

func (r *rollingLimit) newest() int {
	return r.index
}

func (r *rollingLimit) bump() int {
	r.index = r.oldest()
	return r.index
}

func (r *rollingLimit) Oldest() *time.Time {
	return r.times[r.oldest()]
}

func (r *rollingLimit) Push(t time.Time) {
	r.times[r.bump()] = &t
}
