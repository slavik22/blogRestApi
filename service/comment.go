package service

import (
	"context"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/repository"
)

type CommentService struct {
	ctx   context.Context
	store *repository.Store
}

func NewCommentService(ctx context.Context, store *repository.Store) *CommentService {
	return &CommentService{
		ctx:   ctx,
		store: store,
	}
}

func (s *CommentService) GetComments() ([]model.Comment, error) {
	return s.store.Comment.GetComments(s.ctx)
}

func (s *CommentService) GetComment(commentId uint) (*model.Comment, error) {
	return s.store.Comment.GetComment(s.ctx, commentId)
}

func (s *CommentService) CreateComment(comment model.Comment, userId uint) (uint, error) {
	comment.UserId = userId
	return s.store.Comment.CreateComment(s.ctx, &comment)
}

func (s *CommentService) DeleteComment(commentId uint, userId uint) error {
	return s.store.Comment.DeleteComment(s.ctx, userId, commentId)
}

func (s *CommentService) UpdateComment(comment model.Comment) (*model.Comment, error) {
	return s.store.Comment.UpdateComment(s.ctx, &comment)
}
