package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
)

// UsersUsecase contains business logic related to user entity
type UsersUsecase struct {
	userRepo      UserRepository
	sharedCryptor common.SharedCryptorIface
}

// UsersUsecaseIface exported interface for UsersUsecase
type UsersUsecaseIface interface {
	GetMyProfile(ctx context.Context) (*GetMyProfileOutput, error)
	GetTherapistData(ctx context.Context) ([]GetTherapistDataOutput, error)
	UpdateMyProfile(ctx context.Context, input UpdateMyProfileInput) (*UpdateMyProfileOutput, error)
}

// NewUsersUsecase create new UsersUsecase instance
func NewUsersUsecase(userRepo UserRepository, sharedCryptor common.SharedCryptorIface) *UsersUsecase {
	return &UsersUsecase{userRepo: userRepo, sharedCryptor: sharedCryptor}
}

// decryptUserData decrypts sensitive fields on user and returns plain values.
func (u *UsersUsecase) decryptUserData(user *model.User) (string, *string, *string, error) {
	decryptedEmail := ""

	if user.Email != "" {
		de, err := u.sharedCryptor.Decrypt(user.Email)
		if err != nil {
			return "", nil, nil, err
		}

		decryptedEmail = de
	}

	var phonePtr *string

	if user.PhoneNumber.Valid {
		p, err := u.sharedCryptor.Decrypt(user.PhoneNumber.String)
		if err != nil {
			return "", nil, nil, err
		}

		phonePtr = &p
	}

	var addressPtr *string

	if user.Address.Valid {
		a, err := u.sharedCryptor.Decrypt(user.Address.String)
		if err != nil {
			return "", nil, nil, err
		}

		addressPtr = &a
	}

	return decryptedEmail, phonePtr, addressPtr, nil
}

// GetMyProfileOutput output
type GetMyProfileOutput struct {
	ID          uuid.UUID
	Username    string
	IsActive    bool
	Roles       model.Roles
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	PhoneNumber *string
	Address     *string
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

	decryptedEmail, phonePtr, addressPtr, decErr := u.decryptUserData(user)
	if decErr != nil {
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &GetMyProfileOutput{
		ID:          user.ID,
		Username:    user.Username,
		IsActive:    user.IsActive,
		Roles:       user.Roles,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       decryptedEmail,
		PhoneNumber: phonePtr,
		Address:     addressPtr,
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

// UpdateMyProfileInput input
type UpdateMyProfileInput struct {
	Username    string  `validate:"required"`
	PhoneNumber *string `validate:"required,e164"`
	Address     *string `validate:"required,max=256"`
}

func (i *UpdateMyProfileInput) validate() error {
	if i.PhoneNumber != nil {
		trimmed := strings.ReplaceAll(*i.PhoneNumber, " ", "")
		i.PhoneNumber = &trimmed
	}

	if i.Address != nil {
		addr := strings.TrimSpace(*i.Address)
		i.Address = &addr
	}

	return common.Validator.Struct(i)
}

// UpdateMyProfileOutput output
type UpdateMyProfileOutput struct {
	Message string
}

// UpdateMyProfile allows user to update their own profile
func (u *UsersUsecase) UpdateMyProfile(ctx context.Context, input UpdateMyProfileInput) (*UpdateMyProfileOutput, error) {
	requester := model.GetUserFromCtx(ctx)
	if requester == nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: ErrUnauthorized.Error(),
		}
	}

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	// convert to sql nulls then encrypt values
	encPhone := sqlNullFromPtr(input.PhoneNumber)
	encAddr := sqlNullFromPtr(input.Address)

	if encPhone.Valid {
		phoneEncrypted, err := u.sharedCryptor.Encrypt(encPhone.String)
		if err != nil {
			return nil, UsecaseError{
				ErrType: ErrInternal,
				Message: ErrInternal.Error(),
			}
		}

		encPhone.String = phoneEncrypted
	}

	if encAddr.Valid {
		addressEncrypted, err := u.sharedCryptor.Encrypt(encAddr.String)
		if err != nil {
			return nil, UsecaseError{
				ErrType: ErrInternal,
				Message: ErrInternal.Error(),
			}
		}

		encAddr.String = addressEncrypted
	}

	_, err := u.userRepo.UpdateProfile(ctx, requester.ID, RepoUpdateUserProfileInput{
		Username:    input.Username,
		PhoneNumber: encPhone,
		Address:     encAddr,
	})

	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &UpdateMyProfileOutput{
		Message: "ok",
	}, nil
}
