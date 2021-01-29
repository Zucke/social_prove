package service

import (
	"context"
	"time"

	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/post"
	"github.com/Zucke/social_prove/pkg/post/repository"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const waitTime = 10

// PostService the post service.
type PostService struct {
	repository post.Repository
	log        logger.Logger
}

// Create create a new post.
func (ps *PostService) Create(ctx context.Context, p *post.Post) error {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if err := ps.repository.Create(ctx, p); err != nil {
		ps.log.Error(err)
		return response.ErrCouldNotInsert
	}
	return nil
}

// GetByID returns a post by ID.
func (ps *PostService) GetByID(ctx context.Context, id string) (post.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrInvalidID
	}

	p, err := ps.repository.GetByID(ctx, objectID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, err
	}

	return p, nil
}

// GetAllForUser return all post of a user.
func (ps *PostService) GetAllForUser(ctx context.Context, userID string) ([]post.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		ps.log.Error(err)
		return nil, response.ErrInvalidID
	}

	posts, err := ps.repository.GetAllForUser(ctx, objectUserID)
	if err != nil {
		ps.log.Error(err)
		return nil, err
	}

	return posts, nil
}

// GetAll returns all stored posts.
func (ps *PostService) GetAll(ctx context.Context) ([]post.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	posts, err := ps.repository.GetAll(ctx)
	if err != nil {
		ps.log.Error(err)
		return nil, err
	}

	return posts, nil
}

// Update post by ID.
func (ps *PostService) Update(ctx context.Context, toUpdateID string, currendUserID string, role user.Role, p *post.Post) (post.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, waitTime*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(toUpdateID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrInvalidID
	}

	if role == user.Client {
		vPost, err := ps.GetByID(ctx, toUpdateID)
		if err != nil {
			ps.log.Error(err)
			return post.Post{}, err
		}
		if currendUserID != vPost.UserID.Hex() {
			return post.Post{}, response.ErrorUnauthorized
		}
	}

	err = ps.repository.Update(ctx, objectID, p)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrorInternalServerError
	}
	updatedPost, err := ps.GetByID(ctx, toUpdateID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, err
	}

	return updatedPost, nil

}

// Delete remove a post by ID.
func (ps *PostService) Delete(ctx context.Context, toDeleteID string, currendUserID string, role user.Role) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(toDeleteID)
	if err != nil {
		ps.log.Error(err)
		return response.ErrInvalidID
	}

	if role == user.Client {
		vPost, err := ps.GetByID(ctx, toDeleteID)
		if err != nil {
			ps.log.Error(err)
			return err
		}
		if currendUserID != vPost.UserID.Hex() {
			return response.ErrorUnauthorized
		}
	}

	err = ps.repository.Delete(ctx, objectID)
	if err != nil {
		ps.log.Error(err)
		return err
	}
	return nil
}

// AddLike add a like to user.
func (ps *PostService) AddLike(ctx context.Context, fanID, postID string) (post.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectFanID, err := primitive.ObjectIDFromHex(fanID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrInvalidID
	}
	objectPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrInvalidID
	}

	err = ps.repository.AddLike(ctx, objectFanID, objectPostID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, err
	}

	updatedPost, err := ps.GetByID(ctx, postID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, err
	}

	return updatedPost, nil
}

// DeleteLike delete a like from user.
func (ps *PostService) DeleteLike(ctx context.Context, fanID, postID string) (post.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectFanID, err := primitive.ObjectIDFromHex(fanID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrInvalidID
	}
	objectPostID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, response.ErrInvalidID
	}

	err = ps.repository.DeleteLike(ctx, objectFanID, objectPostID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, err
	}

	updatedPost, err := ps.GetByID(ctx, postID)
	if err != nil {
		ps.log.Error(err)
		return post.Post{}, err
	}

	return updatedPost, nil
}

// WithPagination returns users with a pagination limit.
func (ps *PostService) WithPagination(p []post.Post, page int, limit int) ([]post.Post, int) {
	if limit < 0 {
		limit = 0
	}

	if page < 0 {
		page = 0
	}

	total := len(p)
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

	newUsers := p[start:end]

	return newUsers, total
}

// New create and configure user services.
func New(coll *mongo.Collection, log logger.Logger) post.Service {
	return &PostService{
		repository: repository.Mongo(coll, log),
		log:        log,
	}
}
