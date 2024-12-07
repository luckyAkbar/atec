package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryIface interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, input CreateUserInput, txController ...*gorm.DB) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, userID uuid.UUID, input UpdateUserInput) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

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

type CreateUserInput struct {
	Email    string
	Password string
	Username string
	IsActive bool
	Roles    model.Roles
}

func (r *UserRepository) Create(ctx context.Context, input CreateUserInput, txController ...*gorm.DB) (*model.User, error) {
	tx := r.db.WithContext(ctx)
	if len(txController) > 0 {
		tx = txController[0]
	}

	user := &model.User{
		Email:    input.Email,
		Password: input.Password,
		IsActive: input.IsActive,
		Roles:    input.Roles,
		Username: input.Username,
	}

	err := tx.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

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

type UpdateUserInput struct {
	Email    string
	Password string `json:"-"`
	Username string
	IsActive *bool
}

func (uui UpdateUserInput) ToUpdatedFields() map[string]interface{} {
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

func (r *UserRepository) Update(ctx context.Context, userID uuid.UUID, input UpdateUserInput) (*model.User, error) {
	user := &model.User{}
	err := r.db.WithContext(ctx).Model(user).
		Clauses(clause.Returning{}).Where("id = ?", userID).
		Updates(input.ToUpdatedFields()).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
