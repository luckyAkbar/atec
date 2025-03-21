// Code generated by mockery v2.52.2. DO NOT EDIT.

package usecase

import mock "github.com/stretchr/testify/mock"

// TransactionController is an autogenerated mock type for the TransactionController type
type TransactionController struct {
	mock.Mock
}

type TransactionController_Expecter struct {
	mock *mock.Mock
}

func (_m *TransactionController) EXPECT() *TransactionController_Expecter {
	return &TransactionController_Expecter{mock: &_m.Mock}
}

// Begin provides a mock function with no fields
func (_m *TransactionController) Begin() any {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Begin")
	}

	var r0 any
	if rf, ok := ret.Get(0).(func() any); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(any)
		}
	}

	return r0
}

// TransactionController_Begin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Begin'
type TransactionController_Begin_Call struct {
	*mock.Call
}

// Begin is a helper method to define mock.On call
func (_e *TransactionController_Expecter) Begin() *TransactionController_Begin_Call {
	return &TransactionController_Begin_Call{Call: _e.mock.On("Begin")}
}

func (_c *TransactionController_Begin_Call) Run(run func()) *TransactionController_Begin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TransactionController_Begin_Call) Return(_a0 any) *TransactionController_Begin_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionController_Begin_Call) RunAndReturn(run func() any) *TransactionController_Begin_Call {
	_c.Call.Return(run)
	return _c
}

// Commit provides a mock function with no fields
func (_m *TransactionController) Commit() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransactionController_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type TransactionController_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
func (_e *TransactionController_Expecter) Commit() *TransactionController_Commit_Call {
	return &TransactionController_Commit_Call{Call: _e.mock.On("Commit")}
}

func (_c *TransactionController_Commit_Call) Run(run func()) *TransactionController_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TransactionController_Commit_Call) Return(_a0 error) *TransactionController_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionController_Commit_Call) RunAndReturn(run func() error) *TransactionController_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with no fields
func (_m *TransactionController) Rollback() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransactionController_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type TransactionController_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
func (_e *TransactionController_Expecter) Rollback() *TransactionController_Rollback_Call {
	return &TransactionController_Rollback_Call{Call: _e.mock.On("Rollback")}
}

func (_c *TransactionController_Rollback_Call) Run(run func()) *TransactionController_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TransactionController_Rollback_Call) Return(_a0 error) *TransactionController_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionController_Rollback_Call) RunAndReturn(run func() error) *TransactionController_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransactionController creates a new instance of TransactionController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionController(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionController {
	mock := &TransactionController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
