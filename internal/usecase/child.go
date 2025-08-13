package usecase

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
)

// ChildUsecase child usecase
type ChildUsecase struct {
	childRepo  ChildRepository
	userRepo   UserRepository
	resultRepo ResultRepository
}

// ChildUsecaseIface interface
type ChildUsecaseIface interface {
	Register(ctx context.Context, input RegisterChildInput) (*RegisterChildOutput, error)
	Update(ctx context.Context, input UpdateChildInput) (*UpdateChildOutput, error)
	GetRegisteredChildren(ctx context.Context, input GetRegisteredChildrenInput) ([]GetRegisteredChildrenOutput, error)
	Search(ctx context.Context, input SearchChildInput) ([]SearchChildOutput, error)
	HandleGetStatistic(ctx context.Context, input GetStatisticInput) (*GetStatisticOutput, error)
}

// NewChildUsecase create new ChildUsecase instance
func NewChildUsecase(childRepo ChildRepository, resultRepo ResultRepository, userRepo UserRepository) *ChildUsecase {
	return &ChildUsecase{
		childRepo:  childRepo,
		resultRepo: resultRepo,
		userRepo:   userRepo,
	}
}

// RegisterChildInput input
type RegisterChildInput struct {
	DateOfBirth  time.Time `validate:"required"`
	Gender       bool
	Name         string `validate:"required"`
	GuardianName *string
}

// Validate validate RegisterChildInput
func (rci RegisterChildInput) Validate() error {
	return common.Validator.Struct(rci)
}

