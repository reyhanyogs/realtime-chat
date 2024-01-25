package user

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/reyhanyogs/realtime-chat/domain"
	"github.com/reyhanyogs/realtime-chat/util"
)

const (
	secretKey = "secret"
)

type service struct {
	domain.Repository
	timeout time.Duration
}

type JWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewService(repository domain.Repository) domain.Service {
	return &service{
		Repository: repository,
		timeout:    time.Duration(5) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *domain.CreateUserReq) (*domain.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &domain.CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) Login(c context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	signedString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	return &domain.LoginUserRes{
		AccessToken: signedString,
		ID:          strconv.Itoa(int(u.ID)),
		Username:    u.Username,
	}, nil
}
