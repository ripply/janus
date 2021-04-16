package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestHashrateRequest(t *testing.T) {
	//preparing the request
	requestParams := []json.RawMessage{} //eth_hashrate has no params
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	getHashrateResponse := qtum.GetHashrateResponse{Difficulty: float64(457134700)}
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodGetStakingInfo, getHashrateResponse)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ProxyETHHashrate{qtumClient}
	got, err := proxyEth.Request(request)
	if err != nil {
		t.Fatal(err)
	}

	want := eth.HashrateResponse("0x1b3f526c")
	if !reflect.DeepEqual(got, &want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			*request,
			want,
			got,
		)
	}
}
