package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"personal_erp/backend/internal/models"
)

type CategoryStore struct{ db *sql.DB }

func NewCategoryStore(db *sql.DB) *CategoryStore { return &CategoryStore{db: db} }

func (s *CategoryStore) List(ctx context.Context) ([]models.Category, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, color, icon, created_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Color, &c.Icon, &c.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (s *CategoryStore) Create(ctx context.Context, name, color string, icon *string) (models.Category, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO categories (name, color, icon) VALUES (?, ?, ?)`, name, color, icon)
	if err != nil {
		return models.Category{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Category{}, err
	}
	return s.Get(ctx, id)
}

func (s *CategoryStore) Get(ctx context.Context, id int64) (models.Category, error) {
	var c models.Category
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, color, icon, created_at FROM categories WHERE id = ?`, id).
		Scan(&c.ID, &c.Name, &c.Color, &c.Icon, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, ErrNotFound
	}
	return c, err
}

func (s *CategoryStore) Update(ctx context.Context, id int64, name, color string, icon *string) (models.Category, error) {
	res, err := s.db.ExecContext(ctx,
		`UPDATE categories SET name = ?, color = ?, icon = ? WHERE id = ?`, name, color, icon, id)
	if err != nil {
		return models.Category{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return models.Category{}, err
	}
	if affected == 0 {
		return models.Category{}, ErrNotFound
	}
	return s.Get(ctx, id)
}

func (s *CategoryStore) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *CategoryStore) Exists(ctx context.Context, id int64) (bool, error) {
	var one int
	err := s.db.QueryRowContext(ctx, `SELECT 1 FROM categories WHERE id = ?`, id).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (s *CategoryStore) Name(ctx context.Context, id int64) (string, error) {
	var name string
	err := s.db.QueryRowContext(ctx, `SELECT name FROM categories WHERE id = ?`, id).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("category %d: %w", id, err)
	}
	return name, nil
}
