package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Zucke/social_prove/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Client struct {
	*mongo.Client
	log logger.Logger
}

// DBName Database name.
const DBName = "draid"

// Collections.
const (
	UserCollection = "users"
	PostCollection = "posts"
	TripCollection = "trips"
)

// Errors.
var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrCouldNotInsert  = errors.New("could not insert")
	ErrCouldNotUpdate  = errors.New("could not update")
	ErrCouldNotDelete  = errors.New("could not delete")
	ErrCouldNotFound   = errors.New("could not found")
	ErrCouldNotParseID = errors.New("could not parse id")
)

// Close disconnect the database.
func (c *Client) Close(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}

// Collection returns a MongoDB's collection from a name.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.Client.Database(DBName).Collection(name)
}

// Start initialize the mongo client and set de index.
func (c *Client) Start(ctx context.Context) error {
	if err := c.indexes(ctx); err != nil {
		c.log.Errorf("cannot create indexes: %v", err)
		return err
	}

	return nil
}

// indexes set index to all models.
func (c *Client) indexes(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	database := c.Client.Database(DBName)
	indexOpts := options.CreateIndexes().SetMaxTime(time.Second * 10)

	// User Indexes.
	userEmailIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true).SetUnique(true),
		Keys:    bsonx.MDoc{"email": bsonx.Int32(1)},
	}

	userUIDIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true),
		Keys:    bsonx.MDoc{"uid": bsonx.Int32(1)},
	}

	userRoleIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true),
		Keys:    bsonx.MDoc{"role": bsonx.Int32(1)},
	}

	geolocationIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true),
		Keys:    bsonx.MDoc{"location": bsonx.String("2dsphere")},
	}

	userIndexes := database.Collection(UserCollection).Indexes()
	_, err := userIndexes.CreateMany(
		ctx,
		[]mongo.IndexModel{
			userEmailIndexModel,
			userRoleIndexModel,
			userUIDIndexModel,
		},
		indexOpts,
	)
	if err != nil {
		return err
	}

	// Trip indexes.
	userIDIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true),
		Keys:    bsonx.MDoc{"user_id": bsonx.Int32(1)},
	}

	tripIndexes := database.Collection(TripCollection).Indexes()
	_, err = tripIndexes.CreateOne(ctx, userIDIndexModel, indexOpts)
	if err != nil {
		return err
	}

	// Post indexes.

	postIndexes := database.Collection(PostCollection).Indexes()
	_, err = postIndexes.CreateMany(ctx, []mongo.IndexModel{userIDIndexModel, geolocationIndexModel}, indexOpts)
	if err != nil {
		return err
	}

	return nil
}

// NewClient returns a new client for mongo.
func NewClient(ctx context.Context, log logger.Logger, source string) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(source))
	if err != nil {
		log.Errorf("cannot create mongodb connection: %v", err)
		return nil, err
	}
	return &Client{Client: client, log: log}, nil
}
