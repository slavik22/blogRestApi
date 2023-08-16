package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/slavik22/blogRestApi/model"
	"gorm.io/gorm"
)

// PostMysqlRepo ...
type PostMysqlRepo struct {
	db *gorm.DB
}

// NewPostRepo ...
func NewPostMysqlRepo(db *gorm.DB) *PostMysqlRepo {
	return &PostMysqlRepo{db: db}
}

func (repo *PostMysqlRepo) GetPosts(context.Context) ([]model.Post, error) {
	var posts []model.Post
	err := repo.db.Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("error while fetching posts %v", err)
	}

	return posts, nil
}

func (repo *PostMysqlRepo) GetPost(ctx context.Context, userId uint, postId uint) (*model.Post, error) {
	var post model.Post
	err := repo.db.First(&post, "userId = ? AND postId = ?", userId, postId).Error
	if err != nil {
		return nil, fmt.Errorf("no post found %v", err)
	}

	return &post, nil
}

func (repo *PostMysqlRepo) CreatePost(ctx context.Context, post *model.Post) (uint, error) {
	if post == nil {
		return 0, errors.New("No post provided")
	}
	err := repo.db.Create(post).Error
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (repo *PostMysqlRepo) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	err := repo.db.Where("id = ? AND user_id = ?", post.ID, post.UserId).
		Updates(model.Post{Title: post.Title, Body: post.Body}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no post found %v", err)
		}
		return nil, err
	}

	return post, nil
}

func (repo *PostMysqlRepo) DeletePost(ctx context.Context, userId uint, postId uint) error {
	err := repo.db.Where("id = ? AND user_id = ?", postId, userId).Delete(model.Post{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no post found %v", err)
		}
		return err
	}
	return nil
}
