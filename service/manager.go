package service

import (
	"context"
	"errors"
	"github.com/slavik22/blogRestApi/repository"
)

// Manager is just a collection of all services we have in the project
type Manager struct {
	UserService UserServ
	PostService PostServ
}

// NewManager creates new service manager
func NewManager(ctx context.Context, store *repository.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No repository provided")
	}
	return &Manager{
		UserService: NewUserService(ctx, store),
		PostService: NewPostService(ctx, store),
	}, nil
}
