package transformer

import (
	"encoding/json"
	"testing"
)

func TestGetUncleByBlockHashAndIndexReturnsNil(t *testing.T) {
	// request body doesn't matter, there is no QTUM object to proxy calls to
	requestParams := []json.RawMessage{}
	request, err := PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ETHGetUncleByBlockHashAndIndex{}
	got, err := proxyEth.Request(request, nil)
	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Errorf(
			"error\ninput: %s\nwant: nil\ngot: %s",
			*request,
			got,
		)
	}
}
