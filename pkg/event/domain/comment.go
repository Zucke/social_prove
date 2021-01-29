package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Comment is the commend model
type Comment struct {
	ID        primitive.ObjectID
	PostID    primitive.ObjectID
	UserID    primitive.ObjectID
	Body      string
	Likes     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
