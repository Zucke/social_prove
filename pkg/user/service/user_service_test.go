package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	fmock "github.com/Zucke/social_prove/pkg/auth/mock"
	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
	mock "github.com/Zucke/social_prove/pkg/user/mock"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)

	userWithPassword := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	userWithoutPassword := user.User{
		Email:     "user@example.com",
		UID:       "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	userWithIncorrectEmail := user.User{
		Email:     "userexample.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name        string
		user        user.User
		err         error
		active      bool
		hasPassword bool
		times       int
	}{
		{
			name:        "With password success",
			user:        userWithPassword,
			err:         nil,
			active:      true,
			hasPassword: true,
			times:       1,
		},
		{
			name:        "With password failure",
			user:        userWithPassword,
			hasPassword: true,
			active:      true,
			err:         response.ErrCouldNotInsert,
			times:       1,
		},
		{
			name:   "Without password success",
			user:   userWithoutPassword,
			active: true,
			err:    nil,
			times:  1,
		},
		{
			name:   "Without password failure",
			user:   userWithoutPassword,
			active: true,
			err:    response.ErrCouldNotInsert,
			times:  1,
		},
		{
			name:  "With invalid email",
			user:  userWithIncorrectEmail,
			err:   response.ErrInvalidEmail,
			times: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Create(gomock.Any(), gomock.Eq(&test.user)).
				Return(test.err).
				Times(test.times)

			s := UserService{
				repository: m,
				log:        l,
			}

			err := s.Create(ctx, &test.user)
			assert.Equal(t, err, test.err)
			assert.Equal(t, test.active, test.user.Active)
			if test.hasPassword {
				assert.NotNil(t, test.user.HashPassword)
			} else {
				assert.Nil(t, test.user.HashPassword)
			}
		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)

	us := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name string
		user user.User
		err  error
	}{
		{
			name: "Success",
			user: us,
			err:  nil,
		},
		{
			name: "Failure",
			user: user.User{},
			err:  response.ErrorNotFound,
		},
	}
	email := "test@example.com"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetByEmail(gomock.Any(), email).
				Return(test.user, test.err).
				Times(1)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, err := s.GetByEmail(ctx, email)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, test.user)
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)

	validUser := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	validUser.EncryptPassword()
	userWithBadEmail := user.User{
		Email:     "userexample.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	userWithBadEmail.EncryptPassword()
	userBadPassword := user.User{
		Email:     "user@example.com",
		Password:  "02123",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	userBadPassword.EncryptPassword()
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name       string
		user       user.User
		resultUser user.User
		err        error
		times      int
	}{
		{
			name:       "Success",
			user:       validUser,
			resultUser: validUser,
			err:        nil,
			times:      1,
		},
		{
			name:       "Invalid mail",
			user:       userWithBadEmail,
			resultUser: user.User{},
			err:        response.ErrorBadEmailOrPassword,
			times:      0,
		},
		{
			name:       "Bad password",
			user:       userBadPassword,
			resultUser: user.User{},
			err:        response.ErrorBadEmailOrPassword,
			times:      1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			m.
				EXPECT().
				GetByEmail(gomock.Any(), test.user.Email).
				Return(validUser, nil).
				Times(test.times)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, tokenString, err := s.LoginUser(ctx, &test.user)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, &test.resultUser)
			if err == nil {
				assert.NotEmpty(t, tokenString)
			}
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)
	id := primitive.NewObjectID()
	us := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		id    string
		user  user.User
		err   error
		times int
	}{
		{
			name:  "Success",
			id:    id.Hex(),
			user:  us,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure",
			id:    id.Hex(),
			user:  user.User{},
			err:   response.ErrorNotFound,
			times: 1,
		},
		{
			name:  "With invalid id",
			id:    "123",
			user:  user.User{},
			err:   response.ErrInvalidID,
			times: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetByID(gomock.Any(), id).
				Return(test.user, test.err).
				Times(test.times)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, err := s.GetByID(ctx, test.id)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, test.user)
		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)
	su := []user.User{
		{
			Email:     "user@example.com",
			Password:  "123456",
			FirstName: "user",
			LastName:  "test",
			Role:      user.Client,
			Active:    true,
		},
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		users []user.User
		err   error
	}{
		{
			name:  "Success",
			users: su,
			err:   nil,
		},
		{
			name:  "Failure",
			users: nil,
			err:   response.ErrorNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetAll(gomock.Any()).
				Return(test.users, test.err).
				Times(1)

			s := UserService{
				repository: m,
				log:        l,
			}

			users, err := s.GetAll(ctx)
			assert.Equal(t, err, test.err)
			assert.Equal(t, len(users), len(test.users))
		})
	}
}

