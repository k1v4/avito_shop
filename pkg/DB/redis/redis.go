package redis

import (
	"context"
	"fmt"
	"github.com/k1v4/avito_shop/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Port        string        `env:"REDIS_PORT" envDefault:"6379"`
	Host        string        `env:"REDIS_HOST" envDefault:"localhost"`
	Password    string        `env:"REDIS_PASSWORD" envDefault:"123"`
	User        string        `env:"REDIS_USER" envDefault:"root"`
	DB          int           `env:"REDIS_DB" envDefault:"0"`
	MaxRetries  int           `env:"REDIS_MAX_RETRIES" envDefault:"3"`
	DialTimeout time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	Timeout     time.Duration `env:"REDIS_TIMEOUT" envDefault:"5s"`
}

func NewClient(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	l := logger.GetLoggerFromContext(ctx)

	pong, err := db.Ping(ctx).Result()
	if err != nil {
		l.Error(ctx, fmt.Sprintf("failed to connect to redis server: %s\n", err.Error()))

		return nil, err
	}

	l.Info(ctx, fmt.Sprintf("%s connected to redis server: %s\n", pong, cfg.Port))

	return db, nil
}
