package midtrans

import (
	"capstone/entities/transaction"
	"strconv"
	"strings"
)

type BankTransfer struct {
	StatusCode        string           `json:"status_code"`
	StatusMessage     string           `json:"status_message"`
	TransactionID     string           `json:"transaction_id"`
	OrderID           string           `json:"order_id"`
	MerchantID        string           `json:"merchant_id"`
	GrossAmount       string           `json:"gross_amount"`
	Currency          string           `json:"currency"`
	PaymentType       string           `json:"payment_type"`
	TransactionTime   string           `json:"transaction_time"`
	TransactionStatus string           `json:"transaction_status"`
	VANumbers         []VirtualAccount `json:"va_numbers"`
	FraudStatus       string           `json:"fraud_status"`
}

type VirtualAccount struct {
	Bank     string `json:"bank"`
	VaNumber string `json:"va_number"`
}

func (r *BankTransfer) ToTransaction(trans *transaction.Transaction) (*transaction.Transaction, error) {
	newGrossAmount := strings.Split(r.GrossAmount, ".")
	price, err := strconv.Atoi(newGrossAmount[0])
	if err != nil {
		return nil, err
	}
	var paymentLink string
	var bank string
	for _, account := range r.VANumbers {
		paymentLink = account.VaNumber
		bank = account.Bank
	}
	return &transaction.Transaction{
		ID:             trans.ID,
		ConsultationID: trans.ConsultationID,
		Consultation:   trans.Consultation,
		Price:          price,
		Status:         r.TransactionStatus,
		PaymentType:    r.PaymentType,
		PaymentLink:    paymentLink,
		Bank:           bank,
	}, nil
}
