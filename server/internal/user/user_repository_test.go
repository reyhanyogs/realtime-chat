package user

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/reyhanyogs/realtime-chat/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name  string
		arg   *domain.User
		idRes int64
		err   bool
		query func(query *sqlmock.ExpectedQuery, arg *domain.User, idRes int64)
	}{
		{
			name: "success",
			arg: &domain.User{
				Username: "test",
				Password: "test",
				Email:    "test@gmail.com",
			},
			idRes: 123,
			err:   false,
			query: func(query *sqlmock.ExpectedQuery, arg *domain.User, idRes int64) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(idRes)
				query.WithArgs(arg.Username, arg.Password, arg.Email).WillReturnRows(rows)
			},
		},
		{
			name:  "error",
			arg:   &domain.User{},
			idRes: 0,
			err:   true,
			query: func(query *sqlmock.ExpectedQuery, arg *domain.User, idRes int64) {
				query.WithArgs(arg.Username, arg.Password, arg.Email).WillReturnError(errors.New("error"))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			query := "INSERT INTO users\\(username, password, email\\) VALUES \\(\\$1, \\$2, \\$3\\) returning id"
			mockQuery := mock.ExpectQuery(query)
			tc.query(mockQuery, tc.arg, tc.idRes)

			repo := NewRepository(db)

			got, err := repo.CreateUser(context.TODO(), tc.arg)

			if tc.err {
				assert.Error(t, err)
				assert.Equal(t, errors.New("error"), err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got.ID)
				assert.Equal(t, tc.idRes, got.ID)
			}

		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	testCases := []struct {
		name  string
		arg   string
		res   *domain.User
		err   bool
		query func(query *sqlmock.ExpectedQuery, arg string, res *domain.User)
	}{
		{
			name: "success",
			arg:  "test@gmail.com",
			res: &domain.User{
				ID:       123,
				Username: "test",
				Email:    "test@gmail.com",
				Password: "test",
			},
			err: false,
			query: func(query *sqlmock.ExpectedQuery, arg string, res *domain.User) {
				row := sqlmock.NewRows([]string{"id", "email", "username", "password"}).AddRow(res.ID, res.Email, res.Username, res.Password)
				query.WithArgs(arg).WillReturnRows(row)
			},
		},
		{
			name: "error",
			arg:  "",
			res:  &domain.User{},
			err:  true,
			query: func(query *sqlmock.ExpectedQuery, arg string, res *domain.User) {
				query.WithArgs(arg).WillReturnError(errors.New("error"))
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			query := "SELECT id, email, username, password FROM users WHERE email = \\$1"
			mockQuery := mock.ExpectQuery(query)
			tc.query(mockQuery, tc.arg, tc.res)

			repo := NewRepository(db)

			got, err := repo.GetUserByEmail(context.TODO(), tc.arg)

			if tc.err {
				assert.Error(t, err)
				assert.Equal(t, err, errors.New("error"))
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.Equal(t, tc.res.Email, got.Email)
				assert.Equal(t, tc.res.ID, got.ID)
				assert.Equal(t, tc.res.Password, got.Password)
				assert.Equal(t, tc.res.Username, got.Username)
			}
		})
	}
}
