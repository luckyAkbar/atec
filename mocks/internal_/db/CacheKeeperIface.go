// Code generated by mockery v2.52.2. DO NOT EDIT.

package db

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	redsync "github.com/go-redsync/redsync/v4"

	time "time"
)

// CacheKeeperIface is an autogenerated mock type for the CacheKeeperIface type
type CacheKeeperIface struct {
	mock.Mock
}

type CacheKeeperIface_Expecter struct {
	mock *mock.Mock
}

func (_m *CacheKeeperIface) EXPECT() *CacheKeeperIface_Expecter {
	return &CacheKeeperIface_Expecter{mock: &_m.Mock}
}

// AcquireLock provides a mock function with given fields: key
func (_m *CacheKeeperIface) AcquireLock(key string) (*redsync.Mutex, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for AcquireLock")
	}

	var r0 *redsync.Mutex
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*redsync.Mutex, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) *redsync.Mutex); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*redsync.Mutex)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CacheKeeperIface_AcquireLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AcquireLock'
type CacheKeeperIface_AcquireLock_Call struct {
	*mock.Call
}

// AcquireLock is a helper method to define mock.On call
//   - key string
func (_e *CacheKeeperIface_Expecter) AcquireLock(key interface{}) *CacheKeeperIface_AcquireLock_Call {
	return &CacheKeeperIface_AcquireLock_Call{Call: _e.mock.On("AcquireLock", key)}
}

func (_c *CacheKeeperIface_AcquireLock_Call) Run(run func(key string)) *CacheKeeperIface_AcquireLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CacheKeeperIface_AcquireLock_Call) Return(_a0 *redsync.Mutex, _a1 error) *CacheKeeperIface_AcquireLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CacheKeeperIface_AcquireLock_Call) RunAndReturn(run func(string) (*redsync.Mutex, error)) *CacheKeeperIface_AcquireLock_Call {
	_c.Call.Return(run)
	return _c
}

// Del provides a mock function with given fields: ctx, key
func (_m *CacheKeeperIface) Del(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Del")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheKeeperIface_Del_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Del'
type CacheKeeperIface_Del_Call struct {
	*mock.Call
}

// Del is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *CacheKeeperIface_Expecter) Del(ctx interface{}, key interface{}) *CacheKeeperIface_Del_Call {
	return &CacheKeeperIface_Del_Call{Call: _e.mock.On("Del", ctx, key)}
}

func (_c *CacheKeeperIface_Del_Call) Run(run func(ctx context.Context, key string)) *CacheKeeperIface_Del_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CacheKeeperIface_Del_Call) Return(_a0 error) *CacheKeeperIface_Del_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheKeeperIface_Del_Call) RunAndReturn(run func(context.Context, string) error) *CacheKeeperIface_Del_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, key
func (_m *CacheKeeperIface) Get(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CacheKeeperIface_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type CacheKeeperIface_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *CacheKeeperIface_Expecter) Get(ctx interface{}, key interface{}) *CacheKeeperIface_Get_Call {
	return &CacheKeeperIface_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *CacheKeeperIface_Get_Call) Run(run func(ctx context.Context, key string)) *CacheKeeperIface_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CacheKeeperIface_Get_Call) Return(_a0 string, _a1 error) *CacheKeeperIface_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CacheKeeperIface_Get_Call) RunAndReturn(run func(context.Context, string) (string, error)) *CacheKeeperIface_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrLock provides a mock function with given fields: ctx, key
func (_m *CacheKeeperIface) GetOrLock(ctx context.Context, key string) (string, *redsync.Mutex, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for GetOrLock")
	}

	var r0 string
	var r1 *redsync.Mutex
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, *redsync.Mutex, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *redsync.Mutex); ok {
		r1 = rf(ctx, key)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*redsync.Mutex)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, key)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CacheKeeperIface_GetOrLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrLock'
type CacheKeeperIface_GetOrLock_Call struct {
	*mock.Call
}

// GetOrLock is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *CacheKeeperIface_Expecter) GetOrLock(ctx interface{}, key interface{}) *CacheKeeperIface_GetOrLock_Call {
	return &CacheKeeperIface_GetOrLock_Call{Call: _e.mock.On("GetOrLock", ctx, key)}
}

