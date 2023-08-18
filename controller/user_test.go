package controller

import (
	"context"
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

//var (
//	mockDB = map[string]*model.User{
//		"jon@labstack.com": {Name: "User", Email: "user@gmail.com"},
//	}
//	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com"}`
//)

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mock_repository.MockUserRepo)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"name":     user.Name,
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(store *mock_repository.MockUserRepo) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user.ID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				str := recorder.Body.String()

				id, err := strconv.Atoi(strings.TrimSpace(str))

				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, uint(id), user.ID)
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

			tc.buildStubs(userRepo)

			store, err := repository.New(context.Background(), &gorm.DB{}, userRepo, postRepo, commentRepo)

			if err != nil {
				t.Error(err)
			}

			serviceManager, c, rec := newTestServer(t, tc.body, store)

			// Init controllers
			userController := NewUserController(context.Background(), serviceManager)

			// Assertions
			_ = userController.SignUp(c)

			tc.checkResponse(rec)
		})
	}
}

func randomUser(t *testing.T) (user model.User, password string) {
	password = util2.RandomString(6)
	hashedPassword, err := util2.HashPassword(password)
	require.NoError(t, err)

	user = model.User{
		ID:       1000,
		Name:     util2.RandomOwner(),
		Password: hashedPassword,
		Email:    util2.RandomEmail(),
	}
	return
}
