package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	db_mock "github.com/luckyAkbar/atec/mocks/internal_/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecaeAdapter_UsecaseErrorUCAdapter(t *testing.T) {
	testCase := []struct {
		input  error
		output error
	}{
		{
			input:  assert.AnError,
			output: usecase.ErrRepoInternal,
		},
		{
			input:  repository.ErrNotFound,
			output: usecase.ErrRepoNotFound,
		},
		{
			input:  repository.ErrTimeout,
			output: usecase.ErrRepoTimeout,
		},
		{
			input:  nil,
			output: nil,
		},
	}

	for _, tc := range testCase {
		t.Run("adapter", func(t *testing.T) {
			result := repository.UsecaseErrorUCAdapter(tc.input)
			assert.Equal(t, result, tc.output)
		})
	}
}

func TestChildRepositoryUCAdapter(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewChildRepository(kit.DB)

	adapter := repository.NewChildRepositoryUCAdapter(repo)

	t.Run("Create", func(t *testing.T) {
		dbMock.ExpectBegin()

		dbMock.ExpectQuery("^INSERT INTO \"children\"").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		dbMock.ExpectCommit()

		_, err := adapter.Create(ctx, usecase.RepoCreateChildInput{})
		assert.NoError(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		id := uuid.New()

		dbMock.ExpectQuery(`^SELECT .+ FROM "children"`).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		_, err := adapter.FindByID(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		offset := 11
		limit := 1

		dbMock.ExpectQuery(`^SELECT .+ FROM "children"`).
			WithArgs(limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()))

		_, err := adapter.Search(ctx, usecase.RepoSearchChildInput{
			Limit:  limit,
			Offset: offset,
		})
		assert.NoError(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		childID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectQuery("^UPDATE \"children\" SET").
			WithArgs(sqlmock.AnyArg(), childID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parent_user_id"}).AddRow(childID, uuid.New()))

		dbMock.ExpectCommit()

		_, err := adapter.Update(ctx, childID, usecase.RepoUpdateChildInput{})
		assert.NoError(t, err)
	})

	t.Run("DeleteAllUserChildren", func(t *testing.T) {
		userID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectExec("^DELETE FROM \"children\"").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbMock.ExpectCommit()

		err := adapter.DeleteAllUserChildren(ctx, usecase.RepoDeleteAllUserChildrenInput{
			UserID:     userID,
			HardDelete: true,
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteAllUserChildren with custom tx", func(t *testing.T) {
		userID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectExec("^DELETE FROM \"children\"").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbMock.ExpectCommit()

		err := adapter.DeleteAllUserChildren(ctx, usecase.RepoDeleteAllUserChildrenInput{
			UserID:     userID,
			HardDelete: true,
		}, kit.DB)
		assert.NoError(t, err)
	})

	t.Run("DeleteAllUserChildren with invalid tx", func(t *testing.T) {
		err := adapter.DeleteAllUserChildren(ctx, usecase.RepoDeleteAllUserChildrenInput{
			UserID:     uuid.New(),
			HardDelete: true,
		}, 1)
		assert.Error(t, err)
	})
}

func TestPackageRepositoryUCAdapter(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	cacher := db_mock.NewCacheKeeperIface(t)
	repo := repository.NewPackageRepo(kit.DB, cacher)

	adapter := repository.NewPackageRepositoryUCAdapter(repo)

	t.Run("Create - no controller", func(t *testing.T) {
		dbMock.ExpectBegin()
		dbMock.ExpectQuery(`^INSERT INTO "packages"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(assert.AnError)
		dbMock.ExpectRollback()

		_, err := adapter.Create(ctx, usecase.RepoCreatePackageInput{})
		require.Error(t, err)
	})
	t.Run("Create - with controller", func(t *testing.T) {
		dbMock.ExpectBegin()
		dbMock.ExpectQuery(`^INSERT INTO "packages"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(assert.AnError)
		dbMock.ExpectRollback()

		_, err := adapter.Create(ctx, usecase.RepoCreatePackageInput{}, kit.DB)
		require.Error(t, err)
	})
	t.Run("Create - invalid controller", func(t *testing.T) {
		_, err := adapter.Create(ctx, usecase.RepoCreatePackageInput{}, 1)
		require.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		id := uuid.New()

		cacher.EXPECT().AcquireLock(mock.Anything).Return(nil, assert.AnError).Once()

		err := adapter.Delete(ctx, id)
		require.Error(t, err)
	})

	t.Run("FindAllActivePackages", func(t *testing.T) {
		cacher.EXPECT().GetOrLock(ctx, string(repository.AllActivePackageCacheKey)).
			Return("", nil, assert.AnError).Once()

		_, err := adapter.FindAllActivePackages(ctx)
		require.Error(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		id := uuid.New()

		cacher.EXPECT().GetOrLock(ctx, mock.Anything).
			Return("", nil, assert.AnError).Once()

		_, err := adapter.FindByID(ctx, id)
		require.Error(t, err)
	})

	t.Run("FindOldestActiveAndLockedPackage", func(t *testing.T) {
		dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
			WithArgs(true, true, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		_, err := adapter.FindOldestActiveAndLockedPackage(ctx)
		require.NoError(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		dbMock.ExpectQuery(`^SELECT .+ FROM "packages"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		_, err := adapter.Search(ctx, usecase.RepoSearchPackageInput{})
		require.NoError(t, err)
	})

	t.Run("Update - no controller", func(t *testing.T) {
		id := uuid.New()

		cacher.EXPECT().AcquireLock(mock.Anything).Return(nil, assert.AnError).Once()

		_, err := adapter.Update(ctx, id, usecase.RepoUpdatePackageInput{})
		require.Error(t, err)
	})

	t.Run("Update - with controller", func(t *testing.T) {
		id := uuid.New()

		cacher.EXPECT().AcquireLock(mock.Anything).Return(nil, assert.AnError).Once()

		_, err := adapter.Update(ctx, id, usecase.RepoUpdatePackageInput{}, kit.DB)
		require.Error(t, err)
	})

	t.Run("Update - err", func(t *testing.T) {
		id := uuid.New()

		_, err := adapter.Update(ctx, id, usecase.RepoUpdatePackageInput{}, 1)
		require.Error(t, err)
	})
}

func TestResultRepositoryUCAdapter(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewResultRepository(kit.DB)

	adapter := repository.NewResultRepositoryUCAdapter(repo)

	t.Run("Create", func(t *testing.T) {
		dbMock.ExpectBegin()

		dbMock.ExpectQuery("^INSERT INTO \"results\"").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		dbMock.ExpectCommit()

		_, err := adapter.Create(ctx, usecase.RepoCreateResultInput{})
		require.NoError(t, err)
	})

	t.Run("FindAllUserHistory", func(t *testing.T) {
		dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()))

		_, err := adapter.FindAllUserHistory(ctx, usecase.RepoFindAllUserHistoryInput{})
		require.NoError(t, err)
	})

	t.Run("FindByID", func(t *testing.T) {
		dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		_, err := adapter.FindByID(ctx, uuid.New())
		require.NoError(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		dbMock.ExpectQuery(`^SELECT .+ FROM "results"`).
			WillReturnError(assert.AnError)

		_, err := adapter.Search(ctx, usecase.RepoSearchResultInput{})
		require.Error(t, err)
	})

	t.Run("DeleteAllUserResults", func(t *testing.T) {
		userID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectExec("^DELETE FROM \"results\"").
			WithArgs(userID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbMock.ExpectCommit()

		err := adapter.DeleteAllUserResults(ctx, usecase.RepoDeleteAllUserResultsInput{
			UserID:     userID,
			HardDelete: true,
		})
		require.NoError(t, err)
	})

	t.Run("DeleteAllUserResults with custom tx", func(t *testing.T) {
		userID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectExec("^DELETE FROM \"results\"").
			WithArgs(userID, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbMock.ExpectCommit()

		err := adapter.DeleteAllUserResults(ctx, usecase.RepoDeleteAllUserResultsInput{
			UserID:     userID,
			HardDelete: true,
		}, kit.DB)
		require.NoError(t, err)
	})

	t.Run("DeleteAllUserResults with invalid tx", func(t *testing.T) {
		err := adapter.DeleteAllUserResults(ctx, usecase.RepoDeleteAllUserResultsInput{
			UserID:     uuid.New(),
			HardDelete: true,
		}, 1)
		assert.Error(t, err)
	})
}

func TestUserRepositoryUCAdapter(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	adapter := repository.NewUserRepositoryUCAdapter(repo)

	t.Run("Create - the tx controller is not gorm", func(t *testing.T) {
		_, err := adapter.Create(ctx, usecase.RepoCreateUserInput{}, &model.User{})
		assert.Error(t, err)
	})

	t.Run("Create - the tx controller is not gorm", func(t *testing.T) {
		_, err := adapter.Create(ctx, usecase.RepoCreateUserInput{}, kit.DB)
		assert.Error(t, err)
	})

	t.Run("Create - tx controller wasn't supplied", func(t *testing.T) {
		dbMock.ExpectBegin()

		dbMock.ExpectQuery("^INSERT INTO \"users\"").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		dbMock.ExpectCommit()

		_, err := adapter.Create(ctx, usecase.RepoCreateUserInput{})
		assert.NoError(t, err)
	})

	t.Run("FindByEmail - ok", func(t *testing.T) {
		email := "email"

		dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		_, err := adapter.FindByEmail(ctx, email)
		assert.NoError(t, err)
	})

	t.Run("FindByID - ok", func(t *testing.T) {
		id := uuid.New()

		dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		_, err := adapter.FindByID(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("IsAdminAccountExists - ok", func(t *testing.T) {
		dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(model.RolesAdministrator, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

		_, err := adapter.IsAdminAccountExists(ctx)
		assert.NoError(t, err)
	})

	t.Run("Search - ok", func(t *testing.T) {
		offset := 11
		limit := 1
		role := model.RolesAdministrator

		dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
			WithArgs(role, limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()))

		_, err := adapter.Search(ctx, usecase.RepoSearchUserInput{
			Role:   role,
			Limit:  limit,
			Offset: offset,
		})
		assert.NoError(t, err)
	})

	t.Run("Update - ok", func(t *testing.T) {
		userID := uuid.New()
		email := "test@email.com"
		password := "password1"
		username := "username"
		isActive := true

		dbMock.ExpectBegin()

		dbMock.ExpectQuery("^UPDATE \"users\" SET").
			WithArgs(email, isActive, password, sqlmock.AnyArg(), userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		dbMock.ExpectCommit()

		_, err := adapter.Update(ctx, userID, usecase.RepoUpdateUserInput{
			Email:    email,
			Password: password,
			Username: username,
			IsActive: &isActive,
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteByID - ok", func(t *testing.T) {
		userID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectExec("^DELETE FROM \"users\"").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbMock.ExpectCommit()

		err := adapter.DeleteByID(ctx, usecase.RepoDeleteUserByIDInput{
			UserID:     userID,
			HardDelete: true,
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteByID with custom tx", func(t *testing.T) {
		userID := uuid.New()

		dbMock.ExpectBegin()

		dbMock.ExpectExec("^DELETE FROM \"users\"").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		dbMock.ExpectCommit()

		err := adapter.DeleteByID(ctx, usecase.RepoDeleteUserByIDInput{
			UserID:     userID,
			HardDelete: true,
		}, kit.DB)
		assert.NoError(t, err)
	})

	t.Run("DeleteByID with invalid tx", func(t *testing.T) {
		err := adapter.DeleteByID(ctx, usecase.RepoDeleteUserByIDInput{
			UserID: uuid.New(),
		}, 1)
		assert.Error(t, err)
	})
}
