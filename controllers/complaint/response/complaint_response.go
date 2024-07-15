package response

type ComplaintResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Age            int    `json:"age"`
	Gender         string `json:"gender"`
	Message        string `json:"message"`
	MedicalHistory string `json:"medical_history"`
}
