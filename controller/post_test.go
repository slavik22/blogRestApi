package controller

import (
	"context"
	"github.com/labstack/gommon/log"
	util2 "github.com/slavik22/blogRestApi/lib/util"
	"github.com/slavik22/blogRestApi/model"
	"github.com/slavik22/blogRestApi/repository"
	mock_repository "github.com/slavik22/blogRestApi/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreatePostAPI(t *testing.T) {
	post := randomPost(t)

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

			serviceManager, c, rec := newTestServer(t, tc.body, store)

			// Init controllers
			postController := NewUPostController(context.Background(), serviceManager)

			// Assertions
			err = postController.CreatePost(c)

			log.Info(err)

			tc.checkResponse(rec)
		})
	}
}

func randomPost(t *testing.T) (post model.Post) {
	post = model.Post{
		ID:    1,
		Title: util2.RandomString(10),
		Body:  util2.RandomString(100),
	}
	return
}
