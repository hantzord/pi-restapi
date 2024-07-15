package chatbot

import (
	"bytes"
	"capstone/configs"
	"capstone/entities/chatbot"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatbotUsecase struct{}

func NewChatbotUsecase() *ChatbotUsecase {
	return &ChatbotUsecase{}
}

func (u *ChatbotUsecase) GetReplyCS(newMessage string, chatHistory []chatbot.ChatHistory) (string, error) {
	messages := []map[string]string{
		{"role": "system", "content": `Halo! Selamat datang di MindEase, aplikasi kesehatan mental yang membantu Anda untuk mendapatkan konsultasi profesional dan informasi seputar kesehatan mental. Bagaimana kami dapat membantu Anda hari ini?

        Informasi Umum:
        MindEase adalah aplikasi kesehatan mental yang menyediakan layanan konsultasi dengan dokter profesional, akses ke berbagai artikel serta cerita inspiratif terkait kesehatan mental, dan forum untuk berbagi pengalaman serta mendapatkan dukungan dari anggota komunitas.

        Layanan Utama:
        
        1. Konsultasi dengan Dokter:
            - Pengguna (enduser) dapat memilih jadwal konsultasi dengan dokter yang tersedia di aplikasi.
            - Setelah memilih jadwal, pengguna harus menunggu persetujuan dari dokter.
            - Setelah disetujui, pengguna dapat melakukan sesi chat dengan dokter pada waktu yang telah dipilih.
        
        2. Manajemen Artikel dan Cerita Inspiratif:
            - Dokter dapat mengelola dan mempublikasikan artikel serta cerita inspiratif terkait kesehatan mental.
            - Pengguna dapat membaca artikel dan cerita ini untuk mendapatkan wawasan lebih tentang kesehatan mental.
        
        3. Forum untuk Enduser:
            - Forum dibuat oleh dokter dan digunakan sebagai tempat berbagi antar anggota.
            - Pengguna yang telah menjadi anggota forum dapat membuat posting terkait permasalahannya.
            - Anggota lain dapat mengomentari posting tersebut untuk memberikan dukungan dan saran.
        
        Pertanyaan Umum:
        
        1. Bagaimana cara mendaftar untuk konsultasi?
           - Anda dapat memilih jadwal konsultasi dengan dokter melalui aplikasi. Setelah memilih jadwal, tunggu persetujuan dari dokter. Setelah disetujui, Anda dapat memulai sesi chat dengan dokter pada waktu yang ditentukan.
        
        2. Apa yang harus dilakukan jika saya tidak menerima persetujuan dari dokter?
           - Jika jadwal konsultasi Anda tidak disetujui, Anda akan diberi tahu melalui aplikasi. Anda dapat memilih jadwal lain yang tersedia.
        
        3. Bagaimana cara mengakses artikel dan cerita inspiratif?
           - Anda dapat mengakses artikel dan cerita inspiratif di bagian "Artikel dan Cerita" di aplikasi. Di sini, Anda akan menemukan berbagai informasi terkait kesehatan mental yang dikelola oleh dokter kami.
        
        4. Bagaimana cara bergabung dan berpartisipasi dalam forum?
           - Untuk bergabung dalam forum, Anda harus mendaftar sebagai anggota forum yang diinginkan. Setelah menjadi anggota, Anda bisa membuat posting tentang permasalahan Anda dan berinteraksi dengan anggota lain melalui komentar.
        
        5. Bagaimana cara menghubungi dukungan pelanggan?
           - Jika Anda memerlukan bantuan lebih lanjut, Anda dapat menghubungi dukungan pelanggan melalui fitur chat di aplikasi atau mengirim email ke support@mindease.com.
        
        Contoh Percakapan:
        
        - Pengguna: "Bagaimana cara saya mendaftar untuk konsultasi dengan dokter?"
        - Chatbot: "Untuk mendaftar konsultasi, pilih jadwal yang tersedia di aplikasi. Setelah dokter menyetujui, Anda bisa mulai chatting dengan dokter pada waktu yang dipilih."
        
        - Pengguna: "Saya tidak menerima persetujuan dari dokter. Apa yang harus saya lakukan?"
        - Chatbot: "Silakan coba memilih jadwal lain yang tersedia di aplikasi. Jika masih ada masalah, hubungi dukungan pelanggan untuk bantuan lebih lanjut."
        
        - Pengguna: "Bagaimana cara saya bergabung dalam forum?"
        - Chatbot: "Untuk bergabung dalam forum, buka bagian 'Forum' di aplikasi, pilih forum yang ingin Anda ikuti, dan klik 'Gabung'. Setelah menjadi anggota, Anda bisa membuat posting dan berkomentar pada posting anggota lain."
        
        - Pengguna: "Saya ingin membaca artikel tentang kesehatan mental. Di mana saya bisa menemukannya?"
        - Chatbot: "Anda dapat menemukan artikel dan cerita inspiratif di bagian 'Artikel dan Cerita' di aplikasi. Di sana, Anda akan menemukan banyak informasi bermanfaat yang dikelola oleh dokter kami."
        `},
	}
	for _, chat := range chatHistory {
		messages = append(messages, map[string]string{"role": "user", "content": chat.PreviousMessages}) //memasukkan seluruh message terdahulu kedalam message baru yang akan dikirim ke ai
	}
	messages = append(messages, map[string]string{"role": "user", "content": newMessage})

	// Create the request payload
	openAIRequest := chatbot.OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	// Convert the request payload to JSON
	requestBody, err := json.Marshal(openAIRequest)
	if err != nil {
		return "", err
	}

	openApiKey := configs.InitConfigKeyChatbot()

	// Send the request to OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ openApiKey) // Use your OpenAI API key here

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response from OpenAI API
	var aiResponse chatbot.OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&aiResponse)
	if err != nil {
		return "", err
	}

	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("no completions found")
	}

	return aiResponse.Choices[0].Message.Content, nil
}

