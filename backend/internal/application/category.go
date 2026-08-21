package application

import (
	"context"

	"personal-erp-backend/internal/domain"
	"personal-erp-backend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Create(ctx context.Context, name, color string, icon *string) (domain.Category, error) {
	if name == "" {
		return domain.Category{}, &domain.ValidationError{Message: "name is required"}
	}
	if color == "" {
		color = "#2563eb"
	}
	return s.repo.Create(ctx, name, color, icon)
}

func (s *CategoryService) Update(ctx context.Context, id int64, name, color string, icon *string) (domain.Category, error) {
	if name == "" {
		return domain.Category{}, &domain.ValidationError{Message: "name is required"}
	}
	if color == "" {
		color = "#2563eb"
	}
	return s.repo.Update(ctx, id, name, color, icon)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *CategoryService) Exists(ctx context.Context, id int64) (bool, error) {
	return s.repo.Exists(ctx, id)
}