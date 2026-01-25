package models

type ClientConfig struct {
	ID                   int    `json: "id"`
	Origin               string `json:"origin"`
	Email                string `json:"email"`
	Name                 string `json:"name"`
	Company              string `json:"company"`
	EmailApiEndpoint     string `json:"email_api_endpoint"`
	EmailApiEnvTokenName string `json:"email_api_env_token_name"`
}
