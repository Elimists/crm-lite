package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateUserToken(
	subscriber string,
	tenant string,
	roles []string,
	scopes []string) (string, error) {

	signingKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    subscriber, // user
		"iss":    "crm-lite", // issuer
		"aud":    "crm-lite", // consumer
		"tenant": tenant,
		"roles":  roles,
		"scopes": scopes,
		"iat":    jwt.NewNumericDate(time.Now()),
		"nbf":    jwt.NewNumericDate(time.Now()),
		"exp":    jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (bool, jwt.MapClaims) {
	signingKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return signingKey, nil
	})
	if err != nil {
		return false, nil
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return false, nil
			}
		}
		return true, claims
	}

	return false, nil
}
