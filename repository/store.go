package repository

import (
	"context"
	"gorm.io/gorm"
)

// Store contains all repositories
type Store struct {
	Db *gorm.DB

	User    UserRepo
	Post    PostRepo
	Comment CommentRepo
}

// New creates new repository
func New(ctx context.Context, db *gorm.DB, userRepo UserRepo, postRepo PostRepo, commentRepo CommentRepo) (*Store, error) {
	//if err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{}); err != nil {
	//	log.Println(err)
	//}

	var store Store

	// Init MySQL repositories
	if db != nil {
		store.Db = db
		store.User = userRepo
		store.Post = postRepo
		store.Comment = commentRepo
	}

	return &store, nil
}
