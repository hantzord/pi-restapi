package utilities

import (
	"bytes"
	"capstone/configs"
	"io"
	"net/http"
)

func PaymentMidtrans(payload []byte, midtrans *configs.Midtrans) ([]byte, error) {
	req, err := http.NewRequest("POST", midtrans.BaseURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+midtrans.ServerKeyBase64)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	return body, nil

}
