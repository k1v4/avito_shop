package v1

import (
	"errors"
	"fmt"
	"github.com/k1v4/avito_shop/internal/entity"
	"github.com/k1v4/avito_shop/internal/usecase"
	"github.com/k1v4/avito_shop/pkg/jwtPkg"
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

	// POST /api/auth
	handler.POST("/auth", r.Auth)

	//GET /api/buy/{item}
	handler.GET("/buy/:item", r.Buy)

	//POST /api/sendCoin"
	handler.POST("/sendCoin", r.SendCoins)

	//GET  /api/info
	handler.GET("/info", r.Info)
}

func (r *conatainerRoutes) Info(c echo.Context) error {
	const op = "handler.Info"

	ctx := c.Request().Context()

	token := jwtPkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userId, err := jwtPkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	info, err := r.t.GetInfo(ctx, userId)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, info)
}

func (r *conatainerRoutes) SendCoins(c echo.Context) error {
	const op = "handler.SendCoins"

	ctx := c.Request().Context()

	u := new(entity.SendCoinRequest)
	if err := c.Bind(u); err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	if u.Amount <= 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, errors.New("amount must be greater than 0"))
	}

	token := jwtPkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userId, err := jwtPkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	err = r.t.SendCoins(ctx, u.ToUserName, userId, u.Amount)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, nil)
}

func (r *conatainerRoutes) Buy(c echo.Context) error {
	const op = "handler.Buy"
	ctx := c.Request().Context()
	itemName := c.Param("item")

	if len(strings.TrimSpace(itemName)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "item name is required")
	}

	token := jwtPkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// получаем ник
	userId, err := jwtPkg.ValidateTokenAndGetUserId(token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "bad request")

		return fmt.Errorf("%s: %s", op, err)
	}

	err = r.t.BuyItem(ctx, userId, itemName)
	if err != nil {
		if errors.Is(err, usecase.ErrNoCoins) {
			errorResponse(c, http.StatusBadRequest, usecase.ErrNoCoins.Error())

			return fmt.Errorf("%s: %s", op, err)
		}

		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, nil)
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
