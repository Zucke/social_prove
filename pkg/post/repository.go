package post

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Repository the post repository
type Repository interface {
	GetAll(ctx context.Context) ([]Post, error)
	GetAllForUser(ctx context.Context, userID primitive.ObjectID) ([]Post, error)
	Create(ctx context.Context, p *Post) error
	GetByID(ctx context.Context, id primitive.ObjectID) (Post, error)
	Update(ctx context.Context, id primitive.ObjectID, p *Post) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddLike(ctx context.Context, fanID, postID primitive.ObjectID) error
	DeleteLike(ctx context.Context, fanID, postID primitive.ObjectID) error
}