// getGuardianName converts optional input GuardianName to sql.NullString safely
func (rci RegisterChildInput) getGuardianName() sql.NullString {
	if rci.GuardianName == nil {
		return sql.NullString{}
	}

	trimmed := strings.TrimSpace(*rci.GuardianName)
	if trimmed == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: trimmed, Valid: true}
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

	child, err := u.childRepo.Create(ctx, RepoCreateChildInput{
		ParentUserID: requester.ID,
		DateOfBirth:  input.DateOfBirth,
		Gender:       input.Gender,
		Name:         input.Name,
		GuardianName: input.getGuardianName(),
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
	ChildID      uuid.UUID `validate:"required"`
	DateOfBirth  *time.Time
	Gender       *bool
	Name         *string
	GuardianName *sql.NullString
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
	case ErrRepoNotFound:
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

	_, err = u.childRepo.Update(ctx, child.ID, RepoUpdateChildInput{
		DateOfBirth:  input.DateOfBirth,
		Gender:       input.Gender,
		Name:         input.Name,
		GuardianName: input.GuardianName,
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
	ID             uuid.UUID
	ParentUserID   uuid.UUID
	ParentUsername string
	DateOfBirth    time.Time
	Gender         bool
	Name           string
	GuardianName   sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
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

	children, err := u.childRepo.Search(ctx, RepoSearchChildInput{
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
	case ErrRepoNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	parent, err := u.userRepo.FindByID(ctx, requester.ID)
	switch err {
	default:
		logrus.WithContext(ctx).WithField("input", helper.Dump(input)).Error("failed to find parent data")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case ErrRepoNotFound:
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
			ID:             child.ID,
			ParentUserID:   child.ParentUserID,
			ParentUsername: parent.Username,
			DateOfBirth:    child.DateOfBirth,
			Gender:         child.Gender,
			Name:           child.Name,
			GuardianName:   child.GuardianName,
			CreatedAt:      child.CreatedAt,
			UpdatedAt:      child.UpdatedAt,
			DeletedAt:      sql.NullTime(child.DeletedAt),
		})
	}

	return output, nil
}

// SearchChildInput input
type SearchChildInput struct {
	ParentUserID *uuid.UUID
	Name         *string
	Gender       *bool
	Limit        int `validate:"required,min=1,max=100"`
	Offset       int `validate:"min=0"`
}

func (sci SearchChildInput) validate() error {
	return common.Validator.Struct(sci)
}

// SearchChildOutput output
type SearchChildOutput struct {
	ID           uuid.UUID
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	Name         string
	GuardianName sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

// Search allow requester to full search registered child data based on multiple search params
func (u *ChildUsecase) Search(ctx context.Context, input SearchChildInput) ([]SearchChildOutput, error) {
	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	if requester.Role != model.RolesTherapist {
		return nil, UsecaseError{
			ErrType: ErrForbidden,
			Message: "insufficient permission to access this feature",
		}
	}

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	children, err := u.childRepo.Search(ctx, RepoSearchChildInput{
		ParentUserID: input.ParentUserID,
		Name:         input.Name,
		Gender:       input.Gender,
		Limit:        input.Limit,
		Offset:       input.Offset,
	})

	switch err {
	default:
		logrus.WithContext(ctx).WithField("input", helper.Dump(input)).Error("failed to search children data")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case ErrRepoNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	output := []SearchChildOutput{}

	for _, child := range children {
		output = append(output, SearchChildOutput{
			ID:           child.ID,
			ParentUserID: child.ParentUserID,
			DateOfBirth:  child.DateOfBirth,
			Gender:       child.Gender,
			Name:         child.Name,
			GuardianName: child.GuardianName,
			CreatedAt:    child.CreatedAt,
			UpdatedAt:    child.UpdatedAt,
			DeletedAt:    sql.NullTime(child.DeletedAt),
		})
	}

	return output, nil
}

// GetStatisticInput input
type GetStatisticInput struct {
	ChildID uuid.UUID `validate:"required"`
}

func (gsi GetStatisticInput) validate() error {
	return common.Validator.Struct(gsi)
}

// StatisticComponent is the single component of statistic. composed of at least
// the total score for a given time of test.
type StatisticComponent struct {
	Total     int                `json:"total"`
	CreatedAt time.Time          `json:"created_at"`
	Detail    model.ResultDetail `json:"detail"`
}

// GetStatisticOutput represent the overall data to build the statistic
type GetStatisticOutput struct {
	Statistic []StatisticComponent `json:"statistic"`
}

// HandleGetStatistic get the statistic of a given child id. It requires the valid
// authorization of the parent or admin role. It will return the overall statistic
// of the child, which is composed of time of test and the total score of the test.
func (u *ChildUsecase) HandleGetStatistic(ctx context.Context, input GetStatisticInput) (*GetStatisticOutput, error) {
	logger := logrus.WithContext(ctx).WithField("input", helper.Dump(input))

	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "getting statistic requires valid authorization",
		}
	}

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	child, err := u.childRepo.FindByID(ctx, input.ChildID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find child data from database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case ErrRepoNotFound:
		return nil, UsecaseError{
			ErrType: ErrNotFound,
			Message: ErrNotFound.Error(),
		}
	case nil:
		break
	}

	if child.ParentUserID != requester.ID && requester.Role != model.RolesTherapist {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "getting statistic for this child must be done either by parent or therapist",
		}
	}

	batchSize := 100
	offset := 0
	isFirst := true
	mustContinue := true
	statComponents := []StatisticComponent{}

	for {
		results, err := u.resultRepo.Search(ctx, RepoSearchResultInput{
			ChildID: input.ChildID,
			Limit:   batchSize,
			Offset:  offset,
		})

		switch err {
		default:
			logger.WithError(err).Error("failed to get statistic from database")

			return nil, UsecaseError{
				ErrType: ErrInternal,
				Message: ErrInternal.Error(),
			}
		case ErrRepoNotFound:
			if isFirst {
				return nil, UsecaseError{
					ErrType: ErrNotFound,
					Message: ErrNotFound.Error(),
				}
			}

			mustContinue = false
		case nil:
		}

		for _, res := range results {
			total := 0

			for _, detail := range res.Result {
				total += detail.Grade
			}

			statComponents = append(statComponents, StatisticComponent{
				Total:     total,
				CreatedAt: res.CreatedAt,
				Detail:    res.Result,
			})
		}

		offset += batchSize
		isFirst = false

		// avoid extra one query if the result is less than batch size
		if len(results) < batchSize {
			mustContinue = false
		}

		if !mustContinue {
			break
		}
	}

	return &GetStatisticOutput{
		Statistic: statComponents,
	}, nil
}
