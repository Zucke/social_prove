package user

import (
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Role to the user on the system.
type Role uint

// Roles to de user.
const (
	Client Role = iota
	Admin
	Super
)

// User is the user model.
type User struct {
	ID             primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Email          string               `json:"email,omitempty" bson:"email,omitempty"`
	Password       string               `json:"password,omitempty" bson:"-"`
	HashPassword   []byte               `json:"-" bson:"password,omitempty"`
	UID            string               `json:"uid,omitempty" bson:"uid,omitempty"`
	FirstName      string               `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName       string               `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Country        string               `json:"country,omitempty" bson:"country,omitempty"`
	State          string               `json:"state,omitempty" bson:"state,omitempty"`
	City           string               `json:"city,omitempty" bson:"city,omitempty"`
	Bio            string               `json:"bio,omitempty" bson:"bio,omitempty"`
	Picture        string               `json:"picture,omitempty" bson:"picture,omitempty"`
	Following      []primitive.ObjectID `json:"following,omitempty" bson:"following,omitempty"`
	Role           Role                 `json:"role,omitempty" bson:"role,omitempty"`
	Active         bool                 `json:"active" bson:"active"`
	NotificationID string               `json:"notification_id,omitempty" bson:"notification_id,omitempty"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// ComparePassword compare the HashPassword with a raw password and return true if they are the same
func (u User) ComparePassword(password string) bool {
	saltedPassword := getSatlForPassword(password)
	err := bcrypt.CompareHashAndPassword(u.HashPassword, []byte(saltedPassword))

	return err == nil
}

// EncryptPassword generate a hashed with the password and put the result in HashPassword.
func (u *User) EncryptPassword() (err error) {
	saltedPassword := getSatlForPassword(u.Password)
	u.HashPassword, err = bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return nil
}

// ValidateEmail confirm valid email format.
func (u User) ValidateEmail() bool {
	const pattern = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	re := regexp.MustCompile(pattern)

	return re.MatchString(u.Email)
}

func getSatlForPassword(password string) string {
	left, right := "", ""
	for i, char := range password {
		if i%2 == 0 {
			left += string(char)
		} else {
			right += string(char)
		}
	}

	return left + password + right
}
