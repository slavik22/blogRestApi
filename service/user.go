package service

import (
	"context"
	"github.com/slavik22/blogRestApi/lib/util"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/repository"
)

type UserService struct {
	ctx   context.Context
	store *repository.Store
}

func NewUserService(ctx context.Context, store *repository.Store) *UserService {
	return &UserService{
		ctx:   ctx,
		store: store,
	}
}

func (s *UserService) CreateUser(user model.User) (uint, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	return s.store.User.CreateUser(s.ctx, &user)
}

func (s *UserService) SignIn(email, password string) (string, error) {
	user, err := s.store.User.GetUser(s.ctx, email)
	if err != nil {
		return "", err
	}

	err = util.CheckPassword(password, user.Password)

	if err != nil {
		return "", err
	}

	return util.GenerateToken(user.ID)

}
