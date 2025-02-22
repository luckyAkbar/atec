package model

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authCtxKey string

var (
	authUserCtxKey authCtxKey = "github.com/luckyAkbar/atec/internal/model:AuthUser"
)

// AuthUser represent authenticated user and will be used to embed value to context
type AuthUser struct {
	ID   uuid.UUID
	Role Roles
}

// SetUserToCtx set user to context
func SetUserToCtx(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, authUserCtxKey, user)
}

// GetUserFromCtx get user from context
func GetUserFromCtx(ctx context.Context) *AuthUser {
	user, ok := ctx.Value(authUserCtxKey).(AuthUser)
	if !ok {
		return nil
	}

	return &user
}

// LoginTokenClaims custom claims to be placed in payload for login jwt token
type LoginTokenClaims struct {
	Role Roles `json:"role"`

	jwt.RegisteredClaims
}

// ChangePasswordTokenQuery is the key in the query parameters to handle
// change password request
const ChangePasswordTokenQuery = "change_password_token"
