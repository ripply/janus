package transformer

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGasPriceRequest(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{}
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGasPrice{qtumClient}
	got, err := proxyEth.Request(request)
	if err != nil {
		t.Fatal(err)
	}

	want := string("0x28") //price is hardcoded inside the implement
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			want,
			got,
		)
	}
}
