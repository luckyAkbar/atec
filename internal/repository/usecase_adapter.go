package repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"gorm.io/gorm"
)

// UsecaseErrorUCAdapter convert repository error to usecase error.
// Default case should be returning usecase.ErrRepoInternal
func UsecaseErrorUCAdapter(repoErr error) error {
	switch repoErr {
	// the default case should always trigger internal server error on the usecase's caller side
	default:
		return usecase.ErrRepoInternal
	case ErrNotFound:
		return usecase.ErrRepoNotFound
	case ErrTimeout:
		return usecase.ErrRepoTimeout
	case nil:
		return nil
	}
}

// ChildRepositoryUCAdapter child repository usecase adapter
type ChildRepositoryUCAdapter struct {
	repo *ChildRepository
}

// NewChildRepositoryUCAdapter create new ChildRepositoryUCAdapter instance
func NewChildRepositoryUCAdapter(repo *ChildRepository) *ChildRepositoryUCAdapter {
	return &ChildRepositoryUCAdapter{
		repo: repo,
	}
}

// Create call the repository's Create method and convert the error to usecase error
func (r *ChildRepositoryUCAdapter) Create(ctx context.Context, input usecase.RepoCreateChildInput) (*model.Child, error) {
	res, err := r.repo.Create(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// FindByID call the repository's FindByID method and convert the error to usecase error
func (r *ChildRepositoryUCAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Child, error) {
	res, err := r.repo.FindByID(ctx, id)

	return res, UsecaseErrorUCAdapter(err)
}

// Search call the repository's Search method and convert the error to usecase error
func (r *ChildRepositoryUCAdapter) Search(ctx context.Context, input usecase.RepoSearchChildInput) ([]model.Child, error) {
	res, err := r.repo.Search(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// Update call the repository's Update method and convert the error to usecase error
func (r *ChildRepositoryUCAdapter) Update(ctx context.Context, id uuid.UUID, input usecase.RepoUpdateChildInput) (*model.Child, error) {
	res, err := r.repo.Update(ctx, id, input)

	return res, UsecaseErrorUCAdapter(err)
}

// DeleteAllUserChildren call the repository's DeleteAllUserChildren method and convert the error to usecase error
func (r *ChildRepositoryUCAdapter) DeleteAllUserChildren(
	ctx context.Context,
	input usecase.RepoDeleteAllUserChildrenInput,
	txController ...any,
) error {
	if len(txController) == 0 {
		return r.repo.DeleteAllUserChildren(ctx, input)
	}

	if tx, ok := txController[0].(*gorm.DB); ok {
		return r.repo.DeleteAllUserChildren(ctx, input, tx)
	}

	return fmt.Errorf("%w: invalid transaction controller, expecting typeof gorm transaction", usecase.ErrRepoInternal)
}

// PackageRepositoryUCAdapter package repository usecase adapter
type PackageRepositoryUCAdapter struct {
	repo *PackageRepo
}

// NewPackageRepositoryUCAdapter create new PackageRepositoryUCAdapter instance
func NewPackageRepositoryUCAdapter(repo *PackageRepo) *PackageRepositoryUCAdapter {
	return &PackageRepositoryUCAdapter{
		repo: repo,
	}
}

// Create call the repository's Create method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) Create(
	ctx context.Context,
	input usecase.RepoCreatePackageInput,
	txControllers ...any,
) (*model.Package, error) {
	if len(txControllers) == 0 {
		return r.repo.Create(ctx, input)
	}

	// need to ensure that if supplied, the transaction controller is a gorm transaction
	if tx, ok := txControllers[0].(*gorm.DB); ok {
		return r.repo.Create(ctx, input, tx)
	}

	return nil, fmt.Errorf("%w: invalid transaction controller, expecting typeof gorm transaction", usecase.ErrRepoInternal)
}

// Delete call the repository's Delete method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.repo.Delete(ctx, id)

	return UsecaseErrorUCAdapter(err)
}

// FindAllActivePackages call the repository's FindAllActivePackages method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) FindAllActivePackages(ctx context.Context) ([]model.Package, error) {
	res, err := r.repo.FindAllActivePackages(ctx)

	return res, UsecaseErrorUCAdapter(err)
}

// FindByID call the repository's FindByID method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	res, err := r.repo.FindByID(ctx, id)

	return res, UsecaseErrorUCAdapter(err)
}

// FindOldestActiveAndLockedPackage call the repository's FindOldestActiveAndLockedPackage method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error) {
	res, err := r.repo.FindOldestActiveAndLockedPackage(ctx)

	return res, UsecaseErrorUCAdapter(err)
}

