// Code generated by mockery v2.52.2. DO NOT EDIT.

package usecase

import (
	usecase "github.com/luckyAkbar/atec/internal/usecase"
	mock "github.com/stretchr/testify/mock"
)

// TransactionControllerFactory is an autogenerated mock type for the TransactionControllerFactory type
type TransactionControllerFactory struct {
	mock.Mock
}

type TransactionControllerFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *TransactionControllerFactory) EXPECT() *TransactionControllerFactory_Expecter {
	return &TransactionControllerFactory_Expecter{mock: &_m.Mock}
}

// New provides a mock function with no fields
func (_m *TransactionControllerFactory) New() *usecase.TxControllerWrapper {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for New")
	}

	var r0 *usecase.TxControllerWrapper
	if rf, ok := ret.Get(0).(func() *usecase.TxControllerWrapper); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*usecase.TxControllerWrapper)
		}
	}

	return r0
}

// TransactionControllerFactory_New_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'New'
type TransactionControllerFactory_New_Call struct {
	*mock.Call
}

// New is a helper method to define mock.On call
func (_e *TransactionControllerFactory_Expecter) New() *TransactionControllerFactory_New_Call {
	return &TransactionControllerFactory_New_Call{Call: _e.mock.On("New")}
}

func (_c *TransactionControllerFactory_New_Call) Run(run func()) *TransactionControllerFactory_New_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TransactionControllerFactory_New_Call) Return(_a0 *usecase.TxControllerWrapper) *TransactionControllerFactory_New_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionControllerFactory_New_Call) RunAndReturn(run func() *usecase.TxControllerWrapper) *TransactionControllerFactory_New_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransactionControllerFactory creates a new instance of TransactionControllerFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionControllerFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionControllerFactory {
	mock := &TransactionControllerFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
