package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	mock_usecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersUsecase_GetMyProfile(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mock_usecase.NewUserRepository(t)
	usersUsecase := usecase.NewUsersUsecase(mockUserRepo)

	t.Run("unauthorized when requester not in context", func(t *testing.T) {
		res, err := usersUsecase.GetMyProfile(ctx)

		require.Error(t, err)
		assert.Nil(t, res)

		ucErr, ok := err.(usecase.UsecaseError)
		require.True(t, ok)
		assert.Equal(t, usecase.ErrUnauthorized, ucErr.ErrType)
	})

	t.Run("repo error translated to internal error", func(t *testing.T) {
		requester := model.AuthUser{ID: uuid.New(), Role: model.RolesParent}
		reqCtx := model.SetUserToCtx(ctx, requester)

		mockUserRepo.EXPECT().FindByID(reqCtx, requester.ID).Return(nil, usecase.ErrRepoInternal).Once()

		res, err := usersUsecase.GetMyProfile(reqCtx)

		require.Error(t, err)
		assert.Nil(t, res)

		ucErr, ok := err.(usecase.UsecaseError)
		require.True(t, ok)
		assert.Equal(t, usecase.ErrInternal, ucErr.ErrType)
	})

	t.Run("repo not found translated to usecase not found", func(t *testing.T) {
		requester := model.AuthUser{ID: uuid.New(), Role: model.RolesParent}
		reqCtx := model.SetUserToCtx(ctx, requester)

		mockUserRepo.EXPECT().FindByID(reqCtx, requester.ID).Return(nil, usecase.ErrRepoNotFound).Once()

		res, err := usersUsecase.GetMyProfile(reqCtx)

		require.Error(t, err)
		assert.Nil(t, res)

		ucErr, ok := err.(usecase.UsecaseError)
		require.True(t, ok)
		assert.Equal(t, usecase.ErrNotFound, ucErr.ErrType)
	})

	t.Run("success", func(t *testing.T) {
		requester := model.AuthUser{ID: uuid.New(), Role: model.RolesParent}
		reqCtx := model.SetUserToCtx(ctx, requester)

		now := time.Now()
		user := &model.User{
			ID:        requester.ID,
			Username:  "user",
			IsActive:  true,
			Roles:     model.RolesParent,
			CreatedAt: now,
			UpdatedAt: now,
		}

		mockUserRepo.EXPECT().FindByID(reqCtx, requester.ID).Return(user, nil).Once()

		res, err := usersUsecase.GetMyProfile(reqCtx)

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, user.ID, res.ID)
		assert.Equal(t, user.Username, res.Username)
		assert.Equal(t, user.IsActive, res.IsActive)
		assert.Equal(t, user.Roles, res.Roles)
		assert.WithinDuration(t, user.CreatedAt, res.CreatedAt, time.Second)
		assert.WithinDuration(t, user.UpdatedAt, res.UpdatedAt, time.Second)
	})
}
