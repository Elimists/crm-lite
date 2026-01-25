package utils

import (
	"crm-lite/models"
	"encoding/json"
	"log"
	"os"
	"strings"
)

var clientsByOrigin map[string]models.ClientConfig

func LoadClientConfigs(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var clients []models.ClientConfig
	if err := json.Unmarshal(data, &clients); err != nil {
		return err
	}

	clientsByOrigin = make(map[string]models.ClientConfig)
	for _, c := range clients {
		origin := normalizeOrigin(c.Origin)
		clientsByOrigin[origin] = c
		log.Printf("id=%d, origin=%s, env_tok=%s\n", c.ID, c.Origin, c.EmailApiEnvTokenName)
	}

	return nil
}

func GetClientByOrigin(origin string) (models.ClientConfig, bool) {
	client, ok := clientsByOrigin[normalizeOrigin(origin)]
	return client, ok
}

func normalizeOrigin(origin string) string {
	return strings.TrimRight(origin, "/")
}

func VerifyRequiredEnv() {
	for _, v := range clientsByOrigin {
		if os.Getenv(v.EmailApiEnvTokenName) == "" {
			log.Fatalf("[INIT_ERR] Environment variable %s is not set. Please set it before running the application.", v.EmailApiEnvTokenName)
		}
	}
}
