package utils

import (
	"context"
	"crm-lite/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func SendEmailNotification(ctx context.Context, c models.Contact) error {
	clientContextModel, _ := ClientFromContext(ctx)

	apiEndpoint := clientContextModel.EmailApiEndpoint
	apiToken := os.Getenv(clientContextModel.EmailApiEnvTokenName)

	fromEmail := "no-reply@proreact.dev"
	fromName := "Proreact - Pran Pandey"

	toEmail := clientContextModel.Email
	toName := clientContextModel.Name

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
	req, err := http.NewRequest("POST", apiEndpoint, payload)
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

	if resp.Status == "200 OK" {
		fmt.Println("email notification sent to: " + toName)
	}
	return nil
}
