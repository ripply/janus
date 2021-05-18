package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func initializeProxyETHGetBlockByNumber(qtumClient *qtum.Qtum) ETHProxy {
	return &ProxyETHGetBlockByNumber{qtumClient}
}

func TestGetBlockByNumberRequest(t *testing.T) {
	TestETHProxyRequest(
		t,
		initializeProxyETHGetBlockByNumber,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockNumber + `"`), []byte(`false`)},
		&getTransactionByHashResponse,
	)
}

func TestGetBlockByNumberWithTransactionsRequest(t *testing.T) {
	TestETHProxyRequest(
		t,
		initializeProxyETHGetBlockByNumber,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockNumber + `"`), []byte(`true`)},
		&getTransactionByHashResponseWithTransactions,
	)
}

func TestGetBlockByNumberUnknownBlockRequest(t *testing.T) {
	requestParams := []json.RawMessage{[]byte(`"` + getTransactionByHashBlockNumber + `"`), []byte(`true`)}
	request, err := PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := NewDoerMappedMock()
	qtumClient, err := CreateMockedClient(mockedClientDoer)

	unknownBlockResponse := qtum.GetErrorResponse(qtum.ErrInvalidParameter)
	err = mockedClientDoer.AddError(qtum.MethodGetBlockHash, unknownBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBlockByNumber{qtumClient}
	got, err := proxyEth.Request(request, nil)
	if err != nil {
		t.Fatal(err)
	}

	if got != (*eth.GetBlockByNumberResponse)(nil) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			string("nil"),
			string(MustMarshalIndent(got, "", "  ")),
		)
	}
}
