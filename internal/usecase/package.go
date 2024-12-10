package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
)

// PackageUsecase usecase for package
type PackageUsecase struct {
	packageRepo repository.PackageRepoIface
}

// PackageUsecaseIface interface
type PackageUsecaseIface interface {
	Create(ctx context.Context, input CreatePackageInput) (*CreatePackageOutput, error)
}

// NewPackageUsecase create new PackageUsecase instance
func NewPackageUsecase(packageRepo *repository.PackageRepo) *PackageUsecase {
	return &PackageUsecase{
		packageRepo: packageRepo,
	}
}

// CreatePackageInput input by embedding direcly model.Questionnaire to simplify the input anotation
type CreatePackageInput struct {
	PackageName   string              `validate:"required"`
	Questionnaire model.Questionnaire `validate:"required"`
}

// Validate validate CreatePackageInput
func (cpi CreatePackageInput) Validate() error {
	if err := common.Validator.Struct(cpi); err != nil {
		return err
	}

	return cpi.Questionnaire.Validate()
}

// CreatePackageOutput output
type CreatePackageOutput struct {
	ID uuid.UUID
}

// Create create package
func (u *PackageUsecase) Create(ctx context.Context, input CreatePackageInput) (*CreatePackageOutput, error) {
	user := model.GetUserFromCtx(ctx)
	if user == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	logger := logrus.WithContext(ctx).WithField("user", helper.Dump(user))

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	pack, err := u.packageRepo.Create(ctx, repository.CreatePackageInput{
		UserID:        user.ID,
		PackageName:   input.PackageName,
		Questionnaire: input.Questionnaire,
	})

	if err != nil {
		logger.WithError(err).Error("failed to write package data to database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &CreatePackageOutput{
		ID: pack.ID,
	}, nil
}
