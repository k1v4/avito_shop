package redis

import (
	"context"
	"github.com/k1v4/avito_shop/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	cfg := RedisConfig{
		Host:        "localhost",
		Port:        "6379",
		Password:    "",
		User:        "root",
		DB:          0,
		MaxRetries:  3,
		DialTimeout: 5 * time.Second,
		Timeout:     5 * time.Second,
	}

	l := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, l)

	client, err := NewClient(ctx, cfg)

	assert.NoError(t, err)
	assert.NotNil(t, client)

	pong, err := client.Ping(ctx).Result()
	assert.NoError(t, err)
	assert.Equal(t, "PONG", pong)

	defer func() {
		if err := client.Close(); err != nil {
			t.Fatalf("could not close redis client: %v", err)
		}
	}()
}