func (_c *CacheKeeperIface_GetOrLock_Call) Run(run func(ctx context.Context, key string)) *CacheKeeperIface_GetOrLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *CacheKeeperIface_GetOrLock_Call) Return(_a0 string, _a1 *redsync.Mutex, _a2 error) *CacheKeeperIface_GetOrLock_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *CacheKeeperIface_GetOrLock_Call) RunAndReturn(run func(context.Context, string) (string, *redsync.Mutex, error)) *CacheKeeperIface_GetOrLock_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value, expiration
func (_m *CacheKeeperIface) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	ret := _m.Called(ctx, key, value, expiration)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Duration) error); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheKeeperIface_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type CacheKeeperIface_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value string
//   - expiration time.Duration
func (_e *CacheKeeperIface_Expecter) Set(ctx interface{}, key interface{}, value interface{}, expiration interface{}) *CacheKeeperIface_Set_Call {
	return &CacheKeeperIface_Set_Call{Call: _e.mock.On("Set", ctx, key, value, expiration)}
}

func (_c *CacheKeeperIface_Set_Call) Run(run func(ctx context.Context, key string, value string, expiration time.Duration)) *CacheKeeperIface_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Duration))
	})
	return _c
}

func (_c *CacheKeeperIface_Set_Call) Return(_a0 error) *CacheKeeperIface_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheKeeperIface_Set_Call) RunAndReturn(run func(context.Context, string, string, time.Duration) error) *CacheKeeperIface_Set_Call {
	_c.Call.Return(run)
	return _c
}

// SetJSON provides a mock function with given fields: ctx, key, value, expiration
func (_m *CacheKeeperIface) SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	ret := _m.Called(ctx, key, value, expiration)

	if len(ret) == 0 {
		panic("no return value specified for SetJSON")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, any, time.Duration) error); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheKeeperIface_SetJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetJSON'
type CacheKeeperIface_SetJSON_Call struct {
	*mock.Call
}

// SetJSON is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value any
//   - expiration time.Duration
func (_e *CacheKeeperIface_Expecter) SetJSON(ctx interface{}, key interface{}, value interface{}, expiration interface{}) *CacheKeeperIface_SetJSON_Call {
	return &CacheKeeperIface_SetJSON_Call{Call: _e.mock.On("SetJSON", ctx, key, value, expiration)}
}

func (_c *CacheKeeperIface_SetJSON_Call) Run(run func(ctx context.Context, key string, value any, expiration time.Duration)) *CacheKeeperIface_SetJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(any), args[3].(time.Duration))
	})
	return _c
}

func (_c *CacheKeeperIface_SetJSON_Call) Return(_a0 error) *CacheKeeperIface_SetJSON_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheKeeperIface_SetJSON_Call) RunAndReturn(run func(context.Context, string, any, time.Duration) error) *CacheKeeperIface_SetJSON_Call {
	_c.Call.Return(run)
	return _c
}

// SetNil provides a mock function with given fields: ctx, key, expiration
func (_m *CacheKeeperIface) SetNil(ctx context.Context, key string, expiration ...time.Duration) error {
	_va := make([]interface{}, len(expiration))
	for _i := range expiration {
		_va[_i] = expiration[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, key)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SetNil")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...time.Duration) error); ok {
		r0 = rf(ctx, key, expiration...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheKeeperIface_SetNil_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetNil'
type CacheKeeperIface_SetNil_Call struct {
	*mock.Call
}

// SetNil is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - expiration ...time.Duration
func (_e *CacheKeeperIface_Expecter) SetNil(ctx interface{}, key interface{}, expiration ...interface{}) *CacheKeeperIface_SetNil_Call {
	return &CacheKeeperIface_SetNil_Call{Call: _e.mock.On("SetNil",
		append([]interface{}{ctx, key}, expiration...)...)}
}

func (_c *CacheKeeperIface_SetNil_Call) Run(run func(ctx context.Context, key string, expiration ...time.Duration)) *CacheKeeperIface_SetNil_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]time.Duration, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(time.Duration)
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *CacheKeeperIface_SetNil_Call) Return(_a0 error) *CacheKeeperIface_SetNil_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheKeeperIface_SetNil_Call) RunAndReturn(run func(context.Context, string, ...time.Duration) error) *CacheKeeperIface_SetNil_Call {
	_c.Call.Return(run)
	return _c
}

// NewCacheKeeperIface creates a new instance of CacheKeeperIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCacheKeeperIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *CacheKeeperIface {
	mock := &CacheKeeperIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
