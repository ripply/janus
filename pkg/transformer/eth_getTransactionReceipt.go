package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/conversion"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

var STATUS_SUCCESS = "0x1"
var STATUS_FAILURE = "0x0"

// ProxyETHGetTransactionReceipt implements ETHProxy
type ProxyETHGetTransactionReceipt struct {
	*qtum.Qtum
}

func (p *ProxyETHGetTransactionReceipt) Method() string {
	return "eth_getTransactionReceipt"
}

func (p *ProxyETHGetTransactionReceipt) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetTransactionReceiptRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}
	if req == "" {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("empty transaction hash")
	}
	var (
		txHash  = utils.RemoveHexPrefix(string(req))
		qtumReq = qtum.GetTransactionReceiptRequest(txHash)
	)
	return p.request(&qtumReq)
}

func (p *ProxyETHGetTransactionReceipt) request(req *qtum.GetTransactionReceiptRequest) (*eth.GetTransactionReceiptResponse, eth.JSONRPCError) {
	qtumReceipts, err := p.Qtum.GetTransactionReceipt(string(*req))
	if err != nil || (qtumReceipts != nil && len(*qtumReceipts) == 0) {
		ethTx, _, getRewardTransactionErr := getRewardTransactionByHash(p.Qtum, string(*req))
		if getRewardTransactionErr != nil {
			errCause := errors.Cause(err)
			if errCause == qtum.EmptyResponseErr {
				return nil, nil
			}
			p.Qtum.GetDebugLogger().Log("msg", "Transaction does not exist", "txid", string(*req))
			return nil, eth.NewCallbackError(err.Error())
		}
		return &eth.GetTransactionReceiptResponse{
			TransactionHash:  ethTx.Hash,
			TransactionIndex: ethTx.TransactionIndex,
			BlockHash:        ethTx.BlockHash,
			BlockNumber:      ethTx.BlockNumber,
			// TODO: This is higher than GasUsed in geth but does it matter?
			CumulativeGasUsed: NonContractVMGasLimit,
			EffectiveGasPrice: "0x0",
			GasUsed:           NonContractVMGasLimit,
			From:              ethTx.From,
			To:                ethTx.To,
			Logs:              []eth.Log{},
			LogsBloom:         eth.EmptyLogsBloom,
			Status:            STATUS_SUCCESS,
		}, nil
	}

	var cumulativeGasUsed uint64
	var gasUsed uint64
	excepted := "None"

	ethReceipt := &eth.GetTransactionReceiptResponse{}
	// if len(*qtumReceipts) > 1 then the transaction has multiple EVM outputs
	// this is fundamentally different from ethereum so this will never map 1 to 1
	// users will generally not be doing multiple EVM outputs in the same transaction
	// to return something to the client that matches up with what web3 libraries expect
	// we take the first output's values and then add up the gas from all of the outputs
	// if any of the outputs revert, we consider the entire transaction a failure
	for i, qtumReceipt := range *qtumReceipts {
		if excepted == "None" && qtumReceipt.Excepted != "None" {
			excepted = qtumReceipt.Excepted
		}

		if i == 0 {
			ethReceipt.TransactionHash = utils.AddHexPrefix(qtumReceipt.TransactionHash)
			ethReceipt.TransactionIndex = hexutil.EncodeUint64(qtumReceipt.TransactionIndex)
			ethReceipt.BlockHash = utils.AddHexPrefix(qtumReceipt.BlockHash)
			ethReceipt.BlockNumber = hexutil.EncodeUint64(qtumReceipt.BlockNumber)
			ethReceipt.ContractAddress = utils.AddHexPrefixIfNotEmpty(qtumReceipt.ContractAddress)
			ethReceipt.From = utils.AddHexPrefixIfNotEmpty(qtumReceipt.From)
			ethReceipt.To = utils.AddHexPrefixIfNotEmpty(qtumReceipt.To)
			// TODO: researching
			// ! Temporary accept this value to be always zero, as it is at eth logs
			ethReceipt.LogsBloom = eth.EmptyLogsBloom
		}

		gasUsed += qtumReceipt.GasUsed
		cumulativeGasUsed += qtumReceipt.CumulativeGasUsed
	}

	ethReceipt.CumulativeGasUsed = hexutil.EncodeUint64(cumulativeGasUsed)
	ethReceipt.GasUsed = hexutil.EncodeUint64(gasUsed)

	status := STATUS_FAILURE
	if excepted == "None" {
		status = STATUS_SUCCESS
	} else {
		p.Qtum.GetDebugLogger().Log("transaction", ethReceipt.TransactionHash, "msg", "transaction excepted", "message", excepted)
	}
	ethReceipt.Status = status

	for _, qtumReceipt := range *qtumReceipts {
		r := qtum.TransactionReceipt(qtumReceipt)
		logs := conversion.ExtractETHLogsFromTransactionReceipt(&r, r.Log)
		if ethReceipt.Logs == nil && logs != nil {
			ethReceipt.Logs = []eth.Log{}
		}
		ethReceipt.Logs = append(ethReceipt.Logs, logs...)
	}

	qtumTx, err := p.Qtum.GetRawTransaction((*qtumReceipts)[0].TransactionHash, false)
	if err != nil {
		p.GetDebugLogger().Log("msg", "couldn't get transaction", "err", err)
		return nil, eth.NewCallbackError("couldn't get transaction")
	}
	decodedRawQtumTx, err := p.Qtum.DecodeRawTransaction(qtumTx.Hex)
	if err != nil {
		p.GetDebugLogger().Log("msg", "couldn't decode raw transaction", "err", err)
		return nil, eth.NewCallbackError("couldn't decode raw transaction")
	}
	if decodedRawQtumTx.IsContractCreation() {
		ethReceipt.To = ""
	} else {
		ethReceipt.ContractAddress = ""
	}

	// TODO: researching
	// - The following code reason is unknown (see original comment)
	// - Code temporary commented, until an error occures
	// ! Do not remove
	// // contractAddress : DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
	// if status != "0x1" {
	// 	// if failure, should return null for contractAddress, instead of the zero address.
	// 	ethTxReceipt.ContractAddress = ""
	// }

	return ethReceipt, nil
}
