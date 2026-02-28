package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func NotifyTenant() error {
	apiEndpoint := os.Getenv("MAIL_ENDPOINT")
	apiToken := os.Getenv("MAIL_TOKEN")

	fromEmail := "no-reply@proreact.dev"
	fromName := "ProReact Notification System"
	toEmail := " "
	toName := " " + "'s admin"

	phone := ""

	payload := map[string]interface{}{
		"to": []map[string]string{
			{"email": toEmail, "name": toName},
		},
		"from": map[string]string{
			"email": fromEmail, "name": fromName,
		},
		"subject":  "New Contact Form Submission",
		"text":     fmt.Sprintf("New contact submitted:\n\nName: %s\nEmail: %s\nPhone: %s\nMessage: %s", "name", "email", phone, "message"),
		"category": "Weather Wizards Contact Form",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal mail payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Token", apiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mail API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[LOG_MAIL] Email request status %s. Email sent to %s (%s)", resp.Status, toName, toEmail)
	return nil
}
