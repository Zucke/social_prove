package post

import (
	"context"

	"github.com/Zucke/social_prove/pkg/user"
)

//Service the post service
type Service interface {
	Create(ctx context.Context, p *Post) error
	GetByID(ctx context.Context, id string) (Post, error)
	GetAll(ctx context.Context) ([]Post, error)
	GetAllForUser(ctx context.Context, userID string) ([]Post, error)
	Update(ctx context.Context, toUpdateid string, currendUserID string, role user.Role, p *Post) (Post, error)
	Delete(ctx context.Context, toDeleteID string, currendUserID string, role user.Role) error
	AddLike(ctx context.Context, fanID, postID string) (Post, error)
	DeleteLike(ctx context.Context, fanID, postID string) (Post, error)
	WithPagination(p []Post, page int, limit int) ([]Post, int)
}
