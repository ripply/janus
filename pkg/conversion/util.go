package conversion

import (
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/utils"
)

func ExtractETHLogsFromTransactionReceipt(receipt *qtum.TransactionReceipt) []eth.Log {
	logs := make([]eth.Log, 0, len(receipt.Log))
	for i, log := range receipt.Log {
		topics := make([]string, 0, len(log.Topics))
		for _, topic := range log.Topics {
			topics = append(topics, utils.AddHexPrefix(topic))
		}
		logs = append(logs, eth.Log{
			TransactionHash:  utils.AddHexPrefix(receipt.TransactionHash),
			TransactionIndex: hexutil.EncodeUint64(receipt.TransactionIndex),
			BlockHash:        utils.AddHexPrefix(receipt.BlockHash),
			BlockNumber:      hexutil.EncodeUint64(receipt.BlockNumber),
			Data:             utils.AddHexPrefix(log.Data),
			Address:          utils.AddHexPrefix(log.Address),
			Topics:           topics,
			LogIndex:         hexutil.EncodeUint64(uint64(i)),
		})
	}
	return logs
}
