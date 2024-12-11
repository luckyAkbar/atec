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
	ChangeActiveStatus(ctx context.Context, input ChangeActiveStatusInput) (*ChangeActiveStatusOutput, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, input UpdatePackageInput) (*UpdatePackageOutput, error)
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

// ChangeActiveStatusInput input
type ChangeActiveStatusInput struct {
	PackageID    uuid.UUID `validate:"required"`
	ActiveStatus bool
}

// Validate validate ChangeActiveStatusInput
func (casi ChangeActiveStatusInput) Validate() error {
	return common.Validator.Struct(casi)
}

// ChangeActiveStatusOutput output
type ChangeActiveStatusOutput struct {
	Message string
}

// ChangeActiveStatus change package active status from its id. If the package is locked, will raise and forbidden error
func (u *PackageUsecase) ChangeActiveStatus(ctx context.Context, input ChangeActiveStatusInput) (*ChangeActiveStatusOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	pack, err := u.packageRepo.FindByID(ctx, input.PackageID)
	switch err {
	default:
		logger.WithError(err).Error("failed to fetch package from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	// early return if no need to change active status
	if pack.IsActive == input.ActiveStatus {
		return &ChangeActiveStatusOutput{
			Message: "ok",
		}, nil
	}

	if pack.IsLocked {
		return nil, UsecaseError{
			ErrType: ErrForbidden,
			Message: "package is already locked",
		}
	}

	_, err = u.packageRepo.Update(ctx, input.PackageID, repository.UpdatePackageInput{
		ActiveStatus: &input.ActiveStatus,
	})

	if err != nil {
		logger.WithError(err).Error("failed to change package activation status to database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &ChangeActiveStatusOutput{
		Message: "ok",
	}, nil
}

// UpdatePackageInput input
type UpdatePackageInput struct {
	PackageID     uuid.UUID           `validate:"required"`
	PackageName   string              `validate:"required"`
	Questionnaire model.Questionnaire `validate:"required"`
}

// Validate validate UpdatePackageInput
func (upi UpdatePackageInput) Validate() error {
	if err := common.Validator.Struct(upi); err != nil {
		return err
	}

	return upi.Questionnaire.Validate()
}

// UpdatePackageOutput output
type UpdatePackageOutput struct {
	Message string
}

// Update update a package based on its id. Only applicable if the package is not yet locked
func (u *PackageUsecase) Update(ctx context.Context, input UpdatePackageInput) (*UpdatePackageOutput, error) {
	logger := logrus.WithContext(ctx).WithField("id", input.PackageID)

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	pack, err := u.packageRepo.FindByID(ctx, input.PackageID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find package to be updated")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case repository.ErrNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	if pack.IsLocked {
		return nil, UsecaseError{
			ErrType: ErrForbidden,
			Message: "package is already locked",
		}
	}

	_, err = u.packageRepo.Update(ctx, input.PackageID, repository.UpdatePackageInput{
		PackageName:   input.PackageName,
		Questionnaire: &input.Questionnaire,
	})

	if err != nil {
		logger.WithError(err).Error("failed to update package to database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &UpdatePackageOutput{
		Message: "ok",
	}, nil
}

// Delete delete a package with its id by using soft delete technique
func (u *PackageUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	logger := logrus.WithContext(ctx).WithField("id", id.String())

	if err := u.packageRepo.Delete(ctx, id); err != nil {
		logger.WithError(err).Error("failed to delete package to database")

		return UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return nil
}
