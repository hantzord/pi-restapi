package midtrans

import (
	"capstone/entities/transaction"
	"strconv"
	"strings"
)

type EWallet struct {
	StatusCode        string   `json:"status_code"`
	StatusMessage     string   `json:"status_message"`
	TransactionID     string   `json:"transaction_id"`
	OrderID           string   `json:"order_id"`
	GrossAmount       string   `json:"gross_amount"`
	Currency          string   `json:"currency"`
	PaymentType       string   `json:"payment_type"`
	TransactionTime   string   `json:"transaction_time"`
	TransactionStatus string   `json:"transaction_status"`
	FraudStatus       string   `json:"fraud_status"`
	Actions           []Action `json:"actions"`
}

type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

func (r *EWallet) ToTransaction(trans *transaction.Transaction) (*transaction.Transaction, error) {
	newGrossAmount := strings.Split(r.GrossAmount, ".")
	price, err := strconv.Atoi(newGrossAmount[0])
	if err != nil {
		return nil, err
	}
	return &transaction.Transaction{
		ID:             trans.ID,
		ConsultationID: trans.ConsultationID,
		Consultation:   trans.Consultation,
		Price:          price,
		Status:         r.TransactionStatus,
		PaymentType:    "gopay",
		PaymentLink:    r.Actions[0].URL,
		Bank:           "ewallet",
	}, nil
}
