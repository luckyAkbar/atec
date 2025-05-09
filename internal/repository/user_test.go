package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
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

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	email := "test@emailc.com"
	password := "password"
	isActive := true
	roles := model.RolesAdministrator
	username := "testuser"

	testCases := []struct {
		name                 string
		input                usecase.RepoCreateUserInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name: "success",
			input: usecase.RepoCreateUserInput{
				Email:    email,
				Password: password,
				IsActive: isActive,
				Roles:    roles,
				Username: username,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^INSERT INTO \"users\"").
					WithArgs(email, password, username, isActive, roles, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

				dbMock.ExpectCommit()
			},
		},
		{
			name: "error",
			input: usecase.RepoCreateUserInput{
				Email:    email,
				Password: password,
				IsActive: isActive,
				Roles:    roles,
				Username: username,
			},
			wantErr: true,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^INSERT INTO \"users\"").
					WithArgs(email, password, username, isActive, roles, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

func TestUserRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	userID := uuid.New()

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
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(userID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
			},
		},
		{
			name:        "error - unknown just pass the error to the caller",
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(userID, 1).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:        "data not found on db",
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(userID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindByID(ctx, userID)

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

func TestUserRepository_Update(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	userID := uuid.New()
	email := "test@email.com"
	password := "password"
	username := "testuser"
	isActive := true

	testCases := []struct {
		name                 string
		input                usecase.RepoUpdateUserInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutput       *model.Child
	}{
		{
			name: "success",
			input: usecase.RepoUpdateUserInput{
				Email:    email,
				Password: password,
				Username: username,
				IsActive: &isActive,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^UPDATE \"users\" SET").
					WithArgs(email, isActive, password, sqlmock.AnyArg(), userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

				dbMock.ExpectCommit()
			},
		},
		{
			name: "error",
			input: usecase.RepoUpdateUserInput{
				Email:    email,
				Password: password,
				Username: username,
				IsActive: &isActive,
			},
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^UPDATE \"users\" SET").
					WithArgs(email, isActive, password, sqlmock.AnyArg(), userID).
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

			res, err := repo.Update(ctx, userID, tc.input)

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

			assert.Equal(t, res.ID, userID)
		})
	}
}

func TestUserRepository_IsAdminAccountExists(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	userID := uuid.New()

	testCases := []struct {
		name                 string
		wantErr              bool
		expectedErr          error
		expectedOutput       bool
		expectedFunctionCall func()
	}{
		{
			name:           "success",
			wantErr:        false,
			expectedOutput: true,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(model.RolesAdministrator, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
			},
		},
		{
			name:        "error - unknown just pass the error to the caller",
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(model.RolesAdministrator, 1).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:           "data not found on db",
			wantErr:        false,
			expectedOutput: false,
			expectedErr:    repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(model.RolesAdministrator, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.IsAdminAccountExists(ctx)

			if tc.wantErr {
				require.Error(t, err)

				if tc.expectedErr != nil {
					assert.Equal(t, tc.expectedErr, err)
				}

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, res)

			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}

func TestUserRepository_Search(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	role := model.RolesAdministrator
	limit := 100
	offset := 10

	testCases := []struct {
		name                 string
		input                usecase.RepoSearchUserInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutputLen    int
	}{
		{
			name:    "database returning unexpected error",
			wantErr: true,
			input: usecase.RepoSearchUserInput{
				Role:   role,
				Limit:  limit,
				Offset: offset,
			},
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(role, limit, offset).
					WillReturnError(assert.AnError)
			},
		},
		{
			name:    "no result found must return not found error",
			wantErr: true,
			input: usecase.RepoSearchUserInput{
				Role:   role,
				Limit:  limit,
				Offset: offset,
			},
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(role, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			name:    "ok",
			wantErr: false,
			input: usecase.RepoSearchUserInput{
				Role:   role,
				Limit:  limit,
				Offset: offset,
			},
			expectedOutputLen: 2,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(`^SELECT .+ FROM "users"`).
					WithArgs(role, limit, offset).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()))
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
