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

	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/post"
	mock "github.com/Zucke/social_prove/pkg/post/mock"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_GetOneHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id := primitive.NewObjectID()

	p := post.Post{
		ID:          id,
		Description: "conted, bla bla bla",
	}

	tests := []struct {
		name  string
		post  post.Post
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			post:  p,
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure not found",
			post:  p,
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
				Return(test.post, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/post/"+id.Hex(), nil)

			mux := chi.NewRouter()
			mux.Get("/post/{id}", h.GetOneHandler)
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
		name  string
		posts []post.Post
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			posts: []post.Post{},
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure not found",
			posts: []post.Post{},
			code:  http.StatusNotFound,
			err:   response.ErrorNotFound,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetAll(gomock.Any()).
				Return(test.posts, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/post/", nil)

			mux := chi.NewRouter()
			mux.Get("/post/", h.GetAllHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
func TestHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()

	p := post.Post{
		Description: "conted, bla bla bla",
	}

	jPost, err := json.Marshal(p)
	if err != nil {
		assert.NotNil(t, err)
	}

	tests := []struct {
		name  string
		post  post.Post
		body  io.Reader
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			post:  p,
			body:  bytes.NewReader(jPost),
			code:  http.StatusCreated,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure bad request ",
			post:  p,
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
				Create(gomock.Any(), &test.post).
				Return(test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/post/", test.body)

			mux := chi.NewRouter()
			mux.Post("/post/", h.CreateHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	role := user.Client

	p := post.Post{
		Description: "conted, bla bla bla",
	}

	jPost, err := json.Marshal(p)
	if err != nil {
		assert.NotNil(t, err)
	}

	tests := []struct {
		name  string
		post  post.Post
		body  io.Reader
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			post:  p,
			body:  bytes.NewReader(jPost),
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure bad request ",
			post:  p,
			body:  strings.NewReader(``),
			code:  http.StatusBadRequest,
			err:   response.ErrorBadRequest,
			times: 0,
		},
		{
			name:  "Failure internal error ",
			post:  p,
			body:  bytes.NewReader(jPost),
			code:  http.StatusNotFound,
			err:   response.ErrorInternalServerError,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Update(gomock.Any(), id1.Hex(), id2.Hex(), role, &test.post).
				Return(test.post, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPut, "/post/"+id1.Hex(), test.body)
			ctx := context.WithValue(r.Context(), auth.RoleKey, role)
			r = r.WithContext(context.WithValue(ctx, auth.IDKey, id2))

			mux := chi.NewRouter()
			mux.Put("/post/{id}", h.UpdateHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	role := user.Client

	tests := []struct {
		name  string
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure internal error",
			code:  http.StatusNotFound,
			err:   response.ErrorBadRequest,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Delete(gomock.Any(), id1.Hex(), id2.Hex(), role).
				Return(test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodDelete, "/post/"+id1.Hex(), nil)
			ctx := context.WithValue(r.Context(), auth.RoleKey, role)
			r = r.WithContext(context.WithValue(ctx, auth.IDKey, id2))

			mux := chi.NewRouter()
			mux.Delete("/post/{id}", h.DeleteHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_AddLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	p := post.Post{
		Description: "conted, bla bla bla",
	}

	tests := []struct {
		name  string
		post  post.Post
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			post:  p,
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure internal error ",
			post:  p,
			code:  http.StatusNotFound,
			err:   response.ErrorInternalServerError,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				AddLike(gomock.Any(), id2.Hex(), id1.Hex()).
				Return(test.post, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/post/"+id1.Hex()+"/like", nil)
			r = r.WithContext(context.WithValue(r.Context(), auth.IDKey, id2))

			mux := chi.NewRouter()
			mux.Post("/post/{id}/like", h.AddLikeHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}

func TestHandler_DeleteLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewMockService(ctrl)
	l := logger.NewMock()
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	p := post.Post{
		Description: "conted, bla bla bla",
	}

	tests := []struct {
		name  string
		post  post.Post
		code  int
		err   error
		times int
	}{
		{
			name:  "Success",
			post:  p,
			code:  http.StatusOK,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure internal error ",
			post:  p,
			code:  http.StatusNotFound,
			err:   response.ErrorInternalServerError,
			times: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				DeleteLike(gomock.Any(), id2.Hex(), id1.Hex()).
				Return(test.post, test.err).
				Times(test.times)

			h := Handler{
				service: m,
				log:     l,
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodDelete, "/post/"+id1.Hex()+"/like", nil)
			r = r.WithContext(context.WithValue(r.Context(), auth.IDKey, id2))

			mux := chi.NewRouter()
			mux.Delete("/post/{id}/like", h.DeleteLikeHandler)
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.code, w.Code)
		})
	}
}
