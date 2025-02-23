// Package usecase contains all the business logic for each app's features
package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis_rate/v10"
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

// AuthUsecase object containing all usecase level logic related to auth process
type AuthUsecase struct {
	sharedCryptor common.SharedCryptorIface
	userRepo      repository.UserRepositoryIface
	mailer        *common.Mailer
	rateLimiter   *redis_rate.Limiter
}

// AuthUsecaseIface interface exported by AuthUsecase to help ease mocking
type AuthUsecaseIface interface {
	HandleSignup(ctx context.Context, input SignupInput) (*SignupOutput, error)
	HandleAccountVerification(ctx context.Context, input AccountVerificationInput) (*AccountVerificationOutput, error)
	HandleLogin(ctx context.Context, input LoginInput) (*LoginOutput, error)
	HandleInitesetPassword(ctx context.Context, input InitResetPasswordInput) (*InitResetPasswordOutput, error)
	HandleResetPassword(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, error)
	AllowAccess(ctx context.Context, input AllowAccessInput) (*AllowAccessOutput, error)
	HandleResendSignupVerification(ctx context.Context, input ResendSignupVerificationInput) (*ResendSignupVerificationOutput, error)
}

// NewAuthUsecase create new instance for AuthUsecase
func NewAuthUsecase(
	sharedCryptor common.SharedCryptorIface,
	userRepo repository.UserRepositoryIface,
	mailer *common.Mailer,
	rateLimiter *redis_rate.Limiter,
) *AuthUsecase {
	return &AuthUsecase{
		sharedCryptor: sharedCryptor,
		userRepo:      userRepo,
		mailer:        mailer,
		rateLimiter:   rateLimiter,
	}
}

// LoginInput input
type LoginInput struct {
	Email    string `vallidate:"required,email"`
	Password string `validate:"required,min=8"`
}

// Validate validate LoginInput's
func (li LoginInput) Validate() error {
	return common.Validator.Struct(li)
}

// LoginOutput output
type LoginOutput struct {
	Token string
}

