package transformer

import (
	"math/big"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHCall implements ETHProxy
type ProxyETHCall struct {
	*qtum.Qtum
}

func (p *ProxyETHCall) Method() string {
	return "eth_call"
}

func (p *ProxyETHCall) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHCall) request(ethreq *eth.CallRequest) (interface{}, error) {
	// eth req -> qtum req
	qtumreq, err := p.ToRequest(ethreq)
	if err != nil {
		return nil, err
	}

	qtumresp, err := p.CallContract(qtumreq)
	if err != nil {
		return nil, err
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*qtum.CallContractRequest, error) {
	from := ethreq.From
	var err error
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, err
		}
	}

	var gasLimit *big.Int
	if ethreq.Gas != nil {
		gasLimit = ethreq.Gas.Int
	}

	return &qtum.CallContractRequest{
		To:       ethreq.To,
		From:     from,
		Data:     ethreq.Data,
		GasLimit: gasLimit,
	}, nil
}

/**
// TODO: Handle reverts
// https://ethereum.stackexchange.com/questions/84545/how-to-get-reason-revert-using-web3-eth-call
// https://web3js.readthedocs.io/en/v1.2.8/web3-eth.html#handlerevert
{
  "error": null,
  "id": "96",
  "result": {
    "address": "068eccc586b673d0920604483507eef738c0de0e",
    "executionResult": {
      "codeDeposit": 0,
      "depositSize": 0,
      "excepted": "Revert",
      "exceptedMessage": "TransparentUpgradeableProxy: admin cannot fallback to proxy target",
      "gasForDeposit": 0,
      "gasRefunded": 0,
      "gasUsed": 22346,
      "newAddress": "068eccc586b673d0920604483507eef738c0de0e",
      "output": "08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000425472616e73706172656e74557067
7261646561626c6550726f78793a2061646d696e2063616e6e6f742066616c6c6261636b20746f2070726f787920746172676574000000000000000000000000000000000000000000000000000000000000"
    },
    "transactionReceipt": {
      "bloom": "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
      "gasUsed": 22346,
      "log": [],
      "stateRoot": "153097c5d84a1fe23655f06feea20ae759499a3ce94e28cc2a3a75441d125a1f",
      "utxoRoot": "399117d8d0cbd117237b285cc08ae7e9083003a4bda614d7a3818f7869506a7f"
    }
  }
}
*/
func (p *ProxyETHCall) ToResponse(qresp *qtum.CallContractResponse) interface{} {

	if qresp.ExecutionResult.Output == "" {

		return &eth.JSONRPCError{
			Message: "Revert: executionResult output is empty",
			Code:    -32000,
		}

	}

	data := utils.AddHexPrefix(qresp.ExecutionResult.Output)
	qtumresp := eth.CallResponse(data)
	return &qtumresp

}
