package server

import (
	"encoding/json"
	stdLog "log"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/notifier"
	"golang.org/x/net/websocket"
)

func httpHandler(c echo.Context) error {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if !ok {
		return errors.New("Could not find myctx")
	}

	var rpcReq *eth.JSONRPCRequest
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(&rpcReq); err != nil {
		return errors.Wrap(err, "json decoder issue")
	}

	cc.rpcReq = rpcReq

	cc.GetLogger().Log("msg", "proxy RPC", "method", rpcReq.Method)

	// level.Debug(cc.logger).Log("msg", "before call transformer#Transform")
	result, err := cc.transformer.Transform(rpcReq, c)
	// level.Debug(cc.logger).Log("msg", "after call transformer#Transform")

	if err != nil {
		err1 := errors.Cause(err)
		if err != err1 {
			cc.GetErrorLogger().Log("err", err.Error())
			return cc.JSONRPCError(&eth.JSONRPCError{
				Code:    100,
				Message: err1.Error(),
			})
		}

		return err
	}

	// Allow transformer to return an explicit JSON error
	if jerr, isJSONErr := result.(*eth.JSONRPCError); isJSONErr {
		return cc.JSONRPCError(jerr)
	}

	return cc.JSONRPCResult(result)
}

/*
// subscription topic name is a random hex number
// unsubscribe
{"id": 1, "method": "eth_unsubscribe", "params": ["0x9cef478923ff08bf67fde6c64013158d"]}
=>
{
   "jsonrpc": "2.0",
   "id": 1,
   "result": true
}
// new heads
{
   "id": 1,
   "method": "eth_subscribe",
   "params": [
      "newHeads"
   ]
}
=>
{
   "jsonrpc": "2.0",
   "id": 2,
   "result": "0x9ce59a13059e417087c02d3236a0b1cc"
}
=>
{
   "jsonrpc": "2.0",
   "method": "eth_subscription",
   "params": {
      "result": {
         "difficulty": "0x15d9223a23aa",
         "extraData": "0xd983010305844765746887676f312e342e328777696e646f7773",
         "gasLimit": "0x47e7c4",
         "gasUsed": "0x38658",
         "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
         "miner": "0xf8b483dba2c3b7176a3da549ad41a48bb3121069",
         "nonce": "0x084149998194cc5f",
         "number": "0x1348c9",
         "parentHash": "0x7736fab79e05dc611604d22470dadad26f56fe494421b5b333de816ce1f25701",
         "receiptRoot": "0x2fab35823ad00c7bb388595cb46652fe7886e00660a01e867824d3dceb1c8d36",
         "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
         "stateRoot": "0xb3346685172db67de536d8765c43c31009d0eb3bd9c501c9be3229203f15f378",
         "timestamp": "0x56ffeff8",
         "transactionsRoot": "0x0167ffa60e3ebc0b080cdb95f7c0087dd6c0e61413140e39d94d3468d7c9689f"
      },
      "subscription": "0x9ce59a13059e417087c02d3236a0b1cc"
   }
}
// logs
{
   "id": 1,
   "method": "eth_subscribe",
   "params": [
      "logs",
      {
         "address": "0x8320fe7702b96808f7bbc0d4a888ed1468216cfd",
         "topics": [
            "0xd78a0cb8bb633d06981248b816e7bd33c2a35a6089241d099fa519e361cab902"
         ]
      }
   ]
}
=>
{
   "jsonrpc": "2.0",
   "id": 2,
   "result": "0x4a8a4c0517381924f9838102c5a4dcb7"
}
=>
{
   "jsonrpc": "2.0",
   "method": "eth_subscription",
   "params": {
      "subscription": "0x4a8a4c0517381924f9838102c5a4dcb7",
      "result": {
         "address": "0x8320fe7702b96808f7bbc0d4a888ed1468216cfd",
         "blockHash": "0x61cdb2a09ab99abf791d474f20c2ea89bf8de2923a2d42bb49944c8c993cbf04",
         "blockNumber": "0x29e87",
         "data": "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003",
         "logIndex": "0x0",
         "topics": [
            "0xd78a0cb8bb633d06981248b816e7bd33c2a35a6089241d099fa519e361cab902"
         ],
         "transactionHash": "0xe044554a0a55067caafd07f8020ab9f2af60bdfe337e395ecd84b4877a3d1ab4",
         "transactionIndex": "0x0"
      }
   }
}
// new pending transactions
{
   "id": 1,
   "method": "eth_subscribe",
   "params": [
      "newPendingTransactions"
   ]
}
=>
{
   "jsonrpc": "2.0",
   "id": 2,
   "result": "0xc3b33aa549fb9a60e95d21862596617c"
}
=>
{
   "jsonrpc": "2.0",
   "method": "eth_subscription",
   "params": {
      "subscription": "0xc3b33aa549fb9a60e95d21862596617c",
      "result": "0xd6fdc5cc41a9959e922f30cb772a9aef46f4daea279307bc5f7024edc4ccd7fa"
   }
}
// syncing
{
   "id": 1,
   "method": "eth_subscribe",
   "params": [
      "syncing"
   ]
}
=>
{
   "jsonrpc": "2.0",
   "id": 2,
   "result": "0xe2ffeb2703bcf602d42922385829ce96"
}
=>
{
   "subscription": "0xe2ffeb2703bcf602d42922385829ce96",
   "result": {
      "syncing": true,
      "status": {
         "startingBlock": 674427,
         "currentBlock": 67400,
         "highestBlock": 674432,
         "pulledStates": 0,
         "knownStates": 0
      }
   }
}
*/

