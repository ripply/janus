package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHNewBlockFilter implements ETHProxy
type ProxyETHNewBlockFilter struct {
	*qtum.Qtum
	filter *eth.FilterSimulator
}

func (p *ProxyETHNewBlockFilter) Method() string {
	return "eth_newBlockFilter"
}

func (p *ProxyETHNewBlockFilter) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHNewBlockFilter) request() (eth.NewBlockFilterResponse, error) {
	blockCount, err := p.GetBlockCount()
	if err != nil {
		return "", err
	}

	if p.Chain() == qtum.ChainRegTest {
		defer func() {
			if _, generateErr := p.Generate(1, nil); generateErr != nil {
				p.GetErrorLogger().Log("Error generating new block", generateErr)
			}
		}()
	}

	filter := p.filter.New(eth.NewBlockFilterTy)
	filter.Data.Store("lastBlockNumber", blockCount.Uint64())

	return eth.NewBlockFilterResponse(hexutil.EncodeUint64(filter.ID)), nil
}
