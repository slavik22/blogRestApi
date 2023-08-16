package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (u *UserController) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authorizationHeader)
		fmt.Println(header)
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "empty auth header")
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth header")
		}

		if len(headerParts[1]) == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "token is empty")
		}

		userId, err := u.services.UserService.ParseToken(headerParts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "token is incorrect")
		}

		c.Set(userCtx, userId)

		return next(c)
	}
}

func getUserId(c echo.Context) (uint, error) {
	id, ok := c.Get(userCtx).(uint)

	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return id, nil
}
