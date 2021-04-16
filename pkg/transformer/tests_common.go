package transformer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

//copy of qtum.doer interface
type doer interface {
	Do(*http.Request) (*http.Response, error)
}

func newDoerMappedMock() *doerMappedMock {
	return &doerMappedMock{
		Responses: make(map[string][]byte),
	}
}

//type for mocking requests to client with request -> response mapping
type doerMappedMock struct {
	mutex     sync.Mutex
	latestId  int
	Responses map[string][]byte
}

func (d doerMappedMock) updateId(id int) {
	if id > d.latestId {
		d.latestId = id
	}
}

func (d doerMappedMock) Do(request *http.Request) (*http.Response, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	requestJSON, err := parseRequestFromBody(request)
	if err != nil {
		return nil, err
	}

	if d.Responses[requestJSON.Method] == nil {
		log.Printf("No mocked response for %s\n", requestJSON.Method)
	}

	responseWriter := ioutil.NopCloser(bytes.NewReader(d.Responses[requestJSON.Method]))
	return &http.Response{
		StatusCode: 200,
		Body:       responseWriter,
	}, nil
}

func prepareEthRPCRequest(id int, params []json.RawMessage) (*eth.JSONRPCRequest, error) {
	requestID, err := json.Marshal(1)
	if err != nil {
		return nil, err
	}

	paramsArrayRaw, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	requestRPC := eth.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_protocolVersion",
		ID:      requestID,
		Params:  paramsArrayRaw,
	}

	return &requestRPC, nil
}

func prepareRawResponse(requestID int, responseResult interface{}, responseError *eth.JSONRPCError) ([]byte, error) {
	requestIDRaw, err := json.Marshal(requestID)
	if err != nil {
		return nil, err
	}

	var responseResultRaw json.RawMessage
	if responseResult != nil {
		responseResultRaw, err = json.Marshal(responseResult)
		if err != nil {
			return nil, err
		}
	}

	responseRPC := &eth.JSONRPCResult{
		JSONRPC:   "2.0",
		RawResult: responseResultRaw,
		Error:     responseError,
		ID:        requestIDRaw,
	}

	responseRPCRaw, err := json.Marshal(responseRPC)

	return responseRPCRaw, err
}

func (d *doerMappedMock) AddRawResponse(requestType string, rawResponse []byte) {
	d.mutex.Lock()
	d.Responses[requestType] = rawResponse
	d.mutex.Unlock()
}

func (d *doerMappedMock) AddResponse(requestType string, responseResult interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	requestID := d.latestId + 1
	responseRaw, err := prepareRawResponse(requestID, responseResult, nil)
	if err != nil {
		return err
	}

	d.updateId(requestID)
	d.Responses[requestType] = responseRaw
	return nil
}

func (d *doerMappedMock) AddResponseWithRequestID(requestID int, requestType string, responseResult interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	responseRaw, err := prepareRawResponse(requestID, responseResult, nil)
	if err != nil {
		return err
	}

	d.updateId(requestID)
	d.Responses[requestType] = responseRaw
	return nil
}

func (d *doerMappedMock) AddError(requestType string, responseError *eth.JSONRPCError) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	requestID := d.latestId + 1
	responseRaw, err := prepareRawResponse(requestID, nil, responseError)
	if err != nil {
		return err
	}

	d.updateId(requestID)
	d.Responses[requestType] = responseRaw
	return nil
}

func (d *doerMappedMock) AddErrorWithRequestID(requestID int, requestType string, responseError *eth.JSONRPCError) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	responseRaw, err := prepareRawResponse(requestID, nil, responseError)
	if err != nil {
		return err
	}

	d.updateId(requestID)
	d.Responses[requestType] = responseRaw
	return nil
}

func parseRequestFromBody(request *http.Request) (*eth.JSONRPCRequest, error) {
	requestJSON := eth.JSONRPCRequest{}
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestBody, &requestJSON)
	if err != nil {
		return nil, err
	}

	return &requestJSON, err
}

func createMockedClient(doerInstance doer) (qtumClient *qtum.Qtum, err error) {
	logger := kitLog.NewLogfmtLogger(os.Stdout)
	if !isDebugEnvironmentVariableSet() {
		logger = level.NewFilter(logger, level.AllowWarn())
	}
	qtumJSONRPC, err := qtum.NewClient(
		true,
		"http://user:pass@mocked",
		qtum.SetDoer(doerInstance),
		qtum.SetDebug(isDebugEnvironmentVariableSet()),
		qtum.SetLogger(logger),
	)
	if err != nil {
		return
	}

	qtumClient, err = qtum.New(qtumJSONRPC, "test")
	return
}

func isDebugEnvironmentVariableSet() bool {
	return strings.ToLower(os.Getenv("DEBUG")) == "true"
}

func mustMarshalIndent(v interface{}, prefix, indent string) []byte {
	res, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return res
}
