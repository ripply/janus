package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestGetStorageAtRequestWithNoLeadingZeros(t *testing.T) {
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
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")] = make(map[string]string)
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")][leftPadStringWithZerosTo64Bytes(index)] = value
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

func TestGetStorageAtRequestWithLeadingZeros(t *testing.T) {
	index := leftPadStringWithZerosTo64Bytes("abcde")
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
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")] = make(map[string]string)
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")][leftPadStringWithZerosTo64Bytes(index)] = value
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
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")] = make(map[string]string)
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")][leftPadStringWithZerosTo64Bytes(index)] = value
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

func TestLeftPadStringWithZerosTo64Bytes(t *testing.T) {
	tests := make(map[string]string)

	tests["1"] = "0000000000000000000000000000000000000000000000000000000000000001"
	tests["01"] = "0000000000000000000000000000000000000000000000000000000000000001"
	tests["001"] = "0000000000000000000000000000000000000000000000000000000000000001"
	tests["1001"] = "0000000000000000000000000000000000000000000000000000000000001001"
	tests["0000000000000000000000000000000000000000000000000000000000001001"] = "0000000000000000000000000000000000000000000000000000000000001001"
	tests["1111111111111111111111111111111111111111111111111111111111111111"] = "1111111111111111111111111111111111111111111111111111111111111111"
	tests["21111111111111111111111111111111111111111111111111111111111111111"] = "21111111111111111111111111111111111111111111111111111111111111111"

	for input, expected := range tests {
		result := leftPadStringWithZerosTo64Bytes(input)
		if result != expected {
			t.Errorf(
				"error\ninput: %s\nwant: %s\ngot: %s",
				input,
				expected,
				result,
			)
		}
	}
}
