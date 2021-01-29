package auth

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/Zucke/social_prove/pkg/claim"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Key to request context.
type Key string

// Context keys.
const (
	RoleKey Key = "role"
	IDKey   Key = "id"
	UIDKey  Key = "uid"
	UserKey Key = "user"
)

//Errors
var (
	ErrUserNotAuthorized      = errors.New("not authorized")
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrIDNotFound             = errors.New("user ID not found")
	ErrIDNoValid              = errors.New("user ID is not valid")
	ErrRoleNotFound           = errors.New("user role not found")
	ErrRoleNoValid            = errors.New("user role is not valid")
)

// GetID returns user ID from the request context.
func GetID(r *http.Request) (id string, err error) {
	iID := r.Context().Value(IDKey)
	if iID == nil {
		err = ErrIDNotFound
		return
	}

	userID, ok := iID.(primitive.ObjectID)
	if !ok {
		err = ErrIDNoValid
		return
	}
	return userID.Hex(), nil
}

// GetRole returns user role from the request context.
func GetRole(r *http.Request) (role user.Role, err error) {
	iRole := r.Context().Value(RoleKey)
	if iRole == nil {
		err = ErrRoleNotFound
		return
	}
	userRole, ok := iRole.(user.Role)
	if !ok {
		err = ErrRoleNoValid
		return
	}
	return userRole, nil
}

// Authenticator is an authentication middleware.
func Authenticator(next http.Handler) http.Handler {
	signingString := os.Getenv("SIGNING_STRING")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := claim.TokenFromAuthorization(r)
		if err != nil {
			_ = response.HTTPError(w, http.StatusUnauthorized, err.Error())
			return
		}

		c, err := claim.GetFromToken(tokenString, signingString)
		if err != nil {
			_ = response.HTTPError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), RoleKey, c.Role)
		ctx = context.WithValue(ctx, IDKey, c.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithRole validate user role from request context.
func WithRole(roles ...user.Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, err := GetRole(r)
			if err != nil {
				_ = response.HTTPError(
					w,
					http.StatusUnauthorized,
					ErrInsufficientPrivileges.Error(),
				)
				return
			}

			isInRoles := func(role user.Role, roles []user.Role) bool {
				for _, r := range roles {
					if r == role {
						return true
					}
				}

				return false
			}

			if !isInRoles(role, roles) {
				_ = response.HTTPError(
					w,
					http.StatusUnauthorized,
					ErrInsufficientPrivileges.Error(),
				)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}

// WithID validate user id from the request context.
func WithID(id ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := GetID(r)
			if err != nil {
				_ = response.HTTPError(
					w,
					http.StatusUnauthorized,
					ErrUserNotAuthorized.Error(),
				)
				return
			}

			urlID := chi.URLParam(r, "id")

			userRole, err := GetRole(r)
			if err != nil {
				_ = response.HTTPError(
					w,
					http.StatusUnauthorized,
					ErrUserNotAuthorized.Error(),
				)
				return
			}

			var checkID string
			if len(id) >= 1 {
				checkID = id[0]
			} else {
				checkID = urlID
			}

			if checkID != userID && userRole < user.Admin {
				_ = response.HTTPError(
					w,
					http.StatusUnauthorized,
					ErrInsufficientPrivileges.Error(),
				)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		})
	}
}
