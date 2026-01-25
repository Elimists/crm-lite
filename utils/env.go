package utils

import (
	"log"
	"os"
)

var requiredEnv = []string{"MAILTRAP_API_TOKEN", "MAILTRAP_ENDPOINT", "TO_EMAIL", "TO_NAME", "FROM_EMAIL", "FROM_NAME"}

func VerifyRequiredEnv() {
	for _, v := range requiredEnv {
		if os.Getenv(v) == "" {
			log.Fatalf("[INIT_ERR] Environment variable %s is not set. Please set it before running the application.", v)
		}
	}
}
