//  nanti mempelajari ini setelah mempelajari ini baru share ke linkedin

package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"personal_erp/backend/internal/models"
	"personal_erp/backend/internal/store"
)

type Server struct {
	activities *store.ActivityStore
	categories *store.CategoryStore
	db         *sql.DB
}

func NewServer(db *sql.DB, activities *store.ActivityStore, categories *store.CategoryStore) *Server {
	return &Server{db: db, activities: activities, categories: categories}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	writeData(w, http.StatusOK, map[string]string{"status": "ok"})
}


func (s *Server) listCategories(w http.ResponseWriter, r *http.Request) {
	items, err := s.categories.List(r.Context())
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, items)
}

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string  `json:"name"`
		Color string  `json:"color"`
		Icon  *string `json:"icon"`
	}
	if err := decodeBody(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Color = strings.TrimSpace(in.Color)
	if in.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if in.Color == "" {
		in.Color = "#2563eb"
	}
	c, err := s.categories.Create(r.Context(), in.Name, in.Color, in.Icon)
	if err != nil {
		if isDupKey(err) {
			writeError(w, http.StatusConflict, "category name already exists")
			return
		}
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusCreated, c)
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	var in struct {
		Name  string  `json:"name"`
		Color string  `json:"color"`
		Icon  *string `json:"icon"`
	}
	if err := decodeBody(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Color = strings.TrimSpace(in.Color)
	if in.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if in.Color == "" {
		in.Color = "#2563eb"
	}
	c, err := s.categories.Update(r.Context(), id, in.Name, in.Color, in.Icon)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		if isDupKey(err) {
			writeError(w, http.StatusConflict, "category name already exists")
			return
		}
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, c)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := s.categories.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, map[string]bool{"deleted": true})
}

// ---- Activities ----

func (s *Server) listActivities(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := models.ActivityFilter{
		Date:   strings.TrimSpace(q.Get("date")),
		Start:  strings.TrimSpace(q.Get("start")),
		End:    strings.TrimSpace(q.Get("end")),
		Status: strings.TrimSpace(q.Get("status")),
		Query:  strings.TrimSpace(q.Get("q")),
	}
	if v := strings.TrimSpace(q.Get("priority")); v != "" {
		if !validPriority(v) {
			writeError(w, http.StatusBadRequest, "priority must be low, medium, or high")
			return
		}
		f.Priority = v
	}
	if v := strings.TrimSpace(q.Get("category_id")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid category_id")
			return
		}
		f.CategoryID = &id
	}
	for _, d := range []string{f.Date, f.Start, f.End} {
		if d != "" {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				writeError(w, http.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
				return
			}
		}
	}
	items, err := s.activities.List(r.Context(), f)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, items)
}

func (s *Server) getActivity(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid activity id")
		return
	}
	a, err := s.activities.Get(r.Context(), id)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, a)
}

type activityInput struct {
	Title         string  `json:"title"`
	Description   *string `json:"description"`
	CategoryID    *int64  `json:"category_id"`
	ActivityDate  string  `json:"activity_date"`
	StartTime     *string `json:"start_time"`
	EndTime       *string `json:"end_time"`
	DurationMin   *int    `json:"duration_min"`
	Status        string  `json:"status"`
	Priority      string  `json:"priority"`
	Notes         *string `json:"notes"`
}

