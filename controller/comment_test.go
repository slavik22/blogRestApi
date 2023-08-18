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

func TestCreateCommentAPI(t *testing.T) {
	comment := randomComment(t)

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

			serviceManager, c, rec := newTestServer(t, tc.body, store)

			// Init controllers
			commentController := NewUCommentController(context.Background(), serviceManager)

			// Assertions
			err = commentController.CreateComment(c)

			log.Info(err)

			tc.checkResponse(rec)
		})
	}
}

func randomComment(t *testing.T) (comment model.Comment) {
	comment = model.Comment{
		ID:     1,
		PostId: 1,
		Title:  util2.RandomString(10),
		Body:   util2.RandomString(100),
	}
	return
}
