import { resolveProperties, Logger, } from "ethers/lib/utils";
import { BigNumber } from "bignumber.js";
import { checkTransactionType, serializeTransaction } from './helpers/utils';
import { GLOBAL_VARS } from './helpers/global-vars';
import { IntermediateWallet } from './helpers/IntermediateWallet';
const logger = new Logger("QtumWallet");
const forwardErrors = [
    Logger.errors.INSUFFICIENT_FUNDS
];
export class QtumWallet extends IntermediateWallet {
    constructor(privateKey, provider, signer) {
        super(privateKey, provider);
        this.signer = signer;
    }
    /**
     * Override to build a raw QTUM transaction signing UTXO's
     */
    async signTransaction(transaction) {
        const tx = await resolveProperties(transaction);
        // Refactored to check TX type (call, create, p2pkh, deploy error) and calculate needed amount
        const { transactionType, neededAmount } = checkTransactionType(tx);
        // Check if the transactionType matches the DEPLOY_ERROR, throw error else continue
        if (transactionType === GLOBAL_VARS.DEPLOY_ERROR) {
            return logger.throwError("You cannot send QTUM while deploying a contract. Try deploying again without a value.", Logger.errors.NOT_IMPLEMENTED, {
                error: "You cannot send QTUM while deploying a contract. Try deploying again without a value.",
            });
        }
        let utxos = [];
        try {
            console.log("Getting utxos, needed", neededAmount);
            // @ts-ignore
            utxos = await this.provider.getUtxos(tx.from, neededAmount);
            // Grab vins for transaction object.
        }
        catch (error) {
            if (forwardErrors.indexOf(error.code) >= 0) {
                throw error;
            }
            return logger.throwError("Needed amount of UTXO's exceed the total you own.", Logger.errors.INSUFFICIENT_FUNDS, {
                error: error,
            });
        }
        const { serializedTransaction, networkFee } = serializeTransaction(utxos, neededAmount, tx, transactionType, this.privateKey, this.publicKey, this.signer);
        if (networkFee !== "") {
            try {
                // Try again with the network fee included
                const updatedNeededAmount = new BigNumber(neededAmount).plus(networkFee);
                console.log("getting utxos2, need", updatedNeededAmount, networkFee);
                // @ts-ignore
                utxos = await this.provider.getUtxos(tx.from, updatedNeededAmount);
                // Grab vins for transaction object.
            }
            catch (error) {
                if (forwardErrors.indexOf(error.code) >= 0) {
                    throw error;
                }
                return logger.throwError("Needed amount of UTXO's exceed the total you own.", Logger.errors.INSUFFICIENT_FUNDS, {
                    error: error,
                });
            }
            const serialized = await serializeTransaction(utxos, neededAmount, tx, transactionType, this.publicKey, this.privateKey, this.signer);
            return serialized.serializedTransaction;
        }
        return serializedTransaction;
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiUXR1bVdhbGxldC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvUXR1bVdhbGxldC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQSxPQUFPLEVBQ0gsaUJBQWlCLEVBQ2pCLE1BQU0sR0FDVCxNQUFNLGtCQUFrQixDQUFDO0FBSzFCLE9BQU8sRUFBRSxTQUFTLEVBQUUsTUFBTSxjQUFjLENBQUE7QUFDeEMsT0FBTyxFQUFFLG9CQUFvQixFQUFFLG9CQUFvQixFQUFFLE1BQU0saUJBQWlCLENBQUE7QUFDNUUsT0FBTyxFQUFFLFdBQVcsRUFBRSxNQUFNLHVCQUF1QixDQUFBO0FBQ25ELE9BQU8sRUFBRSxrQkFBa0IsRUFBRSxNQUFNLDhCQUE4QixDQUFBO0FBRWpFLE1BQU0sTUFBTSxHQUFHLElBQUksTUFBTSxDQUFDLFlBQVksQ0FBQyxDQUFDO0FBQ3hDLE1BQU0sYUFBYSxHQUFHO0lBQ2xCLE1BQU0sQ0FBQyxNQUFNLENBQUMsa0JBQWtCO0NBQ25DLENBQUM7QUFHRixNQUFNLE9BQU8sVUFBVyxTQUFRLGtCQUFrQjtJQUk5QyxZQUFZLFVBQTJELEVBQUUsUUFBbUIsRUFBRSxNQUFlO1FBQ3pHLEtBQUssQ0FBQyxVQUFVLEVBQUUsUUFBUSxDQUFDLENBQUE7UUFDM0IsSUFBSSxDQUFDLE1BQU0sR0FBRyxNQUFNLENBQUM7SUFDekIsQ0FBQztJQUVEOztPQUVHO0lBQ0gsS0FBSyxDQUFDLGVBQWUsQ0FBQyxXQUErQjtRQUNqRCxNQUFNLEVBQUUsR0FBRyxNQUFNLGlCQUFpQixDQUFDLFdBQVcsQ0FBQyxDQUFDO1FBRWhELDhGQUE4RjtRQUM5RixNQUFNLEVBQUUsZUFBZSxFQUFFLFlBQVksRUFBRSxHQUFHLG9CQUFvQixDQUFDLEVBQUUsQ0FBQyxDQUFDO1FBRW5FLG1GQUFtRjtRQUNuRixJQUFJLGVBQWUsS0FBSyxXQUFXLENBQUMsWUFBWSxFQUFFO1lBQzlDLE9BQU8sTUFBTSxDQUFDLFVBQVUsQ0FDcEIsdUZBQXVGLEVBQ3ZGLE1BQU0sQ0FBQyxNQUFNLENBQUMsZUFBZSxFQUM3QjtnQkFDSSxLQUFLLEVBQUUsdUZBQXVGO2FBQ2pHLENBQ0osQ0FBQztTQUNMO1FBRUQsSUFBSSxLQUFLLEdBQUcsRUFBRSxDQUFDO1FBQ2YsSUFBSTtZQUNBLE9BQU8sQ0FBQyxHQUFHLENBQUMsdUJBQXVCLEVBQUUsWUFBWSxDQUFDLENBQUE7WUFDbEQsYUFBYTtZQUNiLEtBQUssR0FBRyxNQUFNLElBQUksQ0FBQyxRQUFRLENBQUMsUUFBUSxDQUFDLEVBQUUsQ0FBQyxJQUFJLEVBQUUsWUFBWSxDQUFDLENBQUM7WUFDNUQsb0NBQW9DO1NBQ3ZDO1FBQUMsT0FBTyxLQUFVLEVBQUU7WUFDakIsSUFBSSxhQUFhLENBQUMsT0FBTyxDQUFDLEtBQUssQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLEVBQUU7Z0JBQ3hDLE1BQU0sS0FBSyxDQUFDO2FBQ2Y7WUFDRCxPQUFPLE1BQU0sQ0FBQyxVQUFVLENBQ3BCLG1EQUFtRCxFQUNuRCxNQUFNLENBQUMsTUFBTSxDQUFDLGtCQUFrQixFQUNoQztnQkFDSSxLQUFLLEVBQUUsS0FBSzthQUNmLENBQ0osQ0FBQztTQUNMO1FBRUQsTUFBTSxFQUFFLHFCQUFxQixFQUFFLFVBQVUsRUFBRSxHQUFHLG9CQUFvQixDQUFDLEtBQUssRUFBRSxZQUFZLEVBQUUsRUFBRSxFQUFFLGVBQWUsRUFBRSxJQUFJLENBQUMsVUFBVSxFQUFFLElBQUksQ0FBQyxTQUFTLEVBQUUsSUFBSSxDQUFDLE1BQU0sQ0FBQyxDQUFDO1FBRTNKLElBQUksVUFBVSxLQUFLLEVBQUUsRUFBRTtZQUNuQixJQUFJO2dCQUNBLDBDQUEwQztnQkFDMUMsTUFBTSxtQkFBbUIsR0FBRyxJQUFJLFNBQVMsQ0FBQyxZQUFZLENBQUMsQ0FBQyxJQUFJLENBQUMsVUFBVSxDQUFDLENBQUM7Z0JBQ3pFLE9BQU8sQ0FBQyxHQUFHLENBQUMsc0JBQXNCLEVBQUUsbUJBQW1CLEVBQUUsVUFBVSxDQUFDLENBQUE7Z0JBQ3BFLGFBQWE7Z0JBQ2IsS0FBSyxHQUFHLE1BQU0sSUFBSSxDQUFDLFFBQVEsQ0FBQyxRQUFRLENBQUMsRUFBRSxDQUFDLElBQUksRUFBRSxtQkFBbUIsQ0FBQyxDQUFDO2dCQUNuRSxvQ0FBb0M7YUFDdkM7WUFBQyxPQUFPLEtBQVUsRUFBRTtnQkFDakIsSUFBSSxhQUFhLENBQUMsT0FBTyxDQUFDLEtBQUssQ0FBQyxJQUFJLENBQUMsSUFBSSxDQUFDLEVBQUU7b0JBQ3hDLE1BQU0sS0FBSyxDQUFDO2lCQUNmO2dCQUNELE9BQU8sTUFBTSxDQUFDLFVBQVUsQ0FDcEIsbURBQW1ELEVBQ25ELE1BQU0sQ0FBQyxNQUFNLENBQUMsa0JBQWtCLEVBQ2hDO29CQUNJLEtBQUssRUFBRSxLQUFLO2lCQUNmLENBQ0osQ0FBQzthQUNMO1lBQ0QsTUFBTSxVQUFVLEdBQUcsTUFBTSxvQkFBb0IsQ0FBQyxLQUFLLEVBQUUsWUFBWSxFQUFFLEVBQUUsRUFBRSxlQUFlLEVBQUUsSUFBSSxDQUFDLFNBQVMsRUFBRSxJQUFJLENBQUMsVUFBVSxFQUFFLElBQUksQ0FBQyxNQUFNLENBQUMsQ0FBQztZQUN0SSxPQUFPLFVBQVUsQ0FBQyxxQkFBcUIsQ0FBQztTQUMzQztRQUVELE9BQU8scUJBQXFCLENBQUM7SUFDakMsQ0FBQztDQUNKIn0=