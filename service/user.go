package service

import (
	"context"
	"fmt"
	"github.com/slavik22/blogRestApi/lib/types"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/store"

	"github.com/pkg/errors"
)

// UserWebService ...
type UserWebService struct {
	ctx   context.Context
	store *store.Store
}

// NewUserWebService creates a new user web service
func NewUserWebService(ctx context.Context, store *store.Store) *UserWebService {
	return &UserWebService{
		ctx:   ctx,
		store: store,
	}
}

// GetUser ...
func (svc *UserWebService) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	userDB, err := svc.store.User.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.GetUser")
	}
	if userDB == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%s' not found", userID))
	}

	return userDB, nil
}

// CreateUser ...
func (svc UserWebService) CreateUser(ctx context.Context, reqUser *model.User) (*model.User, error) {

	user, err := svc.store.User.CreateUser(ctx, reqUser)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.CreateUser error")
	}

	return user, nil
}

// UpdateUser ...
func (svc *UserWebService) UpdateUser(ctx context.Context, reqUser *model.User) (*model.User, error) {
	userDB, err := svc.store.User.GetUser(ctx, reqUser.ID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.GetUser error")
	}
	if userDB == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%v' not found", reqUser.ID))
	}

	// update user
	_, err = svc.store.User.UpdateUser(ctx, reqUser)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.UpdateUser error")
	}

	// get updated user by ID
	updatedDBUser, err := svc.store.User.GetUser(ctx, reqUser.ID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.GetUser error")
	}

	return updatedDBUser, nil
}

// DeleteUser ...
func (svc *UserWebService) DeleteUser(ctx context.Context, userID uint) error {
	// Check if user exists
	userDB, err := svc.store.User.GetUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "svc.user.GetUser error")
	}
	if userDB == nil {
		return errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%v' not found", userID))
	}

	err = svc.store.User.DeleteUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "svc.user.DeleteUser error")
	}

	return nil
}
