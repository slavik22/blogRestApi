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

// PostRepo is a store for posts
//
//go:generate mockery --dir . --name PostRepo --output ./mocks
type PostRepo interface {
	GetPosts(context.Context) ([]model.Post, error)
	GetPost(context.Context, uint, uint) (*model.Post, error)
	CreatePost(context.Context, *model.Post) (uint, error)
	UpdatePost(context.Context, *model.Post) (*model.Post, error)
	DeletePost(context.Context, uint, uint) error
}
