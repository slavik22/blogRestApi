package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/slavik22/blogRestApi/model"
	"gorm.io/gorm"
)

// CommentMysqlRepo ...
type CommentMysqlRepo struct {
	db *gorm.DB
}

// NewCommentRepo ...
func NewCommentMysqlRepo(db *gorm.DB) *CommentMysqlRepo {
	return &CommentMysqlRepo{db: db}
}

func (repo *CommentMysqlRepo) GetComments(context.Context) ([]model.Comment, error) {
	var Comments []model.Comment
	err := repo.db.Find(&Comments).Error
	if err != nil {
		return nil, fmt.Errorf("error while fetching Comments %v", err)
	}

	return Comments, nil
}

func (repo *CommentMysqlRepo) GetComment(ctx context.Context, commentId uint) (*model.Comment, error) {
	var comment model.Comment
	err := repo.db.First(&comment, "id = ?", commentId).Error
	if err != nil {
		return nil, fmt.Errorf("no Comment found %v", err)
	}

	return &comment, nil
}

func (repo *CommentMysqlRepo) CreateComment(ctx context.Context, comment *model.Comment) (uint, error) {
	if comment == nil {
		return 0, errors.New("No Comment provided")
	}
	err := repo.db.Create(comment).Error
	if err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (repo *CommentMysqlRepo) UpdateComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	err := repo.db.Model(&comment).Updates(model.Comment{Title: comment.Title, Body: comment.Body}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no Comment found %v", err)
		}
		return nil, err
	}

	return comment, nil
}

func (repo *CommentMysqlRepo) DeleteComment(ctx context.Context, userId uint, commentId uint) error {
	err := repo.db.Where("id = ? AND user_id = ?", commentId, userId).Delete(model.Comment{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no Comment found %v", err)
		}
		return err
	}
	return nil
}
