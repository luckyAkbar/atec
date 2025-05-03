package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	mockCommon "github.com/luckyAkbar/atec/mocks/internal_/common"
	mockUsecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/sendinblue/APIv3-go-library/v2/lib"
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

func TestAuthUsecase_HandleSignup(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mockUsecase.NewUserRepository(t)
	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)
	mockTxCtrlFactory := mockUsecase.NewTransactionControllerFactory(t)
	mockMailer := mockCommon.NewMailerIface(t)
	mockRateLimiter := mockUsecase.NewRateLimiter(t)

	sampleValidEmail := "valid@email.sample"
	sampleValidPassword := "validPass!!"
	sampleValidUsername := "validUsername"
	sampleEncryptedEmail := "encryptedEmail"
	sampleHashedPassword := "hashedPassword"

	repoCreateUserInput := usecase.RepoCreateUserInput{
		Email:    sampleEncryptedEmail,
		Password: sampleHashedPassword,
		Username: sampleValidUsername,
		IsActive: false,
		Roles:    model.RolesTherapist,
	}

	uc := usecase.NewAuthUsecase(mockSharedCryptor, mockUserRepo, mockTxCtrlFactory, mockMailer, mockRateLimiter)

	testCases := []struct {
		name                 string
		input                usecase.SignupInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.SignupOutput
		expectedFunctionCall func()
	}{
		{
			name: "invalid email format",
			input: usecase.SignupInput{
				Email:    "invalidEmail",
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "email is required",
			input: usecase.SignupInput{
				Email:    "",
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "password must be atleast 8 chars",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: "7chars!",
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "username is required",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: "",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "system failed to encrypt email",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return("", assert.AnError).Once()
			},
		},
		{
			name: "repository failure to fetch user by email",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "user with the same email already exists",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(&model.User{}, nil).Once()
			},
		},
		{
			name: "system failed to hash password",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return("", assert.AnError).Once()
			},
		},
		{
			name: "system failed to create user",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return(sampleHashedPassword, nil).Once()

				underlyingTransaction := mockUsecase.NewTransactionController(t)
				txCtrlWrapper := usecase.NewTxControllerWrapper(underlyingTransaction)

				mockTxCtrlFactory.EXPECT().New().Return(txCtrlWrapper).Once()
				underlyingTransaction.EXPECT().Begin().Return(struct{}{}).Once()
				mockUserRepo.EXPECT().Create(ctx, repoCreateUserInput, mock.Anything).Return(nil, assert.AnError).Once()
				underlyingTransaction.EXPECT().Rollback().Return(nil).Once()
			},
		},
		{
			name: "rollback failure must still return internal error",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return(sampleHashedPassword, nil).Once()

				underlyingTransaction := mockUsecase.NewTransactionController(t)
				txCtrlWrapper := usecase.NewTxControllerWrapper(underlyingTransaction)

				mockTxCtrlFactory.EXPECT().New().Return(txCtrlWrapper).Once()
				underlyingTransaction.EXPECT().Begin().Return(struct{}{}).Once()
				mockUserRepo.EXPECT().Create(ctx, repoCreateUserInput, mock.Anything).Return(nil, assert.AnError).Once()
				underlyingTransaction.EXPECT().Rollback().Return(assert.AnError).Once()
			},
		},
		{
			name: "system failed to create JWT token for signup",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return(sampleHashedPassword, nil).Once()

				underlyingTransaction := mockUsecase.NewTransactionController(t)
				txCtrlWrapper := usecase.NewTxControllerWrapper(underlyingTransaction)

				mockTxCtrlFactory.EXPECT().New().Return(txCtrlWrapper).Once()
				underlyingTransaction.EXPECT().Begin().Return(struct{}{}).Once()

				mockUserRepo.EXPECT().Create(ctx, repoCreateUserInput, mock.Anything).Return(&model.User{}, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("", assert.AnError).Once()

				underlyingTransaction.EXPECT().Rollback().Return(nil).Once()
			},
		},
		{
			name: "failure to send email",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return(sampleHashedPassword, nil).Once()

				underlyingTransaction := mockUsecase.NewTransactionController(t)
				txCtrlWrapper := usecase.NewTxControllerWrapper(underlyingTransaction)

				mockTxCtrlFactory.EXPECT().New().Return(txCtrlWrapper).Once()
				underlyingTransaction.EXPECT().Begin().Return(struct{}{}).Once()

				mockUserRepo.EXPECT().Create(ctx, repoCreateUserInput, mock.Anything).Return(&model.User{}, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("jwt", nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(&lib.CreateSmtpEmail{}, assert.AnError).Once()

				underlyingTransaction.EXPECT().Rollback().Return(nil).Once()
			},
		},
		{
			name: "failure on commit must trigger internal server error",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return(sampleHashedPassword, nil).Once()

				underlyingTransaction := mockUsecase.NewTransactionController(t)
				txCtrlWrapper := usecase.NewTxControllerWrapper(underlyingTransaction)

				mockTxCtrlFactory.EXPECT().New().Return(txCtrlWrapper).Once()
				underlyingTransaction.EXPECT().Begin().Return(struct{}{}).Once()

				mockUserRepo.EXPECT().Create(ctx, repoCreateUserInput, mock.Anything).Return(&model.User{}, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("jwt", nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(&lib.CreateSmtpEmail{}, nil).Once()

				underlyingTransaction.EXPECT().Commit().Return(assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.SignupInput{
				Email:    sampleValidEmail,
				Password: sampleValidPassword,
				Username: sampleValidUsername,
			},
			wantErr: false,
			expectedOutput: &usecase.SignupOutput{
				Message: "email confirmation sent",
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(sampleValidEmail).Return(sampleEncryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, sampleEncryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPassword)).Return(sampleHashedPassword, nil).Once()

				underlyingTransaction := mockUsecase.NewTransactionController(t)
				txCtrlWrapper := usecase.NewTxControllerWrapper(underlyingTransaction)

				mockTxCtrlFactory.EXPECT().New().Return(txCtrlWrapper).Once()
				underlyingTransaction.EXPECT().Begin().Return(struct{}{}).Once()

				mockUserRepo.EXPECT().Create(ctx, repoCreateUserInput, mock.Anything).Return(&model.User{}, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("jwt", nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(&lib.CreateSmtpEmail{}, nil).Once()

				underlyingTransaction.EXPECT().Commit().Return(nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleSignup(ctx, tc.input)

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

func TestAuthUsecase_HandleAccountVerification(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mockUsecase.NewUserRepository(t)
	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)

	uc := usecase.NewAuthUsecase(mockSharedCryptor, mockUserRepo, nil, nil, nil)

	validateJWTOpts := common.ValidateJWTOpts{
		Issuer:  string(usecase.TokenIssuerSystem),
		Subject: string(usecase.SignupVerificationToken),
	}

	sampleVerificationToken := "sampletoken"
	invalidJWTToken := &jwt.Token{
		Valid: false,
	}

	signingKey := []byte("key")
	invalidClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	invalidClaimsToken.Valid = true
	invalidClaimsTokenString, _ := invalidClaimsToken.SignedString(signingKey)

	unknownIssuerToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "unknown",
	})
	unknownIssuerTokenString, _ := unknownIssuerToken.SignedString(signingKey)

	unknownTokenSubject := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": "unknown",
	})
	unknownTokenSubject.Valid = true
	unknownTokenSubjectString, _ := unknownTokenSubject.SignedString(signingKey)

	noAudienceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.SignupVerificationToken),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	noAudienceToken.Valid = true
	noAudienceTokenString, _ := noAudienceToken.SignedString(signingKey)

	moreThanOneAudienceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.SignupVerificationToken),
		"aud": []string{"aud1", "aud2"},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	moreThanOneAudienceToken.Valid = true
	moreThanOneAudienceTokenString, _ := moreThanOneAudienceToken.SignedString(signingKey)

	invalidUUIDAudToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.SignupVerificationToken),
		"aud": []string{"invalid-UUID"},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	invalidUUIDAudToken.Valid = true
	invalidUUIDAudTokenString, _ := invalidUUIDAudToken.SignedString(signingKey)

	userID := uuid.New()
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.SignupVerificationToken),
		"aud": []string{userID.String()},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	validToken.Valid = true
	validTokenString, _ := validToken.SignedString(signingKey)

	testCases := []struct {
		name                 string
		input                usecase.AccountVerificationInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.AccountVerificationOutput
		expectedFunctionCall func()
	}{
		{
			name: "missing verification token",
			input: usecase.AccountVerificationInput{
				VerificationToken: "",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid jwt token must trigger unauthorized",
			input: usecase.AccountVerificationInput{
				VerificationToken: "invalidToken",
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT("invalidToken", validateJWTOpts).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "expired jwt token",
			input: usecase.AccountVerificationInput{
				VerificationToken: sampleVerificationToken,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleVerificationToken, validateJWTOpts).Return(nil, jwt.ErrTokenExpired).Once()
			},
		},
		{
			name: "jwt token is not Valid",
			input: usecase.AccountVerificationInput{
				VerificationToken: sampleVerificationToken,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleVerificationToken, validateJWTOpts).Return(invalidJWTToken, nil).Once()
			},
		},
		{
			name: "invalid claims",
			input: usecase.AccountVerificationInput{
				VerificationToken: invalidClaimsTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidClaimsTokenString, validateJWTOpts).Return(invalidClaimsToken, nil).Once()
			},
		},
		{
			name: "unknown issuer",
			input: usecase.AccountVerificationInput{
				VerificationToken: unknownIssuerTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(unknownIssuerTokenString, validateJWTOpts).Return(unknownIssuerToken, nil).Once()
			},
		},
		{
			name: "unknown subject",
			input: usecase.AccountVerificationInput{
				VerificationToken: unknownTokenSubjectString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(unknownTokenSubjectString, validateJWTOpts).Return(unknownTokenSubject, nil).Once()
			},
		},
		{
			name: "no subject token",
			input: usecase.AccountVerificationInput{
				VerificationToken: noAudienceTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(noAudienceTokenString, validateJWTOpts).Return(noAudienceToken, nil).Once()
			},
		},
		{
			name: "more than one aud",
			input: usecase.AccountVerificationInput{
				VerificationToken: moreThanOneAudienceTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(moreThanOneAudienceTokenString, validateJWTOpts).Return(moreThanOneAudienceToken, nil).Once()
			},
		},
		{
			name: "invalid uuid token",
			input: usecase.AccountVerificationInput{
				VerificationToken: invalidUUIDAudTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidUUIDAudTokenString, validateJWTOpts).Return(invalidUUIDAudToken, nil).Once()
			},
		},
		{
			name: "failed to query user from repository",
			input: usecase.AccountVerificationInput{
				VerificationToken: validTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "user wasn't found from repository",
			input: usecase.AccountVerificationInput{
				VerificationToken: validTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "user already active",
			input: usecase.AccountVerificationInput{
				VerificationToken: validTokenString,
			},
			wantErr: false,
			expectedOutput: &usecase.AccountVerificationOutput{
				Message: "your account has been activated",
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(&model.User{IsActive: true}, nil).Once()
			},
		},
		{
			name: "repository failed to activate user",
			input: usecase.AccountVerificationInput{
				VerificationToken: validTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				truth := true
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(&model.User{ID: userID}, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, userID, usecase.RepoUpdateUserInput{
					IsActive: &truth,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.AccountVerificationInput{
				VerificationToken: validTokenString,
			},
			wantErr: false,
			expectedOutput: &usecase.AccountVerificationOutput{
				Message: "your account has been activated",
			},
			expectedFunctionCall: func() {
				truth := true
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(&model.User{ID: userID}, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, userID, usecase.RepoUpdateUserInput{
					IsActive: &truth,
				}).Return(&model.User{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleAccountVerification(ctx, tc.input)

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

func TestAuthUsecase_HandleInitesetPassword(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mockUsecase.NewUserRepository(t)
	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)
	mockMailer := mockCommon.NewMailerIface(t)

	uc := usecase.NewAuthUsecase(mockSharedCryptor, mockUserRepo, nil, mockMailer, nil)

	validEmail := "valid@email.sample"
	encryptedEmail := "sampleEncryptedEmail"
	sampleJWTToken := "token"
	userID := uuid.New()
	user := &model.User{
		ID:       userID,
		IsActive: true,
	}

	testCases := []struct {
		name                 string
		input                usecase.InitResetPasswordInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.InitResetPasswordOutput
		expectedFunctionCall func()
	}{
		{
			name: "email is required",
			input: usecase.InitResetPasswordInput{
				Email: "",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid email",
			input: usecase.InitResetPasswordInput{
				Email: "invalid@email",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "system failed to encrypt email",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return("", assert.AnError).Once()
			},
		},
		{
			name: "repositori failed to find user by email",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return(encryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedEmail).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "no user found with the supplied email",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return(encryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedEmail).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "inactive user should not be able to use this",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return(encryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedEmail).Return(&model.User{IsActive: false}, nil).Once()
			},
		},
		{
			name: "failed to generate jwt token",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return(encryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedEmail).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("", assert.AnError).Once()
			},
		},
		{
			name: "failed to send email",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return(encryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedEmail).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return(sampleJWTToken, nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.InitResetPasswordInput{
				Email: validEmail,
			},
			wantErr: false,
			expectedOutput: &usecase.InitResetPasswordOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(validEmail).Return(encryptedEmail, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedEmail).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return(sampleJWTToken, nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(&lib.CreateSmtpEmail{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleInitesetPassword(ctx, tc.input)

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

func TestAuthUsecase_HandleResetPassword(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mockUsecase.NewUserRepository(t)
	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)
	mockMailer := mockCommon.NewMailerIface(t)

	uc := usecase.NewAuthUsecase(mockSharedCryptor, mockUserRepo, nil, mockMailer, nil)

	sampleValidPass := "123ValidPass"
	hashedPw := "hashedpw"
	validateJWTOpts := common.ValidateJWTOpts{
		Issuer:  string(usecase.TokenIssuerSystem),
		Subject: string(usecase.ChangePasswordToken),
	}

	sampleToken := "token"
	invalidJWTToken := &jwt.Token{
		Valid: false,
	}

	signingKey := []byte("key")
	invalidClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	invalidClaimsToken.Valid = true
	invalidClaimsTokenString, _ := invalidClaimsToken.SignedString(signingKey)

	unknownIssuerToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "unknown",
	})
	unknownIssuerTokenString, _ := unknownIssuerToken.SignedString(signingKey)

	unknownTokenSubject := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": "unknown",
	})
	unknownTokenSubject.Valid = true
	unknownTokenSubjectString, _ := unknownTokenSubject.SignedString(signingKey)

	noAudienceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.ChangePasswordToken),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	noAudienceToken.Valid = true
	noAudienceTokenString, _ := noAudienceToken.SignedString(signingKey)

	moreThanOneAudienceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.ChangePasswordToken),
		"aud": []string{"aud1", "aud2"},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	moreThanOneAudienceToken.Valid = true
	moreThanOneAudienceTokenString, _ := moreThanOneAudienceToken.SignedString(signingKey)

	invalidUUIDAudToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.ChangePasswordToken),
		"aud": []string{"invalid-UUID"},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	invalidUUIDAudToken.Valid = true
	invalidUUIDAudTokenString, _ := invalidUUIDAudToken.SignedString(signingKey)

	userID := uuid.New()
	user := &model.User{ID: userID}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.ChangePasswordToken),
		"aud": []string{userID.String()},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	validToken.Valid = true
	validTokenString, _ := validToken.SignedString(signingKey)

	testCases := []struct {
		name                 string
		input                usecase.ResetPasswordInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.ResetPasswordOutput
		expectedFunctionCall func()
	}{
		{
			name: "token is required",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: "",
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "password is required",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: sampleToken,
				NewPassword:        "",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "password minimum 8",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: sampleToken,
				NewPassword:        "7chars1",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "validate jwt token return unexpected error must return unauthorized",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: sampleToken,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleToken, validateJWTOpts).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "validate jwt token return unexpected error must return unauthorized",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: sampleToken,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleToken, validateJWTOpts).Return(nil, jwt.ErrTokenExpired).Once()
			},
		},
		{
			name: "jwt token is not Valid",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: sampleToken,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleToken, validateJWTOpts).Return(invalidJWTToken, nil).Once()
			},
		},
		{
			name: "invalid claims",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: invalidClaimsTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidClaimsTokenString, validateJWTOpts).Return(invalidClaimsToken, nil).Once()
			},
		},
		{
			name: "unknown issuer",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: unknownIssuerTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(unknownIssuerTokenString, validateJWTOpts).Return(unknownIssuerToken, nil).Once()
			},
		},
		{
			name: "unknown subject",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: unknownTokenSubjectString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(unknownTokenSubjectString, validateJWTOpts).Return(unknownTokenSubject, nil).Once()
			},
		},
		{
			name: "no subject token",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: noAudienceTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(noAudienceTokenString, validateJWTOpts).Return(noAudienceToken, nil).Once()
			},
		},
		{
			name: "more than one aud",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: moreThanOneAudienceTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(moreThanOneAudienceTokenString, validateJWTOpts).Return(moreThanOneAudienceToken, nil).Once()
			},
		},
		{
			name: "invalid uuid token",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: invalidUUIDAudTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidUUIDAudTokenString, validateJWTOpts).Return(invalidUUIDAudToken, nil).Once()
			},
		},
		{
			name: "failure to fetch user from repository",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: validTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "user not found",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: validTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "failed to hash new password",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: validTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPass)).Return("", assert.AnError).Once()
			},
		},
		{
			name: "failed to update user data on repository",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: validTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPass)).Return(hashedPw, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, userID, usecase.RepoUpdateUserInput{
					Password: hashedPw,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.ResetPasswordInput{
				ResetPasswordToken: validTokenString,
				NewPassword:        sampleValidPass,
			},
			wantErr: false,
			expectedOutput: &usecase.ResetPasswordOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
				mockUserRepo.EXPECT().FindByID(ctx, userID).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().Hash([]byte(sampleValidPass)).Return(hashedPw, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, userID, usecase.RepoUpdateUserInput{
					Password: hashedPw,
				}).Return(user, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleResetPassword(ctx, tc.input)

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

func TestAuthUsecase_AllowAccess(t *testing.T) {
	ctx := context.Background()

	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)

	uc := usecase.NewAuthUsecase(mockSharedCryptor, nil, nil, nil, nil)

	sampleVerificationToken := "sampleToken"
	validateJWTOpts := common.ValidateJWTOpts{
		Issuer:  string(usecase.TokenIssuerSystem),
		Subject: string(usecase.LoginToken),
	}

	invalidJWTToken := &jwt.Token{
		Valid: false,
	}

	signingKey := []byte("key")
	invalidClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	invalidClaimsToken.Valid = true
	invalidClaimsTokenString, _ := invalidClaimsToken.SignedString(signingKey)

	unknownIssuerToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "unknown",
	})
	unknownIssuerTokenString, _ := unknownIssuerToken.SignedString(signingKey)

	unknownTokenSubject := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": "unknown",
	})
	unknownTokenSubject.Valid = true
	unknownTokenSubjectString, _ := unknownTokenSubject.SignedString(signingKey)

	noAudienceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.LoginToken),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	noAudienceToken.Valid = true
	noAudienceTokenString, _ := noAudienceToken.SignedString(signingKey)

	moreThanOneAudienceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.LoginToken),
		"aud": []string{"aud1", "aud2"},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	moreThanOneAudienceToken.Valid = true
	moreThanOneAudienceTokenString, _ := moreThanOneAudienceToken.SignedString(signingKey)

	invalidUUIDAudToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": string(usecase.TokenIssuerSystem),
		"sub": string(usecase.LoginToken),
		"aud": []string{"invalid-UUID"},
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
	})
	invalidUUIDAudToken.Valid = true
	invalidUUIDAudTokenString, _ := invalidUUIDAudToken.SignedString(signingKey)

	userID := uuid.New()
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  string(usecase.TokenIssuerSystem),
		"sub":  string(usecase.LoginToken),
		"aud":  []string{userID.String()},
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
		"role": string(model.RoleUser),
	})
	validToken.Valid = true
	validTokenString, _ := validToken.SignedString(signingKey)

	invalidRoleToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  string(usecase.TokenIssuerSystem),
		"sub":  string(usecase.LoginToken),
		"aud":  []string{userID.String()},
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
		"role": "invalid",
	})
	invalidRoleToken.Valid = true
	invalidRoleTokenString, _ := validToken.SignedString(signingKey)

	invalidRoleStringToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  string(usecase.TokenIssuerSystem),
		"sub":  string(usecase.LoginToken),
		"aud":  []string{userID.String()},
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
		"role": "123",
	})
	invalidRoleStringToken.Valid = true
	invalidRoleStringTokenString, _ := validToken.SignedString(signingKey)

	adminID := uuid.New()
	validAdminToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  string(usecase.TokenIssuerSystem),
		"sub":  string(usecase.LoginToken),
		"aud":  []string{adminID.String()},
		"exp":  jwt.NewNumericDate(time.Now().Add(time.Hour * 1)).Unix(),
		"role": string(model.RolesAdmin),
	})
	validAdminToken.Valid = true
	validAdminTokenString, _ := validToken.SignedString(signingKey)

	testCases := []struct {
		name                 string
		input                usecase.AllowAccessInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.AllowAccessOutput
		expectedFunctionCall func()
	}{
		{
			name: "expired jwt token",
			input: usecase.AllowAccessInput{
				Token: sampleVerificationToken,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleVerificationToken, validateJWTOpts).Return(nil, jwt.ErrTokenExpired).Once()
			},
		},
		{
			name: "jwt token is not Valid",
			input: usecase.AllowAccessInput{
				Token: sampleVerificationToken,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(sampleVerificationToken, validateJWTOpts).Return(invalidJWTToken, nil).Once()
			},
		},
		{
			name: "invalid claims",
			input: usecase.AllowAccessInput{
				Token: invalidClaimsTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidClaimsTokenString, validateJWTOpts).Return(invalidClaimsToken, nil).Once()
			},
		},
		{
			name: "unknown issuer",
			input: usecase.AllowAccessInput{
				Token: unknownIssuerTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(unknownIssuerTokenString, validateJWTOpts).Return(unknownIssuerToken, nil).Once()
			},
		},
		{
			name: "unknown subject",
			input: usecase.AllowAccessInput{
				Token: unknownTokenSubjectString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(unknownTokenSubjectString, validateJWTOpts).Return(unknownTokenSubject, nil).Once()
			},
		},
		{
			name: "no subject token",
			input: usecase.AllowAccessInput{
				Token: noAudienceTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(noAudienceTokenString, validateJWTOpts).Return(noAudienceToken, nil).Once()
			},
		},
		{
			name: "more than one aud",
			input: usecase.AllowAccessInput{
				Token: moreThanOneAudienceTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(moreThanOneAudienceTokenString, validateJWTOpts).Return(moreThanOneAudienceToken, nil).Once()
			},
		},
		{
			name: "invalid uuid token",
			input: usecase.AllowAccessInput{
				Token: invalidUUIDAudTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidUUIDAudTokenString, validateJWTOpts).Return(invalidUUIDAudToken, nil).Once()
			},
		},
		{
			name: "undefined role in the token claims",
			input: usecase.AllowAccessInput{
				Token: invalidUUIDAudTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidUUIDAudTokenString, validateJWTOpts).Return(invalidUUIDAudToken, nil).Once()
			},
		},
		{
			name: "role was not string",
			input: usecase.AllowAccessInput{
				Token: invalidRoleStringTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidRoleStringTokenString, validateJWTOpts).Return(invalidRoleStringToken, nil).Once()
			},
		},
		{
			name: "role was not string",
			input: usecase.AllowAccessInput{
				Token: invalidRoleTokenString,
			},
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(invalidRoleTokenString, validateJWTOpts).Return(invalidRoleToken, nil).Once()
			},
		},
		{
			name: "allowing all unauthorized success for role user",
			input: usecase.AllowAccessInput{
				Token:              validTokenString,
				AllowAllAuthorized: true,
			},
			wantErr: false,
			expectedOutput: &usecase.AllowAccessOutput{
				UserRole: model.RoleUser,
				UserID:   userID,
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
			},
		},
		{
			name: "allowing all unauthorized success for role admin",
			input: usecase.AllowAccessInput{
				Token:              validAdminTokenString,
				AllowAllAuthorized: true,
			},
			wantErr: false,
			expectedOutput: &usecase.AllowAccessOutput{
				UserID:   adminID,
				UserRole: model.RolesAdmin,
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validAdminTokenString, validateJWTOpts).Return(validAdminToken, nil).Once()
			},
		},
		{
			name: "allowing admin only must failed for user",
			input: usecase.AllowAccessInput{
				Token:              validTokenString,
				AllowAllAuthorized: false,
				AllowAdminOnly:     true,
			},
			wantErr:     true,
			expectedErr: usecase.ErrForbidden,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validTokenString, validateJWTOpts).Return(validToken, nil).Once()
			},
		},
		{
			name: "allowing admin only must success for role admin",
			input: usecase.AllowAccessInput{
				Token:              validAdminTokenString,
				AllowAllAuthorized: false,
				AllowAdminOnly:     true,
			},
			wantErr: false,
			expectedOutput: &usecase.AllowAccessOutput{
				UserID:   adminID,
				UserRole: model.RolesAdmin,
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().ValidateJWT(validAdminTokenString, validateJWTOpts).Return(validAdminToken, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.AllowAccess(ctx, tc.input)

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

func TestAuthUSecase_HandleResendSignupVerification(t *testing.T) {
	ctx := context.Background()

	mockSharedCryptor := mockCommon.NewSharedCryptorIface(t)
	mockUserRepo := mockUsecase.NewUserRepository(t)
	mockMailer := mockCommon.NewMailerIface(t)
	mockRateLimiter := mockUsecase.NewRateLimiter(t)

	uc := usecase.NewAuthUsecase(mockSharedCryptor, mockUserRepo, nil, mockMailer, mockRateLimiter)

	userEmail := "email@sample.com"
	encryptedUserEmail := "encryptedUserEmail"
	limiterAllow := redis_rate.Result{
		Allowed: 10,
	}
	user := &model.User{IsActive: false}

	testCases := []struct {
		name                 string
		input                usecase.ResendSignupVerificationInput
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.ResendSignupVerificationOutput
		expectedFunctionCall func()
	}{
		{
			name: "email is required",
			input: usecase.ResendSignupVerificationInput{
				Email: "",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "email is invalid format",
			input: usecase.ResendSignupVerificationInput{
				Email: "invalid@format",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "failed to encrypt the user email",
			input: usecase.ResendSignupVerificationInput{
				Email: "invalid@format",
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "failed to encrypt the user email",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return("", assert.AnError).Once()
			},
		},
		{
			name: "failed to call rate limiter",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "rate limited",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrTooManyRequests,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&redis_rate.Result{
					Allowed: 0,
				}, nil).Once()
			},
		},
		{
			name: "repository failed to fetch user by email",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&limiterAllow, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedUserEmail).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "repository return user not found error",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&limiterAllow, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedUserEmail).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "bad request if user already active",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&limiterAllow, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedUserEmail).Return(&model.User{IsActive: true}, nil).Once()
			},
		},
		{
			name: "failed to generate jwt token",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&limiterAllow, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedUserEmail).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("", assert.AnError).Once()
			},
		},
		{
			name: "failed to send verification email",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&limiterAllow, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedUserEmail).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("token", nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.ResendSignupVerificationInput{
				Email: userEmail,
			},
			wantErr: false,
			expectedOutput: &usecase.ResendSignupVerificationOutput{
				Message: "email confirmation sent",
			},
			expectedFunctionCall: func() {
				mockSharedCryptor.EXPECT().Encrypt(userEmail).Return(encryptedUserEmail, nil).Once()
				mockRateLimiter.EXPECT().Allow(ctx, encryptedUserEmail, mock.Anything).Return(&limiterAllow, nil).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, encryptedUserEmail).Return(user, nil).Once()
				mockSharedCryptor.EXPECT().CreateJWT(mock.Anything).Return("token", nil).Once()
				mockMailer.EXPECT().SendEmail(ctx, mock.Anything).Return(&lib.CreateSmtpEmail{}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleResendSignupVerification(ctx, tc.input)

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
