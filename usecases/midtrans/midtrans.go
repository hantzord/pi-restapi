package midtrans

import (
	"capstone/configs"
	"capstone/constants"
	midtransEntities "capstone/entities/midtrans"
	"capstone/entities/payment"
	"capstone/entities/transaction"
	"capstone/utilities"
	"encoding/json"
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransUseCase struct {
	midtransConfig configs.Midtrans
	envi           midtrans.EnvironmentType
}

func NewMidtransUseCase(m configs.Midtrans) midtransEntities.MidtransUseCase {
	envi := midtrans.Sandbox
	if m.IsProd {
		envi = midtrans.Production
	}
	return &MidtransUseCase{
		midtransConfig: m,
		envi:           envi,
	}
}

func (usecase *MidtransUseCase) GenerateSnapURL(transaction *transaction.Transaction) (*transaction.Transaction, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.ID.String(),
			GrossAmt: int64(transaction.Price),
		},
	}

	var client snap.Client
	client.New(usecase.midtransConfig.ServerKey, usecase.envi)

	snapResponse, err := client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}
	transaction.PaymentLink = snapResponse.RedirectURL
	return transaction, nil
}

func (usecase *MidtransUseCase) VerifyPayment(orderID string) (string, error) {
	var client coreapi.Client
	client.New(usecase.midtransConfig.ServerKey, usecase.envi)

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(orderID)

	if e != nil {
		return constants.Deny, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					//return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" && transactionStatusResp.SettlementTime != "" {
				return constants.Success, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				return constants.Deny, nil
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				return constants.Failed, nil
			} else if transactionStatusResp.TransactionStatus == "pending" {
				return constants.Pending, nil
			}
		}
	}
	return constants.Deny, e
}

func (usecase *MidtransUseCase) BankTransfer(transaction *transaction.Transaction) (*transaction.Transaction, error) {
	payload, err := json.Marshal(payment.ToBankTransfer(transaction))
	if err != nil {
		return nil, err
	}

	body, err := utilities.PaymentMidtrans(payload, &usecase.midtransConfig)
	if err != nil {
		return nil, err
	}

	var responseBody midtransEntities.BankTransfer
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	response, err := responseBody.ToTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (usecase *MidtransUseCase) EWallet(transaction *transaction.Transaction) (*transaction.Transaction, error) {
	payload, err := json.Marshal(payment.ToEWallet(transaction))
	if err != nil {
		return nil, err
	}

	body, err := utilities.PaymentMidtrans(payload, &usecase.midtransConfig)
	if err != nil {
		return nil, err
	}

	var responseBody midtransEntities.EWallet
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}
	fmt.Println(responseBody.OrderID)
	fmt.Println(transaction.ID)
	response, err := responseBody.ToTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return response, nil
}
