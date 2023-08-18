package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/slavik22/blogRestApi/lib/types"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/service"
	"net/http"
)

// UserController ...
type UserController struct {
	ctx      context.Context
	services *service.Manager
}

// NewUsers creates a new user controller.
func NewUserController(ctx context.Context, services *service.Manager) *UserController {
	return &UserController{
		ctx:      ctx,
		services: services,
	}
}

func (u *UserController) SignUp(c echo.Context) error {
	var input model.User

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode user data"))
	}

	err := c.Validate(&input)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	createdUser, err := u.services.UserService.CreateUser(input)

	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create user"))
		}
	}

	return c.JSON(http.StatusCreated, createdUser)
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *UserController) SignIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	token, err := u.services.UserService.SignIn(input.Email, input.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create user"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
