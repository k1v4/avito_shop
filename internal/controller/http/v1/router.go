package v1

import (
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *echo.Echo, l logger.Logger, t usecase.IShopService) {
	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	h := handler.Group("/api")
	{
		newShopRoutes(h, t, l)
	}
}
