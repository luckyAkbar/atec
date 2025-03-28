// Code generated by mockery v2.52.2. DO NOT EDIT.

package usecase

import (
	context "context"

	usecase "github.com/luckyAkbar/atec/internal/usecase"
	mock "github.com/stretchr/testify/mock"
)

// ChildUsecaseIface is an autogenerated mock type for the ChildUsecaseIface type
type ChildUsecaseIface struct {
	mock.Mock
}

type ChildUsecaseIface_Expecter struct {
	mock *mock.Mock
}

func (_m *ChildUsecaseIface) EXPECT() *ChildUsecaseIface_Expecter {
	return &ChildUsecaseIface_Expecter{mock: &_m.Mock}
}

// GetRegisteredChildren provides a mock function with given fields: ctx, input
func (_m *ChildUsecaseIface) GetRegisteredChildren(ctx context.Context, input usecase.GetRegisteredChildrenInput) ([]usecase.GetRegisteredChildrenOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for GetRegisteredChildren")
	}

	var r0 []usecase.GetRegisteredChildrenOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.GetRegisteredChildrenInput) ([]usecase.GetRegisteredChildrenOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.GetRegisteredChildrenInput) []usecase.GetRegisteredChildrenOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]usecase.GetRegisteredChildrenOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.GetRegisteredChildrenInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChildUsecaseIface_GetRegisteredChildren_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRegisteredChildren'
type ChildUsecaseIface_GetRegisteredChildren_Call struct {
	*mock.Call
}

// GetRegisteredChildren is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.GetRegisteredChildrenInput
func (_e *ChildUsecaseIface_Expecter) GetRegisteredChildren(ctx interface{}, input interface{}) *ChildUsecaseIface_GetRegisteredChildren_Call {
	return &ChildUsecaseIface_GetRegisteredChildren_Call{Call: _e.mock.On("GetRegisteredChildren", ctx, input)}
}

func (_c *ChildUsecaseIface_GetRegisteredChildren_Call) Run(run func(ctx context.Context, input usecase.GetRegisteredChildrenInput)) *ChildUsecaseIface_GetRegisteredChildren_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.GetRegisteredChildrenInput))
	})
	return _c
}

func (_c *ChildUsecaseIface_GetRegisteredChildren_Call) Return(_a0 []usecase.GetRegisteredChildrenOutput, _a1 error) *ChildUsecaseIface_GetRegisteredChildren_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChildUsecaseIface_GetRegisteredChildren_Call) RunAndReturn(run func(context.Context, usecase.GetRegisteredChildrenInput) ([]usecase.GetRegisteredChildrenOutput, error)) *ChildUsecaseIface_GetRegisteredChildren_Call {
	_c.Call.Return(run)
	return _c
}

// HandleGetStatistic provides a mock function with given fields: ctx, input
func (_m *ChildUsecaseIface) HandleGetStatistic(ctx context.Context, input usecase.GetStatisticInput) (*usecase.GetStatisticOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for HandleGetStatistic")
	}

	var r0 *usecase.GetStatisticOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.GetStatisticInput) (*usecase.GetStatisticOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.GetStatisticInput) *usecase.GetStatisticOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usecase.GetStatisticOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.GetStatisticInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChildUsecaseIface_HandleGetStatistic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HandleGetStatistic'
type ChildUsecaseIface_HandleGetStatistic_Call struct {
	*mock.Call
}

// HandleGetStatistic is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.GetStatisticInput
func (_e *ChildUsecaseIface_Expecter) HandleGetStatistic(ctx interface{}, input interface{}) *ChildUsecaseIface_HandleGetStatistic_Call {
	return &ChildUsecaseIface_HandleGetStatistic_Call{Call: _e.mock.On("HandleGetStatistic", ctx, input)}
}

func (_c *ChildUsecaseIface_HandleGetStatistic_Call) Run(run func(ctx context.Context, input usecase.GetStatisticInput)) *ChildUsecaseIface_HandleGetStatistic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.GetStatisticInput))
	})
	return _c
}

