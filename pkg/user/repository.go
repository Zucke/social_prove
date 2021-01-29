package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository handle the CRUD operations with Users.
type Repository interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, role Role, id primitive.ObjectID, user *User) error
	GetAll(ctx context.Context) ([]User, error)
	GetAllActive(ctx context.Context) ([]User, error)
	GetByUID(ctx context.Context, uid string) (User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	FollowTo(ctx context.Context, followingID primitive.ObjectID, followerID primitive.ObjectID) error
	UnfollowTo(ctx context.Context, followingID primitive.ObjectID, followerID primitive.ObjectID) error
	Delete(ctx context.Context, role Role, id primitive.ObjectID) error
	GetByRole(ctx context.Context, role Role) ([]User, error)
}
