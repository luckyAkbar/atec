package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/common"
	"github.com/luckyAkbar/atec/internal/db"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	common_mock "github.com/luckyAkbar/atec/mocks/internal_/common"
	db_mock "github.com/luckyAkbar/atec/mocks/internal_/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPackageRepository_Create(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	cacher := db_mock.NewCacheKeeperIface(t)
	redsyncMutex := common_mock.NewRedsyncMutex(t)
	repo := repository.NewPackageRepo(kit.DB, cacher)

	mutexWrapper := common.NewRedsyncMutexWrapper(redsyncMutex)

	testCases := []struct {
		name                 string
		wantErr              bool
		input                usecase.RepoCreatePackageInput
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:        "unexpected db error when trying to submit",
			wantErr:     true,
			input:       usecase.RepoCreatePackageInput{},
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()
				dbMock.ExpectQuery(`^INSERT INTO "packages"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
				dbMock.ExpectRollback()
			},
		},
		{
			name:        "failed to acquire lock doesn't returning error",
			wantErr:     false,
			input:       usecase.RepoCreatePackageInput{},
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()
				dbMock.ExpectQuery(`^INSERT INTO "packages"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				dbMock.ExpectCommit()

				cacher.EXPECT().AcquireLock(mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "failed to set cache to JSON, but OK",
			wantErr: false,
			input:   usecase.RepoCreatePackageInput{},
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()
				dbMock.ExpectQuery(`^INSERT INTO "packages"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				dbMock.ExpectCommit()

				cacher.EXPECT().AcquireLock(mock.Anything).Return(mutexWrapper, nil).Once()
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().SetJSON(ctx, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			name:    "failed to unlock, but OK",
			wantErr: false,
			input:   usecase.RepoCreatePackageInput{},
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()
				dbMock.ExpectQuery(`^INSERT INTO "packages"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				dbMock.ExpectCommit()

				cacher.EXPECT().AcquireLock(mock.Anything).Return(mutexWrapper, nil).Once()
				redsyncMutex.EXPECT().Unlock().Return(false, assert.AnError).Once()
				cacher.EXPECT().SetJSON(ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.Create(ctx, tc.input, kit.DB)

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

func TestPackageRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	// dbMock := kit.DBmock
	cacher := db_mock.NewCacheKeeperIface(t)
	repo := repository.NewPackageRepo(kit.DB, cacher)
	id := uuid.New()

	redsyncMutex := common_mock.NewRedsyncMutex(t)
	mutexWrapper := common.NewRedsyncMutexWrapper(redsyncMutex)

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "failed to get or lock the cache by key",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("!", nil, assert.AnError).Once()
			},
		},
		{
			name:        "cache is set to nil",
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("!", nil, db.ErrCacheNil).Once()
			},
		},
		{
			name:        "cache took too long",
			wantErr:     true,
			expectedErr: repository.ErrTimeout,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("!", nil, db.ErrLockWaitingTooLong).Once()
			},
		},
		{
			name:    "failed unmarshalling the value",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("!", nil, nil).Once()
			},
		},
		{
			name:    "got marshalled value from cache",
			wantErr: false,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("{}", nil, nil).Once()
			},
		},
		{
			name:    "when fetching from db, got error",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("", mutexWrapper, nil).Once()
				kit.DBmock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(id, 1).
					WillReturnError(assert.AnError)
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
			},
		},
		{
			name:    "when fetching from db, got error, also failed to unlock",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("", mutexWrapper, nil).Once()
				kit.DBmock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(id, 1).
					WillReturnError(assert.AnError)
				redsyncMutex.EXPECT().Unlock().Return(false, assert.AnError).Once()
			},
		},
		{
			name:    "ok - all good",
			wantErr: false,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("", mutexWrapper, nil).Once()
				kit.DBmock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().SetJSON(ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name:    "the package wasn't found on db",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("", mutexWrapper, nil).Once()
				kit.DBmock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().SetNil(ctx, mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name:    "the package wasn't found on db - failed to set nil to cache",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("", mutexWrapper, nil).Once()
				kit.DBmock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().SetNil(ctx, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			name:    "ok - failed to set json",
			wantErr: false,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(mock.Anything, mock.Anything).Return("", mutexWrapper, nil).Once()
				kit.DBmock.ExpectQuery(`^SELECT .+ FROM "packages"`).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().SetJSON(ctx, mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindByID(ctx, id)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, res)

				if tc.expectedErr != nil {
					t.Logf("Got error: %v (type: %T), Expected error: %v (type: %T)", err, err, tc.expectedErr, tc.expectedErr)
					assert.Equal(t, tc.expectedErr, err)
				}

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, res)
		})
	}
}

