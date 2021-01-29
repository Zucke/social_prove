package handler

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/pagination"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	"github.com/Zucke/social_prove/pkg/user/service"
)

// Handler is the router of the users.
type Handler struct {
	service user.Service
	log     logger.Logger
}

// LoginHandler response a JWT to authorization.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u, storedUser *user.User
	var tokenString string

	err := json.NewDecoder(r.Body).Decode(&u)
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
		storedUser, tokenString, err = h.service.LoginUser(ctx, u)
	}

	if err != nil {
		h.log.Error(err)
		if errors.Is(err, response.ErrorBadEmailOrPassword) {
			_ = response.HTTPError(w, http.StatusBadRequest, err.Error())
			return
		} else if errors.Is(err, response.ErrorNotFound) {
			_ = response.HTTPError(w, http.StatusNotFound, err.Error())
			return

		} else {
			_ = response.HTTPError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	_ = response.JSON(w, http.StatusOK, response.Map{
		"token": tokenString,
		"user":  storedUser,
	})
}

// FirebaseAuthHandler response a JWT to authorization.
func (h *Handler) FirebaseAuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iUID := ctx.Value(auth.UIDKey)
	uid, ok := iUID.(string)
	if !ok {
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorUUIDNotFound.Error())
		return
	}

	var (
		err         error
		u           *user.User
		tokenString string
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		u, tokenString, err = h.service.FirebaseAuth(ctx, uid)
	}

	if err != nil {
		h.log.Error(err)
		if errors.Is(err, response.ErrorNotFound) {
			_ = response.HTTPError(w, http.StatusNotFound, err.Error())
			return
		} else {
			_ = response.HTTPError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	_ = response.JSON(w, http.StatusOK, response.Map{
		"token": tokenString,
		"user":  u,
	})
}

// GetAllHandler response all the users.
func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	var (
		users []user.User
		err   error
	)
	page, limit, ok := pagination.GetPagination(r)
	all, err := strconv.ParseBool(r.URL.Query().Get("all"))
	if err != nil {
		all = false
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		if all {
			users, err = h.service.GetAll(ctx)
		} else {
			users, err = h.service.GetAllActive(ctx)
		}
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, response.ErrorNotFound.Error())
		return
	}

	total := len(users)
	if !ok {
		_ = response.JSON(w, http.StatusOK, response.Map{
			"users": users,
			"total": total,
		})
		return
	}
	users, total = h.service.WithPagination(users, page, limit)
	_ = response.JSON(w, http.StatusOK, render.M{
		"users": users,
		"total": total,
	})
}

// GetAdminsHandler response all the users.
func (h *Handler) GetAdminsHandler(w http.ResponseWriter, r *http.Request) {

	var (
		users []user.User
		err   error
	)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		users, err = h.service.GetByRole(ctx, user.Admin)

	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, response.ErrorNotFound.Error())
		return
	}

	_ = response.JSON(w, http.StatusOK, render.M{
		"users": users,
	})
}

// GetOneHandler response one user by id.
func (h *Handler) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var (
		u   user.User
		err error
	)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		u, err = h.service.GetByID(ctx, id)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, response.ErrorNotFound.Error())
		return
	}

	_ = response.JSON(w, http.StatusOK, response.Map{"user": u})
}

// CreateHandler Start a new user.
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorParsingUser.Error())
		return
	}

	u.Role = user.Client

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		err = h.service.Create(ctx, &u)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Location", r.URL.String()+u.ID.Hex())
	_ = response.JSON(w, http.StatusCreated, response.Map{"user": u})
}

// CreateAdminHandler Start a new user.
func (h *Handler) CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorParsingUser.Error())
		return
	}

	u.Role = user.Admin

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		err = h.service.Create(ctx, &u)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Location", r.URL.String()+u.ID.Hex())
	_ = response.JSON(w, http.StatusCreated, response.Map{"user": u})
}

// UpdateHandler update a stored user by id.
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var u, updatedUser user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorParsingUser.Error())
		return
	}

	id := chi.URLParam(r, "id")

	cu, err := auth.GetID(r)
	if err != nil {
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}
	role, err := auth.GetRole(r)
	if err != nil {
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
		updatedUser, err = h.service.Update(ctx, id, cu, role, &u)
	}

	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, r, render.M{"user": updatedUser})
}

//FollowToHandler follow to somebody
func (h *Handler) FollowToHandler(w http.ResponseWriter, r *http.Request) {
	var followerUser user.User
	followingID := chi.URLParam(r, "id")
	followerID, err := auth.GetID(r)

	if err != nil {
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
		followerUser, err = h.service.FollowTo(ctx, followingID, followerID)
	}

	if err != nil {
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorInternalServerError.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.M{"user": followerUser})
}

//UnfollowToHandler unfollow to somebody
func (h *Handler) UnfollowToHandler(w http.ResponseWriter, r *http.Request) {
	var followedUser user.User
	followingID := chi.URLParam(r, "id")
	followerID, err := auth.GetID(r)
	if err != nil {
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
		followedUser, err = h.service.UnfollowTo(ctx, followingID, followerID)
	}

	if err != nil {
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorInternalServerError.Error())
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, render.M{"user": followedUser})
}

// DeleteHandler Remove a user by ID.
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var err error

	role, err := auth.GetRole(r)
	if err != nil {
		_ = response.HTTPError(w, http.StatusBadRequest, response.ErrorBadRequest.Error())
		return
	}

	select {
	case <-ctx.Done():
		_ = response.HTTPError(w, http.StatusBadGateway, response.ErrTimeout.Error())
		return
	default:
		err = h.service.Delete(ctx, role, id)
	}
	if err != nil {
		h.log.Error(err)
		_ = response.HTTPError(w, http.StatusNotFound, err.Error())
		return
	}

	render.JSON(w, r, render.M{})
}

//Routes configure and return routes for users
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client, user.Admin, user.Super)).
		Get("/", h.GetAllHandler)

	r.Post("/", h.CreateHandler)

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
		Post("/{id}/follow", h.FollowToHandler)
	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Client)).
		Delete("/{id}/follow", h.UnfollowToHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Super)).
		Get("/admin", h.CreateHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Super)).
		Get("/admin", h.GetAdminsHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Admin, user.Super)).
		Delete("/{id}", h.DeleteHandler)

	r.
		With(auth.Authenticator).
		With(auth.WithRole(user.Admin, user.Super)).
		Put("/{id}", h.UpdateHandler)
	return r

}

// NewUserHandler create and configure a new Handler.
func New(coll *mongo.Collection, log logger.Logger, firebaseRepo auth.Repository) *Handler {
	return &Handler{
		log:     log,
		service: service.New(coll, log, firebaseRepo),
	}
}
