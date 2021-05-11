package transformer

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHSendRawTransaction implements ETHProxy
type ProxyETHSendRawTransaction struct {
	*qtum.Qtum
}

var _ ETHProxy = (*ProxyETHSendRawTransaction)(nil)

func (p *ProxyETHSendRawTransaction) Method() string {
	return "eth_sendRawTransaction"
}

func (p *ProxyETHSendRawTransaction) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var params eth.SendRawTransactionRequest
	if err := unmarshalRequest(req.Params, &params); err != nil {
		return nil, err
	}
	if params[0] == "" {
		return nil, errors.Errorf("invalid parameter: raw transaction hexed string is empty")
	}

	return p.request(params)
}

func (p *ProxyETHSendRawTransaction) request(params eth.SendRawTransactionRequest) (eth.SendRawTransactionResponse, error) {
	var (
		qtumHexedRawTx = utils.RemoveHexPrefix(params[0])
		req            = qtum.SendRawTransactionRequest([1]string{qtumHexedRawTx})
	)

	qtumresp, err := p.Qtum.SendRawTransaction(&req)
	if err != nil {
		return eth.SendRawTransactionResponse{}, err
	}

	if p.Chain() == qtum.ChainRegTest {
		if _, err = p.Generate(1, nil); err != nil {
			p.GetErrorLogger().Log("Error generating new block", err)
		}
	}

	resp := *qtumresp
	ethHexedTxHash := utils.AddHexPrefix(resp.Result)
	return eth.SendRawTransactionResponse([1]string{ethHexedTxHash}), nil
}
