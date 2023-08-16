package service

import (
	"github.com/slavik22/blogRestApi/model"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserServ interface {
	CreateUser(user model.User) (uint, error)
	ParseToken(accessToken string) (uint, error)
	GenerateToken(email, password string) (string, error)
	//UpdateUser(context.Context, *model.User) (*model.User, error)
	//DeleteUser(context.Context, uint) error
}

type PostServ interface {
	GetPosts() ([]model.Post, error)
	GetPost(postId uint, userId uint) (*model.Post, error)
	CreatePost(post model.Post, userId uint) (uint, error)
	UpdatePost(post model.Post) (*model.Post, error)
	DeletePost(postId uint, userId uint) error
}
