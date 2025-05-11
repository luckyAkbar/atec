package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestChildRepository_Create(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewChildRepository(kit.DB)

	parentUserID := uuid.New()
	dateOfBirth := time.Now().Add(-time.Hour * 24 * 365 * 5)
	gender := false
	name := "Jane Doe"

	dbGeneratedUUID := uuid.New()

	testCases := []struct {
		name                 string
		input                usecase.RepoCreateChildInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutput       *model.Child
	}{
		{
			name: "success",
			input: usecase.RepoCreateChildInput{
				ParentUserID: parentUserID,
				DateOfBirth:  dateOfBirth,
				Gender:       gender,
				Name:         name,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^INSERT INTO \"children\"").
					WithArgs(parentUserID, dateOfBirth, gender, name, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(dbGeneratedUUID))

				dbMock.ExpectCommit()
			},
			expectedOutput: &model.Child{
				ID:           dbGeneratedUUID,
				ParentUserID: parentUserID,
				DateOfBirth:  dateOfBirth,
				Gender:       gender,
				Name:         name,
			},
		},
		{
			name: "error",
			input: usecase.RepoCreateChildInput{
				ParentUserID: parentUserID,
				DateOfBirth:  dateOfBirth,
				Gender:       gender,
				Name:         name,
			},
			wantErr: true,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^INSERT INTO \"children\"").
					WithArgs(parentUserID, dateOfBirth, gender, name, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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

				assert.Equal(t, res.ID, tc.expectedOutput.ID)
				assert.Equal(t, res.ParentUserID, tc.expectedOutput.ParentUserID)
				assert.Equal(t, res.DateOfBirth, tc.expectedOutput.DateOfBirth)
				assert.Equal(t, res.Gender, tc.expectedOutput.Gender)
				assert.Equal(t, res.Name, tc.expectedOutput.Name)

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

func TestChildRepository_Update(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewChildRepository(kit.DB)

	childID := uuid.New()
	parentUserID := uuid.New()
	dateOfBirth := time.Now().Add(-time.Hour * 24 * 365 * 5)
	gender := false
	name := "Jane Doe"

	testCases := []struct {
		name                 string
		input                usecase.RepoUpdateChildInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutput       *model.Child
	}{
		{
			name: "success",
			input: usecase.RepoUpdateChildInput{
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^UPDATE \"children\" SET").
					WithArgs(dateOfBirth, gender, name, sqlmock.AnyArg(), childID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "parent_user_id"}).AddRow(childID, parentUserID))

				dbMock.ExpectCommit()
			},
			expectedOutput: &model.Child{
				ID:           childID,
				ParentUserID: parentUserID,
				Gender:       gender,
				Name:         name,
				DateOfBirth:  dateOfBirth,
			},
		},
		{
			name: "error",
			input: usecase.RepoUpdateChildInput{
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("^UPDATE \"children\" SET").
					WithArgs(dateOfBirth, gender, name, sqlmock.AnyArg(), childID).
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

			res, err := repo.Update(ctx, childID, tc.input)

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

			assert.Equal(t, res.ID, tc.expectedOutput.ID)
			assert.Equal(t, res.ParentUserID, tc.expectedOutput.ParentUserID)
			assert.Equal(t, res.DateOfBirth, tc.expectedOutput.DateOfBirth)
			assert.Equal(t, res.Gender, tc.expectedOutput.Gender)
			assert.Equal(t, res.Name, tc.expectedOutput.Name)
		})
	}
}

func TestChildRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewChildRepository(kit.DB)

	childID := uuid.New()
	parentUserID := uuid.New()
	randomError := errors.New("random error")

	testCases := []struct {
		name                 string
		input                uuid.UUID
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutput       *model.Child
	}{
		{
			name:    "success",
			input:   childID,
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "children" WHERE id = $1 AND "children"."deleted_at" IS NULL LIMIT $2`)).
					WithArgs(childID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "parent_user_id"}).AddRow(childID, parentUserID))
			},
			expectedOutput: &model.Child{
				ID:           childID,
				ParentUserID: parentUserID,
			},
		},
		{
			name:        "error - unknown just pass the error to the caller",
			input:       childID,
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "children" WHERE id = $1 AND "children"."deleted_at" IS NULL LIMIT $2`)).
					WithArgs(childID, 1).
					WillReturnError(assert.AnError)
			},
			expectedOutput: &model.Child{
				ID:           childID,
				ParentUserID: parentUserID,
			},
		},
		{
			name:        "error - unknown just pass the error to the calle - 2",
			input:       childID,
			wantErr:     true,
			expectedErr: randomError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "children" WHERE id = $1 AND "children"."deleted_at" IS NULL LIMIT $2`)).
					WithArgs(childID, 1).
					WillReturnError(randomError)
			},
			expectedOutput: &model.Child{
				ID:           childID,
				ParentUserID: parentUserID,
			},
		},
		{
			name:        "data not found on db",
			input:       childID,
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "children" WHERE id = $1 AND "children"."deleted_at" IS NULL LIMIT $2`)).
					WithArgs(childID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedOutput: &model.Child{
				ID:           childID,
				ParentUserID: parentUserID,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.FindByID(ctx, tc.input)

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

			assert.Equal(t, res.ID, tc.expectedOutput.ID)
			assert.Equal(t, res.ParentUserID, tc.expectedOutput.ParentUserID)
		})
	}
}

func TestChildRepository_Search(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewChildRepository(kit.DB)

	childID := uuid.New()
	parentUserID := uuid.New()
	name := "mary currie"
	gender := true
	limit := 111
	offset := 222

	testCases := []struct {
		name                 string
		input                usecase.RepoSearchChildInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutputLen    int
	}{
		{
			name: "success 1",
			input: usecase.RepoSearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        limit,
				Offset:       offset,
			},
			wantErr:           false,
			expectedOutputLen: 1,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery("SELECT .+ FROM \"children\"").
					WithArgs(parentUserID, fmt.Sprintf("%%%s%%", name), gender, limit, offset).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id"},
						).AddRow(childID),
					)
			},
		},
		{
			name: "success 5",
			input: usecase.RepoSearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        limit,
				Offset:       offset,
			},
			wantErr:           false,
			expectedOutputLen: 5,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery("SELECT .+ FROM \"children\"").
					WithArgs(parentUserID, fmt.Sprintf("%%%s%%", name), gender, limit, offset).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id"},
						).AddRow(childID).AddRow(childID).
							AddRow(childID).AddRow(childID).
							AddRow(childID),
					)
			},
		},
		{
			name: "error db",
			input: usecase.RepoSearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        limit,
				Offset:       offset,
			},
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery("SELECT .+ FROM \"children\"").
					WithArgs(parentUserID, fmt.Sprintf("%%%s%%", name), gender, limit, offset).
					WillReturnError(assert.AnError)
			},
		},
		{
			name: "no rows returned must trigger not found error",
			input: usecase.RepoSearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        limit,
				Offset:       offset,
			},
			wantErr:     true,
			expectedErr: repository.ErrNotFound,
			expectedFunctionCall: func() {
				dbMock.ExpectQuery("SELECT .+ FROM \"children\"").
					WithArgs(parentUserID, fmt.Sprintf("%%%s%%", name), gender, limit, offset).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id"},
						),
					)
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

func TestChildRepository_DeleteAllUserChildren(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	repo := repository.NewChildRepository(kit.DB)

	userID := uuid.New()

	testCases := []struct {
		name                 string
		input                usecase.RepoDeleteAllUserChildrenInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
	}{
		{
			name: "success",
			input: usecase.RepoDeleteAllUserChildrenInput{
				UserID:     userID,
				HardDelete: true,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectExec("^DELETE FROM \"children\"").
					WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				dbMock.ExpectCommit()
			},
		},
		{
			name: "error - hard delete",
			input: usecase.RepoDeleteAllUserChildrenInput{
				UserID:     userID,
				HardDelete: true,
			},
			wantErr:     true,
			expectedErr: assert.AnError,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectExec("^DELETE FROM \"children\"").
					WithArgs(userID).
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

			err := repo.DeleteAllUserChildren(ctx, tc.input, kit.DB)

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
