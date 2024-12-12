package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
)

// ResultRepository result repository
type ResultRepository struct {
	db *gorm.DB
}

// ResultRepositoryIface result repository interface
type ResultRepositoryIface interface {
	Create(ctx context.Context, input CreateResultInput) (*model.Result, error)
}

// NewResultRepository create new instance of ResultRepository
func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{
		db: db,
	}
}

// CreateResultInput create result input
type CreateResultInput struct {
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Answer    model.AnswerDetail
	Result    model.ResultDetail
}

// Create insert new record of results to the database
func (r *ResultRepository) Create(ctx context.Context, input CreateResultInput) (*model.Result, error) {
	result := &model.Result{
		PackageID: input.PackageID,
		ChildID:   input.ChildID,
		CreatedBy: input.CreatedBy,
		Answer:    input.Answer,
		Result:    input.Result,
	}

	if err := r.db.WithContext(ctx).Create(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
