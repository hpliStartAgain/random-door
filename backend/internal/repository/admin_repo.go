package repository

import (
	"context"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *AdminRepo {
	return &AdminRepo{DB: db}
}

func (r *AdminRepo) ListCities(ctx context.Context) ([]model.City, error) {
	var cities []model.City
	err := r.DB.WithContext(ctx).Order("id ASC").Find(&cities).Error
	return cities, err
}

func (r *AdminRepo) FindCityByID(ctx context.Context, id int64) (*model.City, error) {
	var city model.City
	err := r.DB.WithContext(ctx).First(&city, id).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *AdminRepo) ListTags(ctx context.Context, cityID int64) ([]model.CityTag, error) {
	var tags []model.CityTag
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&tags).Error
	return tags, err
}

func (r *AdminRepo) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	var rows []model.Landmark
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&rows).Error
	return rows, err
}

func (r *AdminRepo) ListFoods(ctx context.Context, cityID int64) ([]model.Food, error) {
	var rows []model.Food
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&rows).Error
	return rows, err
}

func (r *AdminRepo) ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error) {
	var rows []model.Character
	err := r.DB.WithContext(ctx).Where("city_id = ?", cityID).Order("id ASC").Find(&rows).Error
	return rows, err
}

func (r *AdminRepo) CreateLandmark(ctx context.Context, row *model.Landmark) error {
	return r.DB.WithContext(ctx).Create(row).Error
}

func (r *AdminRepo) CreateFood(ctx context.Context, row *model.Food) error {
	return r.DB.WithContext(ctx).Create(row).Error
}

func (r *AdminRepo) CreateCharacter(ctx context.Context, row *model.Character) error {
	return r.DB.WithContext(ctx).Create(row).Error
}

func (r *AdminRepo) UpdateCity(ctx context.Context, id int64, fields map[string]any, tags *[]string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var city model.City
		if err := tx.First(&city, id).Error; err != nil {
			return err
		}
		if len(fields) > 0 {
			if err := tx.Model(&model.City{}).Where("id = ?", id).Updates(fields).Error; err != nil {
				return err
			}
		}
		if tags != nil {
			if err := tx.Where("city_id = ?", id).Delete(&model.CityTag{}).Error; err != nil {
				return err
			}
			for _, tag := range *tags {
				row := model.CityTag{CityID: id, Tag: tag}
				if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *AdminRepo) UpdateLandmark(ctx context.Context, id int64, fields map[string]any) error {
	return updateByID[model.Landmark](ctx, r.DB, id, fields)
}

func (r *AdminRepo) UpdateFood(ctx context.Context, id int64, fields map[string]any) error {
	return updateByID[model.Food](ctx, r.DB, id, fields)
}

func (r *AdminRepo) UpdateCharacter(ctx context.Context, id int64, fields map[string]any) error {
	return updateByID[model.Character](ctx, r.DB, id, fields)
}

func (r *AdminRepo) DeleteLandmark(ctx context.Context, id int64) error {
	return deleteByID[model.Landmark](ctx, r.DB, id, "landmark")
}

func (r *AdminRepo) DeleteFood(ctx context.Context, id int64) error {
	return deleteByID[model.Food](ctx, r.DB, id, "food")
}

func (r *AdminRepo) DeleteCharacter(ctx context.Context, id int64) error {
	return deleteByID[model.Character](ctx, r.DB, id, "character")
}

func updateByID[T any](ctx context.Context, db *gorm.DB, id int64, fields map[string]any) error {
	var row T
	if err := db.WithContext(ctx).First(&row, id).Error; err != nil {
		return err
	}
	if len(fields) == 0 {
		return nil
	}
	return db.WithContext(ctx).Model(&row).Updates(fields).Error
}

func deleteByID[T any](ctx context.Context, db *gorm.DB, id int64, targetType string) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row T
		if err := tx.First(&row, id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&row).Error; err != nil {
			return err
		}
		return tx.Where("target_type = ? AND target_id = ?", targetType, id).Delete(&model.Comment{}).Error
	})
}
