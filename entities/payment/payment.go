package payment

import (
	"capstone/constants"
	transactionEntities "capstone/entities/transaction"
	"github.com/google/uuid"
)

type Payment struct {
	PaymentType        string            `json:"payment_type"`
	TransactionDetails TransactionDetail `json:"transaction_details"`
	BankTransfer       BankTransfer      `json:"bank_transfer"`
}

type TransactionDetail struct {
	OrderID     uuid.UUID `json:"order_id"`
	GrossAmount int       `json:"gross_amount"`
}

type EWalletPayment struct {
	PaymentType        string                   `json:"payment_type"`
	TransactionDetails EWalletTransactionDetail `json:"transaction_details"`
}

type EWalletTransactionDetail struct {
	OrderID     uuid.UUID `json:"order_id"`
	GrossAmount int       `json:"gross_amount"`
}

type BankTransfer struct {
	Bank string `json:"bank"`
}

func ToEWallet(transaction *transactionEntities.Transaction) *EWalletPayment {
	return &EWalletPayment{
		PaymentType: constants.GoPay,
		TransactionDetails: EWalletTransactionDetail{
			OrderID:     transaction.ID,
			GrossAmount: transaction.Price,
		},
	}
}

func ToBankTransfer(transaction *transactionEntities.Transaction) *Payment {
	return &Payment{
		PaymentType: constants.BankTransfer,
		TransactionDetails: TransactionDetail{
			OrderID:     transaction.ID,
			GrossAmount: transaction.Price,
		},
		BankTransfer: BankTransfer{
			Bank: transaction.Bank,
		},
	}
}
