package auth

import (
	"context"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/Zucke/social_prove/pkg/user"
)

//Repository firebase repository
type Repository interface {
	GetFirebaseUser(ctx context.Context, uid string) (user.User, error)
	WithToken(next http.Handler) http.Handler
	GetFirebaseClient() *auth.Client
}
