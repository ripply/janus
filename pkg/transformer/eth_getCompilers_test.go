package transformer

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetCompilersReturnsEmptyArray(t *testing.T) {
	//preparing the request
	requestParams := []json.RawMessage{} //eth_getCompilers has no params
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ETHGetCompilers{}
	got, err := proxyEth.Request(request)
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%v", got) != "[]" {
		t.Errorf(
			"error\nwant: '[]'\ngot: '%v'",
			got,
		)
	}
}
