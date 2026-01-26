package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func VerifyHMAC(body []byte, signature string, secret string) bool {

	if secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	expectedMAC := mac.Sum(nil)

	receivedMAC, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(receivedMAC, expectedMAC)
}
