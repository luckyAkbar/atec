package usecase

import (
	"context"
	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
)

// RepositoryError custom specialized error type expected to be returned
// by the repository implementation
type RepositoryError error

var (
	// ErrRepoInternal represent unhandled error happening from the database server
	ErrRepoInternal RepositoryError = errors.New("internal database server error")

	// ErrRepoNotFound represent data not found on the database
	ErrRepoNotFound RepositoryError = errors.New("data not found")

	// ErrRepoTimeout represent timeout or process took too long
	ErrRepoTimeout RepositoryError = errors.New("timeout")
)

// TransactionControllerFactory transaction controller factory
type TransactionControllerFactory interface {
	New() *TxControllerWrapper
}

// TransactionController transaction controller capabilities
type TransactionController interface {
	Commit() error
	Rollback() error

	// Begin return the transaction object
	Begin() any
}

// TxControllerWrapper will wrap the underlying implementation of the transaction controller
// to confronts with dependency rules without violating ireturn lint rules
type TxControllerWrapper struct {
	impl TransactionController
}

// NewTxControllerWrapper create new TxControllerWrapper instance
func NewTxControllerWrapper(impl TransactionController) *TxControllerWrapper {
	return &TxControllerWrapper{impl}
}

// NewTransactionControllerFactory create new TransactionControllerFactory instance
func NewTransactionControllerFactory(impl TransactionController) *TxControllerWrapper {
	return &TxControllerWrapper{impl}
}

// Commit call the underlying implementation of the transaction controller's Commit
func (t *TxControllerWrapper) Commit() error {
	return t.impl.Commit()
}

// Rollback call the underlying implementation of the transaction controller's Rollback
func (t *TxControllerWrapper) Rollback() error {
	return t.impl.Rollback()
}

// Begin call the underlying implementation of the transaction controller's Begin
func (t *TxControllerWrapper) Begin() any {
	return t.impl.Begin()
}

// RepoCreateChildInput input
type RepoCreateChildInput struct {
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	Name         string
}

// RepoSearchChildInput input to search child data. everything marked as pointer to a datatype means it is optional
type RepoSearchChildInput struct {
	ParentUserID *uuid.UUID
	Name         *string
	Gender       *bool
	Limit        int
	Offset       int
}

// ChildRepository interface
type ChildRepository interface {
	Create(ctx context.Context, input RepoCreateChildInput) (*model.Child, error)
	Update(ctx context.Context, id uuid.UUID, input RepoUpdateChildInput) (*model.Child, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Child, error)
	Search(ctx context.Context, input RepoSearchChildInput) ([]model.Child, error)
}

// RepoRegisterChildInput input
type RepoRegisterChildInput struct {
	DateOfBirth time.Time `validate:"required"`
	Gender      bool
	Name        string `validate:"required"`
}

// RepoUpdateChildInput input
type RepoUpdateChildInput struct {
	DateOfBirth *time.Time
	Gender      *bool
	Name        *string
}

// RepoUpdateUserInput options to update user record
type RepoUpdateUserInput struct {
	Email    string
	Password string `json:"-"`
	Username string
	IsActive *bool
}

// RepoCreateUserInput input to create a new user data
type RepoCreateUserInput struct {
	Email    string
	Password string
	Username string
	IsActive bool
	Roles    model.Roles
}

// RepoSearchUserInput options to search users
type RepoSearchUserInput struct {
	Role   model.Roles
	Limit  int
	Offset int
}

// UserRepository interface exported by UserRepository to help ease mocking
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, input RepoCreateUserInput, txController ...any) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, userID uuid.UUID, input RepoUpdateUserInput) (*model.User, error)
	Search(ctx context.Context, input RepoSearchUserInput) ([]model.User, error)
	IsAdminAccountExists(ctx context.Context) (bool, error)
}

// RepoCreateResultInput create result input
type RepoCreateResultInput struct {
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Answer    model.AnswerDetail
	Result    model.ResultDetail
}

// RepoSearchResultInput search result input
type RepoSearchResultInput struct {
	ID        uuid.UUID
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Limit     int
	Offset    int
}

// RepoFindAllUserHistoryInput input
type RepoFindAllUserHistoryInput struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

// ResultRepository result repository interface
type ResultRepository interface {
	Create(ctx context.Context, input RepoCreateResultInput) (*model.Result, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Result, error)
	Search(ctx context.Context, input RepoSearchResultInput) ([]model.Result, error)
	FindAllUserHistory(ctx context.Context, input RepoFindAllUserHistoryInput) ([]model.Result, error)
}

// RepoCreatePackageInput input
type RepoCreatePackageInput struct {
	UserID                  uuid.UUID
	PackageName             string
	Questionnaire           model.Questionnaire
	IndicationCategories    model.IndicationCategories
	ImageResultAttributeKey model.ImageResultAttributeKey
}

// RepoUpdatePackageInput input
type RepoUpdatePackageInput struct {
	ActiveStatus *bool
	LockStatus   *bool

	Questionnaire *model.Questionnaire
	PackageName   string
}

// RepoSearchPackageInput input to search package. any fields typed with a pointer means it is optional
type RepoSearchPackageInput struct {
	IsActive *bool
	Limit    int
}

// PackageRepo interface for PackageRepo
type PackageRepo interface {
	Create(ctx context.Context, input RepoCreatePackageInput, txControllers ...any) (*model.Package, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error)
	Update(ctx context.Context, id uuid.UUID, input RepoUpdatePackageInput, txControllers ...any) (*model.Package, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, input RepoSearchPackageInput) ([]model.Package, error)
	FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error)
	FindAllActivePackages(ctx context.Context) ([]model.Package, error)
}
