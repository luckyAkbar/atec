package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPackageRepository_Search(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	isActive := true
	limit := 10

	dbMock := kit.DBmock
	repo := repository.NewPackageRepo(kit.DB, nil)

	unexpectedErr := errors.New("unexpected db error")

	testCases := []struct {
		name                 string
		wantErr              bool
		input                usecase.RepoSearchPackageInput
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "unexpected db error",
			wantErr: true,
			input: usecase.RepoSearchPackageInput{
				IsActive: &isActive,
				Limit:    limit,
			},
			expectedErr: unexpectedErr,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(isActive, limit).
					WillReturnError(unexpectedErr)
			},
		},
		{
			name:    "when no package found, must return not found error",
			wantErr: true,
			input: usecase.RepoSearchPackageInput{
				IsActive: &isActive,
				Limit:    limit,
			},
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(isActive, limit).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			name:    "ok",
			wantErr: false,
			input: usecase.RepoSearchPackageInput{
				IsActive: &isActive,
				Limit:    limit,
			},
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(isActive, limit).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.Search(ctx, tc.input)

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

func TestPackageRepository_FindOldestActiveAndLockedPackage(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewPackageRepo(kit.DB, nil)

	unexpectedErr := errors.New("unexpected db error")

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:        "unexpected db error",
			wantErr:     true,
			expectedErr: unexpectedErr,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(true, true, 1).
					WillReturnError(unexpectedErr)
			},
		},
		{
			name:        "when no package found, must return not found error",
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(true, true, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:    "ok",
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(true, true, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindOldestActiveAndLockedPackage(ctx)

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
