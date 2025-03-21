package repository

import (
	"github.com/luckyAkbar/atec/internal/usecase"
	"gorm.io/gorm"
)

// TransactionControllerFactory transaction controller factory
type TransactionControllerFactory struct {
	db *gorm.DB
}

// NewTransactionControllerFactory create new TransactionControllerFactory instance
func NewTransactionControllerFactory(db *gorm.DB) *TransactionControllerFactory {
	return &TransactionControllerFactory{
		db: db,
	}
}

// New return the transaction controller by immediately starting the transaction
func (t *TransactionControllerFactory) New() *usecase.TxControllerWrapper {
	tx := t.db.Begin()

	return usecase.NewTransactionControllerFactory(&TransactionController{
		db: tx,
	})
}

// TransactionController transaction controller
type TransactionController struct {
	db *gorm.DB
}

// NewTransactionController create new TransactionController instance
func NewTransactionController(db *gorm.DB) *TransactionController {
	return &TransactionController{
		db: db,
	}
}

// Begin return the transaction object
func (t *TransactionController) Begin() any {
	return t.db
}

// Commit commit the transaction
func (t *TransactionController) Commit() error {
	return t.db.Commit().Error
}

// Rollback rollback the transaction
func (t *TransactionController) Rollback() error {
	return t.db.Rollback().Error
}
