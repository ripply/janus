// @ts-nocheck
import { resolveProperties, } from "ethers/lib/utils";
import { checkTransactionType } from './helpers/utils';
import { providers } from "ethers";
import { QtumProvider } from "./QtumProvider";
export class QtumWeb3Provider extends providers.Web3Provider {
    constructor(provider, qtumWallet, network) {
        super(provider, network);
        this.qtumWallet = qtumWallet;
    }
    async send(method, params) {
        console.log("send", method, params);
        if (method == "eth_sendTransaction") {
            console.log("qtumRpcProvider", this.qtumRpcProvider);
            const param = params[0];
            if (!param) {
                throw new Error("Expected param to eth_sendTransaction", params);
            }
            if (!param.from) {
                throw new Error("Expected from parameter to eth_sendTransaction", param);
            }
            // const utxos = await this.qtumRpcProvider.getUtxos(param.from, 1);
            // console.log("utxos------", utxos);
            // await this.signTx(param);
            // console.log("Exiting...")
            param.gasLimit = "0x999999";
            param.gasPrice = "0x111";
            const signed = await this.qtumWallet.signTransaction(param);
            console.log("signed", signed);
            return;
        }
        if (typeof method != "string") {
            params = method.params;
            method = method.method;
        }
        return this.jsonRpcFetchFunc(method, params);
    }
    async signTx(transaction) {
        console.log("--------signTx-------");
        const tx = await resolveProperties(transaction);
        console.log("tx", tx);
        const { transactionType, neededAmount } = checkTransactionType(tx);
        console.log("transactionType", transactionType, "neededAmount", neededAmount);
    }
}
QtumWeb3Provider.prototype.sendTransaction = QtumProvider.prototype.sendTransaction;
QtumWeb3Provider.prototype.prepareRequest = QtumProvider.prototype.prepareRequest;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiUXR1bVdlYjNQcm92aWRlci5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvUXR1bVdlYjNQcm92aWRlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxjQUFjO0FBQ2QsT0FBTyxFQUNILGlCQUFpQixHQUNwQixNQUFNLGtCQUFrQixDQUFDO0FBQzFCLE9BQU8sRUFBRSxvQkFBb0IsRUFBd0IsTUFBTSxpQkFBaUIsQ0FBQTtBQUM1RSxPQUFPLEVBQUUsU0FBUyxFQUFFLE1BQU0sUUFBUSxDQUFDO0FBR25DLE9BQU8sRUFBRSxZQUFZLEVBQUUsTUFBTSxnQkFBZ0IsQ0FBQztBQUc5QyxNQUFNLE9BQU8sZ0JBQWlCLFNBQVEsU0FBUyxDQUFDLFlBQVk7SUFFeEQsWUFDSSxRQUE2QyxFQUM3QyxVQUFzQixFQUN0QixPQUFvQjtRQUVwQixLQUFLLENBQUMsUUFBUSxFQUFFLE9BQU8sQ0FBQyxDQUFDO1FBQ3pCLElBQUksQ0FBQyxVQUFVLEdBQUcsVUFBVSxDQUFDO0lBQ2pDLENBQUM7SUFFRCxLQUFLLENBQUMsSUFBSSxDQUFDLE1BQWMsRUFBRSxNQUFrQjtRQUN6QyxPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sRUFBRSxNQUFNLEVBQUUsTUFBTSxDQUFDLENBQUE7UUFDbkMsSUFBSSxNQUFNLElBQUkscUJBQXFCLEVBQUU7WUFDakMsT0FBTyxDQUFDLEdBQUcsQ0FBQyxpQkFBaUIsRUFBRSxJQUFJLENBQUMsZUFBZSxDQUFDLENBQUM7WUFDckQsTUFBTSxLQUFLLEdBQUcsTUFBTSxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQ3hCLElBQUksQ0FBQyxLQUFLLEVBQUU7Z0JBQ1IsTUFBTSxJQUFJLEtBQUssQ0FBQyx1Q0FBdUMsRUFBRSxNQUFNLENBQUMsQ0FBQzthQUNwRTtZQUNELElBQUksQ0FBQyxLQUFLLENBQUMsSUFBSSxFQUFFO2dCQUNiLE1BQU0sSUFBSSxLQUFLLENBQUMsZ0RBQWdELEVBQUUsS0FBSyxDQUFDLENBQUM7YUFDNUU7WUFDRCxvRUFBb0U7WUFDcEUscUNBQXFDO1lBQ3JDLDRCQUE0QjtZQUM1Qiw0QkFBNEI7WUFDNUIsS0FBSyxDQUFDLFFBQVEsR0FBRyxVQUFVLENBQUE7WUFDM0IsS0FBSyxDQUFDLFFBQVEsR0FBRyxPQUFPLENBQUE7WUFDeEIsTUFBTSxNQUFNLEdBQUcsTUFBTSxJQUFJLENBQUMsVUFBVSxDQUFDLGVBQWUsQ0FBQyxLQUFLLENBQUMsQ0FBQztZQUM1RCxPQUFPLENBQUMsR0FBRyxDQUFDLFFBQVEsRUFBRSxNQUFNLENBQUMsQ0FBQztZQUM5QixPQUFNO1NBQ1Q7UUFDRCxJQUFJLE9BQU8sTUFBTSxJQUFJLFFBQVEsRUFBRTtZQUMzQixNQUFNLEdBQUcsTUFBTSxDQUFDLE1BQU0sQ0FBQTtZQUN0QixNQUFNLEdBQUcsTUFBTSxDQUFDLE1BQU0sQ0FBQTtTQUN6QjtRQUNELE9BQU8sSUFBSSxDQUFDLGdCQUFnQixDQUFDLE1BQU0sRUFBRSxNQUFNLENBQUMsQ0FBQztJQUNqRCxDQUFDO0lBRUQsS0FBSyxDQUFDLE1BQU0sQ0FBQyxXQUFnQjtRQUN6QixPQUFPLENBQUMsR0FBRyxDQUFDLHVCQUF1QixDQUFDLENBQUM7UUFDckMsTUFBTSxFQUFFLEdBQUcsTUFBTSxpQkFBaUIsQ0FBQyxXQUFXLENBQUMsQ0FBQztRQUNoRCxPQUFPLENBQUMsR0FBRyxDQUFDLElBQUksRUFBRSxFQUFFLENBQUMsQ0FBQTtRQUNyQixNQUFNLEVBQUUsZUFBZSxFQUFFLFlBQVksRUFBRSxHQUFHLG9CQUFvQixDQUFDLEVBQUUsQ0FBQyxDQUFDO1FBQ25FLE9BQU8sQ0FBQyxHQUFHLENBQUMsaUJBQWlCLEVBQUUsZUFBZSxFQUFFLGNBQWMsRUFBRSxZQUFZLENBQUMsQ0FBQTtJQUNqRixDQUFDO0NBQ0o7QUFFRCxnQkFBZ0IsQ0FBQyxTQUFTLENBQUMsZUFBZSxHQUFHLFlBQVksQ0FBQyxTQUFTLENBQUMsZUFBZSxDQUFDO0FBQ3BGLGdCQUFnQixDQUFDLFNBQVMsQ0FBQyxjQUFjLEdBQUcsWUFBWSxDQUFDLFNBQVMsQ0FBQyxjQUFjLENBQUMifQ==