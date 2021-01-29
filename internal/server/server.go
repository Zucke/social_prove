package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/Zucke/social_prove/internal/db/mongo"
	v1 "github.com/Zucke/social_prove/internal/server/v1"
	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
)

// Server is a base server configuration.
type Server struct {
	server *http.Server
	log    logger.Logger
	port   string
	debug  bool
}

func (serv *Server) getRoutes(client *mongo.Client, fa auth.Repository) (http.Handler, error) {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Google-Token",
			"X-Google-client",
			"c-Control",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	v1Routes, err := v1.New(serv.log, client, fa)
	if err != nil {
		return nil, err
	}

	r.Mount("/api/v1", v1Routes)
	r.Handle(
		"/docs/*",
		http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))),
	)

	return r, nil
}

// New initialize a new server with configuration.
func New(port string, debug bool, client *mongo.Client, log logger.Logger, fa auth.Repository) (*Server, error) {
	serv := &Server{
		port:  port,
		debug: debug,
		log:   log,
	}

	r, err := serv.getRoutes(client, fa)
	if err != nil {
		return nil, err
	}

	serv.server = &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return serv, nil
}

// Close server resources.
func (serv *Server) Close(ctx context.Context) error {
	return nil
}

// Start the server.
func (serv *Server) Start() {
	serv.log.Infof("Server running on http://localhost:%s", serv.port)
	serv.log.Error(serv.server.ListenAndServe())
	os.Exit(1)
}
