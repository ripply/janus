package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
	"github.com/shopspring/decimal"
)

type ETHProxyInitializer = func(*qtum.Qtum) ETHProxy

var (
	getTransactionByHashBlockNumber  = "0xf8f"
	getTransactionByHashBlockHash    = "bba11e1bacc69ba535d478cf1f2e542da3735a517b0b8eebaf7e6bb25eeb48c5"
	getTransactionByHashBlockHexHash = utils.AddHexPrefix(getTransactionByHashBlockHash)
	getTransactionByHashResponseData = eth.GetTransactionByHashResponse{
		BlockHash:        getTransactionByHashBlockHexHash,
		BlockNumber:      getTransactionByHashBlockNumber,
		TransactionIndex: "0x2",
		Hash:             "0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		Nonce:            "0x0",
		Value:            "0x0",
		Input:            "0x",
		From:             "0x0000000000000000000000000000000000000000",
		To:               "0x0000000000000000000000000000000000000000",
		Gas:              "0x0",
		GasPrice:         "0x0",
	}

	getTransactionByHashResponse = eth.GetBlockByHashResponse{
		Number:           getTransactionByHashBlockNumber,
		Hash:             getTransactionByHashBlockHexHash,
		ParentHash:       "0x6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Miner:            "0x0000000000000000000000000000000000000000",
		Size:             "0x26c",
		Nonce:            "0x0",
		TransactionsRoot: "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		ReceiptsRoot:     "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		StateRoot:        "0x3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		Difficulty:       "0x4",
		TotalDifficulty:  "0x4",
		LogsBloom:        EmptyLogsBloom,
		ExtraData:        "0x0",
		GasLimit:         utils.AddHexPrefix(qtum.DefaultBlockGasLimit),
		GasUsed:          "0x0",
		Timestamp:        "0x5b95ebd0",
		Transactions: []interface{}{"0x3208dc44733cbfa11654ad5651305428de473ef1e61a1ec07b0c1a5f4843be91",
			"0x8fcd819194cce6a8454b2bec334d3448df4f097e9cdc36707bfd569900268950"},
		Sha3Uncles: DefaultSha3Uncles,
		Uncles:     []string{},
	}

	getTransactionByHashResponseWithTransactions = eth.GetBlockByHashResponse{
		Number:           getTransactionByHashBlockNumber,
		Hash:             getTransactionByHashBlockHexHash,
		ParentHash:       "0x6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Miner:            "0x0000000000000000000000000000000000000000",
		Size:             "0x26c",
		Nonce:            "0x0",
		TransactionsRoot: "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		ReceiptsRoot:     "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		StateRoot:        "0x3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		Difficulty:       "0x4",
		TotalDifficulty:  "0x4",
		LogsBloom:        EmptyLogsBloom,
		ExtraData:        "0x0",
		GasLimit:         utils.AddHexPrefix(qtum.DefaultBlockGasLimit),
		GasUsed:          "0x0",
		Timestamp:        "0x5b95ebd0",
		Transactions: []interface{}{
			getTransactionByHashResponseData,
			getTransactionByHashResponseData,
		},
		Sha3Uncles: DefaultSha3Uncles,
		Uncles:     []string{},
	}

	getTransactionByBlockResponse = eth.GetBlockByNumberResponse{
		Number:           getTransactionByHashBlockNumber,
		Hash:             getTransactionByHashBlockHexHash,
		ParentHash:       "0x6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Miner:            "0x0000000000000000000000000000000000000000",
		Size:             "0x26c",
		Nonce:            "0x0",
		TransactionsRoot: "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		ReceiptsRoot:     "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		StateRoot:        "0x3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		Difficulty:       "0x4",
		TotalDifficulty:  "0x4",
		LogsBloom:        EmptyLogsBloom,
		ExtraData:        "0x0",
		GasLimit:         utils.AddHexPrefix(qtum.DefaultBlockGasLimit),
		GasUsed:          "0x0",
		Timestamp:        "0x5b95ebd0",
		Transactions: []interface{}{"0x3208dc44733cbfa11654ad5651305428de473ef1e61a1ec07b0c1a5f4843be91",
			"0x8fcd819194cce6a8454b2bec334d3448df4f097e9cdc36707bfd569900268950"},
		Sha3Uncles: DefaultSha3Uncles,
		Uncles:     []string{},
	}

	getTransactionByBlockResponseWithTransactions = eth.GetBlockByNumberResponse{
		Number:           getTransactionByHashBlockNumber,
		Hash:             getTransactionByHashBlockHexHash,
		ParentHash:       "0x6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Miner:            "0x0000000000000000000000000000000000000000",
		Size:             "0x26c",
		Nonce:            "0x0",
		TransactionsRoot: "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		ReceiptsRoot:     "0x0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		StateRoot:        "0x3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		Difficulty:       "0x4",
		TotalDifficulty:  "0x4",
		LogsBloom:        EmptyLogsBloom,
		ExtraData:        "0x0",
		GasLimit:         utils.AddHexPrefix(qtum.DefaultBlockGasLimit),
		GasUsed:          "0x0",
		Timestamp:        "0x5b95ebd0",
		Transactions: []interface{}{
			getTransactionByHashResponseData,
			getTransactionByHashResponseData,
		},
		Sha3Uncles: DefaultSha3Uncles,
		Uncles:     []string{},
	}
)

