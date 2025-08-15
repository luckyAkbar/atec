//nolint:thelper // it is okay for test file to have helper functions
package usecase_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	mockCommon "github.com/luckyAkbar/atec/mocks/internal_/common"
	mock_usecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersUsecase_GetMyProfile(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mock_usecase.NewUserRepository(t)
	mockCryptor := mockCommon.NewSharedCryptorIface(t)
	usersUsecase := usecase.NewUsersUsecase(mockUserRepo, mockCryptor)

	now := time.Now()

	// define contexts used in table tests
	ctxRepoErr := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxNotFound := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxDecEmailFail := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxDecPhoneFail := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxDecAddressFail := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxNoPhone := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxNoAddress := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()
	ctxSuccess := func() context.Context {
		return model.SetUserToCtx(ctx, model.AuthUser{ID: uuid.New(), Role: model.RolesParent})
	}()

	testCases := []struct {
		name                 string
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		assertOutput         func(t *testing.T, res *usecase.GetMyProfileOutput)
	}{
		{
			name:        "unauthorized when requester not in context",
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name:        "repo error translated to internal error",
			ctx:         ctxRepoErr,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxRepoErr)
				mockUserRepo.EXPECT().FindByID(ctxRepoErr, requester.ID).Return(nil, usecase.ErrRepoInternal).Once()
			},
		},
		{
			name:        "repo not found translated to usecase not found",
			ctx:         ctxNotFound,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxNotFound)
				mockUserRepo.EXPECT().FindByID(ctxNotFound, requester.ID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name:        "what happen when email decryption fail?",
			ctx:         ctxDecEmailFail,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxDecEmailFail)
				user := &model.User{
					ID:        requester.ID,
					Email:     "enc-email",
					Username:  "user",
					IsActive:  true,
					Roles:     model.RolesParent,
					CreatedAt: now,
					UpdatedAt: now,
				}
				mockUserRepo.EXPECT().FindByID(ctxDecEmailFail, requester.ID).Return(user, nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-email").Return("", assert.AnError).Once()
			},
		},
		{
			name:        "what happen when phone number decryption fail?",
			ctx:         ctxDecPhoneFail,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxDecPhoneFail)
				user := &model.User{
					ID:          requester.ID,
					Email:       "enc-email",
					Username:    "user",
					IsActive:    true,
					Roles:       model.RolesParent,
					CreatedAt:   now,
					UpdatedAt:   now,
					PhoneNumber: sql.NullString{String: "enc-phone", Valid: true},
				}
				mockUserRepo.EXPECT().FindByID(ctxDecPhoneFail, requester.ID).Return(user, nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-email").Return("plain@email", nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-phone").Return("", assert.AnError).Once()
			},
		},
		{
			name:        "what happen when address decryption fail?",
			ctx:         ctxDecAddressFail,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxDecAddressFail)
				user := &model.User{
					ID:        requester.ID,
					Email:     "enc-email",
					Username:  "user",
					IsActive:  true,
					Roles:     model.RolesParent,
					CreatedAt: now,
					UpdatedAt: now,
					Address:   sql.NullString{String: "enc-addr", Valid: true},
				}
				mockUserRepo.EXPECT().FindByID(ctxDecAddressFail, requester.ID).Return(user, nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-email").Return("plain@email", nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-addr").Return("", assert.AnError).Once()
			},
		},
		{
			name: "what happen when the user didn't have phone number?",
			ctx:  ctxNoPhone,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxNoPhone)
				user := &model.User{
					ID:          requester.ID,
					Email:       "enc-email",
					Username:    "user",
					IsActive:    true,
					Roles:       model.RolesParent,
					CreatedAt:   now,
					UpdatedAt:   now,
					PhoneNumber: sql.NullString{Valid: false},
				}
				mockUserRepo.EXPECT().FindByID(ctxNoPhone, requester.ID).Return(user, nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-email").Return("plain@email", nil).Once()
			},
			assertOutput: func(t *testing.T, res *usecase.GetMyProfileOutput) {
				require.NotNil(t, res)
				assert.Nil(t, res.PhoneNumber)
			},
		},
		{
			name: "what happen when the user didn't have address?",
			ctx:  ctxNoAddress,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxNoAddress)
				user := &model.User{
					ID:        requester.ID,
					Email:     "enc-email",
					Username:  "user",
					IsActive:  true,
					Roles:     model.RolesParent,
					CreatedAt: now,
					UpdatedAt: now,
					Address:   sql.NullString{Valid: false},
				}
				mockUserRepo.EXPECT().FindByID(ctxNoAddress, requester.ID).Return(user, nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-email").Return("plain@email", nil).Once()
			},
			assertOutput: func(t *testing.T, res *usecase.GetMyProfileOutput) {
				require.NotNil(t, res)
				assert.Nil(t, res.Address)
			},
		},
		{
			name: "success",
			ctx:  ctxSuccess,
			expectedFunctionCall: func() {
				requester := model.GetUserFromCtx(ctxSuccess)
				user := &model.User{
					ID:          requester.ID,
					Email:       "enc-email",
					Username:    "user",
					IsActive:    true,
					Roles:       model.RolesParent,
					CreatedAt:   now,
					UpdatedAt:   now,
					PhoneNumber: sql.NullString{String: "enc-phone", Valid: true},
					Address:     sql.NullString{String: "enc-addr", Valid: true},
				}
				mockUserRepo.EXPECT().FindByID(ctxSuccess, requester.ID).Return(user, nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-email").Return("plain@email", nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-phone").Return("+62812", nil).Once()
				mockCryptor.EXPECT().Decrypt("enc-addr").Return("street 1", nil).Once()
			},
			assertOutput: func(t *testing.T, res *usecase.GetMyProfileOutput) {
				require.NotNil(t, res)
				assert.Equal(t, "user", res.Username)
				assert.True(t, res.IsActive)
				assert.Equal(t, model.RolesParent, res.Roles)
				require.NotNil(t, res.PhoneNumber)
				require.NotNil(t, res.Address)
				assert.Equal(t, "+62812", *res.PhoneNumber)
				assert.Equal(t, "street 1", *res.Address)
				assert.WithinDuration(t, now, res.CreatedAt, time.Second)
				assert.WithinDuration(t, now, res.UpdatedAt, time.Second)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := usersUsecase.GetMyProfile(tc.ctx)

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

			if tc.assertOutput != nil {
				tc.assertOutput(t, res)
			}
		})
	}
}

func TestUsersUsecase_GetTherapistData(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := mock_usecase.NewUserRepository(t)
	usersUsecase := usecase.NewUsersUsecase(mockUserRepo, mockCommon.NewSharedCryptorIface(t))

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
