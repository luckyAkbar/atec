package model_test

import (
	"context"
	"testing"

	"github.com/luckyAkbar/atec/internal/model"

	"github.com/google/uuid"
)

func TestSetUserToCtx(t *testing.T) {
	user := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesAdministrator,
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

func TestGetUserFromCtx(t *testing.T) {
	user := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesAdministrator,
	}

	ctx := context.Background()
	newCtx := model.SetUserToCtx(ctx, user)

	t.Run("should correctly read the data", func(t *testing.T) {
		got := model.GetUserFromCtx(newCtx)
		if got == nil {
			t.Errorf("user data is already set, should be able to read it")
			t.FailNow()
		}

		if got.ID != user.ID {
			t.Errorf("expected user id to be %v, but got %v", user.ID, got.ID)
		}

		if got.Role != user.Role {
			t.Errorf("expected user role to be %v, but got %v", user.Role, got.Role)
		}
	})

	t.Run("if the user data is not in the ctx, must return nil", func(t *testing.T) {
		emptyCtx := context.Background()

		got := model.GetUserFromCtx(emptyCtx)
		if got != nil {
			t.Errorf("on empty context value, should return nil")
		}
	})
}
