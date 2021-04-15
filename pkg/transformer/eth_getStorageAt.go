package transformer

import (
	"fmt"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHGetStorageAt implements ETHProxy
type ProxyETHGetStorageAt struct {
	*qtum.Qtum
}

func (p *ProxyETHGetStorageAt) Method() string {
	return "eth_getStorageAt"
}

func (p *ProxyETHGetStorageAt) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetStorageRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	qtumAddress := utils.RemoveHexPrefix(req.Address)
	blockNumber, err := getBlockNumberByParam(p.Qtum, req.BlockNumber, false)
	if err != nil {
		p.GetDebugLogger().Log("msg", fmt.Sprintf("Failed to get block number by param for '%s'", req.BlockNumber), "err", err)
		return nil, err
	}

	return p.request(&qtum.GetStorageRequest{
		Address:     qtumAddress,
		BlockNumber: blockNumber,
	}, utils.RemoveHexPrefix(req.Index))
}

func (p *ProxyETHGetStorageAt) request(ethreq *qtum.GetStorageRequest, index string) (*eth.GetStorageResponse, error) {
	qtumresp, err := p.Qtum.GetStorage(ethreq)
	if err != nil {
		return nil, err
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp, index), nil
}

func (p *ProxyETHGetStorageAt) ToResponse(qtumresp *qtum.GetStorageResponse, index string) *eth.GetStorageResponse {
	// the value for unknown anything
	storageData := eth.GetStorageResponse("0x0000000000000000000000000000000000000000000000000000000000000000")
	for _, outerValue := range *qtumresp {
		qtumStorageData, ok := outerValue[index]
		if ok {
			storageData = eth.GetStorageResponse(utils.AddHexPrefix(qtumStorageData))
			return &storageData
		}
	}

	return &storageData
}
