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

func TestCreatePostAPI(t *testing.T) {
	post := randomPost(t, uint(util2.RandomInt(0, 100)))

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mock_repository.MockPostRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"body":  post.Body,
				"title": post.Title,
			},
			buildStubs: func(store *mock_repository.MockPostRepo) {
				store.EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(post.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				str := recorder.Body.String()

				id, err := strconv.Atoi(strings.TrimSpace(str))

				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, uint(id), post.ID)
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

			tc.buildStubs(postRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			e.Validator = validator.NewValidator()

			json, err := json.Marshal(tc.body)

			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/api/posts", strings.NewReader(string(json)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userId", post.UserId)

			serviceManager, err := service.NewManager(context.Background(), store)
			if err != nil {
				t.Error(err)
			}

			// Init controllers
			postController := NewUPostController(context.Background(), serviceManager)

			// Assertions
			err = postController.CreatePost(c)

			log.Info(err)

			tc.checkResponse(rec)
		})
	}
}
func TestUpdatePostAPI(t *testing.T) {
	post := randomPost(t, uint(util2.RandomInt(0, 100)))

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mock_repository.MockPostRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"body":  post.Body,
				"title": post.Title,
			},
			buildStubs: func(store *mock_repository.MockPostRepo) {
				store.EXPECT().
					UpdatePost(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&post, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, post)
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

			tc.buildStubs(postRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			e.Validator = validator.NewValidator()
			json, err := json.Marshal(tc.body)

			if err != nil {
				t.Error(err)
			}

			url := fmt.Sprintf("/v1/api/posts/%d", post.ID)
			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(string(json)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userId", post.UserId)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(post.ID)))

			serviceManager, err := service.NewManager(context.Background(), store)
			if err != nil {
				t.Error(err)
			}

			postController := NewUPostController(context.Background(), serviceManager)
			err = postController.UpdatePost(c)

			log.Info(err)

			tc.checkResponse(rec)
		})
	}
}
func TestGetPostAPI(t *testing.T) {
	user, _ := randomUser(t)
	post := randomPost(t, user.ID)

	testCases := []struct {
		name          string
		postId        uint
		buildStubs    func(store *mock_repository.MockPostRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			postId: post.ID,
			buildStubs: func(store *mock_repository.MockPostRepo) {
				store.EXPECT().
					GetPost(gomock.Any(), gomock.Eq(post.UserId), gomock.Eq(post.ID)).
					Times(1).
					Return(&post, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchPost(t, recorder.Body, post)
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

			tc.buildStubs(postRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			e := echo.New()
			e.Validator = validator.NewValidator()

			url := fmt.Sprintf("/v1/api/posts/%d", tc.postId)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.Set("userId", post.UserId)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(post.ID)))

			serviceManager, err := service.NewManager(context.Background(), store)

			if err != nil {
				t.Error(err)
			}

			postController := NewUPostController(context.Background(), serviceManager)
			err = postController.GetPostById(c)

			if err != nil {
				t.Error(err)
			}

			tc.checkResponse(rec)
		})
	}
}
func TestGetPostsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 10
	posts := make([]model.Post, n)

	for i := 0; i < n; i++ {
		posts[i] = randomPost(t, user.ID)
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mock_repository.MockPostRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(store *mock_repository.MockPostRepo) {
				store.EXPECT().
					GetPosts(gomock.Any()).
					Times(1).
					Return(posts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAPosts(t, recorder.Body, posts)
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

			tc.buildStubs(postRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			e := echo.New()
			e.Validator = validator.NewValidator()

			req := httptest.NewRequest(http.MethodGet, "/v1/api/posts", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.Set("userId", user.ID)

			serviceManager, err := service.NewManager(context.Background(), store)

			if err != nil {
				t.Error(err)
			}

			postController := NewUPostController(context.Background(), serviceManager)
			err = postController.GetAllPosts(c)

			if err != nil {
				t.Error(err)
			}

			tc.checkResponse(rec)
		})
	}
}

func randomPost(t *testing.T, userId uint) (post model.Post) {
	post = model.Post{
		ID:     1,
		Title:  util2.RandomString(10),
		Body:   util2.RandomString(100),
		UserId: userId,
	}
	return
}

func requireBodyMatchPost(t *testing.T, body *bytes.Buffer, post model.Post) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotPost model.Post
	err = json.Unmarshal(data, &gotPost)
	require.NoError(t, err)
	require.Equal(t, post, gotPost)
}

func requireBodyMatchAPosts(t *testing.T, body *bytes.Buffer, posts []model.Post) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotPosts []model.Post
	err = json.Unmarshal(data, &gotPosts)
	require.NoError(t, err)
	require.Equal(t, posts, gotPosts)
}
