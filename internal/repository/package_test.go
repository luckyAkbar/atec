package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redsync/redsync/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	mockdb "github.com/luckyAkbar/atec/mocks/internal_/db"
)

func TestPackageRepository_Create(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	dbMock := kit.DBmock
	cacherMock := mockdb.NewCacheKeeperIface(t)

	repo := repository.NewPackageRepo(kit.DB, cacherMock)

	userID := uuid.New()
	generatedPackageID := uuid.New()
	name := "name"
	indicationCategories := model.IndicationCategories{
		{
			MinimumScore: 1,
			MaximumScore: 1,
			Name:         "category 1",
			Detail:       "detail 1",
		},
	}
	questionnaire := model.Questionnaire{
		1: model.ChecklistGroup{
			CustomName: "custom name",
			Questions:  []string{"question 1", "question 2"},
			Options: []model.AnswerOption{
				{
					ID:          1,
					Description: "option 1",
					Score:       1,
				},
			},
		},
	}
	imageAttributeKey := model.ImageResultAttributeKey{
		Title:       "title",
		Total:       "total",
		Indication:  "indication",
		ResultID:    "result_id",
		SubmittedAt: time.Now().String(),
	}
	cacheKey := repository.CacheKeyForPackage(model.Package{
		ID: generatedPackageID,
	})

	testCases := []struct {
		name                 string
		input                usecase.RepoCreatePackageInput
		wantErr              bool
		expectedErr          error
		expectedFunctionCall func()
		expectedOutput       *model.Package
		txControllers        []*gorm.DB
	}{
		{
			name: "success",
			input: usecase.RepoCreatePackageInput{
				UserID:                  userID,
				PackageName:             name,
				Questionnaire:           questionnaire,
				IndicationCategories:    indicationCategories,
				ImageResultAttributeKey: imageAttributeKey,
			},
			wantErr: false,
			expectedFunctionCall: func() {
				dbMock.ExpectBegin()

				dbMock.ExpectQuery("INSERT INTO \"packages\"").
					WithArgs(
						userID,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						name,
						false,
						false,
						sqlmock.AnyArg(),
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(generatedPackageID))

				dbMock.ExpectCommit()

				cacherMock.EXPECT().AcquireLock(cacheKey).Return(&redsync.Mutex{}, nil).Once()

				cacherMock.EXPECT().SetJSON(ctx, cacheKey, gomock.Any(), sqlmock.AnyArg()).Return(nil).Once()
			},
			expectedOutput: &model.Package{
				ID: generatedPackageID,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := repo.Create(ctx, tc.input, tc.txControllers...)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.NotNil(t, res)

				assert.Equal(t, res.ID, tc.expectedOutput.ID)

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
