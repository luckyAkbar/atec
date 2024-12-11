package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
)

// ChildUsecase child usecase
type ChildUsecase struct {
	childRepo repository.ChildRepositoryIface
}

// ChildUsecaseIface interface
type ChildUsecaseIface interface {
	Register(ctx context.Context, input RegisterChildInput) (*RegisterChildOutput, error)
	Update(ctx context.Context, input UpdateChildInput) (*UpdateChildOutput, error)
	GetRegisteredChildren(ctx context.Context, input GetRegisteredChildrenInput) ([]GetRegisteredChildrenOutput, error)
}

// NewChildUsecase create new ChildUsecase instance
func NewChildUsecase(childRepo *repository.ChildRepository) *ChildUsecase {
	return &ChildUsecase{
		childRepo: childRepo,
	}
}

// RegisterChildInput input
type RegisterChildInput struct {
	DateOfBirth time.Time `validate:"required"`
	Gender      bool
	Name        string `validate:"required"`
}

// Validate validate RegisterChildInput
func (rci RegisterChildInput) Validate() error {
	return common.Validator.Struct(rci)
}

// RegisterChildOutput output
type RegisterChildOutput struct {
	ID uuid.UUID
}

// Register will register a child and assign the requester id as the parent ID
func (u *ChildUsecase) Register(ctx context.Context, input RegisterChildInput) (*RegisterChildOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	child, err := u.childRepo.Create(ctx, repository.CreateChildInput{
		ParentUserID: requester.ID,
		DateOfBirth:  input.DateOfBirth,
		Gender:       input.Gender,
		Name:         input.Name,
	})

	if err != nil {
		logger.WithError(err).Error("failed to insert child data to database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &RegisterChildOutput{
		ID: child.ID,
	}, nil
}

// UpdateChildInput input
type UpdateChildInput struct {
	ChildID     uuid.UUID `validate:"required"`
	DateOfBirth *time.Time
	Gender      *bool
	Name        *string
}

// Validate validate UpdateChildInput
func (uci UpdateChildInput) Validate() error {
	return common.Validator.Struct(uci)
}

// UpdateChildOutput output
type UpdateChildOutput struct {
	Message string
}

// Update update child data and can only be done by parents aka the child register
func (u *ChildUsecase) Update(ctx context.Context, input UpdateChildInput) (*UpdateChildOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	child, err := u.childRepo.FindByID(ctx, input.ChildID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find ")

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

	if child.ParentUserID != requester.ID {
		return nil, UsecaseError{
			ErrType: ErrForbidden,
			Message: "only the child's parent should be able to update child data",
		}
	}

	_, err = u.childRepo.Update(ctx, child.ID, repository.UpdateChildInput{
		DateOfBirth: input.DateOfBirth,
		Gender:      input.Gender,
		Name:        input.Name,
	})

	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &UpdateChildOutput{
		Message: "ok",
	}, nil
}

// GetRegisteredChildrenInput input
type GetRegisteredChildrenInput struct {
	Limit  int `validate:"min=1,max=100"`
	Offset int `validate:"min=0"`
}

func (grci GetRegisteredChildrenInput) validate() error {
	return common.Validator.Struct(grci)
}

// GetRegisteredChildrenOutput output
type GetRegisteredChildrenOutput struct {
	ID           uuid.UUID
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

// GetRegisteredChildren get registered children by the requester account
func (u *ChildUsecase) GetRegisteredChildren(ctx context.Context, input GetRegisteredChildrenInput) ([]GetRegisteredChildrenOutput, error) {
	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	children, err := u.childRepo.Search(ctx, repository.SearchChildInput{
		ParentUserID: &requester.ID,
		Limit:        input.Limit,
		Offset:       input.Offset,
	})

	switch err {
	default:
		logrus.WithContext(ctx).WithField("input", helper.Dump(input)).Error("failed to search registered children data")

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

	output := []GetRegisteredChildrenOutput{}

	for _, child := range children {
		output = append(output, GetRegisteredChildrenOutput{
			ID:           child.ID,
			ParentUserID: child.ParentUserID,
			DateOfBirth:  child.DateOfBirth,
			Gender:       child.Gender,
			Name:         child.Name,
			CreatedAt:    child.CreatedAt,
			UpdatedAt:    child.UpdatedAt,
			DeletedAt:    sql.NullTime(child.DeletedAt),
		})
	}

	return output, nil
}
