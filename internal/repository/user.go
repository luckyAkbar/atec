package repository

import (
	"context"

	"github.com/luckyAkbar/atec/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryIface interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, input CreateUserInput, txController ...*gorm.DB) (*model.User, error)
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
