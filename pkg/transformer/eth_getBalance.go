package transformer

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHGetBalance implements ETHProxy
type ProxyETHGetBalance struct {
	*qtum.Qtum
}

func (p *ProxyETHGetBalance) Method() string {
	return "eth_getBalance"
}

func (p *ProxyETHGetBalance) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetBalanceRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	addr := utils.RemoveHexPrefix(req.Address)
	{
		// is address a contract or an account?
		qtumreq := qtum.GetAccountInfoRequest(addr)
		qtumresp, err := p.GetAccountInfo(&qtumreq)

		// the address is a contract
		if err == nil {
			// the unit of the balance Satoshi
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "is a contract")
			return hexutil.EncodeUint64(uint64(qtumresp.Balance)), nil
		}
	}

	{
		// try account
		base58Addr, err := p.FromHexAddress(addr)
		if err != nil {
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "error parsing address", "error", err)
			return nil, err
		}

		qtumreq := qtum.GetAddressBalanceRequest{Address: base58Addr}
		qtumresp, err := p.GetAddressBalance(&qtumreq)
		if err != nil {
			if err == qtum.ErrInvalidAddress {
				// invalid address should return 0x0
				return "0x0", nil
			}
			return nil, err
		}

		resp := *qtumresp
		balance := resp.Balance

		// 1 QTUM = 10 ^ 8 Satoshi
		floatBalance, exact := balance.Float64()

		if exact != true {
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "precision loss", "original", balance.String(), "after", floatBalance)
			return exact, errors.New("precision error:  float64 value does not represent the original decimal precisely")
		}

		return hexutil.EncodeUint64(uint64(floatBalance)), nil

	}
}
