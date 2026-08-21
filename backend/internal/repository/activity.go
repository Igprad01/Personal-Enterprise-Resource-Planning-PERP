package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"personal-erp-backend/internal/domain"
)

type ActivityRepository struct{ db *gorm.DB }

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

const activityColumns = `a.id, a.title, a.description, a.category_id, a.activity_date,
	a.start_time, a.end_time, a.duration_min, a.status, a.priority, a.notes,
	a.created_at, a.updated_at, c.name AS category_name, c.color AS category_color`

func (r *ActivityRepository) List(ctx context.Context, f domain.ActivityFilter) ([]domain.Activity, error) {
	var conds []string
	var args []any

	if f.Date != "" {
		conds = append(conds, "a.activity_date = ?")
		args = append(args, f.Date)
	}
	if f.Start != "" {
		conds = append(conds, "a.activity_date >= ?")
		args = append(args, f.Start)
	}
	if f.End != "" {
		conds = append(conds, "a.activity_date <= ?")
		args = append(args, f.End)
	}
	if f.CategoryID != nil {
		conds = append(conds, "a.category_id = ?")
		args = append(args, *f.CategoryID)
	}
	if f.Status != "" {
		conds = append(conds, "a.status = ?")
		args = append(args, f.Status)
	}
	if f.Priority != "" {
		conds = append(conds, "a.priority = ?")
		args = append(args, f.Priority)
	}
	if f.Query != "" {
		conds = append(conds, "(a.title LIKE ? OR a.description LIKE ? OR a.notes LIKE ?)")
		like := "%" + f.Query + "%"
		args = append(args, like, like, like)
	}

	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}

	query := fmt.Sprintf(
		`SELECT %s FROM activities a LEFT JOIN categories c ON c.id = a.category_id%s
		 ORDER BY a.activity_date DESC, a.start_time IS NULL, a.start_time ASC, a.id DESC`,
		activityColumns, where)

	var items []domain.Activity
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ActivityRepository) Get(ctx context.Context, id int64) (domain.Activity, error) {
	query := fmt.Sprintf(`SELECT %s FROM activities a
		LEFT JOIN categories c ON c.id = a.category_id WHERE a.id = ?`, activityColumns)
	var a domain.Activity
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&a).Error
	if err == gorm.ErrRecordNotFound {
		return domain.Activity{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Activity{}, err
	}
	if a.ID == 0 {
		return domain.Activity{}, domain.ErrNotFound
	}
	return a, nil
}

func (r *ActivityRepository) Create(ctx context.Context, a *domain.Activity) (domain.Activity, error) {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return domain.Activity{}, err
	}
	return r.Get(ctx, a.ID)
}

func (r *ActivityRepository) Update(ctx context.Context, id int64, a *domain.Activity) (domain.Activity, error) {
	res := r.db.WithContext(ctx).Model(&domain.Activity{}).Where("id = ?", id).Updates(map[string]any{
		"title":         a.Title,
		"description":   a.Description,
		"category_id":   a.CategoryID,
		"activity_date": a.ActivityDate,
		"start_time":    a.StartTime,
		"end_time":      a.EndTime,
		"duration_min":  a.DurationMin,
		"status":        a.Status,
		"priority":      a.Priority,
		"notes":         a.Notes,
	})
	if res.Error != nil {
		return domain.Activity{}, res.Error
	}
	if res.RowsAffected == 0 {
		return domain.Activity{}, domain.ErrNotFound
	}
	return r.Get(ctx, id)
}

func (r *ActivityRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&domain.Activity{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ActivityRepository) Stats(ctx context.Context, start, end string) (domain.Stats, error) {
	var conds []string
	var args []any
	if start != "" {
		conds = append(conds, "activity_date >= ?")
		args = append(args, start)
	}
	if end != "" {
		conds = append(conds, "activity_date <= ?")
		args = append(args, end)
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}

	var stats domain.Stats
	countsQuery := `SELECT COUNT(*),
		SUM(status = 'done'),
		SUM(status = 'cancelled'),
		SUM(status = 'in_progress'),
		SUM(status = 'planned'),
		COALESCE(SUM(duration_min), 0)
		FROM activities` + where
	if err := r.db.WithContext(ctx).Raw(countsQuery, args...).Scan(&stats).Error; err != nil {
		return stats, err
	}

	stats.ByCategory = []domain.CategoryStat{}
	err := r.db.WithContext(ctx).Raw(
		`SELECT COALESCE(c.name, 'Uncategorized') AS category, COUNT(*) AS count, COALESCE(SUM(a.duration_min), 0) AS minutes
		 FROM activities a LEFT JOIN categories c ON c.id = a.category_id`+where+`
		 GROUP BY COALESCE(c.name, 'Uncategorized')
		 ORDER BY COUNT(*) DESC, COALESCE(c.name, 'Uncategorized')`, args...).
		Scan(&stats.ByCategory).Error
	return stats, err
}