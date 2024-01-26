package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/reyhanyogs/realtime-chat/domain"
	"github.com/reyhanyogs/realtime-chat/domain/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUserHandler(t *testing.T) {

	testCases := []struct {
		name          string
		user          *domain.CreateUserReq
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		stubs         func(service *mocks.MockService)
	}{
		{
			name: "success",
			user: &domain.CreateUserReq{
				Username: "test",
				Email:    "test@gmail.com",
				Password: "test",
			},
			stubs: func(service *mocks.MockService) {
				service.EXPECT().CreateUser(gomock.Any(), &domain.CreateUserReq{
					Username: "test",
					Email:    "test@gmail.com",
					Password: "test",
				}).Times(1).Return(&domain.CreateUserRes{
					ID:       "123",
					Username: "test",
					Email:    "test@gmail.com",
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCreateUser(t, recorder.Body, &domain.CreateUserRes{
					ID:       "123",
					Username: "test",
					Email:    "test@gmail.com",
				})
			},
		},
		{
			name: "error create user",
			user: &domain.CreateUserReq{
				Username: "test",
				Email:    "test@gmail.com",
				Password: "test",
			},
			stubs: func(service *mocks.MockService) {
				service.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(&domain.CreateUserRes{}, errors.New("error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatchCreateUser(t, recorder.Body, &domain.CreateUserRes{})
			},
		},
		{
			name: "bad request",
			user: &domain.CreateUserReq{},
			stubs: func(service *mocks.MockService) {
				service.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchCreateUser(t, recorder.Body, &domain.CreateUserRes{})
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockService(ctrl)
			tc.stubs(service)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			NewHandler(router, service)

			recorder := httptest.NewRecorder()
			url := "/signup"
			body, err := json.Marshal(tc.user)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	testCases := []struct {
		name          string
		user          *domain.LoginUserReq
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
		stubs         func(service *mocks.MockService)
	}{
		{
			name: "success",
			user: &domain.LoginUserReq{
				Email:    "test@gmail.com",
				Password: "test",
			},
			stubs: func(service *mocks.MockService) {
				service.EXPECT().Login(gomock.Any(), &domain.LoginUserReq{
					Email:    "test@gmail.com",
					Password: "test",
				}).Times(1).Return(&domain.LoginUserRes{
					ID:       "123",
					Username: "test",
				}, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchLoginUser(t, recorder.Body, &domain.LoginUserRes{
					ID:       "123",
					Username: "test",
				})
			},
		},
		{
			name: "bad request",
			user: &domain.LoginUserReq{},
			stubs: func(service *mocks.MockService) {
				service.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchLoginUser(t, recorder.Body, &domain.LoginUserRes{})
			},
		},
		{
			name: "error login user",
			user: &domain.LoginUserReq{
				Email:    "test@gmail.com",
				Password: "test",
			},
			stubs: func(service *mocks.MockService) {
				service.EXPECT().Login(gomock.Any(), &domain.LoginUserReq{
					Email:    "test@gmail.com",
					Password: "test",
				}).Times(1).Return(&domain.LoginUserRes{}, errors.New("error"))
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatchLoginUser(t, recorder.Body, &domain.LoginUserRes{})
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockService(ctrl)
			tc.stubs(service)

			gin.SetMode(gin.ReleaseMode)
			router := gin.New()
			NewHandler(router, service)

			recorder := httptest.NewRecorder()
			url := "/login"
			body, err := json.Marshal(tc.user)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(t, *recorder)
		})
	}
}

func requireBodyMatchCreateUser(t *testing.T, body *bytes.Buffer, user *domain.CreateUserRes) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got *domain.CreateUserRes
	err = json.Unmarshal(data, &got)

	require.NoError(t, err)
	require.Equal(t, user.Username, got.Username)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, user.Email, got.Email)
}

func requireBodyMatchLoginUser(t *testing.T, body *bytes.Buffer, user *domain.LoginUserRes) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got *domain.CreateUserRes
	err = json.Unmarshal(data, &got)

	require.NoError(t, err)
	require.Equal(t, user.Username, got.Username)
	require.Equal(t, user.ID, got.ID)
}
