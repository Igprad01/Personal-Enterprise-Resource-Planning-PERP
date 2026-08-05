package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"personal_erp/backend/internal/models"
)

var ErrNotFound = errors.New("not found")

type ActivityStore struct{ db *sql.DB }

func NewActivityStore(db *sql.DB) *ActivityStore { return &ActivityStore{db: db} }

const activityColumns = `a.id, a.title, a.description, a.category_id, a.activity_date,
	a.start_time, a.end_time, a.duration_min, a.status, a.priority, a.notes,
	a.created_at, a.updated_at, c.name, c.color`

func scanActivity(row interface{ Scan(...any) error }) (models.Activity, error) {
	var a models.Activity
	var (
		description, notes, catName, catColor sql.NullString
		start, end                            sql.NullString
		catID                                 sql.NullInt64
		dur                                   sql.NullInt64
	)
	err := row.Scan(&a.ID, &a.Title, &description, &catID, &a.ActivityDate,
		&start, &end, &dur, &a.Status, &a.Priority, &notes,
		&a.CreatedAt, &a.UpdatedAt, &catName, &catColor)
	if err != nil {
		return models.Activity{}, err
	}
	if description.Valid {
		a.Description = &description.String
	}
	if notes.Valid {
		a.Notes = &notes.String
	}
	if catID.Valid {
		a.CategoryID = &catID.Int64
	}
	if start.Valid {
		t := start.String
		a.StartTime = &t
	}
	if end.Valid {
		t := end.String
		a.EndTime = &t
	}
	if dur.Valid {
		d := int(dur.Int64)
		a.DurationMin = &d
	}
	if catName.Valid {
		a.CategoryName = &catName.String
	}
	if catColor.Valid {
		a.CategoryColor = &catColor.String
	}
	return a, nil
}

func (s *ActivityStore) List(ctx context.Context, f models.ActivityFilter) ([]models.Activity, error) {
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

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.Activity{}
	for rows.Next() {
		a, err := scanActivity(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

func (s *ActivityStore) Get(ctx context.Context, id int64) (models.Activity, error) {
	query := fmt.Sprintf(`SELECT %s FROM activities a
		LEFT JOIN categories c ON c.id = a.category_id WHERE a.id = ?`, activityColumns)
	row := s.db.QueryRowContext(ctx, query, id)
	a, err := scanActivity(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Activity{}, ErrNotFound
	}
	return a, err
}

func (s *ActivityStore) Create(ctx context.Context, a *models.Activity) (models.Activity, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO activities
			(title, description, category_id, activity_date, start_time, end_time, duration_min, status, priority, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Title, a.Description, a.CategoryID, a.ActivityDate,
		a.StartTime, a.EndTime, a.DurationMin, a.Status, a.Priority, a.Notes)
	if err != nil {
		return models.Activity{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Activity{}, err
	}
	return s.Get(ctx, id)
}

func (s *ActivityStore) Update(ctx context.Context, id int64, a *models.Activity) (models.Activity, error) {
	res, err := s.db.ExecContext(ctx,
		`UPDATE activities SET
			title = ?, description = ?, category_id = ?, activity_date = ?,
			start_time = ?, end_time = ?, duration_min = ?, status = ?, priority = ?, notes = ?
		 WHERE id = ?`,
		a.Title, a.Description, a.CategoryID, a.ActivityDate,
		a.StartTime, a.EndTime, a.DurationMin, a.Status, a.Priority, a.Notes, id)
	if err != nil {
		return models.Activity{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return models.Activity{}, err
	}
	if affected == 0 {
		return models.Activity{}, ErrNotFound
	}
	return s.Get(ctx, id)
}

func (s *ActivityStore) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM activities WHERE id = ?`, id)
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

func (s *ActivityStore) Stats(ctx context.Context, start, end string) (models.Stats, error) {
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

	var stats models.Stats
	countsQuery := `SELECT COUNT(*),
		SUM(status = 'done'),
		SUM(status = 'cancelled'),
		SUM(status = 'in_progress'),
		SUM(status = 'planned'),
		COALESCE(SUM(duration_min), 0)
		FROM activities` + where
	if err := s.db.QueryRowContext(ctx, countsQuery, args...).
		Scan(&stats.Total, &stats.Done, &stats.Cancelled, &stats.InProgress, &stats.Planned, &stats.TotalMinutes); err != nil {
		return stats, err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT COALESCE(c.name, 'Uncategorized'), COUNT(*), COALESCE(SUM(a.duration_min), 0)
		 FROM activities a LEFT JOIN categories c ON c.id = a.category_id`+where+`
		 GROUP BY COALESCE(c.name, 'Uncategorized')
		 ORDER BY COUNT(*) DESC, COALESCE(c.name, 'Uncategorized')`, args...)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	stats.ByCategory = []models.CategoryStat{}
	for rows.Next() {
		var cs models.CategoryStat
		if err := rows.Scan(&cs.Category, &cs.Count, &cs.Minutes); err != nil {
			return stats, err
		}
		stats.ByCategory = append(stats.ByCategory, cs)
	}
	return stats, rows.Err()
}
