package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/repository"
	"time"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

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
	user.Password = generatePasswordHash(user.Password)
	return s.store.User.CreateUser(s.ctx, &user)
}

func (s *UserService) GenerateToken(email, password string) (string, error) {
	user, err := s.store.User.GetUser(s.ctx, email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *UserService) ParseToken(accessToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

//// UpdateUser ...
//func (svc *UserWebService) UpdateUser(ctx context.Context, reqUser *model.User) (*model.User, error) {
//	userDB, err := svc.repository.User.GetUser(ctx, reqUser.ID)
//	if err != nil {
//		return nil, errors.Wrap(err, "svc.user.GetUser error")
//	}
//	if userDB == nil {
//		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%v' not found", reqUser.ID))
//	}
//
//	// update user
//	_, err = svc.repository.User.UpdateUser(ctx, reqUser)
//	if err != nil {
//		return nil, errors.Wrap(err, "svc.user.UpdateUser error")
//	}
//
//	// get updated user by ID
//	updatedDBUser, err := svc.repository.User.GetUser(ctx, reqUser.ID)
//	if err != nil {
//		return nil, errors.Wrap(err, "svc.user.GetUser error")
//	}
//
//	return updatedDBUser, nil
//}
//
//// DeleteUser ...
//func (svc *UserWebService) DeleteUser(ctx context.Context, userID uint) error {
//	// Check if user exists
//	userDB, err := svc.repository.User.GetUser(ctx, userID)
//	if err != nil {
//		return errors.Wrap(err, "svc.user.GetUser error")
//	}
//	if userDB == nil {
//		return errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%v' not found", userID))
//	}
//
//	err = svc.repository.User.DeleteUser(ctx, userID)
//	if err != nil {
//		return errors.Wrap(err, "svc.user.DeleteUser error")
//	}
//
//	return nil
//}
