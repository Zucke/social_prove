package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Zucke/social_prove/internal/db/mongo"
	"github.com/Zucke/social_prove/internal/server"
	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/logger"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8000"
	}

	var dbURL string
	if dbURL = os.Getenv("DATABASE_URI"); dbURL == "" {
		dbURL = "mongodb://127.0.0.1:27017"
	}

	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	log := logger.New("draid", !*debug)

	ctx := context.Background()
	dbClient, err := mongo.NewClient(ctx, log, dbURL)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = dbClient.Start(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var fa auth.Repository
	// firebaseCredentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	// fa, err := auth.NewFirebaseAuth(context.Background(), firebaseCredentialsPath)
	// if err != nil {
	// 	log.Error(err)
	// 	os.Exit(1)
	// }

	srv, err := server.New(port, *debug, dbClient, log, fa)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Start the server.
	go srv.Start()

	// Wait for an interrupt.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	dbClient.Close(ctx)
	srv.Close(ctx)
}
