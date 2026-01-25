package utils

import (
	"context"
	"crm-lite/models"
)

type clientContextKeyType struct{}

var clientContextKey = clientContextKeyType{}

func ClientFromContext(ctx context.Context) (models.ClientConfig, bool) {
	client, ok := ctx.Value(clientContextKey).(models.ClientConfig)
	return client, ok
}
