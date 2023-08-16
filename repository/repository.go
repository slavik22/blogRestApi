package repository

import (
	"context"
	"github.com/slavik22/blogRestApi/model"
)

// UserRepo is a repository for users
//
//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	GetUser(context.Context, string, string) (*model.User, error)
	CreateUser(context.Context, *model.User) (uint, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, uint) error
}
