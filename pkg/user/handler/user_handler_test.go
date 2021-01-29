package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	mock "github.com/Zucke/social_prove/pkg/user/mock"
)

func TestHandler_FirebaseAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()

	uid := "12442df"
	ctx := context.WithValue(context.Background(), auth.UIDKey, uid)
	token := "token12321446"

	u := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	if err := u.EncryptPassword(); err != nil {
		assert.Nil(t, err)
		t.Fail()
	}

	tests := []struct {
		name  string
		user  *user.User
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			user:  &u,
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Invalid no found",
			user:  &user.User{},
			code:  http.StatusNotFound,
			err:   response.ErrorNotFound,
			times: 1,
		},
		{
			name:  "invalid internal error",
			user:  &user.User{},
			code:  http.StatusInternalServerError,
			err:   response.ErrorInternalServerError,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				FirebaseAuth(gomock.Any(), uid).
				Return(test.user, token, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/auth/google", nil)
			r = r.WithContext(ctx)

			mux := chi.NewRouter()
			mux.Post("/auth/google", h.FirebaseAuthHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
func TestHandler_LoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	token := "token12321446"
	l := logger.NewMock()

	userLogin := user.User{
		Email:    "user@example.com",
		Password: "123456",
	}

	jsonUserLogin, err := json.Marshal(userLogin)
	assert.NoError(t, err)

	badPasswordUserLogin := user.User{
		Email:    "user@example.com",
		Password: "654321",
	}
	jsonbadPasswordUserLogin, err := json.Marshal(badPasswordUserLogin)
	assert.NoError(t, err)

	u := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	if err := u.EncryptPassword(); err != nil {
		assert.Nil(t, err)
		t.Fail()
	}

	tests := []struct {
		name  string
		user  *user.User
		rUser *user.User
		body  io.Reader
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			user:  &userLogin,
			rUser: &u,
			body:  strings.NewReader(string(jsonUserLogin)),
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Invalid body",
			user:  &user.User{},
			rUser: &user.User{},
			body:  strings.NewReader(``),
			code:  http.StatusBadRequest,
			err:   nil,
			times: 0,
		},
		{
			name:  "Password does not match",
			user:  &badPasswordUserLogin,
			rUser: &user.User{},
			body:  strings.NewReader(string(jsonbadPasswordUserLogin)),
			code:  http.StatusBadRequest,
			err:   response.ErrorBadEmailOrPassword,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				LoginUser(gomock.Any(), test.user).
				Return(test.rUser, token, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/login", test.body)

			mux := chi.NewRouter()
			mux.Post("/login", h.LoginHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_GetOneHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id := primitive.NewObjectID()

	u := user.User{
		ID:        id,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}

	tests := []struct {
		name  string
		user  user.User
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			user:  u,
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure",
			user:  user.User{},
			code:  http.StatusNotFound,
			err:   response.ErrorNotFound,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetByID(gomock.Any(), id.Hex()).
				Return(test.user, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/users/"+id.Hex(), nil)

			mux := chi.NewRouter()
			mux.Get("/users/{id}", h.GetOneHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
func TestHandler_GetAllHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()

	tests := []struct {
		name   string
		users  []user.User
		code   int
		err    error
		times  int
		pTimes int
	}{
		{
			name:   "Success",
			users:  []user.User{},
			code:   http.StatusOK,
			err:    nil,
			times:  1,
			pTimes: 0,
		},
		{
			name:   "Failure",
			users:  []user.User{},
			code:   http.StatusNotFound,
			err:    response.ErrorNotFound,
			times:  1,
			pTimes: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetAllActive(gomock.Any()).
				Return(test.users, test.err).
				Times(test.times)
			m.
				EXPECT().
				WithPagination(test.users, 0, len(test.users)).
				Return(test.users, len(test.users)).
				Times(test.pTimes)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/users/", nil)

			mux := chi.NewRouter()
			mux.Get("/users/", h.GetAllHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}

}

func TestHandler_CreateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()

	u := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}

	j, err := json.Marshal(u)
	if err != nil {
		assert.NotNil(t, err)
	}

	tests := []struct {
		name  string
		user  user.User
		body  io.Reader
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			user:  u,
			body:  bytes.NewReader(j),
			code:  http.StatusCreated,
			err:   nil,
			times: 1,
		},
		{
			name:  "Invalid body",
			user:  u,
			body:  strings.NewReader(``),
			code:  http.StatusBadRequest,
			err:   nil,
			times: 0,
		},
		{
			name:  "User not inserted",
			user:  u,
			body:  bytes.NewReader(j),
			code:  http.StatusBadRequest,
			err:   response.ErrCouldNotInsert,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Create(gomock.Any(), &test.user).
				Return(test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/users", test.body)

			mux := chi.NewRouter()
			mux.Post("/users", h.CreateHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_FollowTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	u := user.User{
		ID:        id1,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}

	tests := []struct {
		name        string
		user        user.User
		code        int
		followingID string
		followerID  string
		err         error
		times       int
	}{
		{
			name:        "Success",
			user:        u,
			followingID: id1.Hex(),
			followerID:  id2.Hex(),
			code:        http.StatusOK,
			err:         nil,
			times:       1,
		},
		{
			name:        "Invalid Id",
			user:        u,
			followingID: "l,lñk",
			followerID:  id2.Hex(),
			code:        http.StatusBadRequest,
			err:         response.ErrInvalidID,
			times:       1,
		},
		{
			name:        "Same Id",
			user:        u,
			followingID: id2.Hex(),
			followerID:  id2.Hex(),
			code:        http.StatusBadRequest,
			err:         response.ErrCantFollowYou,
			times:       1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				FollowTo(gomock.Any(), test.followingID, test.followerID).
				Return(test.user, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), auth.IDKey, id2)
			r := httptest.NewRequest(http.MethodPost, "/users/"+test.followingID, nil).WithContext(ctx)

			mux := chi.NewRouter()
			mux.Post("/users/{id}", h.FollowToHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
func TestHandler_UnfollowTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	u := user.User{
		ID:        id1,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}

	tests := []struct {
		name        string
		user        user.User
		code        int
		followingID string
		followerID  string
		err         error
		times       int
	}{
		{
			name:        "Success",
			user:        u,
			followingID: id1.Hex(),
			followerID:  id2.Hex(),
			code:        http.StatusOK,
			err:         nil,
			times:       1,
		},
		{
			name:        "Invalid Id",
			user:        u,
			followingID: "1234",
			followerID:  id2.Hex(),
			code:        http.StatusBadRequest,
			err:         response.ErrInvalidID,
			times:       1,
		},
		{
			name:        "Same Id",
			user:        u,
			followingID: id2.Hex(),
			followerID:  id2.Hex(),
			code:        http.StatusBadRequest,
			err:         response.ErrCantFollowYou,
			times:       1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				UnfollowTo(gomock.Any(), test.followingID, test.followerID).
				Return(test.user, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			ctx := context.WithValue(context.Background(), auth.IDKey, id2)
			r := httptest.NewRequest(http.MethodDelete, "/users/"+test.followingID, nil).WithContext(ctx)

			mux := chi.NewRouter()
			mux.Delete("/users/{id}", h.UnfollowToHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_DeleteHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	role := user.Admin

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id := primitive.NewObjectID()

	u := user.User{
		ID:        id,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}

	tests := []struct {
		name  string
		user  user.User
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			user:  u,
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure",
			user:  user.User{},
			code:  http.StatusNotFound,
			err:   response.ErrorNotFound,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Delete(gomock.Any(), role, id.Hex()).
				Return(test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodDelete, "/users/"+id.Hex(), nil)
			ctx := context.WithValue(r.Context(), auth.RoleKey, role)
			r = r.WithContext(ctx)

			mux := chi.NewRouter()
			mux.Delete("/users/{id}", h.DeleteHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
func TestHandler_UpdateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	role := user.Admin
	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id := primitive.NewObjectID()

	u := user.User{
		ID:        id,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}

	bodyUser, err := json.Marshal(u)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		user  user.User
		code  int
		body  io.Reader
		err   error
		times int
	}{
		{
			name:  "Success",
			user:  u,
			body:  strings.NewReader(string(bodyUser)),
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure",
			user:  u,
			body:  strings.NewReader(string(bodyUser)),
			code:  http.StatusNotFound,
			err:   response.ErrorNotFound,
			times: 1,
		},
		{
			name:  "Invalid body",
			user:  u,
			body:  strings.NewReader(``),
			code:  http.StatusBadRequest,
			err:   response.ErrorBadRequest,
			times: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Update(gomock.Any(), id.Hex(), id.Hex(), role, &test.user).
				Return(test.user, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPut, "/users/"+id.Hex(), test.body)
			ctx := context.WithValue(r.Context(), auth.RoleKey, role)
			ctx = context.WithValue(ctx, auth.IDKey, id)
			r = r.WithContext(ctx)

			mux := chi.NewRouter()
			mux.Put("/users/{id}", h.UpdateHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
