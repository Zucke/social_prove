package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID          primitive.ObjectID
	Title       string
	Date        time.Time
	Picture     string
	Description string
	Route       [][]int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
