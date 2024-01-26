package user

import (
	"context"
	"errors"
	"testing"

	"github.com/reyhanyogs/realtime-chat/domain"
	"github.com/reyhanyogs/realtime-chat/domain/mocks"
	"github.com/reyhanyogs/realtime-chat/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUserService(t *testing.T) {
	testCases := []struct {
		name    string
		user    *domain.CreateUserReq
		want    *domain.CreateUserRes
		err     bool
		wantErr error
		stubs   func(repo *mocks.MockRepository)
	}{
		{
			name: "success",
			user: &domain.CreateUserReq{
				Username: "test",
				Email:    "test@gmail.com",
				Password: "test",
			},
			want: &domain.CreateUserRes{
				ID: "123",
			},
			err:     false,
			wantErr: nil,
			stubs: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(&domain.User{ID: 123}, nil)
			},
		},
		{
			name: "error on hash password",
			user: &domain.CreateUserReq{
				Username: "test",
				Email:    "test@gmail.com",
				Password: util.RandomString(73),
			},
			want:    &domain.CreateUserRes{},
			err:     true,
			wantErr: errors.New("error"),
			stubs: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:    "error on creating",
			user:    &domain.CreateUserReq{},
			want:    &domain.CreateUserRes{},
			err:     true,
			wantErr: errors.New("error"),
			stubs: func(repo *mocks.MockRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(&domain.User{}, errors.New("error"))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockrepo := mocks.NewMockRepository(ctrl)
			testService := NewService(mockrepo)
			tc.stubs(mockrepo)

			got, err := testService.CreateUser(context.Background(), tc.user)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestLoginService(t *testing.T) {
	testCases := []struct {
		name    string
		user    *domain.LoginUserReq
		want    *domain.LoginUserRes
		err     bool
		wantErr error
		stubs   func(repo *mocks.MockRepository)
	}{
		{
			name: "success",
			user: &domain.LoginUserReq{
				Email:    "test@gmail.com",
				Password: "test",
			},
			want: &domain.LoginUserRes{
				ID:       "123",
				Username: "test",
			},
			err: false,
			stubs: func(repo *mocks.MockRepository) {
				hashPassword, _ := util.HashPassword("test")
				repo.EXPECT().GetUserByEmail(gomock.Any(), "test@gmail.com").Times(1).Return(&domain.User{ID: 123, Username: "test", Password: hashPassword}, nil)
			},
		},
		{
			name: "error password",
			user: &domain.LoginUserReq{
				Email: "test@gmail.com",
			},
			want:    &domain.LoginUserRes{},
			err:     true,
			wantErr: errors.New("error"),
			stubs: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetUserByEmail(gomock.Any(), "test@gmail.com").Times(1).Return(&domain.User{ID: 123, Username: "test"}, nil)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepository(ctrl)
			testService := NewService(mockRepo)
			tc.stubs(mockRepo)

			got, err := testService.Login(context.Background(), tc.user)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.ID, got.ID)
				assert.Equal(t, tc.want.Username, got.Username)
			}
		})
	}
}
