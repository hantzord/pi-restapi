package midtrans

import "capstone/entities/transaction"

type MidtransUseCase interface {
	GenerateSnapURL(transaction *transaction.Transaction) (*transaction.Transaction, error)
	VerifyPayment(orderID string) (string, error)
	BankTransfer(transaction *transaction.Transaction) (*transaction.Transaction, error)
	EWallet(transaction *transaction.Transaction) (*transaction.Transaction, error)
}
