import { providers } from "ethers";
import { Networkish } from "@ethersproject/networks";
import { ExternalProvider, JsonRpcFetchFunc } from "@ethersproject/providers";
export declare class QtumWeb3Provider extends providers.Web3Provider {
    constructor(provider: ExternalProvider | JsonRpcFetchFunc, network?: Networkish);
    getUtxos(): Promise<void>;
    send(method: string, params: Array<any>): Promise<any>;
}
