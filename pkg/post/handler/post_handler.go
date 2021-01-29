package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/pagination"
	"github.com/Zucke/social_prove/pkg/post"
	"github.com/Zucke/social_prove/pkg/post/service"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
)

// Handler is the router of the post.
type Handler struct {
	service post.Service
	log     logger.Logger
}

// GetAllHandler response all the posts.
func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	var (
		posts []post.Post
		err   error
	)
	page, limit, ok := pagination.GetPagination(r)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		posts, err = h.service.GetAll(ctx)

	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, response.ErrorNotFound.Error())
		return
	}

	total := len(posts)
	if !ok {
		_ = response.JSON(w, http.StatusOK, response.Map{
			"posts": posts,
			"total": total,
		})
		return
	}
	posts, total = h.service.WithPagination(posts, page, limit)
	_ = response.JSON(w, http.StatusOK, render.M{
		"posts": posts,
		"total": total,
	})
}

// GetOneHandler response one post by id.
func (h *Handler) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var (
		p   post.Post
		err error
	)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		p, err = h.service.GetByID(ctx, id)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, response.ErrorNotFound.Error())
		return
	}

	_ = response.JSON(w, http.StatusOK, render.M{"post": p})
}

// CreateHandler create a new post.
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Post

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorParsingUser.Error())
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		err = h.service.Create(ctx, &p)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusCreated, render.M{"post": p})
}

// UpdateHandler update a stored post by id.
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var p, updatedPost post.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorParsingUser.Error())
		return
	}

	id := chi.URLParam(r, "id")

	lID, err := auth.GetID(r)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}
	role, err := auth.GetRole(r)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		updatedPost, err = h.service.Update(ctx, id, lID, role, &p)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, r, render.M{"post": updatedPost})
}

// DeleteHandler Remove a user by ID.
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lID, err := auth.GetID(r)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	role, err := auth.GetRole(r)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		err = h.service.Delete(ctx, id, lID, role)
	}
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, r, render.M{})
}

// AddLikeHandler add like to post.
func (h *Handler) AddLikeHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Post
	id := chi.URLParam(r, "id")

	lID, err := auth.GetID(r)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	select {
	case <-ctx.Done():
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		p, err = h.service.AddLike(ctx, lID, id)
	}
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, r, render.M{"post": p})
}

// DeleteLikeHandler add like to post.
func (h *Handler) DeleteLikeHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Post
	id := chi.URLParam(r, "id")

	lID, err := auth.GetID(r)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	select {
	case <-ctx.Done():
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		p, err = h.service.DeleteLike(ctx, lID, id)
	}
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, r, render.M{"post": p})
}

//Routes configure and return routes for users
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client, user.Admin, user.Super)).
		Get("/", h.GetAllHandler)
	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client)).
		Post("/", h.CreateHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client, user.Admin, user.Super)).
		Get("/{id}", h.GetOneHandler)
	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client, user.Admin, user.Super)).
		Put("/{id}", h.UpdateHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client)).
		Post("/{id}/like", h.AddLikeHandler)
	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client)).
		Delete("/{id}/like", h.DeleteLikeHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client, user.Admin, user.Super)).
		Delete("/{id}", h.DeleteHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client, user.Admin, user.Super)).
		Put("/{id}", h.UpdateHandler)
	return r

}

// NewPostHandler create and configure a new Handler.
func New(coll *mongo.Collection, log logger.Logger) *Handler {
	return &Handler{
		log:     log,
		service: service.New(coll, log),
	}
}
