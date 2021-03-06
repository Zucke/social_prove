// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Zucke/social_prove/pkg/post (interfaces: Repository)

// Package mock_post is a generated GoMock package.
package mock_post

import (
	context "context"
	post "github.com/Zucke/social_prove/pkg/post"
	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddLike mocks base method
func (m *MockRepository) AddLike(arg0 context.Context, arg1, arg2 primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLike", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLike indicates an expected call of AddLike
func (mr *MockRepositoryMockRecorder) AddLike(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLike", reflect.TypeOf((*MockRepository)(nil).AddLike), arg0, arg1, arg2)
}

// Create mocks base method
func (m *MockRepository) Create(arg0 context.Context, arg1 *post.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockRepository) Delete(arg0 context.Context, arg1 primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), arg0, arg1)
}

// DeleteLike mocks base method
func (m *MockRepository) DeleteLike(arg0 context.Context, arg1, arg2 primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLike", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLike indicates an expected call of DeleteLike
func (mr *MockRepositoryMockRecorder) DeleteLike(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLike", reflect.TypeOf((*MockRepository)(nil).DeleteLike), arg0, arg1, arg2)
}

// GetAll mocks base method
func (m *MockRepository) GetAll(arg0 context.Context) ([]post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockRepositoryMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepository)(nil).GetAll), arg0)
}

// GetAllForUser mocks base method
func (m *MockRepository) GetAllForUser(arg0 context.Context, arg1 primitive.ObjectID) ([]post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllForUser", arg0, arg1)
	ret0, _ := ret[0].([]post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllForUser indicates an expected call of GetAllForUser
func (mr *MockRepositoryMockRecorder) GetAllForUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllForUser", reflect.TypeOf((*MockRepository)(nil).GetAllForUser), arg0, arg1)
}

// GetByID mocks base method
func (m *MockRepository) GetByID(arg0 context.Context, arg1 primitive.ObjectID) (post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockRepositoryMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), arg0, arg1)
}

// Update mocks base method
func (m *MockRepository) Update(arg0 context.Context, arg1 primitive.ObjectID, arg2 *post.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockRepositoryMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), arg0, arg1, arg2)
}
