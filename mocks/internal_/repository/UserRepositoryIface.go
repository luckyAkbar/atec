// Code generated by mockery v2.52.2. DO NOT EDIT.

package repository

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "github.com/luckyAkbar/atec/internal/model"

	repository "github.com/luckyAkbar/atec/internal/repository"

	uuid "github.com/google/uuid"
)

// UserRepositoryIface is an autogenerated mock type for the UserRepositoryIface type
type UserRepositoryIface struct {
	mock.Mock
}

type UserRepositoryIface_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepositoryIface) EXPECT() *UserRepositoryIface_Expecter {
	return &UserRepositoryIface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, input, txController
func (_m *UserRepositoryIface) Create(ctx context.Context, input repository.CreateUserInput, txController ...*gorm.DB) (*model.User, error) {
	_va := make([]interface{}, len(txController))
	for _i := range txController {
		_va[_i] = txController[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.CreateUserInput, ...*gorm.DB) (*model.User, error)); ok {
		return rf(ctx, input, txController...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.CreateUserInput, ...*gorm.DB) *model.User); ok {
		r0 = rf(ctx, input, txController...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.CreateUserInput, ...*gorm.DB) error); ok {
		r1 = rf(ctx, input, txController...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryIface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserRepositoryIface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - input repository.CreateUserInput
//   - txController ...*gorm.DB
func (_e *UserRepositoryIface_Expecter) Create(ctx interface{}, input interface{}, txController ...interface{}) *UserRepositoryIface_Create_Call {
	return &UserRepositoryIface_Create_Call{Call: _e.mock.On("Create",
		append([]interface{}{ctx, input}, txController...)...)}
}

func (_c *UserRepositoryIface_Create_Call) Run(run func(ctx context.Context, input repository.CreateUserInput, txController ...*gorm.DB)) *UserRepositoryIface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*gorm.DB, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*gorm.DB)
			}
		}
		run(args[0].(context.Context), args[1].(repository.CreateUserInput), variadicArgs...)
	})
	return _c
}

func (_c *UserRepositoryIface_Create_Call) Return(_a0 *model.User, _a1 error) *UserRepositoryIface_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryIface_Create_Call) RunAndReturn(run func(context.Context, repository.CreateUserInput, ...*gorm.DB) (*model.User, error)) *UserRepositoryIface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepositoryIface) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmail")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryIface_FindByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByEmail'
type UserRepositoryIface_FindByEmail_Call struct {
	*mock.Call
}

// FindByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *UserRepositoryIface_Expecter) FindByEmail(ctx interface{}, email interface{}) *UserRepositoryIface_FindByEmail_Call {
	return &UserRepositoryIface_FindByEmail_Call{Call: _e.mock.On("FindByEmail", ctx, email)}
}

func (_c *UserRepositoryIface_FindByEmail_Call) Run(run func(ctx context.Context, email string)) *UserRepositoryIface_FindByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepositoryIface_FindByEmail_Call) Return(_a0 *model.User, _a1 error) *UserRepositoryIface_FindByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryIface_FindByEmail_Call) RunAndReturn(run func(context.Context, string) (*model.User, error)) *UserRepositoryIface_FindByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// FindByID provides a mock function with given fields: ctx, id
func (_m *UserRepositoryIface) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryIface_FindByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByID'
type UserRepositoryIface_FindByID_Call struct {
	*mock.Call
}

// FindByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *UserRepositoryIface_Expecter) FindByID(ctx interface{}, id interface{}) *UserRepositoryIface_FindByID_Call {
	return &UserRepositoryIface_FindByID_Call{Call: _e.mock.On("FindByID", ctx, id)}
}

func (_c *UserRepositoryIface_FindByID_Call) Run(run func(ctx context.Context, id uuid.UUID)) *UserRepositoryIface_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *UserRepositoryIface_FindByID_Call) Return(_a0 *model.User, _a1 error) *UserRepositoryIface_FindByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryIface_FindByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*model.User, error)) *UserRepositoryIface_FindByID_Call {
	_c.Call.Return(run)
	return _c
}

// IsAdminAccountExists provides a mock function with given fields: ctx
func (_m *UserRepositoryIface) IsAdminAccountExists(ctx context.Context) (bool, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for IsAdminAccountExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (bool, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryIface_IsAdminAccountExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsAdminAccountExists'
type UserRepositoryIface_IsAdminAccountExists_Call struct {
	*mock.Call
}

// IsAdminAccountExists is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserRepositoryIface_Expecter) IsAdminAccountExists(ctx interface{}) *UserRepositoryIface_IsAdminAccountExists_Call {
	return &UserRepositoryIface_IsAdminAccountExists_Call{Call: _e.mock.On("IsAdminAccountExists", ctx)}
}

func (_c *UserRepositoryIface_IsAdminAccountExists_Call) Run(run func(ctx context.Context)) *UserRepositoryIface_IsAdminAccountExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserRepositoryIface_IsAdminAccountExists_Call) Return(_a0 bool, _a1 error) *UserRepositoryIface_IsAdminAccountExists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryIface_IsAdminAccountExists_Call) RunAndReturn(run func(context.Context) (bool, error)) *UserRepositoryIface_IsAdminAccountExists_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, userID, input
func (_m *UserRepositoryIface) Update(ctx context.Context, userID uuid.UUID, input repository.UpdateUserInput) (*model.User, error) {
	ret := _m.Called(ctx, userID, input)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, repository.UpdateUserInput) (*model.User, error)); ok {
		return rf(ctx, userID, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, repository.UpdateUserInput) *model.User); ok {
		r0 = rf(ctx, userID, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, repository.UpdateUserInput) error); ok {
		r1 = rf(ctx, userID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryIface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type UserRepositoryIface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - input repository.UpdateUserInput
func (_e *UserRepositoryIface_Expecter) Update(ctx interface{}, userID interface{}, input interface{}) *UserRepositoryIface_Update_Call {
	return &UserRepositoryIface_Update_Call{Call: _e.mock.On("Update", ctx, userID, input)}
}

func (_c *UserRepositoryIface_Update_Call) Run(run func(ctx context.Context, userID uuid.UUID, input repository.UpdateUserInput)) *UserRepositoryIface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(repository.UpdateUserInput))
	})
	return _c
}

func (_c *UserRepositoryIface_Update_Call) Return(_a0 *model.User, _a1 error) *UserRepositoryIface_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryIface_Update_Call) RunAndReturn(run func(context.Context, uuid.UUID, repository.UpdateUserInput) (*model.User, error)) *UserRepositoryIface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepositoryIface creates a new instance of UserRepositoryIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepositoryIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepositoryIface {
	mock := &UserRepositoryIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