func TestUserService_GetAllActive(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)
	su := []user.User{
		{
			Email:     "user@example.com",
			Password:  "123456",
			FirstName: "user",
			LastName:  "test",
			Role:      user.Client,
			Active:    true,
		},
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		users []user.User
		err   error
	}{
		{
			name:  "Success",
			users: su,
			err:   nil,
		},
		{
			name:  "Failure",
			users: nil,
			err:   response.ErrorNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetAllActive(gomock.Any()).
				Return(test.users, test.err).
				Times(1)

			s := UserService{
				repository: m,
				log:        l,
			}

			users, err := s.GetAllActive(ctx)
			assert.Equal(t, err, test.err)
			assert.Equal(t, len(users), len(test.users))
		})
	}
}

func TestUserService_GetByUID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	uid := "123456"
	m := mock.NewMockRepository(ctrl)
	us := user.User{
		Email:     "user@example.com",
		UID:       uid,
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name string
		uid  string
		user user.User
		err  error
	}{
		{
			name: "Success",
			uid:  uid,
			user: us,
			err:  nil,
		},
		{
			name: "Failure",
			uid:  uid,
			user: user.User{},
			err:  response.ErrorNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetByUID(gomock.Any(), test.uid).
				Return(test.user, test.err).
				Times(1)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, err := s.GetByUID(ctx, test.uid)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, test.user)
		})
	}
}

