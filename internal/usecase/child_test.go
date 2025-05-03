package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mockUsecase "github.com/luckyAkbar/atec/mocks/internal_/usecase"
)

func TestRegisterChildInputValidate(t *testing.T) {
	t.Run("Valid RegisterChildInput", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			DateOfBirth: time.Now(),
			Gender:      true,
			Name:        "John Doe",
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Invalid RegisterChildInput - Missing DateOfBirth", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			Gender: true,
			Name:   "John Doe",
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid RegisterChildInput - Missing Name", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			DateOfBirth: time.Now(),
			Gender:      true,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Invalid RegisterChildInput - Missing DateOfBirth and Name", func(t *testing.T) {
		input := usecase.RegisterChildInput{
			Gender: true,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})
}

func TestUpdateChildInputValidate(t *testing.T) {
	name := "Fulan"

	t.Run("Valid UpdateChildInput", func(t *testing.T) {
		dateOfBirth := time.Now()
		gender := true

		input := usecase.UpdateChildInput{
			ChildID:     uuid.New(),
			DateOfBirth: &dateOfBirth,
			Gender:      &gender,
			Name:        &name,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Invalid UpdateChildInput - Missing ChildID", func(t *testing.T) {
		dateOfBirth := time.Now()
		gender := true
		input := usecase.UpdateChildInput{
			DateOfBirth: &dateOfBirth,
			Gender:      &gender,
			Name:        &name,
		}

		if err := input.Validate(); err == nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("Valid UpdateChildInput - Only ChildID", func(t *testing.T) {
		input := usecase.UpdateChildInput{
			ChildID: uuid.New(),
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Valid UpdateChildInput - ChildID and DateOfBirth", func(t *testing.T) {
		dateOfBirth := time.Now()
		input := usecase.UpdateChildInput{
			ChildID:     uuid.New(),
			DateOfBirth: &dateOfBirth,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Valid UpdateChildInput - ChildID and Gender", func(t *testing.T) {
		gender := true
		input := usecase.UpdateChildInput{
			ChildID: uuid.New(),
			Gender:  &gender,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})

	t.Run("Valid UpdateChildInput - ChildID and Name", func(t *testing.T) {
		input := usecase.UpdateChildInput{
			ChildID: uuid.New(),
			Name:    &name,
		}

		if err := input.Validate(); err != nil {
			t.Errorf("Validate() error = %v, wantErr %v", err, false)
		}
	})
}

func TestChildUsecase_Register(t *testing.T) {
	ctx := context.Background()

	user := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesTherapist,
	}

	userCtx := model.SetUserToCtx(ctx, user)

	mockChildRepo := mockUsecase.NewChildRepository(t)

	uc := usecase.NewChildUsecase(mockChildRepo, nil)

	childID := uuid.New()
	dateOfBirth := time.Now()
	gender := false
	name := "John Doe"

	testCases := []struct {
		name                 string
		input                usecase.RegisterChildInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.RegisterChildOutput
		expectedFunctionCall func()
	}{
		{
			name:        "empty user in ctx",
			input:       usecase.RegisterChildInput{},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name:        "invalid input",
			input:       usecase.RegisterChildInput{},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "repository failed to create child",
			input: usecase.RegisterChildInput{
				DateOfBirth: dateOfBirth,
				Gender:      gender,
				Name:        name,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Create(userCtx, usecase.RepoCreateChildInput{
					DateOfBirth:  dateOfBirth,
					ParentUserID: user.ID,
					Gender:       gender,
					Name:         name,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.RegisterChildInput{
				DateOfBirth: dateOfBirth,
				Gender:      gender,
				Name:        name,
			},
			ctx:     userCtx,
			wantErr: false,
			expectedOutput: &usecase.RegisterChildOutput{
				ID: childID,
			},
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Create(userCtx, usecase.RepoCreateChildInput{
					DateOfBirth:  dateOfBirth,
					ParentUserID: user.ID,
					Gender:       gender,
					Name:         name,
				}).Return(&model.Child{ID: childID}, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.Register(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, res)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestChildUsecase_Update(t *testing.T) {
	ctx := context.Background()

	user := model.AuthUser{
		ID:   uuid.New(),
		Role: model.RolesTherapist,
	}

	userCtx := model.SetUserToCtx(ctx, user)

	mockChildRepo := mockUsecase.NewChildRepository(t)

	uc := usecase.NewChildUsecase(mockChildRepo, nil)

	childID := uuid.New()
	dateOfBirth := time.Now()
	gender := false
	name := "John Doe"
	child := &model.Child{
		ID:           childID,
		ParentUserID: user.ID,
	}

	testCases := []struct {
		name                 string
		input                usecase.UpdateChildInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutput       *usecase.UpdateChildOutput
		expectedFunctionCall func()
	}{
		{
			name:        "empty user in ctx",
			input:       usecase.UpdateChildInput{},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name:        "invalid input",
			input:       usecase.UpdateChildInput{},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "repository failed to find user by id",
			input: usecase.UpdateChildInput{
				ChildID:     childID,
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "repository return not found",
			input: usecase.UpdateChildInput{
				ChildID:     childID,
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "only the parent able to update children data",
			input: usecase.UpdateChildInput{
				ChildID:     childID,
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrForbidden,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(&model.Child{ParentUserID: uuid.New()}, nil).Once()
			},
		},
		{
			name: "repository failed to update child",
			input: usecase.UpdateChildInput{
				ChildID:     childID,
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(child, nil).Once()
				mockChildRepo.EXPECT().Update(userCtx, childID, usecase.RepoUpdateChildInput{
					DateOfBirth: &dateOfBirth,
					Gender:      &gender,
					Name:        &name,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "ok",
			input: usecase.UpdateChildInput{
				ChildID:     childID,
				DateOfBirth: &dateOfBirth,
				Gender:      &gender,
				Name:        &name,
			},
			ctx:     userCtx,
			wantErr: false,
			expectedOutput: &usecase.UpdateChildOutput{
				Message: "ok",
			},
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(child, nil).Once()
				mockChildRepo.EXPECT().Update(userCtx, childID, usecase.RepoUpdateChildInput{
					DateOfBirth: &dateOfBirth,
					Gender:      &gender,
					Name:        &name,
				}).Return(child, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.Update(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, res)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestAuthUsecase_GetRegisteredChildren(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	user := model.AuthUser{
		ID:   userID,
		Role: model.RolesTherapist,
	}

	userCtx := model.SetUserToCtx(ctx, user)

	mockChildRepo := mockUsecase.NewChildRepository(t)

	uc := usecase.NewChildUsecase(mockChildRepo, nil)

	children := []model.Child{
		{
			ID:           uuid.New(),
			ParentUserID: userID,
			DateOfBirth:  time.Now(),
			Gender:       true,
			Name:         "John Does Nothing",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	testCases := []struct {
		name                 string
		input                usecase.GetRegisteredChildrenInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutputLen    int
		expectedFunctionCall func()
	}{
		{
			name: "empty user in ctx",
			input: usecase.GetRegisteredChildrenInput{
				Limit:  10,
				Offset: 1,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name: "invalid input limit undefined",
			input: usecase.GetRegisteredChildrenInput{
				Offset: 1,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid input offset negative",
			input: usecase.GetRegisteredChildrenInput{
				Limit:  10,
				Offset: -1,
			},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "repository search return error",
			input: usecase.GetRegisteredChildrenInput{
				Limit:  20,
				Offset: 1,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Search(userCtx, usecase.RepoSearchChildInput{
					ParentUserID: &userID,
					Limit:        20,
					Offset:       1,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "repository return not found error",
			input: usecase.GetRegisteredChildrenInput{
				Limit:  20,
				Offset: 1,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Search(userCtx, usecase.RepoSearchChildInput{
					ParentUserID: &userID,
					Limit:        20,
					Offset:       1,
				}).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "ok",
			input: usecase.GetRegisteredChildrenInput{
				Limit:  20,
				Offset: 1,
			},
			ctx:               userCtx,
			wantErr:           false,
			expectedOutputLen: 1,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Search(userCtx, usecase.RepoSearchChildInput{
					ParentUserID: &userID,
					Limit:        20,
					Offset:       1,
				}).Return(children, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.GetRegisteredChildren(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Len(t, res, tc.expectedOutputLen)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestChildUsecase_Search(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	user := model.AuthUser{
		ID:   userID,
		Role: model.RolesTherapist,
	}

	userCtx := model.SetUserToCtx(ctx, user)

	mockChildRepo := mockUsecase.NewChildRepository(t)

	uc := usecase.NewChildUsecase(mockChildRepo, nil)

	parentUserID := uuid.New()
	name := "Jane Doe"
	gender := false

	children := []model.Child{
		{
			ID:           uuid.New(),
			ParentUserID: userID,
			DateOfBirth:  time.Now(),
			Gender:       gender,
			Name:         name,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	testCases := []struct {
		name                 string
		input                usecase.SearchChildInput
		ctx                  context.Context
		wantErr              bool
		expectedErr          error
		expectedOutputLen    int
		expectedFunctionCall func()
	}{
		{
			name:        "invalid input limit undefined",
			input:       usecase.SearchChildInput{},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "invalid input: negative value on offset",
			input: usecase.SearchChildInput{
				Limit:  10,
				Offset: -1,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "search function returning unexpected error",
			input: usecase.SearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        10,
				Offset:       1,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Search(userCtx, usecase.RepoSearchChildInput{
					ParentUserID: &parentUserID,
					Name:         &name,
					Gender:       &gender,
					Limit:        10,
					Offset:       1,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "search function returning not found",
			input: usecase.SearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        10,
				Offset:       1,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Search(userCtx, usecase.RepoSearchChildInput{
					ParentUserID: &parentUserID,
					Name:         &name,
					Gender:       &gender,
					Limit:        10,
					Offset:       1,
				}).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "ok",
			input: usecase.SearchChildInput{
				ParentUserID: &parentUserID,
				Name:         &name,
				Gender:       &gender,
				Limit:        10,
				Offset:       1,
			},
			ctx:               userCtx,
			wantErr:           false,
			expectedOutputLen: 1,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().Search(userCtx, usecase.RepoSearchChildInput{
					ParentUserID: &parentUserID,
					Name:         &name,
					Gender:       &gender,
					Limit:        10,
					Offset:       1,
				}).Return(children, nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.Search(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Len(t, res, tc.expectedOutputLen)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}

func TestChildUseccase_HandleGetStatistic(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	user := model.AuthUser{
		ID:   userID,
		Role: model.RolesParent,
	}

	therapistID := uuid.New()
	therapist := model.AuthUser{
		ID:   therapistID,
		Role: model.RolesTherapist,
	}

	nonParentID := uuid.New()
	nonParent := model.AuthUser{
		ID:   nonParentID,
		Role: model.RolesParent,
	}

	userCtx := model.SetUserToCtx(ctx, user)
	therapistCtx := model.SetUserToCtx(ctx, therapist)
	nonParentCtx := model.SetUserToCtx(ctx, nonParent)

	mockChildRepo := mockUsecase.NewChildRepository(t)
	mockResultRepo := mockUsecase.NewResultRepository(t)

	uc := usecase.NewChildUsecase(mockChildRepo, mockResultRepo)

	childID := uuid.New()
	child := &model.Child{
		ID:           childID,
		ParentUserID: userID,
	}

	batchSize := 100

	genResult := func(num int) []model.Result {
		var results []model.Result

		for range num {
			results = append(results, model.Result{
				ID:        uuid.New(),
				ChildID:   childID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}

		return results
	}

	testCases := []struct {
		name                 string
		input                usecase.GetStatisticInput
		wantErr              bool
		ctx                  context.Context
		expectedErr          error
		expectedOutputLen    int
		expectedFunctionCall func()
	}{
		{
			name:        "requester in context is empty",
			input:       usecase.GetStatisticInput{},
			ctx:         ctx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
		},
		{
			name:        "invalid input childID is empty",
			input:       usecase.GetStatisticInput{},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrBadRequest,
		},
		{
			name: "find child by id failed",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "find child return not found",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "the requester wasn't the parent nor the admin",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:         nonParentCtx,
			wantErr:     true,
			expectedErr: usecase.ErrUnauthorized,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(nonParentCtx, childID).Return(child, nil).Once()
			},
		},
		{
			name: "admin: search from repository returning unexpected error",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:         therapistCtx,
			wantErr:     true,
			expectedErr: usecase.ErrInternal,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(therapistCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Search(therapistCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  0,
				}).Return(nil, assert.AnError).Once()
			},
		},
		{
			name: "parent: when child has no result yet, return not found",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:         userCtx,
			wantErr:     true,
			expectedErr: usecase.ErrNotFound,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Search(userCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  0,
				}).Return(nil, usecase.ErrRepoNotFound).Once()
			},
		},
		{
			name: "parent: 1 query only if result less than batch size",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:               userCtx,
			wantErr:           false,
			expectedOutputLen: 99,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Search(userCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  0,
				}).Return(genResult(99), nil).Once()
			},
		},
		{
			name: "parent: more query needed if result returned is equal to batch size",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:               userCtx,
			wantErr:           false,
			expectedOutputLen: 199,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(userCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Search(userCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  0,
				}).Return(genResult(100), nil).Once()
				mockResultRepo.EXPECT().Search(userCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  100,
				}).Return(genResult(99), nil).Once()
			},
		},
		{
			name: "admin: more query needed if result returned is equal to batch size",
			input: usecase.GetStatisticInput{
				ChildID: childID,
			},
			ctx:               therapistCtx,
			wantErr:           false,
			expectedOutputLen: 299,
			expectedFunctionCall: func() {
				mockChildRepo.EXPECT().FindByID(therapistCtx, childID).Return(child, nil).Once()
				mockResultRepo.EXPECT().Search(therapistCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  0,
				}).Return(genResult(100), nil).Once()
				mockResultRepo.EXPECT().Search(therapistCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  100,
				}).Return(genResult(100), nil).Once()
				mockResultRepo.EXPECT().Search(therapistCtx, usecase.RepoSearchResultInput{
					ChildID: childID,
					Limit:   batchSize,
					Offset:  200,
				}).Return(genResult(99), nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedFunctionCall != nil {
				tc.expectedFunctionCall()
			}

			res, err := uc.HandleGetStatistic(tc.ctx, tc.input)

			if !tc.wantErr {
				require.NoError(t, err)
				assert.Len(t, res.Statistic, tc.expectedOutputLen)

				return
			}

			require.Error(t, err)

			switch e := err.(type) {
			default:
				t.Errorf("expecting usecase error but got %T", err)
			case usecase.UsecaseError:
				assert.Equal(t, tc.expectedErr, e.ErrType)
			}
		})
	}
}
