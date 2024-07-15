package response

type TransactionCount struct {
	TotalTransaction   int64 `json:"total_transaction"`
	TransactionSuccess int64 `json:"transaction_success"`
	TransactionFailed  int64 `json:"transaction_Failed"`
	TransactionToday   int64 `json:"transaction_today"`
	TransactionWeek    int64 `json:"transaction_week"`
	TransactionMonth   int64 `json:"transaction_month"`
	TransactionYear    int64 `json:"transaction_year"`
}

func ToTransactionCount(totalTransaction, transactionSuccess, transactionFailed, transactionToday, transactionWeek, transactionMonth, transactionYear int64) *TransactionCount {
	return &TransactionCount{
		TotalTransaction:   totalTransaction,
		TransactionSuccess: transactionSuccess,
		TransactionFailed:  transactionFailed,
		TransactionToday:   transactionToday,
		TransactionWeek:    transactionWeek,
		TransactionMonth:   transactionMonth,
		TransactionYear:    transactionYear,
	}
}