func TestUserService_FollowTo(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	id := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	m := mock.NewMockRepository(ctrl)
	userFollowing := user.User{
		ID:        id,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
		Following: []primitive.ObjectID{id2},
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name        string
		user        user.User
		err         error
		id          string
		idtimes     int
		followTimes int
	}{
		{
			name:        "succes following",
			user:        userFollowing,
			err:         nil,
			id:          id.Hex(),
			idtimes:     1,
			followTimes: 1,
		},
		{
			name:        "failure, user is same has following",
			user:        user.User{},
			err:         response.ErrCantFollowYou,
			id:          id.Hex(),
			followTimes: 1,
			idtimes:     0,
		},
		{
			name:        "failure, bad id",
			user:        user.User{},
			err:         response.ErrInvalidID,
			id:          "",
			followTimes: 0,
			idtimes:     0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				FollowTo(gomock.Any(), id2, id).
				Return(test.err).
				Times(test.followTimes)
			m.
				EXPECT().
				GetByID(gomock.Any(), id).
				Return(test.user, nil).
				Times(test.idtimes)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, err := s.FollowTo(ctx, id2.Hex(), test.id)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, test.user)
		})
	}
}
func TestUserService_UnfollowTo(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	id := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	m := mock.NewMockRepository(ctrl)
	userFollowing := user.User{
		ID:        id,
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
		Following: []primitive.ObjectID{id2},
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name        string
		user        user.User
		err         error
		id          string
		idtimes     int
		followTimes int
	}{
		{
			name:        "succes Unfollowing",
			user:        userFollowing,
			err:         nil,
			id:          id.Hex(),
			idtimes:     1,
			followTimes: 1,
		},
		{
			name:        "failure, user is same has following",
			user:        user.User{},
			err:         response.ErrCantFollowYou,
			id:          id.Hex(),
			followTimes: 1,
			idtimes:     0,
		},
		{
			name:        "failure, bad id",
			user:        user.User{},
			err:         response.ErrInvalidID,
			id:          "",
			followTimes: 0,
			idtimes:     0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				FollowTo(gomock.Any(), id2, id).
				Return(test.err).
				Times(test.followTimes)
			m.
				EXPECT().
				GetByID(gomock.Any(), id).
				Return(test.user, nil).
				Times(test.idtimes)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, err := s.FollowTo(ctx, id2.Hex(), test.id)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, test.user)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)
	id := primitive.NewObjectID()
	us := user.User{
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		id    string
		user  user.User
		err   error
		times int
	}{
		{
			name:  "Success",
			id:    id.Hex(),
			user:  us,
			err:   nil,
			times: 1,
		},
		{
			name:  "Failure",
			id:    id.Hex(),
			user:  user.User{},
			err:   response.ErrorNotFound,
			times: 1,
		},
		{
			name:  "With invalid id",
			id:    "123",
			user:  user.User{},
			err:   response.ErrInvalidID,
			times: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Delete(gomock.Any(), test.user.Role, id).
				Return(test.err).
				Times(test.times)

			s := UserService{
				repository: m,
				log:        l,
			}

			err := s.Delete(ctx, test.user.Role, test.id)
			assert.Equal(t, err, test.err)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	cliendUser := user.User{
		ID:        id1,
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	adminUser := user.User{
		ID:        id1,
		Email:     "user@example.com",
		Password:  "123456",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Admin,
		Active:    true,
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name    string
		id      string
		user    user.User
		rUser   user.User
		err     error
		times   int
		IDtimes int
	}{
		{
			name:    "Success cliend",
			id:      id1.Hex(),
			user:    cliendUser,
			rUser:   cliendUser,
			err:     nil,
			times:   1,
			IDtimes: 1,
		},
		{
			name:    "Success admin",
			id:      id1.Hex(),
			user:    adminUser,
			rUser:   adminUser,
			err:     nil,
			times:   1,
			IDtimes: 1,
		},
		{
			name:    "Failure cliend id unauthorized",
			id:      id2.Hex(),
			user:    cliendUser,
			rUser:   user.User{},
			err:     response.ErrorUnauthorized,
			times:   0,
			IDtimes: 0,
		},
		{
			name:    "Cliend with invalid id",
			id:      "123",
			user:    cliendUser,
			rUser:   user.User{},
			err:     response.ErrorUnauthorized,
			times:   0,
			IDtimes: 0,
		},
		{
			name:    "Admin with invalid id",
			id:      "123",
			user:    adminUser,
			rUser:   user.User{},
			err:     response.ErrInvalidID,
			times:   0,
			IDtimes: 0,
		},
		{
			name:    "Failed to update",
			id:      id1.Hex(),
			user:    adminUser,
			rUser:   user.User{},
			err:     response.ErrorInternalServerError,
			times:   1,
			IDtimes: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Update(gomock.Any(), test.user.Role, gomock.Any(), &test.user).
				Return(test.err).
				Times(test.times)
			m.
				EXPECT().
				GetByID(gomock.Any(), gomock.Any()).
				Return(test.user, nil).
				Times(test.IDtimes)

			s := UserService{
				repository: m,
				log:        l,
			}

			u, err := s.Update(ctx, test.id, test.user.ID.Hex(), test.user.Role, &test.user)
			assert.Equal(t, err, test.err)
			assert.Equal(t, u, test.rUser)
		})
	}
}
func TestUserService_FirebaseAuth(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)
	fm := fmock.NewMockRepository(ctrl)

	uid := "2134gh"
	us := user.User{
		Email:     "user@example.com",
		FirstName: "user",
		LastName:  "test",
		Role:      user.Client,
		Active:    true,
	}
	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name   string
		user   user.User
		err    error
		fErr   error
		times  int
		fTimes int
	}{
		{
			name:   "Success",
			user:   us,
			err:    nil,
			times:  1,
			fTimes: 0,
		},
		{
			name:   "Succes found in fb",
			user:   us,
			err:    response.ErrorNotFound,
			fErr:   nil,
			times:  1,
			fTimes: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			m.
				EXPECT().
				GetByUID(gomock.Any(), uid).
				Return(test.user, test.err).
				Times(test.times)

			m.
				EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(test.fTimes)
			fm.
				EXPECT().
				GetFirebaseUser(gomock.Any(), uid).
				Return(test.user, test.fErr).
				Times(test.fTimes)

			s := UserService{
				repository:   m,
				log:          l,
				firebaseRepo: fm,
			}

			u, tokenString, err := s.FirebaseAuth(ctx, uid)
			if err != nil {
				assert.Equal(t, err, test.err)

			} else {
				assert.Equal(t, err, test.fErr)

			}
			assert.NotEmpty(t, tokenString)
			test.user.CreatedAt = u.CreatedAt
			test.user.ID = u.ID
			test.user.UpdatedAt = u.UpdatedAt
			test.user.Active = u.Active
			assert.Equal(t, u, &test.user)

		})
	}
}
