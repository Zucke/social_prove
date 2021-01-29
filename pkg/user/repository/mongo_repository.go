package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
)

// Repository storage to the user model.
type Repository struct {
	coll *mongo.Collection
	log  logger.Logger
}

// Create create a new user.
func (r *Repository) Create(ctx context.Context, u *user.User) error {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

// GetByEmail returns a user by email address.
func (r *Repository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	u := user.User{}
	result := r.coll.FindOne(ctx, bson.M{"email": email})

	if result.Err() != nil {
		r.log.Error(result.Err().Error())
		return user.User{}, response.ErrorNotFound

	}
	err := result.Decode(&u)
	if err != nil {
		r.log.Error(err)
		return user.User{}, response.ErrorInternalServerError
	}

	return u, nil
}

// GetByUID returns a user by UID.
func (r *Repository) GetByUID(ctx context.Context, uid string) (user.User, error) {
	u := user.User{}
	result := r.coll.FindOne(ctx, bson.M{"uid": uid})
	if result.Err() != nil {
		r.log.Error(result.Err().Error())
		return user.User{}, response.ErrorNotFound

	}
	err := result.Decode(&u)
	if err != nil {
		r.log.Error(err)
		return user.User{}, response.ErrorInternalServerError
	}
	return u, nil
}

// GetByID returns a user by ID.
func (r *Repository) GetByID(ctx context.Context, objectID primitive.ObjectID) (user.User, error) {
	u := user.User{}
	result := r.coll.FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		r.log.Error(result.Err().Error())
		return user.User{}, response.ErrorNotFound

	}
	err := result.Decode(&u)
	if err != nil {
		r.log.Error(err)
		return user.User{}, response.ErrorInternalServerError
	}
	return u, nil
}

// GetAll returns all stored users.
func (r *Repository) GetAll(ctx context.Context) ([]user.User, error) {
	opt := options.Find().SetProjection(bson.M{"password": 0})
	users := make([]user.User, 0)

	cursor, err := r.coll.Find(ctx, bson.M{"role": user.Client}, opt)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return users, response.ErrorNotFound
	}

	if err != nil {
		r.log.Error(err)
		return nil, response.ErrorInternalServerError
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		u := user.User{}
		if err := cursor.Decode(&u); err != nil {
			r.log.Error(err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// GetAllActive returns all active stored users.
func (r *Repository) GetAllActive(ctx context.Context) ([]user.User, error) {
	opt := options.Find().SetProjection(bson.M{"password": 0})
	users := make([]user.User, 0)

	cursor, err := r.coll.Find(ctx, bson.M{"role": user.Client, "active": true}, opt)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return users, response.ErrorNotFound
	}

	if err != nil {
		r.log.Error(err)
		return nil, response.ErrorInternalServerError
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		u := user.User{}
		if err := cursor.Decode(&u); err != nil {
			r.log.Error(err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// GetByRole returns stored users by role.
func (r *Repository) GetByRole(ctx context.Context, role user.Role) ([]user.User, error) {
	opt := options.Find().SetProjection(bson.M{"password": 0})

	users := make([]user.User, 0)
	cursor, err := r.coll.Find(ctx, bson.M{"role": role}, opt)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return users, response.ErrorNotFound
	}
	if err != nil {
		r.log.Error(err)
		return nil, response.ErrorInternalServerError
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		u := user.User{}
		if err := cursor.Decode(&u); err != nil {
			r.log.Error(err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

// Update user by ID.
func (r *Repository) Update(ctx context.Context, role user.Role, id primitive.ObjectID, u *user.User) error {
	var filter bson.M

	switch role {
	case user.Client:
		filter = bson.M{
			"_id":  id,
			"role": user.Client,
		}

	case user.Admin:
		filter = bson.M{
			"_id":  id,
			"role": user.Client,
		}

	case user.Super:
		filter = bson.M{
			"_id": id,
		}
	}

	update := bson.M{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"country":    u.Country,
		"state":      u.State,
		"city":       u.City,
		"bio":        u.Bio,
		"picture":    u.Picture,
		"updated_at": time.Now(),
	}

	sr := r.coll.FindOneAndUpdate(ctx, filter, bson.M{"$set": update})
	if err := sr.Err(); err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

// FollowTo add a user id to the following array if not exist.
func (r *Repository) FollowTo(ctx context.Context, followingID primitive.ObjectID, followerID primitive.ObjectID) error {
	filter := bson.M{
		"_id": followerID, "following": bson.M{"$nin": followingID},
	}

	update := bson.M{
		"$push": bson.M{"following": followingID},
	}
	result := r.coll.FindOneAndUpdate(ctx, filter, update)
	if err := result.Err(); err != nil {
		r.log.Error(err)
		return response.ErrorInternalServerError
	}

	return nil
}

// UnfollowTo delete a user id to the following array if exist.
func (r *Repository) UnfollowTo(ctx context.Context, followingID primitive.ObjectID, followerID primitive.ObjectID) error {
	filter := bson.M{
		"_id": followerID, "following": bson.M{"$in": followingID},
	}

	update := bson.M{
		"$pull": bson.M{"following": followingID},
	}
	result := r.coll.FindOneAndUpdate(ctx, filter, update)
	if err := result.Err(); err != nil {
		r.log.Error(err)
		return response.ErrorInternalServerError
	}

	return nil
}

// Delete remove a user by ID.
func (r *Repository) Delete(ctx context.Context, role user.Role, id primitive.ObjectID) error {
	var filter bson.M
	switch role {
	case user.Admin:
		filter = bson.M{"_id": id, "role": user.Client}

	case user.Super:
		filter = bson.M{"_id": id}

	}

	_, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		r.log.Error(err)
		return response.ErrorInternalServerError
	}

	return nil
}

// Mongo create a new Repository.
func Mongo(coll *mongo.Collection, log logger.Logger) user.Repository {
	return &Repository{
		coll: coll,
		log:  log,
	}
}
