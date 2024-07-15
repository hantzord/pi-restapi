package payment

import "capstone/entities/transaction"

type Method interface {
	EWallet(transaction *transaction.Transaction) (*transaction.Transaction, error)
	BankTransfer(transaction *transaction.Transaction) (*transaction.Transaction, error)
}
