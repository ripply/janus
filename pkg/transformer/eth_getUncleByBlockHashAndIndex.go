package transformer

import (
	"github.com/qtumproject/janus/pkg/eth"
)

type ETHGetUncleByBlockHashAndIndex struct {
}

func (p *ETHGetUncleByBlockHashAndIndex) Method() string {
	return "eth_getUncleByBlockHashAndIndex"
}

func (p *ETHGetUncleByBlockHashAndIndex) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	// hardcoded to nil
	return nil, nil
}
