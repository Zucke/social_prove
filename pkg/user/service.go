package user

import "context"

// Service the user service.
type Service interface {
	Create(ctx context.Context, u *User) error
	LoginUser(ctx context.Context, u *User) (*User, string, error)
	Update(ctx context.Context, toUpdateid string, currendUserID string, role Role, u *User) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUID(ctx context.Context, uid string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetAllActive(ctx context.Context) ([]User, error)
	FollowTo(ctx context.Context, followingID string, followerID string) (User, error)
	UnfollowTo(ctx context.Context, followingID string, followerID string) (User, error)
	Delete(ctx context.Context, role Role, id string) error
	GetByRole(ctx context.Context, role Role) ([]User, error)
	FirebaseAuth(ctx context.Context, uid string) (*User, string, error)
	WithPagination(users []User, page int, limit int) ([]User, int)
}
