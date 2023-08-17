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

	User    UserRepo
	Post    PostRepo
	Comment CommentRepo
}

// New creates new repository
func New(ctx context.Context, db *gorm.DB) (*Store, error) {
	if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}); err != nil {
		log.Println(err)
	}

	var store Store

	// Init MySQL repositories
	if db != nil {
		store.Db = db
		store.User = NewUserMysqlRepo(db)
		store.Post = NewPostMysqlRepo(db)
		store.Comment = NewCommentMysqlRepo(db)
	}

	return &store, nil
}
