package response

type ConsultationCount struct {
	TotalConsultation    int64 `json:"total_consultation"`
	TodayConsultation    int64 `json:"today_consultation"`
	ActiveConsultation   int64 `json:"active_consultation"`
	DoneConsultation     int64 `json:"done_consultation"`
	RejectedConsultation int64 `json:"rejected_consultation"`
	IncomingConsultation int64 `json:"incoming_consultation"`
	PendingConsultation  int64 `json:"pending_consultation"`
}
