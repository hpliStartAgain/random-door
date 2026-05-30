package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) FindByAnonymousID(ctx context.Context, anonymousID string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).Where("anonymous_id = ?", anonymousID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) UpdateCurrentCity(ctx context.Context, userID, cityID int64) error {
	return r.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).
		Update("current_city_id", cityID).Error
}
