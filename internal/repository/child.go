package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ChildRepository child repository
type ChildRepository struct {
	db *gorm.DB
}

// ChildRepositoryIface interface
type ChildRepositoryIface interface {
	Create(ctx context.Context, input CreateChildInput) (*model.Child, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateChildInput) (*model.Child, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Child, error)
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

// UpdateChildInput input
type UpdateChildInput struct {
	DateOfBirth *time.Time
	Gender      *bool
	Name        *string
}

// ToUpdateFields converts UpdateChildInput to dynamic gorm update fields
func (uci UpdateChildInput) ToUpdateFields() map[string]interface{} {
	fields := map[string]interface{}{}

	if uci.DateOfBirth != nil {
		fields["date_of_birth"] = *uci.DateOfBirth
	}

	if uci.Gender != nil {
		fields["gender"] = *uci.Gender
	}

	if uci.Name != nil {
		fields["name"] = *uci.Name
	}

	return fields
}

// Update update child records on database based on id
func (r *ChildRepository) Update(ctx context.Context, id uuid.UUID, input UpdateChildInput) (*model.Child, error) {
	child := &model.Child{}

	err := r.db.WithContext(ctx).Model(child).
		Clauses(clause.Returning{}).Where("id = ?", id).
		Updates(input.ToUpdateFields()).Error

	if err != nil {
		return nil, err
	}

	return child, nil
}

// FindByID find child by id
func (r *ChildRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Child, error) {
	child := &model.Child{}

	err := r.db.WithContext(ctx).Take(child, "id = ?", id).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return child, nil
	}
}
