// Code generated by mockery v2.52.2. DO NOT EDIT.

package common

import mock "github.com/stretchr/testify/mock"

// RedsyncMutex is an autogenerated mock type for the RedsyncMutex type
type RedsyncMutex struct {
	mock.Mock
}

type RedsyncMutex_Expecter struct {
	mock *mock.Mock
}

func (_m *RedsyncMutex) EXPECT() *RedsyncMutex_Expecter {
	return &RedsyncMutex_Expecter{mock: &_m.Mock}
}

// Lock provides a mock function with no fields
func (_m *RedsyncMutex) Lock() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Lock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RedsyncMutex_Lock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Lock'
type RedsyncMutex_Lock_Call struct {
	*mock.Call
}

// Lock is a helper method to define mock.On call
func (_e *RedsyncMutex_Expecter) Lock() *RedsyncMutex_Lock_Call {
	return &RedsyncMutex_Lock_Call{Call: _e.mock.On("Lock")}
}

func (_c *RedsyncMutex_Lock_Call) Run(run func()) *RedsyncMutex_Lock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *RedsyncMutex_Lock_Call) Return(_a0 error) *RedsyncMutex_Lock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RedsyncMutex_Lock_Call) RunAndReturn(run func() error) *RedsyncMutex_Lock_Call {
	_c.Call.Return(run)
	return _c
}

// Unlock provides a mock function with no fields
func (_m *RedsyncMutex) Unlock() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Unlock")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RedsyncMutex_Unlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unlock'
type RedsyncMutex_Unlock_Call struct {
	*mock.Call
}

// Unlock is a helper method to define mock.On call
func (_e *RedsyncMutex_Expecter) Unlock() *RedsyncMutex_Unlock_Call {
	return &RedsyncMutex_Unlock_Call{Call: _e.mock.On("Unlock")}
}

func (_c *RedsyncMutex_Unlock_Call) Run(run func()) *RedsyncMutex_Unlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *RedsyncMutex_Unlock_Call) Return(_a0 bool, _a1 error) *RedsyncMutex_Unlock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RedsyncMutex_Unlock_Call) RunAndReturn(run func() (bool, error)) *RedsyncMutex_Unlock_Call {
	_c.Call.Return(run)
	return _c
}

// NewRedsyncMutex creates a new instance of RedsyncMutex. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRedsyncMutex(t interface {
	mock.TestingT
	Cleanup(func())
}) *RedsyncMutex {
	mock := &RedsyncMutex{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