func (_c *ChildUsecaseIface_HandleGetStatistic_Call) Return(_a0 *usecase.GetStatisticOutput, _a1 error) *ChildUsecaseIface_HandleGetStatistic_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChildUsecaseIface_HandleGetStatistic_Call) RunAndReturn(run func(context.Context, usecase.GetStatisticInput) (*usecase.GetStatisticOutput, error)) *ChildUsecaseIface_HandleGetStatistic_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: ctx, input
func (_m *ChildUsecaseIface) Register(ctx context.Context, input usecase.RegisterChildInput) (*usecase.RegisterChildOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 *usecase.RegisterChildOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.RegisterChildInput) (*usecase.RegisterChildOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.RegisterChildInput) *usecase.RegisterChildOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usecase.RegisterChildOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.RegisterChildInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChildUsecaseIface_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type ChildUsecaseIface_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.RegisterChildInput
func (_e *ChildUsecaseIface_Expecter) Register(ctx interface{}, input interface{}) *ChildUsecaseIface_Register_Call {
	return &ChildUsecaseIface_Register_Call{Call: _e.mock.On("Register", ctx, input)}
}

func (_c *ChildUsecaseIface_Register_Call) Run(run func(ctx context.Context, input usecase.RegisterChildInput)) *ChildUsecaseIface_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.RegisterChildInput))
	})
	return _c
}

func (_c *ChildUsecaseIface_Register_Call) Return(_a0 *usecase.RegisterChildOutput, _a1 error) *ChildUsecaseIface_Register_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChildUsecaseIface_Register_Call) RunAndReturn(run func(context.Context, usecase.RegisterChildInput) (*usecase.RegisterChildOutput, error)) *ChildUsecaseIface_Register_Call {
	_c.Call.Return(run)
	return _c
}

// Search provides a mock function with given fields: ctx, input
func (_m *ChildUsecaseIface) Search(ctx context.Context, input usecase.SearchChildInput) ([]usecase.SearchChildOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 []usecase.SearchChildOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.SearchChildInput) ([]usecase.SearchChildOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.SearchChildInput) []usecase.SearchChildOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]usecase.SearchChildOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.SearchChildInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChildUsecaseIface_Search_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Search'
type ChildUsecaseIface_Search_Call struct {
	*mock.Call
}

// Search is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.SearchChildInput
func (_e *ChildUsecaseIface_Expecter) Search(ctx interface{}, input interface{}) *ChildUsecaseIface_Search_Call {
	return &ChildUsecaseIface_Search_Call{Call: _e.mock.On("Search", ctx, input)}
}

func (_c *ChildUsecaseIface_Search_Call) Run(run func(ctx context.Context, input usecase.SearchChildInput)) *ChildUsecaseIface_Search_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.SearchChildInput))
	})
	return _c
}

func (_c *ChildUsecaseIface_Search_Call) Return(_a0 []usecase.SearchChildOutput, _a1 error) *ChildUsecaseIface_Search_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChildUsecaseIface_Search_Call) RunAndReturn(run func(context.Context, usecase.SearchChildInput) ([]usecase.SearchChildOutput, error)) *ChildUsecaseIface_Search_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, input
func (_m *ChildUsecaseIface) Update(ctx context.Context, input usecase.UpdateChildInput) (*usecase.UpdateChildOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *usecase.UpdateChildOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.UpdateChildInput) (*usecase.UpdateChildOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.UpdateChildInput) *usecase.UpdateChildOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usecase.UpdateChildOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.UpdateChildInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChildUsecaseIface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ChildUsecaseIface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - input usecase.UpdateChildInput
func (_e *ChildUsecaseIface_Expecter) Update(ctx interface{}, input interface{}) *ChildUsecaseIface_Update_Call {
	return &ChildUsecaseIface_Update_Call{Call: _e.mock.On("Update", ctx, input)}
}

func (_c *ChildUsecaseIface_Update_Call) Run(run func(ctx context.Context, input usecase.UpdateChildInput)) *ChildUsecaseIface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.UpdateChildInput))
	})
	return _c
}

func (_c *ChildUsecaseIface_Update_Call) Return(_a0 *usecase.UpdateChildOutput, _a1 error) *ChildUsecaseIface_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChildUsecaseIface_Update_Call) RunAndReturn(run func(context.Context, usecase.UpdateChildInput) (*usecase.UpdateChildOutput, error)) *ChildUsecaseIface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewChildUsecaseIface creates a new instance of ChildUsecaseIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChildUsecaseIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChildUsecaseIface {
	mock := &ChildUsecaseIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
