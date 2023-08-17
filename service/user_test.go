package service

// TestGetUser runs tests for GetUser service
//func TestGetUser(t *testing.T) {
//	input := &model.User{
//		ID:   0,
//		Name: "Joe Smith",
//	}
//
//	tests := []struct {
//		name         string
//		expectations func(userRepo *mocks.UserRepo)
//		input        *model.User
//		err          error
//	}{
//		{
//			name: "valid and found",
//			expectations: func(userRepo *mocks.UserRepo) {
//				userRepo.On("GetUser", context.Background(), input.ID).Return(input.ToDB(), nil)
//			},
//			input: input,
//		},
//		{
//			name: "valid user ID but not found",
//			expectations: func(userRepo *mocks.UserRepo) {
//				userRepo.On("GetUser", context.Background(), input.ID).Return(nil, nil)
//			},
//			input: input,
//			err:   errors.New("User '7a2f922c-073a-11eb-adc1-0242ac120002' not found: resource not found"),
//		},
//		{
//			name: "repository error",
//			expectations: func(userRepo *mocks.UserRepo) {
//				userRepo.On("GetUser", context.Background(), input.ID).Return(nil, errors.New("some error"))
//			},
//			input: input,
//			err:   errors.New("svc.user.GetUser: some error"),
//		},
//	}
//	for _, test := range tests {
//		t.Logf("running: %s", test.name)
//
//		ctx := context.Background()
//
//		userRepo := &mocks.UserRepo{}
//		svc := NewUserWebService(context.Background(), &repository.Store{User: userRepo})
//		test.expectations(userRepo)
//
//		_, err := svc.GetUser(ctx, test.input.ID)
//		if err != nil {
//			if test.err != nil {
//				assert.Equal(t, test.err.Error(), err.Error())
//			} else {
//				t.Errorf("Expected no error, found: %s", err.Error())
//			}
//		}
//		userRepo.AssertExpectations(t)
//	}
//}
