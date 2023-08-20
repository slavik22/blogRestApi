package controller

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/service"
	"net/http"
	"strconv"
)

type PostController struct {
	ctx      context.Context
	services *service.Manager
}

func NewUPostController(ctx context.Context, services *service.Manager) *PostController {
	return &PostController{
		ctx:      ctx,
		services: services,
	}
}

// GetAllPosts godoc
//
//	@Summary		Get All Posts
//	@Security		ApiKeyAuth
//	@Tags			Posts
//	@Description	get all Posts
//	@ID				get-all-Posts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]model.Post
//	@Router			/api/v1/posts [get]
func (h *PostController) GetAllPosts(c echo.Context) error {
	posts, err := h.services.PostService.GetPosts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, posts)
}

// GetPostById godoc
//
//	@Summary		Get Post By ID
//	@Security		ApiKeyAuth
//	@Tags			Posts
//	@Description	get model.Post by id
//	@ID				get-Post-by-id
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Post
//	@Router			/api/v1/posts/:id [get]
func (h *PostController) GetPostById(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	postId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "post id is incorrect"))
	}

	post, err := h.services.PostService.GetPost(uint(postId), userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, post)
}

// CreatePost godoc
//
//	@Summary		Create Post
//	@Security		ApiKeyAuth
//	@Tags			Post
//	@Description	create Post
//	@Description	create model.Post
//	@ID				create-Post
//	@Accept			json
//	@Produce		json
//	@Success		200	{uint}	id
//	@Router			/api/v1/posts [post]
func (h *PostController) CreatePost(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	var post model.Post

	if err := c.Bind(&post); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	id, err := h.services.PostService.CreatePost(post, userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, id)
}

// UpdatePost godoc
//
//	@Summary		Update Post
//	@Security		ApiKeyAuth
//	@Tags			Post
//	@Description	update model.Post
//	@ID				update-Post
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.Post
//	@Router			/api/v1/posts [put]
func (h *PostController) UpdatePost(c echo.Context) error {
	userId, err := getUserId(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	postId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "post id is incorrect"))
	}

	var post model.Post

	if err := c.Bind(&post); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	post.ID = uint(postId)
	post.UserId = userId

	updatedPost, err := h.services.PostService.UpdatePost(post)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, *updatedPost)
}

// DeletePost godoc
//
//	@Summary		Delete Post
//	@Security		ApiKeyAuth
//	@Tags			Post
//	@Description	delete model.Post
//	@ID				delete-Post
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/posts [delete]
func (h *PostController) DeletePost(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	postId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "post id is incorrect"))
	}

	err = h.services.PostService.DeletePost(uint(postId), userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "post deleted")
}