func (s *Server) validateActivity(ctx context.Context, in *activityInput) (*models.Activity, string) {
	in.Title = strings.TrimSpace(in.Title)
	in.ActivityDate = strings.TrimSpace(in.ActivityDate)

	if in.Title == "" {
		return nil, "title is required"
	}
	if len(in.Title) > 255 {
		return nil, "title must be at most 255 characters"
	}
	if _, err := time.Parse("2006-01-02", in.ActivityDate); err != nil {
		return nil, "activity_date must be a valid date (YYYY-MM-DD)"
	}

	status := strings.ToLower(strings.TrimSpace(in.Status))
	if status == "" {
		status = "planned"
	}
	if !validStatus(status) {
		return nil, "status must be planned, in_progress, done, or cancelled"
	}

	priority := strings.ToLower(strings.TrimSpace(in.Priority))
	if priority == "" {
		priority = "medium"
	}
	if !validPriority(priority) {
		return nil, "priority must be low, medium, or high"
	}

	if in.CategoryID != nil {
		exists, err := s.categories.Exists(ctx, *in.CategoryID)
		if err != nil {
			return nil, "category lookup failed"
		}
		if !exists {
			return nil, "category_id does not exist"
		}
	}

	var start, end *string
	if in.StartTime != nil && strings.TrimSpace(*in.StartTime) != "" {
		t := strings.TrimSpace(*in.StartTime)
		if !validTime(t) {
			return nil, "start_time must be HH:MM or HH:MM:SS"
		}
		start = &t
	}
	if in.EndTime != nil && strings.TrimSpace(*in.EndTime) != "" {
		t := strings.TrimSpace(*in.EndTime)
		if !validTime(t) {
			return nil, "end_time must be HH:MM or HH:MM:SS"
		}
		end = &t
	}

	duration := in.DurationMin
	if start != nil && end != nil {
		if *end < *start {
			return nil, "end_time must be after start_time"
		}
		d := computeMinutes(*start, *end)
		duration = &d
	}
	if duration != nil && *duration < 0 {
		return nil, "duration_min cannot be negative"
	}

	return &models.Activity{
		Title:        in.Title,
		Description:  in.Description,
		CategoryID:   in.CategoryID,
		ActivityDate: in.ActivityDate,
		StartTime:    start,
		EndTime:      end,
		DurationMin:  duration,
		Status:       status,
		Priority:     priority,
		Notes:        in.Notes,
	}, ""
}

func (s *Server) createActivity(w http.ResponseWriter, r *http.Request) {
	var in activityInput
	if err := decodeBody(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	activity, msg := s.validateActivity(r.Context(), &in)
	if msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	a, err := s.activities.Create(r.Context(), activity)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusCreated, a)
}

func (s *Server) updateActivity(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid activity id")
		return
	}
	var in activityInput
	if err := decodeBody(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	activity, msg := s.validateActivity(r.Context(), &in)
	if msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	a, err := s.activities.Update(r.Context(), id, activity)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "activity not found")
			return
		}
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, a)
}

func (s *Server) deleteActivity(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid activity id")
		return
	}
	if err := s.activities.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "activity not found")
			return
		}
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (s *Server) activityStats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	start := strings.TrimSpace(q.Get("start"))
	end := strings.TrimSpace(q.Get("end"))
	if d := strings.TrimSpace(q.Get("date")); d != "" {
		start = d
		end = d
	}
	for _, d := range []string{start, end} {
		if d != "" {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				writeError(w, http.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
				return
			}
		}
	}
	stats, err := s.activities.Stats(r.Context(), start, end)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeData(w, http.StatusOK, stats)
}

// ---- helpers ----

func validStatus(v string) bool {
	switch v {
	case "planned", "in_progress", "done", "cancelled":
		return true
	}
	return false
}

func validPriority(v string) bool {
	switch v {
	case "low", "medium", "high":
		return true
	}
	return false
}

func validTime(v string) bool {
	for _, layout := range []string{"15:04", "15:04:05"} {
		if _, err := time.Parse(layout, v); err == nil {
			return true
		}
	}
	return false
}

func computeMinutes(start, end string) int {
	layout := "15:04:05"
	parse := func(v string) time.Time {
		if len(v) == 5 {
			t, err := time.Parse("15:04", v)
			if err != nil {
				return time.Time{}
			}
			return t
		}
		t, _ := time.Parse(layout, v)
		return t
	}
	s := parse(start)
	e := parse(end)
	d := e.Sub(s)
	if d < 0 {
		return 0
	}
	return int(d.Minutes())
}

func pathID(r *http.Request) (int64, error) {
	raw := r.PathValue("id")
	return strconv.ParseInt(raw, 10, 64)
}

func isDupKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