// HandleLogin contains logic to handle login request
func (u *AuthUsecase) HandleLogin(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	logger := logrus.WithContext(ctx).WithField("email", input.Email)

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	emailEnc, err := u.sharedCryptor.Encrypt(input.Email)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: "encryption process failed",
		}
	}

	user, err := u.userRepo.FindByEmail(ctx, emailEnc)
	switch err {
	default:
		logger.WithError(err).Error("failed to find user by email")

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

	if !user.IsActive {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "this account still not activated",
		}
	}

	pwDecoded, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		logger.WithError(err).Error("failed to decode base 64 string")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	err = u.sharedCryptor.CompareHash(pwDecoded, []byte(input.Password))
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid password",
		}
	}

	loginToken, err := u.sharedCryptor.CreateJWT(model.LoginTokenClaims{
		Role: user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    string(TokenIssuerSystem),
			Subject:   string(LoginToken),
			Audience:  []string{user.ID.String()},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginTokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	if err != nil {
		logger.WithError(err).Error("failed to generate jwt token after login")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &LoginOutput{
		Token: loginToken,
	}, nil
}

// SignupInput input
type SignupInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	Username string `validate:"required"`
}

// Validate validate SignupInput's fields
func (si SignupInput) Validate() error {
	return common.Validator.Struct(si)
}

// SignupOutput output
type SignupOutput struct {
	Message string
}

// JWTTokenType known jwt token type for field sub
type JWTTokenType string

// known JWTTokenType
//
//nolint:gosec
var (
	SignupVerificationToken JWTTokenType = "signup-verification-token"
	LoginToken              JWTTokenType = "login-token"
	ChangePasswordToken     JWTTokenType = "change-password"
)

// JWTTokenIssuer known jwt token issuer for field iss
type JWTTokenIssuer string

// known JWTToken issuer
var (
	TokenIssuerSystem JWTTokenIssuer = "system"
)

// HandleSignup contains logic to handle signup request
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
		HTMLContent:   accountVerificationEmailTemplate(token),
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

// AccountVerificationInput input
type AccountVerificationInput struct {
	VerificationToken string `validate:"required"`
}

// Validate validate AccountVerificationInput's fields
func (avi AccountVerificationInput) Validate() error {
	return common.Validator.Struct(avi)
}

// AccountVerificationOutput output
type AccountVerificationOutput struct {
	Message string
}

// HandleAccountVerification after receiving the email from signup process, if the user receive the email, the token
// inside the email can be used to activate user's account. Then the user can do the login
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

// InitResetPasswordInput input
type InitResetPasswordInput struct {
	Email string `validate:"required,email"`
}

// Validate validate InitResetPasswordInput
func (irpi InitResetPasswordInput) Validate() error {
	return common.Validator.Struct(irpi)
}

// InitResetPasswordOutput output
type InitResetPasswordOutput struct {
	Message string
}

// HandleInitesetPassword will handle change password request by sending email containing the necessary
// data to change password
func (u *AuthUsecase) HandleInitesetPassword(ctx context.Context, input InitResetPasswordInput) (*InitResetPasswordOutput, error) {
	logger := logrus.WithContext(ctx).WithField("email", input.Email)

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	emailEnc, err := u.sharedCryptor.Encrypt(input.Email)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: "encryption process failed",
		}
	}

	user, err := u.userRepo.FindByEmail(ctx, emailEnc)
	switch err {
	default:
		logger.WithError(err).Error("failed to find user by email")

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

	if !user.IsActive {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "this account active status is disabled",
		}
	}

	changePassToken, err := u.sharedCryptor.CreateJWT(jwt.RegisteredClaims{
		Issuer:    string(TokenIssuerSystem),
		Subject:   string(ChangePasswordToken),
		Audience:  []string{user.ID.String()},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.ChangePasswordTokenExpiry())),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	if err != nil {
		logger.WithError(err).Error("failed to create JWT token for change password")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	_, err = u.mailer.SendEmail(ctx, common.SendEmailInput{
		ReceiverName:  user.Username,
		ReceiverEmail: input.Email,
		Subject:       "Permintaan Reset Kata Sandi",
		//nolint:lll
		HTMLContent: fmt.Sprintf(`
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
						<h1>Permintaan Reset Kata Sandi</h1>
					</div>
					<div class="content">
						<p>Baru saja sistem menerima permintaan reset kata sandi terhadap akun Anda. Apabila Anda merasa melakukannya, silahkan klik tombol di bawah ini:</p>
						<div class="btn-container">
							<a href="%s?%s=%s" class="btn">Reset Kata Sandi</a>
						</div>
					</div>
					<div class="footer">
						<p>Jika Anda tidak merasa melakukannya, silahkan abaikan email ini dan akun anda tidak akan diganti kata sandinya</p>
					</div>
				</div>
			</body>
			</html>
			`, config.ServerResetPasswordBaseURL(), model.ChangePasswordTokenQuery, changePassToken),
	})

	if err != nil {
		logger.WithError(err).Error("failed to send verification email to user's mail")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &InitResetPasswordOutput{
		Message: "ok",
	}, nil
}

// ResetPasswordInput input
type ResetPasswordInput struct {
	ResetPasswordToken string `validate:"required"`
	NewPassword        string `validate:"required,min=8"`
}

// Validate validate ResetPasswordInput
func (rpi ResetPasswordInput) Validate() error {
	return common.Validator.Struct(rpi)
}

// ResetPasswordOutput output
type ResetPasswordOutput struct {
	Message string
}

// HandleResetPassword will handle change password if the supplied token is valid
func (u *AuthUsecase) HandleResetPassword(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, error) {
	logger := logrus.WithContext(ctx)

	if err := input.Validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	_, claims, err := u.parseJWTToken(input.ResetPasswordToken, parseJWTTokenInput{
		expectedIssuer:      TokenIssuerSystem,
		expectedSubject:     ChangePasswordToken,
		expectedAudiences:   nil,
		expectedAudienceLen: 1,
	})

	if err != nil {
		return nil, err
	}

	// no need to check the err here, because it's already checked
	// when calling the parseJWTToken
	audiences, _ := claims.GetAudience()
	targetUserID := audiences[0]

	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: "invalid value of user id",
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

	newPasswordHashed, err := u.sharedCryptor.Hash([]byte(input.NewPassword))
	if err != nil {
		logger.WithError(err).Error("failed to hash new user password")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	_, err = u.userRepo.Update(ctx, user.ID, repository.UpdateUserInput{
		Password: newPasswordHashed,
	})

	if err != nil {
		logger.WithError(err).Error("failed to update new user password to database")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &ResetPasswordOutput{
		Message: "ok",
	}, nil
}

// AllowAccessInput input
type AllowAccessInput struct {
	Token              string
	AllowAllAuthorized bool
	AllowAdminOnly     bool
}

// AllowAccessOutput output
type AllowAccessOutput struct {
	UserID   uuid.UUID
	UserRole model.Roles
}

// AllowAccess will perform validation and checking for supplied jwt token.
// Based on the supplied params, you can determine whether to allow the request or not based on the
// returned value. If the error is not nil, safe to assume that you should not let the request pass.
func (u *AuthUsecase) AllowAccess(_ context.Context, input AllowAccessInput) (*AllowAccessOutput, error) {
	jwtToken, _, err := u.parseJWTToken(input.Token, parseJWTTokenInput{
		expectedIssuer:      TokenIssuerSystem,
		expectedSubject:     LoginToken,
		expectedAudiences:   nil,
		expectedAudienceLen: 1,
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !(ok && jwtToken.Valid) {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid jwt payload or the token is invalid",
		}
	}

	// no need to check the err here, because it's already checked
	// when calling the parseJWTToken
	audiences, _ := claims.GetAudience()
	targetUserID := audiences[0]

	userID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid value of user id",
		}
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "undefined role on auth token",
		}
	}

	var userRole model.Roles

	switch role {
	default:
		return nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid role on auth token",
		}
	case string(model.RoleUser):
		userRole = model.RoleUser
	case string(model.RolesAdmin):
		userRole = model.RolesAdmin
	}

	// early return if not specified for admin only
	if input.AllowAllAuthorized {
		return &AllowAccessOutput{
			UserID:   userID,
			UserRole: userRole,
		}, nil
	}

	if input.AllowAdminOnly && userRole != model.RolesAdmin {
		return nil, UsecaseError{
			ErrType: ErrForbidden,
			Message: "access denied due to invalid required access role",
		}
	}

	return &AllowAccessOutput{
		UserID:   userID,
		UserRole: userRole,
	}, nil
}

