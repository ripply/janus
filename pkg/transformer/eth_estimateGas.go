package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHEstimateGas struct {
	*ProxyETHCall
}

func (p *ProxyETHEstimateGas) Method() string {
	return "eth_estimateGas"
}

func (p *ProxyETHEstimateGas) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var ethreq eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &ethreq); err != nil {
		return nil, err
	}

	// eth req -> qtum req
	qtumreq, err := p.ToRequest(&ethreq)
	if err != nil {
		return nil, err
	}

	// qtum [code: -5] Incorrect address occurs here
	qtumresp, err := p.CallContract(qtumreq)
	if err != nil {
		return nil, err
	}

	return p.toResp(qtumresp)
}

func (p *ProxyETHEstimateGas) toResp(qtumresp *qtum.CallContractResponse) (*eth.EstimateGasResponse, error) {
	gas := eth.EstimateGasResponse(hexutil.EncodeUint64(uint64(qtumresp.ExecutionResult.GasUsed)))
	p.GetDebugLogger().Log(p.Method(), gas)
	return &gas, nil
}
