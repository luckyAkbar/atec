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
	Create(ctx context.Context, input CreatePackageInput) (*model.Package, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Package, error)
	Update(ctx context.Context, id uuid.UUID, input UpdatePackageInput) (*model.Package, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewPackageRepo create new package repo instance
func NewPackageRepo(db *gorm.DB) *PackageRepo {
	return &PackageRepo{
		db: db,
	}
}

// CreatePackageInput input
type CreatePackageInput struct {
	UserID        uuid.UUID
	PackageName   string
	Questionnaire model.Questionnaire
}

// Create insert new record of packages to the database
func (r *PackageRepo) Create(ctx context.Context, input CreatePackageInput) (*model.Package, error) {
	pack := &model.Package{
		CreatedBy:     input.UserID,
		Questionnaire: input.Questionnaire,
		Name:          input.PackageName,
	}

	if err := r.db.WithContext(ctx).Create(pack).Error; err != nil {
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
func (r *PackageRepo) Update(ctx context.Context, id uuid.UUID, input UpdatePackageInput) (*model.Package, error) {
	pack := &model.Package{}

	err := r.db.WithContext(ctx).Model(pack).
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
