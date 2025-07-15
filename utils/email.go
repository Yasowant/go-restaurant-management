package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// EmailPayload represents the Resend email body
type EmailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

// SendEmail sends email via Resend API (sandbox: only to your email)
func SendEmail(to, subject, htmlBody string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("❌ RESEND_API_KEY is not set in .env")
	}

	// You can only send to your own email if using onboarding@resend.dev
	payload := EmailPayload{
		From:    "onboarding@resend.dev", // 🔒 Resend sandbox sender
		To:      to,
		Subject: subject,
		Html:    htmlBody,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("❌ Failed to encode payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("❌ Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("❌ Request error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return fmt.Errorf("❌ Email sending failed: %s - %s", resp.Status, string(body))
	}

	// Optional: Log in development
	fmt.Printf("✅ Email sent to %s\n", to)
	return nil
}
