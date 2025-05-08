package repository_test

import (
	"context"
	"testing"

	"github.com/luckyAkbar/atec/internal/repository"
	"github.com/luckyAkbar/atec/internal/usecase"
	"github.com/stretchr/testify/assert"
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

func TestUserRepositoryUCAdapter_Create(t *testing.T) {
	ctx := context.Background()
	kit, closer := InitializeRepoTestKit(t)

	defer closer()

	// dbMock := kit.DBmock
	repo := repository.NewUserRepository(kit.DB)

	adapter := repository.NewUserRepositoryUCAdapter(repo)

	t.Run("the tx controller is not gorm", func(t *testing.T) {
		_, err := adapter.Create(ctx, usecase.RepoCreateUserInput{}, 1)
		assert.Error(t, err)
	})
}
