import { Provider, TransactionRequest } from "@ethersproject/abstract-provider";
import { ExternallyOwnedAccount, Signer } from "@ethersproject/abstract-signer";
import { SigningKey } from "@ethersproject/signing-key";
import { BytesLike } from "@ethersproject/bytes";
import { IntermediateWallet } from './helpers/IntermediateWallet';
export declare class QtumWallet extends IntermediateWallet {
    readonly signer?: Signer;
    constructor(privateKey: BytesLike | ExternallyOwnedAccount | SigningKey, provider?: Provider, signer?: Signer);
    /**
     * Override to build a raw QTUM transaction signing UTXO's
     */
    signTransaction(transaction: TransactionRequest): Promise<string>;
}
