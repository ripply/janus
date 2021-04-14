package transformer

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHGetTransactionByBlockNumberAndIndex implements ETHProxy
type ProxyETHGetTransactionByBlockNumberAndIndex struct {
	*qtum.Qtum
}

func (p *ProxyETHGetTransactionByBlockNumberAndIndex) Method() string {
	return "eth_getTransactionByBlockNumberAndIndex"
}

func (p *ProxyETHGetTransactionByBlockNumberAndIndex) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.GetTransactionByBlockNumberAndIndex
	if err := json.Unmarshal(rawreq.Params, &req); err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal request")
	}
	if req.TransactionHash == "" {
		return nil, errors.New("invalid argument 0: empty hex string")
	}

	return p.request(&req)
}

func (p *ProxyETHGetTransactionByBlockNumberAndIndex) request(req *eth.GetTransactionByBlockNumberAndIndex) (interface{}, error) {
	transactionIndex, err := hexutil.DecodeUint64(req.TransactionIndex)
	if err != nil {
		return nil, errors.Wrap(err, "invalid argument 1")
	}

	// Proxy eth_getBlockByNumber and return the transaction at requested index
	getBlockByNumber := ProxyETHGetBlockByNumber{p.Qtum}
	blockByNumber, err := getBlockByNumber.request(&eth.GetBlockByNumberRequest{BlockNumber: json.RawMessage([]byte(`"` + req.TransactionHash + `"`)), FullTransaction: true})

	if err != nil {
		return nil, err
	}

	if blockByNumber == nil {
		return nil, nil
	}

	if len(blockByNumber.Transactions) < int(transactionIndex) {
		return nil, nil
	}

	return blockByNumber.Transactions[int(transactionIndex)], nil
}
