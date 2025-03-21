// Code generated by mockery v2.52.2. DO NOT EDIT.

package usecase

import (
	context "context"

	model "github.com/luckyAkbar/atec/internal/model"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/luckyAkbar/atec/internal/usecase"

	uuid "github.com/google/uuid"
)

// PackageRepo is an autogenerated mock type for the PackageRepo type
type PackageRepo struct {
	mock.Mock
}

type PackageRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *PackageRepo) EXPECT() *PackageRepo_Expecter {
	return &PackageRepo_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, input, txControllers
func (_m *PackageRepo) Create(ctx context.Context, input usecase.RepoCreatePackageInput, txControllers ...any) (*model.Package, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, txControllers...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *model.Package
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.RepoCreatePackageInput, ...any) (*model.Package, error)); ok {
		return rf(ctx, input, txControllers...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.RepoCreatePackageInput, ...any) *model.Package); ok {
		r0 = rf(ctx, input, txControllers...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.RepoCreatePackageInput, ...any) error); ok {
		r1 = rf(ctx, input, txControllers...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageRepo_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type PackageRepo_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.RepoCreatePackageInput
//   - txControllers ...any
func (_e *PackageRepo_Expecter) Create(ctx interface{}, input interface{}, txControllers ...interface{}) *PackageRepo_Create_Call {
	return &PackageRepo_Create_Call{Call: _e.mock.On("Create",
		append([]interface{}{ctx, input}, txControllers...)...)}
}

func (_c *PackageRepo_Create_Call) Run(run func(ctx context.Context, input usecase.RepoCreatePackageInput, txControllers ...any)) *PackageRepo_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]any, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(any)
			}
		}
		run(args[0].(context.Context), args[1].(usecase.RepoCreatePackageInput), variadicArgs...)
	})
	return _c
}

func (_c *PackageRepo_Create_Call) Return(_a0 *model.Package, _a1 error) *PackageRepo_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PackageRepo_Create_Call) RunAndReturn(run func(context.Context, usecase.RepoCreatePackageInput, ...any) (*model.Package, error)) *PackageRepo_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *PackageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PackageRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type PackageRepo_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *PackageRepo_Expecter) Delete(ctx interface{}, id interface{}) *PackageRepo_Delete_Call {
	return &PackageRepo_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *PackageRepo_Delete_Call) Run(run func(ctx context.Context, id uuid.UUID)) *PackageRepo_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *PackageRepo_Delete_Call) Return(_a0 error) *PackageRepo_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PackageRepo_Delete_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *PackageRepo_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllActivePackages provides a mock function with given fields: ctx
func (_m *PackageRepo) FindAllActivePackages(ctx context.Context) ([]model.Package, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindAllActivePackages")
	}

	var r0 []model.Package
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]model.Package, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []model.Package); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Package)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageRepo_FindAllActivePackages_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllActivePackages'
type PackageRepo_FindAllActivePackages_Call struct {
	*mock.Call
}

// FindAllActivePackages is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PackageRepo_Expecter) FindAllActivePackages(ctx interface{}) *PackageRepo_FindAllActivePackages_Call {
	return &PackageRepo_FindAllActivePackages_Call{Call: _e.mock.On("FindAllActivePackages", ctx)}
}

func (_c *PackageRepo_FindAllActivePackages_Call) Run(run func(ctx context.Context)) *PackageRepo_FindAllActivePackages_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PackageRepo_FindAllActivePackages_Call) Return(_a0 []model.Package, _a1 error) *PackageRepo_FindAllActivePackages_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PackageRepo_FindAllActivePackages_Call) RunAndReturn(run func(context.Context) ([]model.Package, error)) *PackageRepo_FindAllActivePackages_Call {
	_c.Call.Return(run)
	return _c
}

// FindByID provides a mock function with given fields: ctx, id
func (_m *PackageRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindByID")
	}

	var r0 *model.Package
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.Package, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.Package); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageRepo_FindByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByID'
type PackageRepo_FindByID_Call struct {
	*mock.Call
}

// FindByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *PackageRepo_Expecter) FindByID(ctx interface{}, id interface{}) *PackageRepo_FindByID_Call {
	return &PackageRepo_FindByID_Call{Call: _e.mock.On("FindByID", ctx, id)}
}

func (_c *PackageRepo_FindByID_Call) Run(run func(ctx context.Context, id uuid.UUID)) *PackageRepo_FindByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *PackageRepo_FindByID_Call) Return(_a0 *model.Package, _a1 error) *PackageRepo_FindByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PackageRepo_FindByID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*model.Package, error)) *PackageRepo_FindByID_Call {
	_c.Call.Return(run)
	return _c
}

