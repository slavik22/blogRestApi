package store

import (
	"context"
	"github.com/slavik22/blogRestApi/model"
)

// UserRepo is a store for users
//
//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	GetUser(context.Context, uint) (*model.User, error)
	CreateUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, uint) error
}