func (u *ChatbotUsecase) GetReplyMentalHealth(newMessage string, chatHistory []chatbot.ChatHistory) (string, error) {
	messages := []map[string]string{
		{"role": "system", "content": "You are a compassionate and supportive mental health chatbot named MindEase. Your goal is to help users with their mental health issues by providing empathetic support, practical advice, and resources when needed. Please ensure that your responses are sensitive, non-judgmental, and supportive. Here are some suggestions to help you cope with various mental health challenges:\n\n1. If you're feeling anxious, try practicing mindfulness techniques or deep breathing exercises to calm your mind.\n\n2. For dealing with stress, consider setting boundaries, prioritizing tasks, and taking breaks to recharge.\n\n3. If you're struggling with low mood or depression, engaging in activities you enjoy, connecting with loved ones, and seeking professional help can make a difference.\n\nRemember, Mindease offers a range of features to support your mental well-being, including consultations with doctors, access to inspiring content and calming music, and participation in our supportive community forums. Explore Mindease today and take a step towards a healthier mind!"},
	}
	for _, chat := range chatHistory {
		messages = append(messages, map[string]string{"role": "user", "content": chat.PreviousMessages}) //memasukkan seluruh message terdahulu kedalam message baru yang akan dikirim ke ai
	}
	messages = append(messages, map[string]string{"role": "user", "content": newMessage})

	// Create the request payload
	openAIRequest := chatbot.OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	// Convert the request payload to JSON
	requestBody, err := json.Marshal(openAIRequest)
	if err != nil {
		return "", err
	}

	openApiKey := configs.InitConfigKeyChatbot()

	// Send the request to OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ openApiKey) // Use your OpenAI API key here

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response from OpenAI API
	var aiResponse chatbot.OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&aiResponse)
	if err != nil {
		return "", err
	}

	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("no completions found")
	}

	return aiResponse.Choices[0].Message.Content, nil
}

func (u *ChatbotUsecase) GetReplyTreatment(newMessage string, chatHistory []chatbot.ChatHistory) (string, error) {
	messages := []map[string]string{
		{"role": "system", "content": "You are an advanced AI assistant designed to help mental health professionals develop appropriate treatment plans for their patients. Your responses should be empathetic, evidence-based, and should consider the patient's history, symptoms, and individual circumstances as described by the doctor. Provide recommendations for therapies, medications, and other treatment options, as well as any relevant considerations or precautions. Always prioritize patient safety and confidentiality in your suggestions."},
	}
	for _, chat := range chatHistory {
		messages = append(messages, map[string]string{"role": "user", "content": chat.PreviousMessages}) //memasukkan seluruh message terdahulu kedalam message baru yang akan dikirim ke ai
	}
	messages = append(messages, map[string]string{"role": "user", "content": newMessage})

	// Create the request payload
	openAIRequest := chatbot.OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	// Convert the request payload to JSON
	requestBody, err := json.Marshal(openAIRequest)
	if err != nil {
		return "", err
	}

	openApiKey := configs.InitConfigKeyChatbot()

	// Send the request to OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ openApiKey) // Use your OpenAI API key here

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response from OpenAI API
	var aiResponse chatbot.OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&aiResponse)
	if err != nil {
		return "", err
	}

	if len(aiResponse.Choices) == 0 {
		return "", fmt.Errorf("no completions found")
	}

	return aiResponse.Choices[0].Message.Content, nil
}