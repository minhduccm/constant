package transaction

import (
	"github.com/ninjadotorg/constant/common"
)

type TxBuyBackRequest struct {
	*BuyBackRequestInfo
	*Tx // fee
	// TODO: signature?
}

type BuyBackRequestInfo struct {
	BuyBackFromTxID *common.Hash
	VoutIndex       int
}

// CreateTxBuyBackRequest
// senderKey and paymentInfo is for paying fee
func CreateTxBuyBackRequest(
	feeArgs FeeArgs,
	buyBackRequestInfo *BuyBackRequestInfo,
) (*TxBuyBackRequest, error) {
	// Create tx for fee &
	tx, err := CreateTx(
		feeArgs.SenderKey,
		feeArgs.PaymentInfo,
		feeArgs.Rts,
		feeArgs.UsableTx,
		feeArgs.Commitments,
		feeArgs.Fee,
		feeArgs.SenderChainID,
		false,
	)
	if err != nil {
		return nil, err
	}

	txBuyBackRequest := &TxBuyBackRequest{
		BuyBackRequestInfo: buyBackRequestInfo,
		Tx:                 tx,
	}
	txBuyBackRequest.Type = common.TxBuyBackRequest
	return txBuyBackRequest, nil
}

func (tx *TxBuyBackRequest) Hash() *common.Hash {
	// get hash of tx
	record := tx.Tx.Hash().String()
	record += tx.BuyBackFromTxID.String()
	record += string(tx.VoutIndex)

	// final hash
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

func (tx *TxBuyBackRequest) ValidateTransaction() bool {
	// validate for normal tx
	if !tx.Tx.ValidateTransaction() {
		return false
	}
	return true
}

func (tx *TxBuyBackRequest) GetType() string {
	return tx.Tx.Type
}

func (tx *TxBuyBackRequest) GetTxVirtualSize() uint64 {
	// TODO: calculate
	return 0
}

func (tx *TxBuyBackRequest) GetSenderAddrLastByte() byte {
	return tx.Tx.AddressLastByte
}

func (tx *TxBuyBackRequest) GetTxFee() uint64 {
	return tx.Tx.Fee
}
