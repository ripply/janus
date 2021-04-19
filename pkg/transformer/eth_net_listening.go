package transformer

import (
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHGetCode implements ETHProxy
type ProxyNetListening struct {
	*qtum.Qtum
}

func (p *ProxyNetListening) Method() string {
	return "net_listening"
}

func (p *ProxyNetListening) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	networkInfo, err := p.GetNetworkInfo()
	if err != nil {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "Failed to query network info", "err", err)
		return false, err
	}

	p.GetDebugLogger().Log("method", p.Method(), "network active", networkInfo.NetworkActive)
	return networkInfo.NetworkActive, nil
}
