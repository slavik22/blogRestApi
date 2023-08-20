package repository

import (
	"context"
	"errors"
	"github.com/slavik22/blogRestApi/model"
	"gorm.io/gorm"
)

// UserMysqlRepo ...
type UserMysqlRepo struct {
	db *gorm.DB
}

// NewUserRepo ...
func NewUserMysqlRepo(db *gorm.DB) *UserMysqlRepo {
	return &UserMysqlRepo{db: db}
}

// GetUser retrieves user from Postgres
func (repo *UserMysqlRepo) GetUser(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := repo.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates user in Postgres
func (repo *UserMysqlRepo) CreateUser(ctx context.Context, user *model.User) (uint, error) {
	if user == nil {
		return 0, errors.New("No user provided")
	}
	err := repo.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// UpdateUser updates user in Postgres
func (repo *UserMysqlRepo) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := repo.db.Save(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //not found
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes user in Postgres
func (repo *UserMysqlRepo) DeleteUser(ctx context.Context, id uint) error {
	if id < 0 {
		return errors.New("No user ID provided")
	}
	err := repo.db.Where("id = ?", id).Delete(model.User{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return nil
}
