package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestResultRepository_Create(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewResultRepository(kit.DB)

	packageID := uuid.New()
	childID := uuid.New()
	CreatedByID := uuid.New()

	dbGeneratedUUID := uuid.New()

	testCases := []struct {
		name                 string
		input                usecase.RepoCreateResultInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name: "success",
			input: usecase.RepoCreateResultInput{
				PackageID: packageID,
				ChildID:   childID,
				CreatedBy: CreatedByID,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^INSERT INTO \"results\"").
					WithArgs(packageID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), childID, CreatedByID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(dbGeneratedUUID))

				dbMock.ExpectCommit()
			},
		},
		{
			name: "error",
			input: usecase.RepoCreateResultInput{
				PackageID: packageID,
				ChildID:   childID,
				CreatedBy: CreatedByID,
			},
			wantErr: true,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^INSERT INTO \"results\"").
					WithArgs(packageID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), childID, CreatedByID).
					WillReturnError(assert.AnError)

				dbMock.ExpectRollback()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.Create(ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.NotNil(t, res)

				return
			}

			require.Error(t, err)
			require.Nil(t, res)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			}
		})
	}
}

func TestResultRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewResultRepository(kit.DB)

	resultID := uuid.New()

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "success",
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(resultID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(resultID))
			},
		},
		{
			name:        "error - unknown just pass the error to the caller",
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(resultID, 1).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:        "data not found on db",
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(resultID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindByID(ctx, resultID)

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

func TestResultRepository_Search(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewResultRepository(kit.DB)

	resultID := uuid.New()
	packageID := uuid.New()
	childID := uuid.New()
	createdByID := uuid.New()
	limit := 100
	offset := 10

	testCases := []struct {
		name                 string
		input                usecase.RepoSearchResultInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutputLen    int
	}{
		{
			name:    "database returning unexpected error",
			wantErr: true,
			input: usecase.RepoSearchResultInput{
				ID:        resultID,
				PackageID: packageID,
				ChildID:   childID,
				CreatedBy: createdByID,
				Limit:     limit,
				Offset:    offset,
			},
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(resultID, packageID, childID, createdByID, limit, offset).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:    "no result found must return not found error",
			wantErr: true,
			input: usecase.RepoSearchResultInput{
				ID:        resultID,
				PackageID: packageID,
				ChildID:   childID,
				CreatedBy: createdByID,
				Limit:     limit,
				Offset:    offset,
			},
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(resultID, packageID, childID, createdByID, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			name:    "ok",
			wantErr: false,
			input: usecase.RepoSearchResultInput{
				ID:        resultID,
				PackageID: packageID,
				ChildID:   childID,
				CreatedBy: createdByID,
				Limit:     limit,
				Offset:    offset,
			},
			expectedOutputLen: 2,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(resultID, packageID, childID, createdByID, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(resultID).AddRow(uuid.New()))
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

			assert.Len(t, res, tc.expectedOutputLen)
		})
	}
}

func TestResultRepository_FindAllUserHistory(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewResultRepository(kit.DB)

	userID := uuid.New()
	limit := 100
	offset := 10

	testCases := []struct {
		name                 string
		input                usecase.RepoFindAllUserHistoryInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutputLen    int
	}{
		{
			name:    "database returning unexpected error",
			wantErr: true,
			input: usecase.RepoFindAllUserHistoryInput{
				UserID: userID,
				Limit:  limit,
				Offset: offset,
			},
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(userID, userID, limit, offset).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:    "no result found must return not found error",
			wantErr: true,
			input: usecase.RepoFindAllUserHistoryInput{
				UserID: userID,
				Limit:  limit,
				Offset: offset,
			},
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(userID, userID, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			name:    "ok",
			wantErr: false,
			input: usecase.RepoFindAllUserHistoryInput{
				UserID: userID,
				Limit:  limit,
				Offset: offset,
			},
			expectedOutputLen: 2,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
					WithArgs(userID, userID, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindAllUserHistory(ctx, tc.input)

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

			assert.Len(t, res, tc.expectedOutputLen)
		})
	}
}
