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

// CommentController ...
type CommentController struct {
	ctx      context.Context
	services *service.Manager
}

// NewUCommentController creates a new Comment controller.
func NewUCommentController(ctx context.Context, services *service.Manager) *CommentController {
	return &CommentController{
		ctx:      ctx,
		services: services,
	}
}

// GetAllComments @Summary Get All Comments
// @Security ApiKeyAuth
// @Tags Comments
// @Description get all Comments
// @ID get-all-Comments
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.Comment
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router v1/api/Comments [get]
func (h *CommentController) GetAllComments(c echo.Context) error {
	Comments, err := h.services.CommentService.GetComments()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, Comments)
}

// GetCommentById @Summary Get Comment By ID
// @Security ApiKeyAuth
// @Tags Comments
// @Description get model.Comment by id
// @ID get-Comment-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Comment
// @Router v1/api/Comments/:id [get]
func (h *CommentController) GetCommentById(c echo.Context) error {
	CommentId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Comment id is incorrect"))
	}

	Comment, err := h.services.CommentService.GetComment(uint(CommentId))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, Comment)
}

// CreateComment @Summary Create Comment
// @Security ApiKeyAuth
// @Tags Comment
// @Description create Comment
// @Description create model.Comment
// @ID create-Comment
// @Accept  json
// @Produce  json
// @Success 200 {object} createCommentResponse
// @Router v1/api/Comments [Comment]
func (h *CommentController) CreateComment(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	var comment model.Comment

	if err := c.Bind(&comment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if comment.PostId == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "post id is incorrect"))
	}

	id, err := h.services.CommentService.CreateComment(comment, userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, id)
}

// UpdateComment @Summary Update Comment
// @Security ApiKeyAuth
// @Tags Comment
// @Description update model.Comment
// @ID update-Comment
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Comment
// @Router v1/api/Comments [update]
func (h *CommentController) UpdateComment(c echo.Context) error {
	userId, err := getUserId(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	CommentId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Comment id is incorrect"))
	}

	var Comment model.Comment

	if err := c.Bind(&Comment); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	Comment.ID = uint(CommentId)
	Comment.UserId = userId

	updatedComment, err := h.services.CommentService.UpdateComment(Comment)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, *updatedComment)
}

//type CommentUpdateRequest struct {
//	Title string
//	Body  string
//}

// DeleteComment @Summary Delete Comment
// @Security ApiKeyAuth
// @Tags Comment
// @Description delete model.Comment
// @ID delete-Comment
// @Accept  json
// @Produce  json
// @Success 200
// @Router v1/api/Comments [delete]
func (h *CommentController) DeleteComment(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.Wrap(err, "user is not authorized"))
	}

	CommentId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "Comment id is incorrect"))
	}

	err = h.services.CommentService.DeleteComment(uint(CommentId), userId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Comment deleted")
}
