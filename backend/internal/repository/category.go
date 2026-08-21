package repository

import (
	"context"

	"gorm.io/gorm"

	"personal-erp-backend/internal/domain"
)

type CategoryRepository struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	var items []domain.Category
	err := r.db.WithContext(ctx).Order("name").Find(&items).Error
	return items, err
}

func (r *CategoryRepository) Create(ctx context.Context, name, color string, icon *string) (domain.Category, error) {
	c := domain.Category{Name: name, Color: color, Icon: icon}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return domain.Category{}, err
	}
	return c, nil
}

func (r *CategoryRepository) Get(ctx context.Context, id int64) (domain.Category, error) {
	var c domain.Category
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err == gorm.ErrRecordNotFound {
		return domain.Category{}, domain.ErrNotFound
	}
	return c, err
}

func (r *CategoryRepository) Update(ctx context.Context, id int64, name, color string, icon *string) (domain.Category, error) {
	res := r.db.WithContext(ctx).Model(&domain.Category{}).Where("id = ?", id).
		Updates(map[string]any{"name": name, "color": color, "icon": icon})
	if res.Error != nil {
		return domain.Category{}, res.Error
	}
	if res.RowsAffected == 0 {
		return domain.Category{}, domain.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&domain.Category{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}