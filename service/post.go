package service

import (
	"context"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/repository"
)

type PostService struct {
	ctx   context.Context
	store *repository.Store
}

func NewPostService(ctx context.Context, store *repository.Store) *PostService {
	return &PostService{
		ctx:   ctx,
		store: store,
	}
}

func (s *PostService) GetPosts() ([]model.Post, error) {
	return s.store.Post.GetPosts(s.ctx)
}

func (s *PostService) GetPost(postId uint, userId uint) (*model.Post, error) {
	return s.store.Post.GetPost(s.ctx, userId, postId)
}

func (s *PostService) CreatePost(post model.Post, userId uint) (uint, error) {
	post.UserId = userId
	return s.store.Post.CreatePost(s.ctx, &post)
}

func (s *PostService) DeletePost(postId uint, userId uint) error {
	return s.store.Post.DeletePost(s.ctx, userId, postId)
}

func (s *PostService) UpdatePost(post model.Post) (*model.Post, error) {
	return s.store.Post.UpdatePost(s.ctx, &post)
}