func setupGetBlockByHashResponses(t *testing.T, mockedClientDoer *doerMappedMock) {
	//preparing answer to "getblockhash"
	getBlockHashResponse := qtum.GetBlockHashResponse(getTransactionByHashBlockHexHash)
	err := mockedClientDoer.AddResponse(qtum.MethodGetBlockHash, getBlockHashResponse)
	if err != nil {
		t.Fatal(err)
	}

	getBlockHeaderResponse := qtum.GetBlockHeaderResponse{
		Hash:              getTransactionByHashBlockHash,
		Confirmations:     1,
		Height:            3983,
		Version:           536870912,
		VersionHex:        "20000000",
		Merkleroot:        "0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		Time:              1536551888,
		Mediantime:        1536551728,
		Nonce:             0,
		Bits:              "207fffff",
		Difficulty:        4.656542373906925,
		Chainwork:         "0000000000000000000000000000000000000000000000000000000000001f20",
		HashStateRoot:     "3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		HashUTXORoot:      "130a3e712d9f8b06b83f5ebf02b27542fb682cdff3ce1af1c17b804729d88a47",
		Previousblockhash: "6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Flags:             "proof-of-stake",
		Proofhash:         "15bd6006ecbab06708f705ecf68664b78b388e4d51416cdafb019d5b90239877",
		Modifier:          "a79c00d1d570743ca8135a173d535258026d26bafbc5f3d951c3d33486b1f120",
	}
	err = mockedClientDoer.AddResponse(qtum.MethodGetBlockHeader, getBlockHeaderResponse)
	if err != nil {
		t.Fatal(err)
	}

	getBlockResponse := qtum.GetBlockResponse{
		Hash:              getTransactionByHashBlockHash,
		Confirmations:     1,
		Strippedsize:      584,
		Size:              620,
		Weight:            2372,
		Height:            3983,
		Version:           536870912,
		VersionHex:        "20000000",
		Merkleroot:        "0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		Time:              1536551888,
		Mediantime:        1536551728,
		Nonce:             0,
		Bits:              "207fffff",
		Difficulty:        4.656542373906925,
		Chainwork:         "0000000000000000000000000000000000000000000000000000000000001f20",
		HashStateRoot:     "3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		HashUTXORoot:      "130a3e712d9f8b06b83f5ebf02b27542fb682cdff3ce1af1c17b804729d88a47",
		Previousblockhash: "6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Flags:             "proof-of-stake",
		Proofhash:         "15bd6006ecbab06708f705ecf68664b78b388e4d51416cdafb019d5b90239877",
		Modifier:          "a79c00d1d570743ca8135a173d535258026d26bafbc5f3d951c3d33486b1f120",
		Txs: []string{"3208dc44733cbfa11654ad5651305428de473ef1e61a1ec07b0c1a5f4843be91",
			"8fcd819194cce6a8454b2bec334d3448df4f097e9cdc36707bfd569900268950"},
		Nextblockhash: "d7758774cfdd6bab7774aa891ae035f1dc5a2ff44240784b5e7bdfd43a7a6ec1",
		Signature:     "3045022100a6ab6c2b14b1f73e734f1a61d4d22385748e48836492723a6ab37cdf38525aba022014a51ecb9e51f5a7a851641683541fec6f8f20205d0db49e50b2a4e5daed69d2",
	}
	err = mockedClientDoer.AddResponse(qtum.MethodGetBlock, getBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	getTransactionResponse := qtum.GetTransactionResponse{
		Amount:            decimal.NewFromFloat(0.20689141),
		Fee:               decimal.NewFromFloat(-0.2012),
		Confirmations:     2,
		BlockHash:         getTransactionByHashBlockHash,
		BlockIndex:        2,
		BlockTime:         1533092896,
		ID:                "11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		Time:              1533092879,
		ReceivedAt:        1533092879,
		Bip125Replaceable: "no",
		Details: []*qtum.TransactionDetail{{Account: "",
			Category:  "send",
			Amount:    decimal.NewFromInt(0),
			Vout:      0,
			Fee:       decimal.NewFromFloat(-0.2012),
			Abandoned: false}},
		Hex: "020000000159c0514feea50f915854d9ec45bc6458bb14419c78b17e7be3f7fd5f563475b5010000006a473044022072d64a1f4ea2d54b7b05050fc853ab192c91cc5ca17e23007867f92f2ab59d9202202b8c9ab9348c8edbb3b98b1788382c8f37642ec9bd6a4429817ab79927319200012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140feffffff02000000000000000063010403400d0301644440c10f190000000000000000000000006b22910b1e302cf74803ffd1691c2ecb858d3712000000000000000000000000000000000000000000000000000000000000000a14be528c8378ff082e4ba43cb1baa363dbf3f577bfc260e66272970100001976a9146b22910b1e302cf74803ffd1691c2ecb858d371288acb00f0000",
	}
	err = mockedClientDoer.AddResponse(qtum.MethodGetTransaction, getTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	decodedRawTransactionResponse := qtum.DecodedRawTransactionResponse{
		ID:       "11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		Hash:     "d0fe0caa1b798c36da37e9118a06a7d151632d670b82d1c7dc3985577a71880f",
		Size:     552,
		Vsize:    552,
		Version:  2,
		Locktime: 608,
		Vins: []*qtum.DecodedRawTransactionInV{{
			TxID: "7f5350dc474f2953a3f30282c1afcad2fb61cdcea5bd949c808ecc6f64ce1503",
			Vout: 0,
			ScriptSig: struct {
				Asm string `json:"asm"`
				Hex string `json:"hex"`
			}{
				Asm: "3045022100af4de764705dbd3c0c116d73fe0a2b78c3fab6822096ba2907cfdae2bb28784102206304340a6d260b364ef86d6b19f2b75c5e55b89fb2f93ea72c05e09ee037f60b[ALL] 03520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140",
				Hex: "483045022100af4de764705dbd3c0c116d73fe0a2b78c3fab6822096ba2907cfdae2bb28784102206304340a6d260b364ef86d6b19f2b75c5e55b89fb2f93ea72c05e09ee037f60b012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140",
			},
		}},
		Vouts: []*qtum.DecodedRawTransactionOutV{},
	}
	err = mockedClientDoer.AddResponse(qtum.MethodDecodeRawTransaction, decodedRawTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	getTransactionReceiptResponse := qtum.GetTransactionReceiptResponse{}
	err = mockedClientDoer.AddResponse(qtum.MethodGetTransactionReceipt, &getTransactionReceiptResponse)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Get an actual response for this (only addresses are used in this test though)
	getRawTransactionResponse := qtum.GetRawTransactionResponse{
		Vouts: []qtum.RawTransactionVout{
			{
				Details: struct {
					Addresses []string `json:"addresses"`
					Asm       string   `json:"asm"`
					Hex       string   `json:"hex"`
					// ReqSigs   interface{} `json:"reqSigs"`
					Type string `json:"type"`
				}{
					Addresses: []string{
						"7926223070547d2d15b2ef5e7383e541c338ffe9",
					},
				},
			},
		},
	}
	err = mockedClientDoer.AddResponse(qtum.MethodGetRawTransaction, &getRawTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}
}

func initializeProxyETHGetBlockByHash(qtumClient *qtum.Qtum) ETHProxy {
	return &ProxyETHGetBlockByHash{qtumClient}
}

func TestGetBlockByHashRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByHash,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockHexHash + `"`), []byte(`false`)},
		&getTransactionByHashResponse,
	)
}

func TestGetBlockByHashTransactionsRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByHash,
		[]json.RawMessage{[]byte(`"` + getTransactionByHashBlockHexHash + `"`), []byte(`true`)},
		&getTransactionByHashResponseWithTransactions,
	)
}

func testETHProxyRequest(t *testing.T, initializer ETHProxyInitializer, requestParams []json.RawMessage, want interface{}) {
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)

	setupGetBlockByHashResponses(t, mockedClientDoer)

	//preparing proxy & executing request
	proxyEth := initializer(qtumClient)
	got, err := proxyEth.Request(request)
	if err != nil {
		t.Fatalf("Failed to process request on %T.Request(%s): %s", proxyEth, requestParams, err)
	}

	if !reflect.DeepEqual(got, want) {
		wantString := string(mustMarshalIndent(want, "", "  "))
		gotString := string(mustMarshalIndent(got, "", "  "))
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			wantString,
			gotString,
		)
		if wantString == gotString {
			t.Errorf("Want and Got are equal strings but !DeepEqual, probably differ in types (%T ?= %T)", want, got)
		}
	}
}
