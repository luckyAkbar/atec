package repository

import (
	"context"
	"fmt"
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
	Search(ctx context.Context, input SearchChildInput) ([]model.Child, error)
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

// SearchChildInput input to search child data. everything marked as pointer to a datatype means it is optional
type SearchChildInput struct {
	ParentUserID *uuid.UUID
	Name         *string
	Gender       *bool
	Limit        int
	Offset       int
}

func (sci SearchChildInput) buildSearchField(cursor *gorm.DB) *gorm.DB {
	if sci.ParentUserID != nil {
		cursor = cursor.Where("parent_user_id = ?", sci.ParentUserID)
	}

	if sci.Name != nil {
		cursor = cursor.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", *sci.Name))
	}

	if sci.Gender != nil {
		cursor = cursor.Where("gender = ?", *sci.Gender)
	}

	if sci.Limit > 0 {
		cursor = cursor.Limit(sci.Limit)
	}

	if sci.Offset > 0 {
		cursor = cursor.Offset(sci.Offset)
	}

	return cursor
}

// Search search children data based on provided search parameters
func (r *ChildRepository) Search(ctx context.Context, input SearchChildInput) ([]model.Child, error) {
	children := []model.Child{}

	cursor := r.db.WithContext(ctx)
	query := input.buildSearchField(cursor)

	err := query.Find(&children).Error
	if err != nil {
		return nil, err
	}

	if len(children) == 0 {
		return nil, ErrNotFound
	}

	return children, nil
}
