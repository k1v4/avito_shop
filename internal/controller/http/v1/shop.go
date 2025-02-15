package v1

import (
	"errors"
	"fmt"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type conatainerRoutes struct {
	t usecase.IShopService
	l logger.Logger
}

func newShopRoutes(handler *echo.Group, t usecase.IShopService, l logger.Logger) {
	r := &conatainerRoutes{t, l}

	// GET /api/auth
	handler.POST("/auth", r.Auth)
}

func (r *conatainerRoutes) Auth(c echo.Context) error {
	const op = "handler.Login"

	ctx := c.Request().Context()

	u := new(entity.AuthRequest)
	if err := c.Bind(u); err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	if len(strings.TrimSpace(u.Username)) == 0 || len(strings.TrimSpace(u.Password)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := r.t.Login(ctx, u.Username, u.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrWrongPassword) {
			errorResponse(c, http.StatusUnauthorized, "bad request")

			return fmt.Errorf("%s: %w", op, err)
		}

		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, entity.AuthResponse{
		Token: token,
	})
}
