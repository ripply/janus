package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestGetStorageAtRequest(t *testing.T) {
	index := "abcde"
	blockNumber := "0x1234"
	requestParams := []json.RawMessage{[]byte(`"` + getTransactionByHashBlockNumber + `"`), []byte(`"0x` + index + `"`), []byte(`"` + blockNumber + `"`)}
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)

	value := "0x012341231441234123412343211234abcde12342332100000223030004005000"

	getStorageResponse := qtum.GetStorageResponse{}
	getStorageResponse["12345"] = make(map[string]string)
	getStorageResponse["12345"][index] = value
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodGetStorage, getStorageResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetStorageAt{qtumClient}
	got, err := proxyEth.Request(request)
	if err != nil {
		t.Fatal(err)
	}

	want := eth.GetStorageResponse(value)

	if !reflect.DeepEqual(got, &want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}

func TestGetStorageAtUnknownFieldRequest(t *testing.T) {
	index := "abcde"
	blockNumber := "0x1234"
	requestParams := []json.RawMessage{[]byte(`"` + getTransactionByHashBlockNumber + `"`), []byte(`"0x1234"`), []byte(`"` + blockNumber + `"`)}
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)

	unknownValue := "0x0000000000000000000000000000000000000000000000000000000000000000"
	value := "0x012341231441234123412343211234abcde12342332100000223030004005000"

	getStorageResponse := qtum.GetStorageResponse{}
	getStorageResponse["12345"] = make(map[string]string)
	getStorageResponse["12345"][index] = value
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodGetStorage, getStorageResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetStorageAt{qtumClient}
	got, err := proxyEth.Request(request)
	if err != nil {
		t.Fatal(err)
	}

	want := eth.GetStorageResponse(unknownValue)

	if !reflect.DeepEqual(got, &want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}
