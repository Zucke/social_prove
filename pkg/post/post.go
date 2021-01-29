package post

import (
	"os/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Post is the post model
type Post struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	User        *user.User         `json:"user,omitempty" bson:"user,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Badge       string             `json:"badge,omitempty" bson:"badge,omitempty"`
	Pictures    []string           `json:"pictures,omitempty" bson:"pictures,omitempty"`
	Likes       []string           `json:"likes,omitempty" bson:"likes,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