// Search call the repository's Search method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) Search(ctx context.Context, input usecase.RepoSearchPackageInput) ([]model.Package, error) {
	res, err := r.repo.Search(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// Update call the repository's Update method and convert the error to usecase error
func (r *PackageRepositoryUCAdapter) Update(
	ctx context.Context,
	id uuid.UUID,
	input usecase.RepoUpdatePackageInput,
	txControllers ...any,
) (*model.Package, error) {
	if len(txControllers) == 0 {
		return r.repo.Update(ctx, id, input)
	}

	if tx, ok := txControllers[0].(*gorm.DB); ok {
		return r.repo.Update(ctx, id, input, tx)
	}

	return nil, fmt.Errorf("%w: invalid transaction controller, expecting typeof gorm transaction", usecase.ErrRepoInternal)
}

// ResultRepositoryUCAdapter result repository usecase adapter
type ResultRepositoryUCAdapter struct {
	repo *ResultRepository
}

// NewResultRepositoryUCAdapter create new ResultRepositoryUCAdapter instance
func NewResultRepositoryUCAdapter(repo *ResultRepository) *ResultRepositoryUCAdapter {
	return &ResultRepositoryUCAdapter{
		repo: repo,
	}
}

// Create call the repository's Create method and convert the error to usecase error
func (r *ResultRepositoryUCAdapter) Create(ctx context.Context, input usecase.RepoCreateResultInput) (*model.Result, error) {
	res, err := r.repo.Create(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// FindAllUserHistory call the repository's FindAllUserHistory method and convert the error to usecase error
func (r *ResultRepositoryUCAdapter) FindAllUserHistory(ctx context.Context, input usecase.RepoFindAllUserHistoryInput) ([]model.Result, error) {
	res, err := r.repo.FindAllUserHistory(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// FindByID call the repository's FindByID method and convert the error to usecase error
func (r *ResultRepositoryUCAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.Result, error) {
	res, err := r.repo.FindByID(ctx, id)

	return res, UsecaseErrorUCAdapter(err)
}

// Search call the repository's Search method and convert the error to usecase error
func (r *ResultRepositoryUCAdapter) Search(ctx context.Context, input usecase.RepoSearchResultInput) ([]model.Result, error) {
	res, err := r.repo.Search(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// DeleteAllUserResults call the repository's DeleteAllUserResults method and convert the error to usecase error
func (r *ResultRepositoryUCAdapter) DeleteAllUserResults(
	ctx context.Context,
	input usecase.RepoDeleteAllUserResultsInput,
	txController ...any,
) error {
	if len(txController) == 0 {
		return r.repo.DeleteAllUserResults(ctx, input)
	}

	if tx, ok := txController[0].(*gorm.DB); ok {
		return r.repo.DeleteAllUserResults(ctx, input, tx)
	}

	return fmt.Errorf("%w: invalid transaction controller, expecting typeof gorm transaction", usecase.ErrRepoInternal)
}

// UserRepositoryUCAdapter user repository usecase adapter
type UserRepositoryUCAdapter struct {
	repo *UserRepository
}

// NewUserRepositoryUCAdapter create new UserRepositoryUCAdapter instance
func NewUserRepositoryUCAdapter(repo *UserRepository) *UserRepositoryUCAdapter {
	return &UserRepositoryUCAdapter{
		repo: repo,
	}
}

// Create call the repository's Create method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) Create(ctx context.Context, input usecase.RepoCreateUserInput, txController ...any) (*model.User, error) {
	if len(txController) == 0 {
		return r.repo.Create(ctx, input)
	}

	if tx, ok := txController[0].(*gorm.DB); ok {
		return r.repo.Create(ctx, input, tx)
	}

	return nil, fmt.Errorf(
		"%w: invalid transaction controller, expecting typeof gorm transaction got: %+v",
		usecase.ErrRepoInternal,
		reflect.TypeOf(txController[0]),
	)
}

// FindByEmail call the repository's FindByEmail method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	res, err := r.repo.FindByEmail(ctx, email)

	return res, UsecaseErrorUCAdapter(err)
}

// FindByID call the repository's FindByID method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	res, err := r.repo.FindByID(ctx, id)

	return res, UsecaseErrorUCAdapter(err)
}

// IsAdminAccountExists call the repository's IsAdminAccountExists method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) IsAdminAccountExists(ctx context.Context) (bool, error) {
	res, err := r.repo.IsAdminAccountExists(ctx)

	return res, UsecaseErrorUCAdapter(err)
}

// Search call the repository's Search method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) Search(ctx context.Context, input usecase.RepoSearchUserInput) ([]model.User, error) {
	res, err := r.repo.Search(ctx, input)

	return res, UsecaseErrorUCAdapter(err)
}

// Update call the repository's Update method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) Update(ctx context.Context, userID uuid.UUID, input usecase.RepoUpdateUserInput) (*model.User, error) {
	res, err := r.repo.Update(ctx, userID, input)

	return res, UsecaseErrorUCAdapter(err)
}

// DeleteByID call the repository's DeleteByID method and convert the error to usecase error
func (r *UserRepositoryUCAdapter) DeleteByID(ctx context.Context, input usecase.RepoDeleteUserByIDInput, txController ...any) error {
	if len(txController) == 0 {
		return r.repo.DeleteByID(ctx, input)
	}

	if tx, ok := txController[0].(*gorm.DB); ok {
		return r.repo.DeleteByID(ctx, input, tx)
	}

	return fmt.Errorf("%w: invalid transaction controller, expecting typeof gorm transaction", usecase.ErrRepoInternal)
}
