package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHBlockNumber implements ETHProxy
type ProxyETHBlockNumber struct {
	*qtum.Qtum
}

func (p *ProxyETHBlockNumber) Method() string {
	return "eth_blockNumber"
}

func (p *ProxyETHBlockNumber) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHBlockNumber) request() (*eth.BlockNumberResponse, eth.JSONRPCError) {
	qtumresp, err := p.Qtum.GetBlockCount()
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHBlockNumber) ToResponse(qtumresp *qtum.GetBlockCountResponse) *eth.BlockNumberResponse {
	hexVal := hexutil.EncodeBig(qtumresp.Int)
	ethresp := eth.BlockNumberResponse(hexVal)
	return &ethresp
}
