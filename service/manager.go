package service

import (
	"context"
	"errors"
	"github.com/slavik22/blogRestApi/store"
)

// Manager is just a collection of all services we have in the project
type Manager struct {
	User UserService
}

// NewManager creates new service manager
func NewManager(ctx context.Context, store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}
	return &Manager{
		User: NewUserWebService(ctx, store),
	}, nil
}
