package service

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Zucke/social_prove/pkg/auth"
	"github.com/Zucke/social_prove/pkg/claim"
	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	"github.com/Zucke/social_prove/pkg/user/repository"
)

const waitTime = 10

// UserService the user service.
type UserService struct {
	repository   user.Repository
	firebaseRepo auth.Repository
	log          logger.Logger
}

// Create create a new user.
func (us *UserService) Create(ctx context.Context, u *user.User) error {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	if !u.ValidateEmail() {
		return response.ErrInvalidEmail
	}

	if u.Password != "" {
		err := u.EncryptPassword()
		if err != nil {
			us.log.Error(err)
			return response.ErrCouldNotInsert
		}
	}

	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}

	u.Active = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if err := us.repository.Create(ctx, u); err != nil {
		us.log.Error(err)
		return response.ErrCouldNotInsert
	}
	u.Password = ""
	return nil
}

//FirebaseAuth service for firebase auth
func (us *UserService) FirebaseAuth(ctx context.Context, uid string) (*user.User, string, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()
	u, err := us.GetByUID(ctx, uid)
	if err != nil {
		us.log.Error(err)
		u, err = us.firebaseRepo.GetFirebaseUser(ctx, uid)
		if err != nil {
			us.log.Error(err)
			return &user.User{}, "", response.ErrorNotFound
		}
		err = us.Create(ctx, &u)

		if err != nil {
			us.log.Error(err)
			return &user.User{}, "", err
		}

	}

	tokenString, err := claim.GenerateToken(os.Getenv("SIGNING_STRING"), u.ID.Hex(), uint(u.Role))
	if err != nil {
		us.log.Error(err)
		return &user.User{}, "", response.ErrorInternalServerError
	}

	return &u, tokenString, nil

}

//LoginUser evaluate a user and return if it a valid login and it token
func (us *UserService) LoginUser(ctx context.Context, u *user.User) (*user.User, string, error) {
	var tokenString string

	if !u.ValidateEmail() {
		return &user.User{}, "", response.ErrorBadEmailOrPassword
	}
	matchUser, err := us.GetByEmail(ctx, u.Email)

	if err != nil {
		us.log.Error(err)
		return &user.User{}, "", err
	}

	tokenString, err = claim.GenerateToken(os.Getenv("SIGNING_STRING"), matchUser.ID.Hex(), uint(matchUser.Role))
	if err != nil {
		us.log.Error(err)
		return &user.User{}, "", response.ErrorInternalServerError
	}

	if !matchUser.ComparePassword(u.Password) {
		return &user.User{}, "", response.ErrorBadEmailOrPassword
	}

	return &matchUser, tokenString, nil

}

// GetByEmail returns a user by email address.
func (us *UserService) GetByEmail(ctx context.Context, email string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	u, err := us.repository.GetByEmail(ctx, email)
	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	return u, nil
}

// GetByID returns a user by ID.
func (us *UserService) GetByID(ctx context.Context, id string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrInvalidID
	}

	u, err := us.repository.GetByID(ctx, objectID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	return u, nil
}

// GetAll returns all stored users.
func (us *UserService) GetAll(ctx context.Context) ([]user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	users, err := us.repository.GetAll(ctx)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	return users, nil
}

// GetByRole return a list of users by role.
func (us *UserService) GetByRole(ctx context.Context, role user.Role) ([]user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	users, err := us.repository.GetByRole(ctx, role)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	return users, nil
}

// GetAllActive returns all active stored users.
func (us *UserService) GetAllActive(ctx context.Context) ([]user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	users, err := us.repository.GetAllActive(ctx)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}

	return users, nil
}

// GetByUID returns a user by UID.
func (us *UserService) GetByUID(ctx context.Context, uid string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	u, err := us.repository.GetByUID(ctx, uid)
	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	return u, nil
}

// Update user by ID.
func (us *UserService) Update(ctx context.Context, toUpdateid string, currendUserID string, role user.Role, u *user.User) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	if role == user.Client {
		if currendUserID != toUpdateid {
			return user.User{}, response.ErrorUnauthorized
		}
	}

	objectID, err := primitive.ObjectIDFromHex(toUpdateid)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrInvalidID
	}

	err = us.repository.Update(ctx, role, objectID, u)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrorInternalServerError
	}
	updatedUser, err := us.GetByID(ctx, toUpdateid)
	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	return updatedUser, nil

}

// FollowTo add user to the following list
func (us *UserService) FollowTo(ctx context.Context, followingID string, followerID string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	var u user.User
	if followingID == followerID {
		return user.User{}, response.ErrCantFollowYou
	}

	followingObjectID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrInvalidID
	}
	followerObjectID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrInvalidID
	}

	err = us.repository.FollowTo(ctx, followingObjectID, followerObjectID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	u, err = us.GetByID(ctx, followerID)

	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	return u, nil
}

// UnfollowTo delete user of the following list
func (us *UserService) UnfollowTo(ctx context.Context, followingID string, followerID string) (user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	var u user.User
	if followingID == followerID {
		return user.User{}, response.ErrCantFollowYou
	}

	followingObjectID, err := primitive.ObjectIDFromHex(followingID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrInvalidID
	}
	followerObjectID, err := primitive.ObjectIDFromHex(followerID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, response.ErrInvalidID
	}

	err = us.repository.UnfollowTo(ctx, followingObjectID, followerObjectID)
	if err != nil {
		us.log.Error(err)
		return user.User{}, err
	}

	u, err = us.GetByID(ctx, followerID)

	if err != nil {
		us.log.Error(err)
		return u, err
	}

	return u, nil
}

// Delete remove a user by ID.
func (us *UserService) Delete(ctx context.Context, role user.Role, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		us.log.Error(err)
		return response.ErrInvalidID
	}
	err = us.repository.Delete(ctx, role, objectID)
	if err != nil {
		us.log.Error(err)
		return err
	}
	return nil
}

// WithPagination returns users with a pagination limit.
func (us *UserService) WithPagination(users []user.User, page int, limit int) ([]user.User, int) {
	if limit < 0 {
		limit = 0
	}

	if page < 0 {
		page = 0
	}

	total := len(users)
	if limit > total {
		limit = total
	}

	start := (page - 1) * limit
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	newUsers := users[start:end]

	return newUsers, total
}

// New create and configure user services.
func New(coll *mongo.Collection, log logger.Logger, firebaseRepo auth.Repository) user.Service {
	return &UserService{
		repository:   repository.Mongo(coll, log),
		log:          log,
		firebaseRepo: firebaseRepo,
	}
}
