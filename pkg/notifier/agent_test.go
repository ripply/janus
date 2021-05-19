package notifier

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestAgentAddSubscription(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedSubscriptionID := "0x08e2af779d38a09e4c11442d9de22413"
	want := `{"subscription":"` + expectedSubscriptionID + `","result":{"address":"0x0000000000000000000000000000000000000000","blockHash":"0xbba11e1bacc69ba535d478cf1f2e542da3735a517b0b8eebaf7e6bb25eeb48c5","blockNumber":"0xf8f","data":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","topics":["0xtopic1"],"transactionHash":"0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5","transactionIndex":"0x2"}}`

	doer := internal.NewDoerMappedMock()
	topic1 := "topic1"

	doer.AddResponse(qtum.MethodWaitForLogs, qtum.WaitForLogsResponse{
		Entries: []qtum.TransactionReceipt{
			internal.QtumTransactionReceipt([]qtum.Log{
				{
					Address: internal.QtumTransactionReceipt(nil).ContractAddress,
					Topics:  []string{topic1},
					Data:    "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				},
			}),
		},
		Count:     1,
		NextBlock: internal.QtumTransactionReceipt(nil).BlockNumber + 1,
	})

	mockedClient, err := internal.CreateMockedClient(doer)
	if err != nil {
		t.Fatal(err)
	}
	agent := NewAgent(ctx, mockedClient)

	notifierContext, cancelNotifierContext := context.WithCancel(ctx)

	sentValuesChannel := make(chan interface{}, 10)
	sentValuesMutex := sync.Mutex{}
	sentValues := []interface{}{}

	send := func(v interface{}) error {
		sentValuesMutex.Lock()
		defer sentValuesMutex.Unlock()
		sentValues = append(sentValues, v)
		defer func() { sentValuesChannel <- nil }()
		return nil
	}

	notifier := NewNotifier(notifierContext, cancelNotifierContext, send, log.NewLogfmtLogger(os.Stdout))

	id, err := agent.NewSubscription(notifier, &eth.EthSubscriptionRequest{
		Method: "logs",
		Params: &eth.EthLogSubscriptionParameter{
			Address: eth.ETHAddress{},
			Topics: []interface{}{
				topic1,
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Fatal("Empty subscription id")
	}

	if len(id) != 34 {
		t.Fatalf("Subscription id incorrect length: %s", id)
	}

	select {
	case <-time.After(150 * time.Millisecond):
		// good
	case <-sentValuesChannel:
		t.Fatal("Received subscription result before sending subscription id to client")
	}

	notifier.ResponseSent()

	select {
	case <-sentValuesChannel:
		gotBytes := sentValues[0].([]byte)
		got := string(gotBytes)
		var receivedEthSubscription eth.EthSubscription
		err = json.Unmarshal(gotBytes, &receivedEthSubscription)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %s: %s", got, err)
		}
		receivedEthSubscription.SubscriptionID = expectedSubscriptionID
		gotBytes, err = json.Marshal(receivedEthSubscription)
		if err != nil {
			t.Fatalf("Failed to marshal: %s", err)
		}
		got = string(gotBytes)
		if got != want {
			t.Fatalf(
				"logs subscription error\nwant: %s\ngot: %s",
				want,
				got,
			)
		}
	case <-time.After(350 * time.Millisecond):
		t.Fatalf("Timed out waiting for subscription")
	}
}
