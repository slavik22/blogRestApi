package repository

import (
	"context"
	"github.com/slavik22/blogRestApi/model"
	"gorm.io/gorm"
	"log"
)

// Store contains all repositories
type Store struct {
	Db *gorm.DB

	User UserRepo
}

// New creates new repository
func New(ctx context.Context, db *gorm.DB) (*Store, error) {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Println(err)
	}

	var store Store

	// Init MySQL repositories
	if db != nil {
		store.Db = db
		store.User = NewUserMysqlRepo(db)
	}

	return &store, nil
}
