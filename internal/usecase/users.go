package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
)

// UsersUsecase contains business logic related to user entity
type UsersUsecase struct {
	userRepo UserRepository
}

// UsersUsecaseIface exported interface for UsersUsecase
type UsersUsecaseIface interface {
	GetMyProfile(ctx context.Context) (*GetMyProfileOutput, error)
	GetTherapistData(ctx context.Context) ([]GetTherapistDataOutput, error)
}

// NewUsersUsecase create new UsersUsecase instance
func NewUsersUsecase(userRepo UserRepository) *UsersUsecase {
	return &UsersUsecase{userRepo: userRepo}
}

// GetMyProfileOutput output
type GetMyProfileOutput struct {
	ID        uuid.UUID
	Username  string
	IsActive  bool
	Roles     model.Roles
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetMyProfile returns currently authenticated user's profile from database
func (u *UsersUsecase) GetMyProfile(ctx context.Context) (*GetMyProfileOutput, error) {
	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	user, err := u.userRepo.FindByID(ctx, requester.ID)
	if err != nil {
		switch err {
		default:
			return nil, UsecaseError{
				ErrType: ErrInternal,
				Message: ErrInternal.Error(),
			}
		case ErrRepoNotFound:
			return nil, UsecaseError{
				ErrType: ErrNotFound,
				Message: ErrNotFound.Error(),
			}
		}
	}

	return &GetMyProfileOutput{
		ID:        user.ID,
		Username:  user.Username,
		IsActive:  user.IsActive,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetTherapistDataOutput output
type GetTherapistDataOutput struct {
	ID        uuid.UUID
	Username  string
	IsActive  bool
	Roles     model.Roles
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetTherapistData returns all users that have therapist role
func (u *UsersUsecase) GetTherapistData(ctx context.Context) ([]GetTherapistDataOutput, error) {
	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	users, err := u.userRepo.GetUsersByRoles(ctx, model.RolesTherapist)
	if err != nil {
		switch err {
		default:
			return nil, UsecaseError{
				ErrType: ErrInternal,
				Message: ErrInternal.Error(),
			}
		case ErrRepoNotFound:
			return nil, UsecaseError{
				ErrType: ErrNotFound,
				Message: ErrNotFound.Error(),
			}
		}
	}

	output := make([]GetTherapistDataOutput, 0, len(users))
	for _, user := range users {
		output = append(output, GetTherapistDataOutput{
			ID:        user.ID,
			Username:  user.Username,
			IsActive:  user.IsActive,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return output, nil
}