// ResendSignupVerificationInput input
type ResendSignupVerificationInput struct {
	Email string `validate:"required,email"`
}

func (rsvi ResendSignupVerificationInput) validate() error {
	return common.Validator.Struct(rsvi)
}

// ResendSignupVerificationOutput output
type ResendSignupVerificationOutput struct {
	Message string `json:"message"`
}

// HandleResendSignupVerification will resend the email verification to the user's email
func (u *AuthUsecase) HandleResendSignupVerification(ctx context.Context, input ResendSignupVerificationInput) (
	*ResendSignupVerificationOutput,
	error,
) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"email": input.Email,
		"func":  "HandleResendSignupVerification",
	})

	if err := input.validate(); err != nil {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: err.Error(),
		}
	}

	emailEnc, err := u.sharedCryptor.Encrypt(input.Email)
	if err != nil {
		logger.WithError(err).Error("failed to encrypt email")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	resendLimit := redis_rate.Limit{
		Rate:   1,
		Burst:  1,
		Period: config.ResendSignupVerificationLimiterDuration(),
	}

	rateLimit, err := u.rateLimiter.Allow(ctx, emailEnc, resendLimit)
	if err != nil {
		logger.WithError(err).Error("failed to perform rate limiter ops")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	if rateLimit.Allowed == 0 {
		return nil, UsecaseError{
			ErrType: ErrTooManyRequests,
			Message: fmt.Sprintf("please retry again after %d", int64(rateLimit.ResetAfter.Seconds())),
		}
	}

	user, err := u.userRepo.FindByEmail(ctx, emailEnc)
	switch err {
	default:
		logger.WithError(err).Error("failed to find user by email")

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

	if user.IsActive {
		return nil, UsecaseError{
			ErrType: ErrBadRequest,
			Message: "this account has been activated",
		}
	}

	token, err := u.sharedCryptor.CreateJWT(jwt.RegisteredClaims{
		Issuer:    string(TokenIssuerSystem),
		Subject:   string(SignupVerificationToken),
		Audience:  []string{user.ID.String()},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.SignupTokenExpiry())),
	})

	if err != nil {
		logger.WithError(err).Error("failed to create JWT token for resend signup verification")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	_, err = u.mailer.SendEmail(ctx, common.SendEmailInput{
		ReceiverName:  user.Username,
		ReceiverEmail: input.Email,
		Subject:       "Verifikasi Akun",
		HTMLContent:   accountVerificationEmailTemplate(token),
	})

	if err != nil {
		logger.WithError(err).Error("failed to send verification email to user's mail")

		return nil, UsecaseError{
			ErrType: ErrInternal,
			Message: ErrInternal.Error(),
		}
	}

	return &ResendSignupVerificationOutput{
		Message: "email confirmation sent",
	}, nil
}

type parseJWTTokenInput struct {
	expectedIssuer      JWTTokenIssuer
	expectedSubject     JWTTokenType
	expectedAudiences   *[]string
	expectedAudienceLen int
}

func (u *AuthUsecase) parseJWTToken(token string, input parseJWTTokenInput) (*jwt.Token, jwt.MapClaims, error) {
	jwtToken, err := u.sharedCryptor.ValidateJWT(token, common.ValidateJWTOpts{
		Issuer:  string(input.expectedIssuer),
		Subject: string(input.expectedSubject),
	})

	switch err {
	default:
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: err.Error(),
		}
	case jwt.ErrTokenExpired:
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "change password token has expired",
		}
	case nil:
		break
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !(ok && jwtToken.Valid) {
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid token claims",
		}
	}

	issuer, err := claims.GetIssuer()
	if err != nil || issuer != string(input.expectedIssuer) {
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "incorrect token issuer used",
		}
	}

	subject, err := claims.GetSubject()
	if err != nil || subject != string(input.expectedSubject) {
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "incorrect token subject used",
		}
	}

	audiences, err := claims.GetAudience()
	if err != nil {
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid audience on token used",
		}
	}

	if len(audiences) != input.expectedAudienceLen {
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "invalid number of audience",
		}
	}

	// early return if no expected audience is supplied
	// some times, this checking is not required to performed here
	// but may still be checked by the caller
	if input.expectedAudiences == nil {
		return jwtToken, claims, nil
	}

	expectedAudiences := *input.expectedAudiences

	if !reflect.DeepEqual(audiences, expectedAudiences) {
		return nil, nil, UsecaseError{
			ErrType: ErrUnauthorized,
			Message: "unexpected audiences value",
		}
	}

	return jwtToken, claims, nil
}

//nolint:lll
func accountVerificationEmailTemplate(token string) string {
	return fmt.Sprintf(`
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
						<a href="%s?validation_token=%s" class="btn">Aktifkan Akun</a>
					</div>
				</div>
				<div class="footer">
					<p>Jika Anda tidak merasa mendaftar, abaikan email ini.</p>
				</div>
			</div>
		</body>
		</html>
		`, config.ServerAccountVerificationBaseURL(), token)
}
