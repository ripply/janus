package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

type ProxyETHChainId struct {
	*qtum.Qtum
}

func (p *ProxyETHChainId) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainId) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var qtumresp *qtum.GetBlockChainInfoResponse
	if err := p.Qtum.Request(qtum.MethodGetBlockChainInfo, nil, &qtumresp); err != nil {
		return nil, err
	}

	var chainId *big.Int
	switch qtumresp.Chain {
	case "regtest":
		chainId = big.NewInt(113)
	default:
		chainId = big.NewInt(81)
		p.GetDebugLogger().Log("method", p.Method(), "msg", "Unknown chain "+qtumresp.Chain)
	}

	return eth.ChainIdResponse(hexutil.EncodeBig(chainId)), nil
}
