// Package repository contains all the functions necessary to interact with databases
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"github.com/luckyAkbar/atec/internal/usecase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepository is an instance containing functions to interact specifically to users database
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository create a new instance of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByEmail find exactly one record from users table with matching email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Take(user, "email = ?", email).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return user, nil
	}
}

// Create insert a new records to users table
func (r *UserRepository) Create(ctx context.Context, input usecase.RepoCreateUserInput, txController ...*gorm.DB) (*model.User, error) {
	tx := r.db.WithContext(ctx)
	if len(txController) > 0 {
		tx = txController[0]
	}

	user := &model.User{
		Email:       input.Email,
		Password:    input.Password,
		IsActive:    input.IsActive,
		Roles:       input.Roles,
		Username:    input.Username,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
	}

	err := tx.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID find exacly one record from users table with matching id
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Take(user, "id = ?", id).Error
	switch err {
	default:
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return user, nil
	}
}

// updateUserInputToUpdatedFields helper function to translate update options to gorm dynamic fields update
func updateUserInputToUpdatedFields(uui usecase.RepoUpdateUserInput) map[string]interface{} {
	fields := map[string]interface{}{}

	if uui.Email != "" {
		fields["email"] = uui.Email
	}

	if uui.Password != "" {
		fields["password"] = uui.Password
	}

	if uui.IsActive != nil {
		fields["is_active"] = *uui.IsActive
	}

	return fields
}

// Update update a users record by its id
func (r *UserRepository) Update(ctx context.Context, userID uuid.UUID, input usecase.RepoUpdateUserInput) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Model(user).
		Clauses(clause.Returning{}).Where("id = ?", userID).
		Updates(updateUserInputToUpdatedFields(input)).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// IsAdminAccountExists check if there is an admin account in the database.
// If atleast 1 admin account is found, will return true, otherwise false
func (r *UserRepository) IsAdminAccountExists(ctx context.Context) (bool, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).First(user, "roles = ?", model.RolesAdministrator).Error
	switch err {
	default:
		return false, err
	case gorm.ErrRecordNotFound:
		return false, nil
	case nil:
		return true, nil
	}
}

func toSearchFields(cursor *gorm.DB, sui usecase.RepoSearchUserInput) *gorm.DB {
	if sui.Role != "" {
		cursor = cursor.Where("roles = ?", sui.Role)
	}

	if sui.Limit > 0 {
		cursor = cursor.Limit(sui.Limit)
	}

	if sui.Offset > 0 {
		cursor = cursor.Offset(sui.Offset)
	}

	return cursor
}

// Search search users with given options
func (r *UserRepository) Search(ctx context.Context, input usecase.RepoSearchUserInput) ([]model.User, error) {
	users := []model.User{}
	cursor := r.db.WithContext(ctx)
	cursor = toSearchFields(cursor, input)

	if err := cursor.Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return users, nil
}

// GetUsersByRoles returns all users with the specified role
func (r *UserRepository) GetUsersByRoles(ctx context.Context, roles model.Roles) ([]model.User, error) {
	users := []model.User{}
	if err := r.db.WithContext(ctx).Where("roles = ?", roles).Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return users, nil
}

// DeleteByID delete a user by its id
func (r *UserRepository) DeleteByID(ctx context.Context, input usecase.RepoDeleteUserByIDInput, txController ...*gorm.DB) error {
	tx := r.db
	if len(txController) > 0 {
		tx = txController[0]
	}

	if input.HardDelete {
		tx = tx.Unscoped()
	}

	err := tx.WithContext(ctx).Where("id = ?", input.UserID).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	return nil
}
