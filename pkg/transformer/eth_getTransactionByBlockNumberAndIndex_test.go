package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/qtum"
)

func initializeProxyETHGetTransactionByBlockNumberAndIndex(qtumClient *qtum.Qtum) ETHProxy {
	return &ProxyETHGetTransactionByBlockNumberAndIndex{qtumClient}
}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	TestETHProxyRequest(
		t,
		initializeProxyETHGetTransactionByBlockNumberAndIndex,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockNumber + `"`), []byte(`"0x0"`)},
		getTransactionByHashResponseData,
	)
}
