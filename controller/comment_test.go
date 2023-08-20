package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	util2 "github.com/slavik22/blogRestApi/lib/util"
	"github.com/slavik22/blogRestApi/lib/validator"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/repository"
	mock_repository "github.com/slavik22/blogRestApi/repository/mock"
	"github.com/slavik22/blogRestApi/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateCommentAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(t, user.ID)
	comment := randomComment(user.ID, post.ID)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mock_repository.MockCommentRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"body":   comment.Body,
				"title":  comment.Title,
				"postId": comment.PostId,
			},
			buildStubs: func(store *mock_repository.MockCommentRepo) {
				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Any()).
					Times(1).
					Return(comment.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				str := recorder.Body.String()

				id, err := strconv.Atoi(strings.TrimSpace(str))

				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, uint(id), comment.ID)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			postRepo := mock_repository.NewMockPostRepo(ctrl)
			commentRepo := mock_repository.NewMockCommentRepo(ctrl)

			tc.buildStubs(commentRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			e.Validator = validator.NewValidator()

			marshal, err := json.Marshal(tc.body)

			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userId", comment.UserId)

			serviceManager, err := service.NewManager(context.Background(), store)
			if err != nil {
				t.Error(err)
			}
			commentController := NewUCommentController(context.Background(), serviceManager)
			err = commentController.CreateComment(c)

			log.Info(err)

			tc.checkResponse(rec)
		})
	}
}
func TestUpdateCommentAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(t, user.ID)
	comment := randomComment(user.ID, post.ID)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mock_repository.MockCommentRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"body":   comment.Body,
				"title":  comment.Title,
				"postId": comment.PostId,
			},
			buildStubs: func(store *mock_repository.MockCommentRepo) {
				store.EXPECT().
					UpdateComment(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&comment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchComment(t, recorder.Body, comment)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			postRepo := mock_repository.NewMockPostRepo(ctrl)
			commentRepo := mock_repository.NewMockCommentRepo(ctrl)

			tc.buildStubs(commentRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			e.Validator = validator.NewValidator()
			marshal, err := json.Marshal(tc.body)

			if err != nil {
				t.Error(err)
			}

			url := fmt.Sprintf("/v1/api/comments/%d", comment.ID)
			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(string(marshal)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userId", comment.UserId)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(comment.ID)))

			serviceManager, err := service.NewManager(context.Background(), store)
			if err != nil {
				t.Error(err)
			}

			commentController := NewUCommentController(context.Background(), serviceManager)
			err = commentController.UpdateComment(c)

			log.Info(err)

			tc.checkResponse(rec)
		})
	}
}
func TestGetCommentAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(t, user.ID)
	comment := randomComment(user.ID, post.ID)

	testCases := []struct {
		name          string
		commentId     uint
		buildStubs    func(store *mock_repository.MockCommentRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			commentId: comment.ID,
			buildStubs: func(store *mock_repository.MockCommentRepo) {
				store.EXPECT().
					GetComment(gomock.Any(), gomock.Eq(comment.ID)).
					Times(1).
					Return(&comment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchComment(t, recorder.Body, comment)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			postRepo := mock_repository.NewMockPostRepo(ctrl)
			commentRepo := mock_repository.NewMockCommentRepo(ctrl)

			tc.buildStubs(commentRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			e := echo.New()
			e.Validator = validator.NewValidator()

			url := fmt.Sprintf("/v1/api/comments/%d", tc.commentId)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.Set("userId", comment.UserId)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(comment.ID)))

			serviceManager, err := service.NewManager(context.Background(), store)

			if err != nil {
				t.Error(err)
			}

			commentController := NewUCommentController(context.Background(), serviceManager)
			err = commentController.GetCommentById(c)

			if err != nil {
				t.Error(err)
			}

			tc.checkResponse(rec)
		})
	}
}
func TestGetCommentsAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(t, user.ID)

	n := 10
	comments := make([]model.Comment, n)

	for i := 0; i < n; i++ {
		comments[i] = randomComment(user.ID, post.ID)
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mock_repository.MockCommentRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mock_repository.MockCommentRepo) {
				store.EXPECT().
					GetComments(gomock.Any()).
					Times(1).
					Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchComments(t, recorder.Body, comments)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			userRepo := mock_repository.NewMockUserRepo(ctrl)
			postRepo := mock_repository.NewMockPostRepo(ctrl)
			commentRepo := mock_repository.NewMockCommentRepo(ctrl)

			tc.buildStubs(commentRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			e := echo.New()
			e.Validator = validator.NewValidator()

			req := httptest.NewRequest(http.MethodGet, "/v1/api/comments", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.Set("userId", user.ID)

			serviceManager, err := service.NewManager(context.Background(), store)

			if err != nil {
				t.Error(err)
			}

			commentController := NewUCommentController(context.Background(), serviceManager)
			err = commentController.GetAllComments(c)

			if err != nil {
				t.Error(err)
			}

			tc.checkResponse(rec)
		})
	}
}

func randomComment(userId uint, postId uint) (comment model.Comment) {
	comment = model.Comment{
		ID:     1,
		PostId: postId,
		Title:  util2.RandomString(10),
		Body:   util2.RandomString(100),
	}
	return
}

func requireBodyMatchComment(t *testing.T, body *bytes.Buffer, comment model.Comment) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotComment model.Comment
	err = json.Unmarshal(data, &gotComment)
	require.NoError(t, err)

	require.Equal(t, comment.ID, gotComment.ID)
	require.Equal(t, comment.Body, gotComment.Body)
	require.Equal(t, comment.Title, gotComment.Title)
	require.Equal(t, comment.PostId, gotComment.PostId)
	require.Equal(t, comment.UserId, gotComment.UserId)
}

func requireBodyMatchComments(t *testing.T, body *bytes.Buffer, comments []model.Comment) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotComments []model.Comment
	err = json.Unmarshal(data, &gotComments)
	require.NoError(t, err)
	require.Equal(t, comments, gotComments)
}
