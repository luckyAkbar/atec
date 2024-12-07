package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	HandleAccountVerification(ctx context.Context, input AccountVerificationInput) (*AccountVerificationOutput, error)
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

type JWTTokenIssuer string

var (
	TokenIssuerSystem JWTTokenIssuer = "system"
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
		Issuer:    string(TokenIssuerSystem),
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
			<!DOCTYPE html>
			<html>
			<head>
				<style>
					/* General styles */
					body {
						font-family: Arial, sans-serif;
						margin: 0;
						padding: 0;
						background-color: #f6f6f6;
						color: #333;
					}
					.email-container {
						max-width: 600px;
						margin: 20px auto;
						background-color: #ffffff;
						border: 1px solid #ddd;
						border-radius: 8px;
						overflow: hidden;
					}
					.header {
						background-color: #4CAF50;
						color: white;
						padding: 20px;
						text-align: center;
					}
					.content {
						padding: 20px;
					}
					.content p {
						margin: 0 0 15px;
						line-height: 1.6;
					}
					.btn-container {
						text-align: center;
						margin: 20px 0;
					}
					.btn {
						display: inline-block;
						background-color: #4CAF50;
						color: white;
						text-decoration: none;
						padding: 10px 20px;
						font-size: 16px;
						border-radius: 5px;
					}
					.btn:hover {
						background-color: #45a049;
					}
					.footer {
						background-color: #f1f1f1;
						text-align: center;
						padding: 10px;
						font-size: 12px;
						color: #666;
					}
				</style>
			</head>
			<body>
				<div class="email-container">
					<div class="header">
						<h1>Konfirmasi Akun</h1>
					</div>
					<div class="content">
						<p>Terimakasih telah mendaftar pada layanan Autism Treatment Evaluation Checklist (ATEC). Untuk mengaktifkan akun Anda, silakan klik tombol berikut:</p>
						<div class="btn-container">
							<a href="%s?verification_token=%s" class="btn">Aktifkan Akun</a>
						</div>
					</div>
					<div class="footer">
						<p>Jika Anda tidak merasa mendaftar, abaikan email ini.</p>
					</div>
				</div>
			</body>
			</html>
			`, config.ServerAccountVerificationBaseURL(), token), // TODO replace the confirmation link once the feature has been developed
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

type AccountVerificationInput struct {
	VerificationToken string `validate:"required"`
}

func (avi AccountVerificationInput) Validate() error {
	return validator.Struct(avi)
}

type AccountVerificationOutput struct {
	Message string
}

func (u *AuthUsecase) HandleAccountVerification(ctx context.Context, input AccountVerificationInput) (*AccountVerificationOutput, error) {
	logger := logrus.WithContext(ctx)

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	token, err := u.sharedCryptor.ValidateJWT(input.VerificationToken, common.ValidateJWTOpts{
		Issuer:  string(TokenIssuerSystem),
		Subject: string(SignupVerificationToken),
	})

	switch err {
	default:
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid token for account verification",
		}
	case jwt.ErrTokenExpired:
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "account validation token has expired",
		}
	case nil:
		break
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid token claims",
		}
	}

	issuer, err := claims.GetIssuer()
	if err != nil || issuer != string(TokenIssuerSystem) {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "incorrect token issuer used",
		}
	}

	subject, err := claims.GetSubject()
	if err != nil || subject != string(SignupVerificationToken) {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "incorrect token subject used",
		}
	}

	audiences, err := claims.GetAudience()
	if err != nil || len(audiences) != 1 {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid audience on token used",
		}
	}

	targetUserID := audiences[0]
	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid user id received from account validation token",
		}
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find user data from db")
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

	// early return if already activated
	if user.IsActive {
		return &AccountVerificationOutput{
			Message: "your account has been activated",
		}, nil
	}

	activeTrue := true
	_, err = u.userRepo.Update(ctx, user.ID, repository.UpdateUserInput{
		IsActive: &activeTrue,
	})

	if err != nil {
		logger.WithError(err).Error("failed to activate user account to database")
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &AccountVerificationOutput{
		Message: "your account has been activated",
	}, nil
}
