package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Zucke/social_prove/pkg/logger"
	"github.com/Zucke/social_prove/pkg/post"
	mock "github.com/Zucke/social_prove/pkg/post/mock"
	"github.com/Zucke/social_prove/pkg/response"
	"github.com/Zucke/social_prove/pkg/user"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	m := mock.NewMockRepository(ctrl)

	p := post.Post{
		Description: "contend bla bla bla, bla",
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		post  post.Post
		err   error
		times int
	}{
		{
			name:  "succes",
			post:  p,
			err:   nil,
			times: 1,
		},
		{
			name:  "failure could't insert",
			post:  p,
			err:   response.ErrCouldNotInsert,
			times: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Create(gomock.Any(), &test.post).
				Return(test.err).
				Times(test.times)

			s := PostService{
				repository: m,
				log:        l,
			}

			err := s.Create(ctx, &test.post)
			assert.Equal(t, err, test.err)

		})
	}
}
func TestUserService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	pID := primitive.NewObjectID()
	m := mock.NewMockRepository(ctrl)

	p := post.Post{
		Description: "contend bla bla bla, bla",
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		post  post.Post
		id    string
		oID   primitive.ObjectID
		err   error
		times int
	}{
		{
			name:  "succes",
			post:  p,
			id:    pID.Hex(),
			oID:   pID,
			err:   nil,
			times: 1,
		},
		{
			name:  "failure bad id",
			post:  post.Post{},
			id:    "1245",
			oID:   primitive.NilObjectID,
			err:   response.ErrInvalidID,
			times: 0,
		},
		{
			name:  "failure internal error",
			post:  post.Post{},
			id:    pID.Hex(),
			oID:   pID,
			err:   response.ErrorInternalServerError,
			times: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetByID(gomock.Any(), test.oID).
				Return(test.post, test.err).
				Times(test.times)

			s := PostService{
				repository: m,
				log:        l,
			}

			resultPost, err := s.GetByID(ctx, test.id)
			assert.Equal(t, err, test.err)
			assert.Equal(t, resultPost, test.post)

		})
	}
}

func TestUserService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	p := []post.Post{
		post.Post{
			ID:          primitive.NewObjectID(),
			Description: "contend bla bla bla, bla",
		},
		post.Post{
			ID:          primitive.NewObjectID(),
			Description: "contend bla bla bla, bla",
		},
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		posts []post.Post
		err   error
		times int
	}{
		{
			name:  "succes",
			posts: p,
			err:   nil,
			times: 1,
		},
		{
			name:  "failure internal error",
			posts: nil,
			err:   response.ErrorInternalServerError,
			times: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetAll(gomock.Any()).
				Return(test.posts, test.err).
				Times(test.times)

			s := PostService{
				repository: m,
				log:        l,
			}

			resultPosts, err := s.GetAll(ctx)
			assert.Equal(t, err, test.err)
			assert.Equal(t, resultPosts, test.posts)

		})
	}
}

