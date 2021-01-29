package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"github.com/Zucke/social_prove/pkg/claim"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
)

// FirebaseAuth is a wrapper to firebase client.
type FirebaseAuth struct {
	client *auth.Client
}

// WithToken Allows you to verify the token provided by FireBase from the Header.
func (fa *FirebaseAuth) WithToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		authString, err := claim.TokenFromAuthorization(r)
		if err != nil {
			_ = response.HTTPError(w, http.StatusForbidden, err.Error())
			return
		}

		token, err := fa.client.VerifyIDTokenAndCheckRevoked(ctx, authString)
		if err != nil {
			_ = response.HTTPError(w, http.StatusForbidden, err.Error())
			return
		}

		uid := token.UID
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UIDKey, uid)))
	})
}

// GetFirebaseClient returns the FireBase's Client.
func (fa *FirebaseAuth) GetFirebaseClient() *auth.Client {
	return fa.client
}

// GetFirebaseUser returns an user from firebase auth.
func (fa *FirebaseAuth) GetFirebaseUser(ctx context.Context, uid string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ur, err := fa.client.GetUser(ctx, uid)
	if err != nil {
		return user.User{}, err
	}

	var firstName, lastName string
	names := strings.SplitAfterN(ur.DisplayName, " ", 2)
	firstName = names[0]
	if len(names) == 2 {
		lastName = names[1]
	}
	u := user.User{
		UID:       uid,
		Email:     ur.Email,
		FirstName: firstName,
		LastName:  lastName,
		Picture:   ur.PhotoURL,
		Role:      user.Client,
	}

	return u, nil
}

// NewFirebaseAuth returns a new FirebaseAuth with configuration.
func NewFirebaseAuth(ctx context.Context, credentialsFilePath string) (*FirebaseAuth, error) {
	opts := option.WithCredentialsFile(credentialsFilePath)

	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	fa := FirebaseAuth{
		client: client,
	}

	return &fa, nil
}
