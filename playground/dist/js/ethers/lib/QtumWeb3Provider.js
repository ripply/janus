// @ts-nocheck
import { providers } from "ethers";
import { QtumProvider } from "./QtumProvider";
export class QtumWeb3Provider extends providers.Web3Provider {
    constructor(provider, network) {
        super(provider, network);
    }
    async getUtxos() {
    }
    send(method, params) {
        console.log("send", method, params);
        if (method == "eth_sendTransaction") {
            this.getUtxos("0", 1);
        }
        if (typeof method != "string") {
            params = method.params;
            method = method.method;
        }
        return this.jsonRpcFetchFunc(method, params);
    }
}
QtumWeb3Provider.prototype.sendTransaction = QtumProvider.prototype.sendTransaction;
QtumWeb3Provider.prototype.prepareRequest = QtumProvider.prototype.prepareRequest;
QtumWeb3Provider.prototype.getUtxos = QtumProvider.prototype.getUtxos;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiUXR1bVdlYjNQcm92aWRlci5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvUXR1bVdlYjNQcm92aWRlci50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxjQUFjO0FBQ2QsT0FBTyxFQUFFLFNBQVMsRUFBRSxNQUFNLFFBQVEsQ0FBQztBQUduQyxPQUFPLEVBQUUsWUFBWSxFQUFFLE1BQU0sZ0JBQWdCLENBQUM7QUFFOUMsTUFBTSxPQUFPLGdCQUFpQixTQUFRLFNBQVMsQ0FBQyxZQUFZO0lBQ3hELFlBQ0ksUUFBNkMsRUFDN0MsT0FBb0I7UUFFcEIsS0FBSyxDQUFDLFFBQVEsRUFBRSxPQUFPLENBQUMsQ0FBQztJQUM3QixDQUFDO0lBRUQsS0FBSyxDQUFDLFFBQVE7SUFDZCxDQUFDO0lBRUQsSUFBSSxDQUFDLE1BQWMsRUFBRSxNQUFrQjtRQUNuQyxPQUFPLENBQUMsR0FBRyxDQUFDLE1BQU0sRUFBRSxNQUFNLEVBQUUsTUFBTSxDQUFDLENBQUE7UUFDbkMsSUFBSSxNQUFNLElBQUkscUJBQXFCLEVBQUU7WUFDakMsSUFBSSxDQUFDLFFBQVEsQ0FBQyxHQUFHLEVBQUUsQ0FBQyxDQUFDLENBQUE7U0FDeEI7UUFDRCxJQUFJLE9BQU8sTUFBTSxJQUFJLFFBQVEsRUFBRTtZQUMzQixNQUFNLEdBQUcsTUFBTSxDQUFDLE1BQU0sQ0FBQTtZQUN0QixNQUFNLEdBQUcsTUFBTSxDQUFDLE1BQU0sQ0FBQTtTQUN6QjtRQUNELE9BQU8sSUFBSSxDQUFDLGdCQUFnQixDQUFDLE1BQU0sRUFBRSxNQUFNLENBQUMsQ0FBQztJQUNqRCxDQUFDO0NBQ0o7QUFFRCxnQkFBZ0IsQ0FBQyxTQUFTLENBQUMsZUFBZSxHQUFHLFlBQVksQ0FBQyxTQUFTLENBQUMsZUFBZSxDQUFDO0FBQ3BGLGdCQUFnQixDQUFDLFNBQVMsQ0FBQyxjQUFjLEdBQUcsWUFBWSxDQUFDLFNBQVMsQ0FBQyxjQUFjLENBQUM7QUFDbEYsZ0JBQWdCLENBQUMsU0FBUyxDQUFDLFFBQVEsR0FBRyxZQUFZLENBQUMsU0FBUyxDQUFDLFFBQVEsQ0FBQyJ9