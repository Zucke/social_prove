package v1

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Zucke/social_prove/internal/db/mongo"
	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
	posthandler "github.com/Zucke/social_prove/pkg/post/handler"
	userhandler "github.com/Zucke/social_prove/pkg/user/handler"
)

// New create and configure routes.
func New(log logger.Logger, dbClient *mongo.Client, fa auth.Repository) (http.Handler, error) {
	r := chi.NewRouter()

	//For User.
	ur := userhandler.New(dbClient.Collection(mongo.UserCollection), log, fa)
	r.Post("/login/", ur.LoginHandler)
	r.Post("/auth/google/", ur.FirebaseAuthHandler)
	r.Mount("/user/", ur.Routes())

	ps := posthandler.New(dbClient.Collection(mongo.PostCollection), log)
	r.Mount("/post/", ps.Routes())

	return r, nil

}
