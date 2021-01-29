package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/post"
	"github.com/Zucke/social_prove/pkg/response"
)

// Repository storage to the post model.
type Repository struct {
	coll *mongo.Collection
	log  logger.Logger
}

var pipeLineColl = "users"

// Create create a new post.
func (r *Repository) Create(ctx context.Context, u *post.Post) error {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

// GetByID returns a post by ID.
func (r *Repository) GetByID(ctx context.Context, objectID primitive.ObjectID) (post.Post, error) {
	p := post.Post{}
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"_id": objectID,
		}}},

		{{"$lookup", bson.D{
			{"from", pipeLineColl},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
	}
	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return post.Post{}, response.ErrorNotFound
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&p); err != nil {
			r.log.Error(err)
			return post.Post{}, response.ErrorInternalServerError
		}
	}
	return p, nil
}

// GetAll returns all stored posts.
func (r *Repository) GetAll(ctx context.Context) ([]post.Post, error) {
	posts := make([]post.Post, 0)

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{}}},

		{{"$lookup", bson.D{
			{"from", pipeLineColl},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return posts, response.ErrorNotFound
	}

	if err != nil {
		r.log.Error(err)
		return nil, response.ErrorInternalServerError
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		p := post.Post{}
		if err := cursor.Decode(&p); err != nil {
			r.log.Error(err)
			continue
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// GetAllForUser returns all store post for a user.
func (r *Repository) GetAllForUser(ctx context.Context, userID primitive.ObjectID) ([]post.Post, error) {
	posts := make([]post.Post, 0)

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"user_id": userID,
		}}},

		{{"$lookup", bson.D{
			{"from", pipeLineColl},
			{"localField", "user_id"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return posts, response.ErrorNotFound
	}

	if err != nil {
		r.log.Error(err)
		return nil, response.ErrorInternalServerError
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		p := post.Post{}
		if err := cursor.Decode(&p); err != nil {
			r.log.Error(err)
			continue
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// AddLike add a like to a user by ID.
func (r *Repository) AddLike(ctx context.Context, fanID, postID primitive.ObjectID) error {
	update := bson.M{
		"$addToSet": bson.M{"likes": fanID},
	}
	result := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": postID}, update)
	if result.Err() != nil {
		r.log.Error(result.Err())
		return response.ErrorInternalServerError
	}

	return nil
}

// DeleteLike delete a like to a user by ID.
func (r *Repository) DeleteLike(ctx context.Context, fanID, postID primitive.ObjectID) error {
	update := bson.M{
		"$pull": bson.M{"likes": fanID},
	}
	result := r.coll.FindOneAndUpdate(ctx, bson.M{"_id": postID}, update)
	if result.Err() != nil {
		r.log.Error(result.Err())
		return response.ErrorInternalServerError
	}

	return nil
}

// Update post by ID.
func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, p *post.Post) error {
	filter := bson.M{
		"_id": id,
	}

	update := bson.M{
		"description": p.Description,
		"badge":       p.Badge,
		"pictures":    p.Pictures,
		"updated_at":  time.Now(),
	}

	sr := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": update})
	if err := sr.Err(); err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

// Delete remove a user by ID.
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		r.log.Error(err)
		return response.ErrorInternalServerError
	}

	return nil
}

// Mongo create a new Repository.
func Mongo(coll *mongo.Collection, log logger.Logger) post.Repository {
	return &Repository{
		coll: coll,
		log:  log,
	}
}
