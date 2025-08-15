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

func TestUsersUsecase_GetTherapistData(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mock_usecase.NewUserRepository(t)
	usersUsecase := usecase.NewUsersUsecase(mockUserRepo)

	now := time.Now()

	ctxInternal := func() context.Context {
		requester := model.AuthUser{ID: uuid.New(), Role: model.RolesParent}

		return model.SetUserToCtx(ctx, requester)
	}()
	ctxNotFound := func() context.Context {
		requester := model.AuthUser{ID: uuid.New(), Role: model.RolesParent}

		return model.SetUserToCtx(ctx, requester)
	}()
	ctxSuccess := func() context.Context {
		requester := model.AuthUser{ID: uuid.New(), Role: model.RolesParent}

		return model.SetUserToCtx(ctx, requester)
	}()

	testCases := []struct {
		name                 string
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutput       []usecase.GetTherapistDataOutput
		expectedFunctionCall func()
	}{
		{
			name:                 "unauthorized when requester not in context",
			ctx:                  ctx,
			wantErr:              true,
			expectedErr:          usecase.ErrUnauthorized,
			expectedFunctionCall: func() {},
		},
		{
			name:        "repo error translated to internal error",
			ctx:         ctxInternal,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockUserRepo.EXPECT().GetUsersByRoles(ctxInternal, model.RolesTherapist).Return(nil, usecase.ErrRepoInternal).Once()
			},
		},
		{
			name:        "repo not found translated to usecase not found",
			ctx:         ctxNotFound,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockUserRepo.EXPECT().GetUsersByRoles(ctxNotFound, model.RolesTherapist).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:    "success",
			ctx:     ctxSuccess,
			wantErr: false,
			expectedFunctionCall: func() {
				users := []model.User{
					{
						ID:        uuid.New(),
						Username:  "t1",
						IsActive:  true,
						Roles:     model.RolesTherapist,
						CreatedAt: now,
						UpdatedAt: now,
					},
					{
						ID:        uuid.New(),
						Username:  "t2",
						IsActive:  false,
						Roles:     model.RolesTherapist,
						CreatedAt: now,
						UpdatedAt: now,
					},
				}

				mockUserRepo.EXPECT().GetUsersByRoles(ctxSuccess, model.RolesTherapist).Return(users, nil).Once()
			},
			expectedOutput: []usecase.GetTherapistDataOutput{
				{
					Username:  "t1",
					IsActive:  true,
					Roles:     model.RolesTherapist,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					Username:  "t2",
					IsActive:  false,
					Roles:     model.RolesTherapist,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}

	for idx := range testCases {
		tc := testCases[idx]

		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := usersUsecase.GetTherapistData(tc.ctx)

			if tc.wantErr {
				require.Error(t, err)

				if tc.expectedErr != nil {
					ucErr, ok := err.(usecase.UsecaseError)
					require.True(t, ok)
					assert.Equal(t, tc.expectedErr, ucErr.ErrType)
				}

				assert.Nil(t, res)

				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			require.Len(t, res, len(tc.expectedOutput))

			for i := range tc.expectedOutput {
				exp := tc.expectedOutput[i]
				got := res[i]

				assert.Equal(t, exp.Username, got.Username)
				assert.Equal(t, exp.IsActive, got.IsActive)
				assert.Equal(t, exp.Roles, got.Roles)
				assert.WithinDuration(t, exp.CreatedAt, got.CreatedAt, time.Second)
				assert.WithinDuration(t, exp.UpdatedAt, got.UpdatedAt, time.Second)
			}
		})
	}
}
