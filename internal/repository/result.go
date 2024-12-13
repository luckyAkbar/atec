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
	FindByID(ctx context.Context, id uuid.UUID) (*model.Result, error)
	Search(ctx context.Context, input SearchResultInput) ([]model.Result, error)
	FindAllUserHistory(ctx context.Context, input FindAllUserHistoryInput) ([]model.Result, error)
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

// FindByID find exactly one record from results table with matching id
func (r *ResultRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Result, error) {
	result := &model.Result{}

	err := r.db.WithContext(ctx).Take(result, "id = ?", id).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return result, nil
	}
}

// SearchResultInput search result input
type SearchResultInput struct {
	ID        uuid.UUID
	PackageID uuid.UUID
	ChildID   uuid.UUID
	CreatedBy uuid.UUID
	Limit     int
	Offset    int
}

func (sri SearchResultInput) toSearchFields(cursor *gorm.DB) *gorm.DB {
	if sri.ID != uuid.Nil {
		cursor = cursor.Where("id = ?", sri.ID)
	}

	if sri.PackageID != uuid.Nil {
		cursor = cursor.Where("package_id = ?", sri.PackageID)
	}

	if sri.ChildID != uuid.Nil {
		cursor = cursor.Where("child_id = ?", sri.ChildID)
	}

	if sri.CreatedBy != uuid.Nil {
		cursor = cursor.Where("created_by = ?", sri.CreatedBy)
	}

	if sri.Limit > 0 {
		cursor = cursor.Limit(sri.Limit)
	}

	if sri.Offset > 0 {
		cursor = cursor.Offset(sri.Offset)
	}

	return cursor
}

// Search search results based on provided search parameters
func (r *ResultRepository) Search(ctx context.Context, input SearchResultInput) ([]model.Result, error) {
	cursor := r.db.WithContext(ctx)
	cursor = input.toSearchFields(cursor)

	results := []model.Result{}

	if err := cursor.Find(&results).Error; err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return results, nil
}

// FindAllUserHistoryInput input
type FindAllUserHistoryInput struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

// FindAllUserHistory find the result made by the userID or the child of the userID
func (r *ResultRepository) FindAllUserHistory(ctx context.Context, input FindAllUserHistoryInput) ([]model.Result, error) {
	results := []model.Result{}

	err := r.db.WithContext(ctx).Where("created_by = ?", input.UserID).
		Or(
			"child_id IN (?)",
			r.db.Model(&model.Child{}).
				Select("id").Where("parent_user_id = ?", input.UserID),
		).Limit(input.Limit).Offset(input.Offset).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return results, nil
}
