package utils

import (
	"crm-lite/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func sendMailtrapEmail(c models.Contact) error {
	apiToken := os.Getenv("MAILTRAP_API_TOKEN")
	toEmail := os.Getenv("TO_EMAIL")
	toName := os.Getenv("TO_NAME")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")

	payloadStr := fmt.Sprintf(`{
		"to": [
			{
				"email":"%s",
				"name":"%s"
			}
		],
		"from": {
			"email":"%s",
			"name":"%s"
		},
		"subject":"New Contact Form Submission",
		"text":"New contact submitted:\n\nName: %s\nEmail: %s\nPhone: %s\nMessage: %s",
		"category":"Weather Wizards Contact Form"
	}`, toEmail, toName, fromEmail, fromName, c.Name, c.Email, c.Phone, c.Message)

	payload := strings.NewReader(payloadStr)
	req, err := http.NewRequest("POST", os.Getenv("MAILTRAP_ENDPOINT"), payload)
	if err != nil {
		fmt.Println("[MAIL_ERROR] Failed to create request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Token", apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[MAIL_ERROR] Failed to send request:", err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[MAIL_ERROR] Failed to read response body:", err)
		return err
	}
	return nil
}
