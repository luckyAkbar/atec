package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"gorm.io/gorm"
)

// ResultRepository result repository
type ResultRepository struct {
	db *gorm.DB
}

// NewResultRepository create new instance of ResultRepository
func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{
		db: db,
	}
}

// Create insert new record of results to the database
func (r *ResultRepository) Create(ctx context.Context, input usecase.RepoCreateResultInput) (*model.Result, error) {
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

func searchResultInputoSearchFields(cursor *gorm.DB, sri usecase.RepoSearchResultInput) *gorm.DB {
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
func (r *ResultRepository) Search(ctx context.Context, input usecase.RepoSearchResultInput) ([]model.Result, error) {
	cursor := r.db.WithContext(ctx)
	cursor = searchResultInputoSearchFields(cursor, input)

	results := []model.Result{}

	if err := cursor.Order("created_at ASC").Find(&results).Error; err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return results, nil
}

// FindAllUserHistory find the result made by the userID or the child of the userID
func (r *ResultRepository) FindAllUserHistory(ctx context.Context, input usecase.RepoFindAllUserHistoryInput) ([]model.Result, error) {
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

// DeleteAllUserResults delete all results made by the userID or submitted for child of the userID
func (r *ResultRepository) DeleteAllUserResults(ctx context.Context, input usecase.RepoDeleteAllUserResultsInput, txController ...*gorm.DB) error {
	tx := r.db
	if len(txController) > 0 {
		tx = txController[0]
	}

	if input.HardDelete {
		tx = tx.Unscoped()
	}

	err := tx.WithContext(ctx).Where("created_by = ?", input.UserID).
		Or(
			"child_id IN (?)",
			r.db.Model(&model.Child{}).
				Select("id").Where("parent_user_id = ?", input.UserID),
		).Delete(&model.Result{}).Error

	if err != nil {
		return err
	}

	return nil
}
