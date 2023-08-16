package service

import (
	"context"
	"github.com/slavik22/blogRestApi/model"
)

// UserService is a service for users
//
//go:generate mockery --dir . --name UserService --output ./mocks
type UserServ interface {
	GetUser(context.Context, uint) (*model.User, error)
	CreateUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, uint) error
}
