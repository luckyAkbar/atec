package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/config"
	"github.com/luckyAkbar/atec/internal/db"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/stdlib/helper"
)

type AuthUsecase struct {
	sharedCryptor common.SharedCryptorIface
	userRepo      repository.UserRepositoryIface
	mailer        *common.Mailer
}

type AuthUsecaseIface interface {
	HandleSignup(ctx context.Context, input SignupInput) (*SignupOutput, error)
}

func NewAuthUsecase(
	sharedCryptor common.SharedCryptorIface,
	userRepo repository.UserRepositoryIface,
	mailer *common.Mailer,
) *AuthUsecase {
	return &AuthUsecase{
		sharedCryptor: sharedCryptor,
		userRepo:      userRepo,
		mailer:        mailer,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

func (u *AuthUsecase) HandleLogin(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	return nil, nil
}

type SignupInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Username string `validate:"required"`
}

func (si SignupInput) Validate() error {
	return validator.Struct(si)
}

type SignupOutput struct {
	Message string
}

type JWTTokenType string

var (
	SignupVerificationToken JWTTokenType = "signup-verification-token"
)

func (u *AuthUsecase) HandleSignup(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"email": helper.Dump(input.Email),
	})

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	emailEncrypted, err := u.sharedCryptor.Encrypt(input.Email)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: "encryption process failed",
		}
	}

	_, err = u.userRepo.FindByEmail(ctx, emailEncrypted)
	switch err {
	default:
		logger.WithError(err).Error("failed to perform query to find user by id")
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	case nil:
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: "your email has been used by another account",
		}
	case repository.ErrNotFound:
		break
	}

	hashedPassword, err := u.sharedCryptor.Hash([]byte(input.Password))
	if err != nil {
		logger.WithError(err).Error("failed to perform hasing password")
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	createUserInput := repository.CreateUserInput{
		Email:    emailEncrypted,
		Password: hashedPassword,
		Username: input.Username,
		IsActive: false,
		Roles:    model.RoleUser,
	}

	tx := db.TxController()

	user, err := u.userRepo.Create(ctx, createUserInput, tx)
	if err != nil {
		logger.WithError(err).Error("failed to create user to database")
		tx.Rollback()
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	token, err := u.sharedCryptor.CreateJWT(jwt.RegisteredClaims{
		Issuer:    "system",
		Subject:   string(SignupVerificationToken),
		Audience:  []string{user.ID.String()},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.SignupTokenExpiry())),
	})

	if err != nil {
		logger.WithError(err).Error("failed to create JWT token for signup verification")
		tx.Rollback()
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	_, err = u.mailer.SendEmail(ctx, common.SendEmailInput{
		ReceiverName:  user.Username,
		ReceiverEmail: input.Email,
		Subject:       "Verifikasi Akun",
		HtmlContent: fmt.Sprintf(`
			<h2>Halo %s!</h2>
			<p>Terimakasih telah mendaftar pada layanan Autism Treatment Evaluation Checklist (ATEC)</p>
			</p> Untuk mengaktifkan akun Anda, silakan klik link berikut: </p> <br>
			<a href="%s">validasi akun</a>	
		`, user.Username, token), // TODO replace the confirmation link once the feature has been developed
	})

	if err != nil {
		logger.WithError(err).Error("failed to send verification email to user's mail")
		tx.Rollback()
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	tx.Commit()

	return &SignupOutput{
		Message: "email confirmation sent",
	}, nil
}
