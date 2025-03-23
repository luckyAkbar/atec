package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	mockCommon "github.com/luckyAkbar/atec/mocks/internal_/common"
	mockUsecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginInput(t *testing.T) {
	t.Run("email is required and in valid format", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "just filling randomly",
			Password: "thisShouldBeVal1dPass!",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("empty email is unacceptable", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "",
			Password: "thisShouldBeVal1dPass!",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must not empty", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "valid@mail.test",
			Password: "",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must be atleast 8 chars", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "valid@mail.test",
			Password: "only7ch",
		}

		err := li.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		li := usecase.LoginInput{
			Email:    "valid@mail.test",
			Password: "only8char",
		}

		err := li.Validate()
		if err != nil {
			t.Errorf("expecting nil err but got %v", err)
		}
	})
}

func TestSignupInput(t *testing.T) {
	t.Run("email is required and in valid format", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "just filling randomly",
			Password: "thisShouldBeVal1dPass!",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("empty email is unacceptable", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "",
			Password: "thisShouldBeVal1dPass!",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must not empty", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("password must be atleast 8 chars", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "1234567",
			Username: "valid",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("username must not empty", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "only8char",
			Username: "",
		}

		err := si.Validate()
		if err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		si := usecase.SignupInput{
			Email:    "valid@mail.test",
			Password: "only8char",
			Username: "valid",
		}

		err := si.Validate()
		if err != nil {
			t.Errorf("expecting nil err but got %v", err)
		}
	})
}

func TestAccountVerificationInput(t *testing.T) {
	t.Run("token is required", func(t *testing.T) {
		vi := usecase.AccountVerificationInput{
			VerificationToken: "",
		}

		if err := vi.Validate(); err == nil {
			t.Error("expecting err but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		vi := usecase.AccountVerificationInput{
			VerificationToken: "ok",
		}

		if err := vi.Validate(); err != nil {
			t.Errorf("expecting nil but got %v", err)
		}
	})
}

func TestInitResetPasswordInput(t *testing.T) {
	t.Run("email is required", func(t *testing.T) {
		irpi := usecase.InitResetPasswordInput{
			Email: "",
		}

		if err := irpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("email must also in valid format", func(t *testing.T) {
		irpi := usecase.InitResetPasswordInput{
			Email: "invalid format",
		}

		if err := irpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		irpi := usecase.InitResetPasswordInput{
			Email: "valid@email.test",
		}

		if err := irpi.Validate(); err != nil {
			t.Errorf("expecting nil but got %v", err)
		}
	})
}

func TestResetPasswordInput(t *testing.T) {
	t.Run("reset password token is required", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "",
			NewPassword:        "validPass!",
		}

		if err := rpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("new password is required", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "validToken",
			NewPassword:        "",
		}

		if err := rpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("new password must be atleast 8 chars", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "validToken",
			NewPassword:        "7chars!",
		}

		if err := rpi.Validate(); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("ok", func(t *testing.T) {
		rpi := usecase.ResetPasswordInput{
			ResetPasswordToken: "validToken",
			NewPassword:        "validPass!",
		}

		if err := rpi.Validate(); err != nil {
			t.Errorf("expecting nil but got %v", err)
		}
	})
}

func TestAuthUsecase_HandleLogin(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := mockUsecase.NewUserRepository(t)
	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)
	uc := usecase.NewAuthUsecase(mockSharedCryptor, mockUserRepo, nil, nil, nil)
	validSampleEmail := "valid@sample.email"
	validSamplePassword := "validPass!!"
	sampleEncryptedEmail := "encryptedEmail"
	sampleJWTToken := "jwtToken"
	user := model.User{
		ID:       uuid.New(),
		IsActive: true,
	}

	testCases := []struct {
		name                 string
		input                usecase.LoginInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.LoginOutput
		expectedFunctionCall func()
	}{
		{
			name: "invalid input, email is empty",
			input: usecase.LoginInput{
				Email:    "",
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid email format",
			input: usecase.LoginInput{
				Email:    "this@isnotemail",
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "password must atleast be 8 chars",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: "missing",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "system failed to encrypt the email",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return("", assert.AnError).Once()
			},
		},
		{
			name: "system failed to query to database for user by email",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "no user was found with the supplied email",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "inactive user shouldn't be able to login",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()

				inactiveUser := user
				inactiveUser.IsActive = false

				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(&inactiveUser, nil).Once()
			},
		},
		{
			name: "inactive user shouldn't be able to login",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()

				inactiveUser := user
				inactiveUser.IsActive = false

				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(&inactiveUser, nil).Once()
			},
		},
		{
			name: "failure to compare hash means unauthorized login",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(&user, nil).Once()
				mockSharedCryptor.EXPECT().CompareHash(mock.Anything, []byte(validSamplePassword)).Return(assert.AnError).Once()
			},
		},
		{
			name: "failure to create JWT must return internal error",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(&user, nil).Once()
				mockSharedCryptor.EXPECT().CompareHash(mock.Anything, []byte(validSamplePassword)).Return(nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("", assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.LoginInput{
				Email:    validSampleEmail,
				Password: validSamplePassword,
			},
			wantErr:     false,
			expectedErr: usecase.ErrInternal,
			expectedOutput: &usecase.LoginOutput{
				Token: sampleJWTToken,
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validSampleEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(&user, nil).Once()
				mockSharedCryptor.EXPECT().CompareHash(mock.Anything, []byte(validSamplePassword)).Return(nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return(sampleJWTToken, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleLogin(ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, res)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}