// FindOldestActiveAndLockedPackage provides a mock function with given fields: ctx
func (_m *PackageRepo) FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for FindOldestActiveAndLockedPackage")
	}

	var r0 *model.Package
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*model.Package, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *model.Package); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageRepo_FindOldestActiveAndLockedPackage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOldestActiveAndLockedPackage'
type PackageRepo_FindOldestActiveAndLockedPackage_Call struct {
	*mock.Call
}

// FindOldestActiveAndLockedPackage is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PackageRepo_Expecter) FindOldestActiveAndLockedPackage(ctx interface{}) *PackageRepo_FindOldestActiveAndLockedPackage_Call {
	return &PackageRepo_FindOldestActiveAndLockedPackage_Call{Call: _e.mock.On("FindOldestActiveAndLockedPackage", ctx)}
}

func (_c *PackageRepo_FindOldestActiveAndLockedPackage_Call) Run(run func(ctx context.Context)) *PackageRepo_FindOldestActiveAndLockedPackage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PackageRepo_FindOldestActiveAndLockedPackage_Call) Return(_a0 *model.Package, _a1 error) *PackageRepo_FindOldestActiveAndLockedPackage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PackageRepo_FindOldestActiveAndLockedPackage_Call) RunAndReturn(run func(context.Context) (*model.Package, error)) *PackageRepo_FindOldestActiveAndLockedPackage_Call {
	_c.Call.Return(run)
	return _c
}

// Search provides a mock function with given fields: ctx, input
func (_m *PackageRepo) Search(ctx context.Context, input usecase.RepoSearchPackageInput) ([]model.Package, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 []model.Package
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.RepoSearchPackageInput) ([]model.Package, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.RepoSearchPackageInput) []model.Package); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Package)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.RepoSearchPackageInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageRepo_Search_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Search'
type PackageRepo_Search_Call struct {
	*mock.Call
}

// Search is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.RepoSearchPackageInput
func (_e *PackageRepo_Expecter) Search(ctx interface{}, input interface{}) *PackageRepo_Search_Call {
	return &PackageRepo_Search_Call{Call: _e.mock.On("Search", ctx, input)}
}

func (_c *PackageRepo_Search_Call) Run(run func(ctx context.Context, input usecase.RepoSearchPackageInput)) *PackageRepo_Search_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.RepoSearchPackageInput))
	})
	return _c
}

func (_c *PackageRepo_Search_Call) Return(_a0 []model.Package, _a1 error) *PackageRepo_Search_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PackageRepo_Search_Call) RunAndReturn(run func(context.Context, usecase.RepoSearchPackageInput) ([]model.Package, error)) *PackageRepo_Search_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, id, input, txControllers
func (_m *PackageRepo) Update(ctx context.Context, id uuid.UUID, input usecase.RepoUpdatePackageInput, txControllers ...any) (*model.Package, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, id, input)
	_ca = append(_ca, txControllers...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *model.Package
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, usecase.RepoUpdatePackageInput, ...any) (*model.Package, error)); ok {
		return rf(ctx, id, input, txControllers...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, usecase.RepoUpdatePackageInput, ...any) *model.Package); ok {
		r0 = rf(ctx, id, input, txControllers...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, usecase.RepoUpdatePackageInput, ...any) error); ok {
		r1 = rf(ctx, id, input, txControllers...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackageRepo_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type PackageRepo_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - input usecase.RepoUpdatePackageInput
//   - txControllers ...any
func (_e *PackageRepo_Expecter) Update(ctx interface{}, id interface{}, input interface{}, txControllers ...interface{}) *PackageRepo_Update_Call {
	return &PackageRepo_Update_Call{Call: _e.mock.On("Update",
		append([]interface{}{ctx, id, input}, txControllers...)...)}
}

func (_c *PackageRepo_Update_Call) Run(run func(ctx context.Context, id uuid.UUID, input usecase.RepoUpdatePackageInput, txControllers ...any)) *PackageRepo_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]any, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(any)
			}
		}
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(usecase.RepoUpdatePackageInput), variadicArgs...)
	})
	return _c
}

func (_c *PackageRepo_Update_Call) Return(_a0 *model.Package, _a1 error) *PackageRepo_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PackageRepo_Update_Call) RunAndReturn(run func(context.Context, uuid.UUID, usecase.RepoUpdatePackageInput, ...any) (*model.Package, error)) *PackageRepo_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewPackageRepo creates a new instance of PackageRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPackageRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *PackageRepo {
	mock := &PackageRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
