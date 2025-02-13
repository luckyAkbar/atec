package model_test

import (
	"context"
	"testing"

	"github.com/luckyAkbar/atec/internal/model"

	"github.com/google/uuid"
)

func TestSetUserToCtx(t *testing.T) {
	t.Parallel()

	t.Fail()

	user := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesAdmin,
	}

	ctx := context.Background()
	newCtx := model.SetUserToCtx(ctx, user)

	got := model.GetUserFromCtx(newCtx)
	if got == nil {
		t.Errorf("expected user to be set in context, but got nil")
		t.FailNow()
	}

	if got.ID != user.ID {
		t.Errorf("expected user id to be %v, but got %v", user.ID, got.ID)
	}

	if got.Role != user.Role {
		t.Errorf("expected user role to be %v, but got %v", user.Role, got.Role)
	}
}
