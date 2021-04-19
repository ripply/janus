package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/qtum"
)

func initializeProxyETHGetBlockByHash(qtumClient *qtum.Qtum) ETHProxy {
	return &ProxyETHGetBlockByHash{qtumClient}
}

func TestGetBlockByHashRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByHash,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockHexHash + `"`), []byte(`false`)},
		&getTransactionByHashResponse,
	)
}

func TestGetBlockByHashTransactionsRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByHash,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockHexHash + `"`), []byte(`true`)},
		&getTransactionByHashResponseWithTransactions,
	)
}