func TestPackageRepository_Update(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	// dbMock := kit.DBmock
	cacher := db_mock.NewCacheKeeperIface(t)
	repo := repository.NewPackageRepo(kit.DB, cacher)
	id := uuid.New()

	redsyncMutex := common_mock.NewRedsyncMutex(t)
	mutexWrapper := common.NewRedsyncMutexWrapper(redsyncMutex)

	activeStatus := true
	lockStatus := true
	packageName := "package name"
	questionnaire := &model.Questionnaire{}

	testCases := []struct {
		name                 string
		wantErr              bool
		input                usecase.RepoUpdatePackageInput
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "failed to acquire lock",
			wantErr: true,
			input:   usecase.RepoUpdatePackageInput{},
			expectedFunctionCall: func() {
				cacher.EXPECT().AcquireLock(mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "when failed to delete stale data from cache, must return error",
			wantErr: true,
			input:   usecase.RepoUpdatePackageInput{},
			expectedFunctionCall: func() {
				cacher.EXPECT().AcquireLock(mock.Anything).Return(mutexWrapper, nil).Once()
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().Del(ctx, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			name:    "when failed to delete stale data from cache, must return error - failed to unlock",
			wantErr: true,
			input:   usecase.RepoUpdatePackageInput{},
			expectedFunctionCall: func() {
				cacher.EXPECT().AcquireLock(mock.Anything).Return(mutexWrapper, nil).Once()
				redsyncMutex.EXPECT().Unlock().Return(false, assert.AnError).Once()
				cacher.EXPECT().Del(ctx, mock.Anything).Return(assert.AnError).Once()
			},
		},
		{
			name:    "got error from db when updating the package",
			wantErr: true,
			input: usecase.RepoUpdatePackageInput{
				ActiveStatus:  &activeStatus,
				LockStatus:    &lockStatus,
				PackageName:   packageName,
				Questionnaire: questionnaire,
			},
			expectedFunctionCall: func() {
				cacher.EXPECT().AcquireLock(mock.Anything).Return(mutexWrapper, nil).Once()
				redsyncMutex.EXPECT().Unlock().Return(true, nil).Once()
				cacher.EXPECT().Del(ctx, mock.Anything).Return(nil).Once()

				kit.DBmock.ExpectQuery(`^UPDATE "packages"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.Update(ctx, id, tc.input, kit.DB)

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

func TestPackageRepository_Delete(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	// dbMock := kit.DBmock
	cacher := db_mock.NewCacheKeeperIface(t)
	repo := repository.NewPackageRepo(kit.DB, cacher)
	id := uuid.New()

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "failed to acquire lock",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().AcquireLock(mock.Anything).Return(nil, assert.AnError).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			err := repo.Delete(ctx, id)

			if tc.wantErr {
				require.Error(t, err)

				if tc.expectedErr != nil {
					assert.Equal(t, tc.expectedErr, err)
				}

				return
			}

			require.NoError(t, err)
		})
	}
}

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

func TestPackageRepository_FindAllActivePackages(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	cacher := db_mock.NewCacheKeeperIface(t)
	repo := repository.NewPackageRepo(kit.DB, cacher)

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name:    "failed to get or lock",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(ctx, mock.Anything).Return("", nil, assert.AnError).Once()
			},
		},
		{
			name:        "cache was set to nil",
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(ctx, mock.Anything).Return("", nil, db.ErrCacheNil).Once()
			},
		},
		{
			name:        "too long waiting for cache",
			wantErr:     true,
			expectedErr: repository.ErrTimeout,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(ctx, mock.Anything).Return("", nil, db.ErrLockWaitingTooLong).Once()
			},
		},
		{
			name:    "failed to unmarshall",
			wantErr: true,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(ctx, mock.Anything).Return("!", nil, nil).Once()
			},
		},
		{
			name:    "ok from cache",
			wantErr: false,
			expectedFunctionCall: func() {
				cacher.EXPECT().GetOrLock(ctx, mock.Anything).Return(`[]`, nil, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindAllActivePackages(ctx)

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
