package main

import (
	"context"
	"fmt"
	"github.com/k1v4/avito_shop/internal/config"
	v1 "github.com/k1v4/avito_shop/internal/controller/http/v1"
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/internal/usecase/repository"
	"github.com/k1v4/avito_shop/pkg/DB/postgres"
	"github.com/k1v4/avito_shop/pkg/DB/redis"
	"github.com/k1v4/avito_shop/pkg/httpserver"
	"github.com/k1v4/avito_shop/pkg/logger"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	ctx := context.Background()
	loggerBack := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, loggerBack)

	loggerBack.Info(ctx, "starting backend")

	cfg := config.MustLoadConfig()
	if cfg == nil {
		loggerBack.Error(ctx, "config is nil")
		return
	}

	clientRedis, err := redis.NewClient(ctx, cfg.RedisConfig)
	if err != nil {
		loggerBack.Error(ctx, "redis client init fail")
	}

	_ = clientRedis

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.UserName,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)

	pg, err := postgres.New(url, postgres.MaxPoolSize(cfg.DBConfig.PoolMax))
	if err != nil {
		loggerBack.Error(ctx, fmt.Sprintf("app - Run - postgres.New: %s", err))
	}
	defer pg.Close()

	loggerBack.Info(ctx, "connected to database successfully")

	containerUseCase := usecase.NewShopUseCase(
		repository.NewShopRepository(pg),
		clientRedis,
	)

	handler := echo.New()
	//handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins: []string{"http://localhost:3000", "http://10.255.196.171:3000"},
	//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	//}))
	v1.NewRouter(handler, loggerBack, containerUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(strconv.Itoa(cfg.RestServerPort)))

	// signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		loggerBack.Info(ctx, "app-Run-signal: "+s.String())
	case err = <-httpServer.Notify():
		loggerBack.Error(ctx, fmt.Sprintf("app-Run-httpServer.Notify: %s", err))
	}

	// shutdown
	err = httpServer.Shutdown()
	if err != nil {
		loggerBack.Error(ctx, fmt.Sprintf("app-Run-httpServer.Shutdown: %s", err))
	}
}