func websocketHandler(c echo.Context) error {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if !ok {
		return errors.New("Could not find myctx")
	}

	cc.GetDebugLogger().Log("msg", "Got websocket request... opening")
	req := c.Request()
	req.Header.Set("Origin", "http://localhost:23888")

	websocket.Handler(func(ws *websocket.Conn) {
		cc.GetDebugLogger().Log("msg", "opened websocket")
		notifier := notifier.NewNotifier(
			c.Request().Context(),
			func() {
				ws.Close()
			}, func(value interface{}) error {
				return websocket.Message.Send(ws, value)
			},
			cc.GetLogger(),
		)
		c.Set("notifier", notifier)
		defer ws.Close()
		defer func() {
			cc.GetDebugLogger().Log("msg", "Websocket connection closed")
		}()
		go notifier.Run()

		for {
			cc.GetDebugLogger().Log("msg", "reading websocket request")
			// Read
			var payload string
			err := websocket.Message.Receive(ws, &payload)
			if err != nil {
				c.Logger().Error(err)
				return
			}

			var rpcReq eth.JSONRPCRequest
			json.Unmarshal([]byte(payload), &rpcReq)

			cc.rpcReq = &rpcReq

			result, err := cc.transformer.Transform(&rpcReq, c)

			response := result

			if err != nil {
				err1 := errors.Cause(err)
				if err != err1 {
					cc.GetErrorLogger().Log("err", err.Error())
					response = cc.GetJSONRPCError(&eth.JSONRPCError{
						Code:    100,
						Message: err1.Error(),
					})
				}
			}

			// Allow transformer to return an explicit JSON error
			if jerr, isJSONErr := response.(*eth.JSONRPCError); isJSONErr {
				response = cc.GetJSONRPCError(jerr)
			} else {
				response, err = cc.GetJSONRPCResult(response)
				if err != nil {
					cc.GetErrorLogger().Log("err", err.Error())
					return
				}
			}

			responseBytes, err := json.Marshal(response)
			if err != nil {
				cc.GetErrorLogger().Log("err", err.Error())
				return
			}

			cc.GetDebugLogger().Log("response", string(responseBytes))

			err = websocket.Message.Send(ws, string(responseBytes))
			if err == nil {
				notifier.ResponseSent()
			} else {
				cc.GetErrorLogger().Log("err", err.Error())
				return
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func errorHandler(err error, c echo.Context) {
	myctx := c.Get("myctx")
	cc, ok := myctx.(*myCtx)
	if ok {
		cc.GetErrorLogger().Log("err", err.Error())
		if err := cc.JSONRPCError(&eth.JSONRPCError{
			Code:    100,
			Message: err.Error(),
		}); err != nil {
			cc.GetErrorLogger().Log("msg", "reply to client", "err", err.Error())
		}
		return
	}

	stdLog.Println("errorHandler", err.Error())
}
