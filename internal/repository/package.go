package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PackageRepo repository for packages
type PackageRepo struct {
	db *gorm.DB
}

// PackageRepoIface interface for PackageRepo
type PackageRepoIface interface {
	Create(ctx context.Context, input CreatePackageInput, txControllers ...*gorm.DB) (*model.Package, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error)
	Update(ctx context.Context, id uuid.UUID, input UpdatePackageInput, txControllers ...*gorm.DB) (*model.Package, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, input SearchPackageInput) ([]model.Package, error)
	FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error)
}

// NewPackageRepo create new package repo instance
func NewPackageRepo(db *gorm.DB) *PackageRepo {
	return &PackageRepo{
		db: db,
	}
}

// CreatePackageInput input
type CreatePackageInput struct {
	UserID                  uuid.UUID
	PackageName             string
	Questionnaire           model.Questionnaire
	IndicationCategories    model.IndicationCategories
	ImageResultAttributeKey model.ImageResultAttributeKey
}

// Create insert new record of packages to the database
func (r *PackageRepo) Create(ctx context.Context, input CreatePackageInput, txControllers ...*gorm.DB) (*model.Package, error) {
	tx := r.db
	if len(txControllers) > 0 {
		tx = txControllers[0]
	}

	pack := &model.Package{
		CreatedBy:               input.UserID,
		Questionnaire:           input.Questionnaire,
		Name:                    input.PackageName,
		IndicationCategories:    input.IndicationCategories,
		ImageResultAttributeKey: input.ImageResultAttributeKey,
	}

	if err := tx.WithContext(ctx).Create(pack).Error; err != nil {
		return nil, err
	}

	return pack, nil
}

// FindByID find package by id
func (r *PackageRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error) {
	pack := &model.Package{}

	err := r.db.WithContext(ctx).Take(pack, "id = ?", id).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return pack, nil
	}
}

// UpdatePackageInput input
type UpdatePackageInput struct {
	ActiveStatus *bool

	Questionnaire *model.Questionnaire
	PackageName   string
}

// ToUpdateFields convert the update params to gorm dynamic update fields
func (upi UpdatePackageInput) ToUpdateFields() map[string]interface{} {
	fields := map[string]interface{}{}

	if upi.ActiveStatus != nil {
		fields["is_active"] = *upi.ActiveStatus
	}

	if upi.PackageName != "" {
		fields["name"] = upi.PackageName
	}

	if upi.Questionnaire != nil {
		fields["questionnaire"] = upi.Questionnaire
	}

	return fields
}

// Update update package record by its id
func (r *PackageRepo) Update(ctx context.Context, id uuid.UUID, input UpdatePackageInput, txControllers ...*gorm.DB) (*model.Package, error) {
	tx := r.db
	if len(txControllers) > 0 {
		tx = txControllers[0]
	}

	pack := &model.Package{}

	err := tx.WithContext(ctx).Model(pack).
		Clauses(clause.Returning{}).Where("id = ?", id).
		Updates(input.ToUpdateFields()).Error

	if err != nil {
		return nil, err
	}

	return pack, nil
}

// Delete use soft delete to delett a package record by its id
func (r *PackageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Package{ID: id}).Error
}

// SearchPackageInput input to search package. any fields typed with a pointer means it is optional
type SearchPackageInput struct {
	IsActive *bool
}

func (spi SearchPackageInput) toSearchFields(cursor *gorm.DB) *gorm.DB {
	if spi.IsActive != nil {
		cursor = cursor.Where("is_active = ?", *spi.IsActive)
	}

	return cursor
}

// Search search package based on provided parameters
func (r *PackageRepo) Search(ctx context.Context, input SearchPackageInput) ([]model.Package, error) {
	packages := []model.Package{}

	conn := r.db.WithContext(ctx)
	cursor := input.toSearchFields(conn)

	if err := cursor.Find(&packages).Error; err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return nil, ErrNotFound
	}

	return packages, nil
}

// FindOldestActiveAndLockedPackage get the oldest active and locked package
func (r *PackageRepo) FindOldestActiveAndLockedPackage(ctx context.Context) (*model.Package, error) {
	pack := &model.Package{}

	err := r.db.WithContext(ctx).Where("is_active = ? AND is_locked = ?", true, true).Order("created_at ASC").Take(pack).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return pack, nil
	}
}