func TestUserService_GetAllForUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	oID := primitive.NewObjectID()
	p := []post.Post{
		post.Post{
			ID:          primitive.NewObjectID(),
			Description: "contend bla bla bla, bla",
		},
		post.Post{
			ID:          primitive.NewObjectID(),
			Description: "contend bla bla bla, bla",
		},
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name  string
		posts []post.Post
		err   error
		id    string
		oID   primitive.ObjectID
		times int
	}{
		{
			name:  "succes",
			posts: p,
			err:   nil,
			id:    oID.Hex(),
			oID:   oID,
			times: 1,
		},
		{
			name:  "failure bad id",
			posts: nil,
			err:   response.ErrInvalidID,
			oID:   primitive.NilObjectID,
			id:    "1234",
			times: 0,
		},
		{
			name:  "failure internal error",
			posts: nil,
			err:   response.ErrorInternalServerError,
			id:    oID.Hex(),
			oID:   oID,
			times: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				GetAllForUser(gomock.Any(), test.oID).
				Return(test.posts, test.err).
				Times(test.times)

			s := PostService{
				repository: m,
				log:        l,
			}

			resultPosts, err := s.GetAllForUser(ctx, test.id)
			assert.Equal(t, err, test.err)
			assert.Equal(t, resultPosts, test.posts)

		})
	}
}

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	p := post.Post{
		ID:          primitive.NewObjectID(),
		UserID:      id2,
		Description: "contend bla bla bla, bla",
	}
	otherUserPost := post.Post{
		ID:          primitive.NewObjectID(),
		UserID:      primitive.NewObjectID(),
		Description: "contend bla bla bla, bla",
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name     string
		post     post.Post
		rpost    post.Post
		err      error
		id       string
		oID      primitive.ObjectID
		times    int
		timesID1 int
		timesID2 int
		role     user.Role
	}{
		{
			name:     "succes",
			rpost:    p,
			post:     p,
			err:      nil,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    1,
			timesID2: 1,
			role:     user.Client,
		},
		{
			name:     "failure bad id",
			rpost:    post.Post{},
			post:     post.Post{},
			err:      response.ErrInvalidID,
			id:       "1234",
			oID:      id1,
			timesID1: 0,
			times:    0,
			timesID2: 0,
			role:     user.Client,
		},
		{
			name:     "failure unauthorized",
			post:     post.Post{},
			rpost:    otherUserPost,
			err:      response.ErrorUnauthorized,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    0,
			timesID2: 0,
			role:     user.Client,
		},
		{
			name:     "succes deferend userID with admin",
			post:     p,
			rpost:    otherUserPost,
			err:      nil,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 0,
			times:    1,
			timesID2: 1,
			role:     user.Admin,
		},
		{
			name:     "failure internal error",
			post:     post.Post{},
			rpost:    p,
			err:      response.ErrorInternalServerError,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    1,
			timesID2: 0,
			role:     user.Client,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Update(gomock.Any(), test.oID, &test.post).
				Return(test.err).
				Times(test.times)
			m.
				EXPECT().
				GetByID(gomock.Any(), test.oID).
				Return(test.rpost, nil).
				Times(test.timesID1)
			m.
				EXPECT().
				GetByID(gomock.Any(), test.oID).
				Return(test.post, nil).
				Times(test.timesID2)

			s := PostService{
				repository: m,
				log:        l,
			}

			resultPosts, err := s.Update(ctx, test.id, id2.Hex(), test.role, &test.post)
			assert.Equal(t, err, test.err)
			assert.Equal(t, resultPosts, test.post)

		})
	}
}
func TestUserService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	p := post.Post{
		ID:          primitive.NewObjectID(),
		UserID:      id2,
		Description: "contend bla bla bla, bla",
	}
	otherUserPost := post.Post{
		ID:          primitive.NewObjectID(),
		UserID:      primitive.NewObjectID(),
		Description: "contend bla bla bla, bla",
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name     string
		post     post.Post
		err      error
		id       string
		oID      primitive.ObjectID
		times    int
		timesID1 int
		role     user.Role
	}{
		{
			name:     "succes",
			post:     p,
			err:      nil,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    1,
			role:     user.Client,
		},
		{
			name:     "failure bad id",
			post:     p,
			err:      response.ErrInvalidID,
			id:       "1234",
			oID:      id1,
			timesID1: 0,
			times:    0,
			role:     user.Client,
		},
		{
			name:     "failure unauthorized",
			post:     otherUserPost,
			err:      response.ErrorUnauthorized,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    0,
			role:     user.Client,
		},
		{
			name:     "succes deferend userID with admin",
			post:     otherUserPost,
			err:      nil,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 0,
			times:    1,
			role:     user.Admin,
		},
		{
			name:     "failure internal error",
			post:     p,
			err:      response.ErrorInternalServerError,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    1,
			role:     user.Client,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				Delete(gomock.Any(), test.oID).
				Return(test.err).
				Times(test.times)
			m.
				EXPECT().
				GetByID(gomock.Any(), test.oID).
				Return(test.post, nil).
				Times(test.timesID1)

			s := PostService{
				repository: m,
				log:        l,
			}

			err := s.Delete(ctx, test.id, id2.Hex(), test.role)
			assert.Equal(t, err, test.err)

		})
	}
}
func TestUserService_AddLike(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	p := post.Post{
		ID:          primitive.NewObjectID(),
		UserID:      id2,
		Description: "contend bla bla bla, bla",
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name     string
		post     post.Post
		err      error
		id       string
		oID      primitive.ObjectID
		times    int
		timesID1 int
		role     user.Role
	}{
		{
			name:     "succes",
			post:     p,
			err:      nil,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    1,
			role:     user.Client,
		},
		{
			name:     "failure bad id",
			post:     post.Post{},
			err:      response.ErrInvalidID,
			id:       "1234",
			oID:      id1,
			timesID1: 0,
			times:    0,
			role:     user.Client,
		},
		{
			name:     "failure internal error",
			post:     post.Post{},
			err:      response.ErrorInternalServerError,
			id:       id1.Hex(),
			oID:      id1,
			times:    1,
			timesID1: 0,
			role:     user.Client,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				AddLike(gomock.Any(), test.oID, id2).
				Return(test.err).
				Times(test.times)
			m.
				EXPECT().
				GetByID(gomock.Any(), id2).
				Return(test.post, nil).
				Times(test.timesID1)

			s := PostService{
				repository: m,
				log:        l,
			}

			resultPosts, err := s.AddLike(ctx, test.id, id2.Hex())
			assert.Equal(t, err, test.err)
			assert.Equal(t, resultPosts, test.post)

		})
	}
}
func TestUserService_DeleteLike(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()
	m := mock.NewMockRepository(ctrl)
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	p := post.Post{
		ID:          primitive.NewObjectID(),
		UserID:      id2,
		Description: "contend bla bla bla, bla",
	}

	ctx := context.Background()
	l := logger.NewMock()

	tests := []struct {
		name     string
		post     post.Post
		err      error
		id       string
		oID      primitive.ObjectID
		times    int
		timesID1 int
		role     user.Role
	}{
		{
			name:     "succes",
			post:     p,
			err:      nil,
			id:       id1.Hex(),
			oID:      id1,
			timesID1: 1,
			times:    1,
			role:     user.Client,
		},
		{
			name:     "failure bad id",
			post:     post.Post{},
			err:      response.ErrInvalidID,
			id:       "1234",
			oID:      id1,
			timesID1: 0,
			times:    0,
			role:     user.Client,
		},
		{
			name:     "failure internal error",
			post:     post.Post{},
			err:      response.ErrorInternalServerError,
			id:       id1.Hex(),
			oID:      id1,
			times:    1,
			timesID1: 0,
			role:     user.Client,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m.
				EXPECT().
				DeleteLike(gomock.Any(), test.oID, id2).
				Return(test.err).
				Times(test.times)
			m.
				EXPECT().
				GetByID(gomock.Any(), id2).
				Return(test.post, nil).
				Times(test.timesID1)

			s := PostService{
				repository: m,
				log:        l,
			}

			resultPosts, err := s.DeleteLike(ctx, test.id, id2.Hex())
			assert.Equal(t, err, test.err)
			assert.Equal(t, resultPosts, test.post)

		})
	}
}
