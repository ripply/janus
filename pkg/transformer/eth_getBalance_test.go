package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestGetBalanceRequestAccount(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960"`), []byte(`"123"`)}
	requestRPC, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare account
	account, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	qtumClient.Accounts = append(qtumClient.Accounts, account)

	//prepare responses
	fromHexAddressResponse := qtum.FromHexAddressResponse("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	getAddressBalanceResponse := qtum.GetAddressBalanceResponse{Balance: 10010000000000}
	err = mockedClientDoer.AddResponseWithRequestID(3, qtum.MethodGetAddressBalance, getAddressBalanceResponse)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Need getaccountinfo to return an account for unit test
	// if getaccountinfo returns nil
	// then address is contract, else account

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBalance{qtumClient}
	got, err := proxyEth.Request(requestRPC)
	if err != nil {
		t.Fatal(err)
	}

	want := string("0x91aa27e8400") //(100000+100)*10^8 == 10010000000000 == 0x91AA27E8400
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			requestRPC,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}

func TestGetBalanceRequestContract(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960"`), []byte(`"123"`)}
	requestRPC, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare account
	account, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	qtumClient.Accounts = append(qtumClient.Accounts, account)

	//prepare responses
	getAccountInfoResponse := qtum.GetAccountInfoResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		Balance: 12431243,
		// Storage json.RawMessage `json:"storage"`,
		// Code    string          `json:"code"`,
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, qtum.MethodGetAccountInfo, getAccountInfoResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBalance{qtumClient}
	got, err := proxyEth.Request(requestRPC)
	if err != nil {
		t.Fatal(err)
	}

	want := string("0xbdaf8b")
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			requestRPC,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}
