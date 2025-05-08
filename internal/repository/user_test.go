package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository_FindByEmail(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	email := "email@test.com"

	testCases := []struct {
		name                 string
		email                string
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "success",
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(email, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
		},
		{
			name:        "error - unknown just pass the error to the caller",
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(email, 1).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:        "data not found on db",
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(email, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindByEmail(ctx, email)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, res)

				if tc.expectedErr != nil {
					assert.Equal(t, tc.expectedErr, err)
				}

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}
