import { providers } from "ethers";
import { Networkish } from "@ethersproject/networks";
import { ExternalProvider, JsonRpcFetchFunc } from "@ethersproject/providers";
import { QtumWallet } from "./QtumWallet";
export declare class QtumWeb3Provider extends providers.Web3Provider {
    readonly qtumWallet: QtumWallet;
    constructor(provider: ExternalProvider | JsonRpcFetchFunc, qtumWallet: QtumWallet, network?: Networkish);
    send(method: string, params: Array<any>): Promise<any>;
    signTx(transaction: any): Promise<String>;
}
