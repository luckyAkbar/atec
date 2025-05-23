// Code generated by mockery v2.52.2. DO NOT EDIT.

package common

import (
	context "context"

	common "github.com/luckyAkbar/atec/internal/common"

	mock "github.com/stretchr/testify/mock"
)

// DistributedLockerIface is an autogenerated mock type for the DistributedLockerIface type
type DistributedLockerIface struct {
	mock.Mock
}

type DistributedLockerIface_Expecter struct {
	mock *mock.Mock
}

func (_m *DistributedLockerIface) EXPECT() *DistributedLockerIface_Expecter {
	return &DistributedLockerIface_Expecter{mock: &_m.Mock}
}

// GetLock provides a mock function with given fields: key
func (_m *DistributedLockerIface) GetLock(key string) (*common.RedsyncMutexWrapper, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for GetLock")
	}

	var r0 *common.RedsyncMutexWrapper
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*common.RedsyncMutexWrapper, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) *common.RedsyncMutexWrapper); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.RedsyncMutexWrapper)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DistributedLockerIface_GetLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLock'
type DistributedLockerIface_GetLock_Call struct {
	*mock.Call
}

// GetLock is a helper method to define mock.On call
//   - key string
func (_e *DistributedLockerIface_Expecter) GetLock(key interface{}) *DistributedLockerIface_GetLock_Call {
	return &DistributedLockerIface_GetLock_Call{Call: _e.mock.On("GetLock", key)}
}

func (_c *DistributedLockerIface_GetLock_Call) Run(run func(key string)) *DistributedLockerIface_GetLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DistributedLockerIface_GetLock_Call) Return(_a0 *common.RedsyncMutexWrapper, _a1 error) *DistributedLockerIface_GetLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DistributedLockerIface_GetLock_Call) RunAndReturn(run func(string) (*common.RedsyncMutexWrapper, error)) *DistributedLockerIface_GetLock_Call {
	_c.Call.Return(run)
	return _c
}

// IsLocked provides a mock function with given fields: ctx, key
func (_m *DistributedLockerIface) IsLocked(ctx context.Context, key string) (bool, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for IsLocked")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DistributedLockerIface_IsLocked_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsLocked'
type DistributedLockerIface_IsLocked_Call struct {
	*mock.Call
}

// IsLocked is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *DistributedLockerIface_Expecter) IsLocked(ctx interface{}, key interface{}) *DistributedLockerIface_IsLocked_Call {
	return &DistributedLockerIface_IsLocked_Call{Call: _e.mock.On("IsLocked", ctx, key)}
}

func (_c *DistributedLockerIface_IsLocked_Call) Run(run func(ctx context.Context, key string)) *DistributedLockerIface_IsLocked_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *DistributedLockerIface_IsLocked_Call) Return(_a0 bool, _a1 error) *DistributedLockerIface_IsLocked_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DistributedLockerIface_IsLocked_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *DistributedLockerIface_IsLocked_Call {
	_c.Call.Return(run)
	return _c
}

// NewDistributedLockerIface creates a new instance of DistributedLockerIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDistributedLockerIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DistributedLockerIface {
	mock := &DistributedLockerIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
