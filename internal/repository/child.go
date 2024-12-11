package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
)

// ChildRepository child repository
type ChildRepository struct {
	db *gorm.DB
}

// ChildRepositoryIface interface
type ChildRepositoryIface interface {
	Create(ctx context.Context, input CreateChildInput) (*model.Child, error)
}

// NewChildRepository create new instance of ChildRepository
func NewChildRepository(db *gorm.DB) *ChildRepository {
	return &ChildRepository{
		db: db,
	}
}

// CreateChildInput input
type CreateChildInput struct {
	ParentUserID uuid.UUID
	DateOfBirth  time.Time
	Gender       bool
	Name         string
}

// Create create a new record on children table
func (r *ChildRepository) Create(ctx context.Context, input CreateChildInput) (*model.Child, error) {
	child := &model.Child{
		ParentUserID: input.ParentUserID,
		DateOfBirth:  input.DateOfBirth,
		Gender:       input.Gender,
		Name:         input.Name,
	}

	err := r.db.WithContext(ctx).Create(child).Error
	if err != nil {
		return nil, err
	}

	return child, nil
}
