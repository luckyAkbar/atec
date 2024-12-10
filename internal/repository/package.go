package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
)

// PackageRepo repository for packages
type PackageRepo struct {
	db *gorm.DB
}

// PackageRepoIface interface for PackageRepo
type PackageRepoIface interface {
	Create(ctx context.Context, input CreatePackageInput) (*model.Package, error)
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
