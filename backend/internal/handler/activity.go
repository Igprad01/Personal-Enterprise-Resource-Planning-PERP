package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"personal-erp-backend/internal/application"
	"personal-erp-backend/internal/domain"
	"personal-erp-backend/pkg/response"
)

type ActivityHandler struct {
	service *application.ActivityService
}

func NewActivityHandler(service *application.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

type activityInput struct {
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	CategoryID   *int64  `json:"category_id"`
	ActivityDate string  `json:"activity_date"`
	StartTime    *string `json:"start_time"`
	EndTime      *string `json:"end_time"`
	DurationMin  *int    `json:"duration_min"`
	Status       string  `json:"status"`
	Priority     string  `json:"priority"`
	Notes        *string `json:"notes"`
}

func (h *ActivityHandler) List(c *gin.Context) {
	q := c.Request.URL.Query()
	f := domain.ActivityFilter{
		Date:   strings.TrimSpace(q.Get("date")),
		Start:  strings.TrimSpace(q.Get("start")),
		End:    strings.TrimSpace(q.Get("end")),
		Status: strings.TrimSpace(q.Get("status")),
		Query:  strings.TrimSpace(q.Get("q")),
	}
	if v := strings.TrimSpace(q.Get("priority")); v != "" {
		f.Priority = v
	}
	if v := strings.TrimSpace(q.Get("category_id")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid category_id")
			return
		}
		f.CategoryID = &id
	}
	for _, d := range []string{f.Date, f.Start, f.End} {
		if d != "" {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				response.Error(c, http.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
				return
			}
		}
	}
	items, err := h.service.List(c.Request.Context(), f)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, items)
}

func (h *ActivityHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid activity id")
		return
	}
	item, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, item)
}

func (h *ActivityHandler) Create(c *gin.Context) {
	var in activityInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	a := &domain.Activity{
		Title:        in.Title,
		Description:  in.Description,
		CategoryID:   in.CategoryID,
		ActivityDate: in.ActivityDate,
		StartTime:    in.StartTime,
		EndTime:      in.EndTime,
		DurationMin:  in.DurationMin,
		Status:       in.Status,
		Priority:     in.Priority,
		Notes:        in.Notes,
	}
	item, err := h.service.Create(c.Request.Context(), a)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusCreated, item)
}

func (h *ActivityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid activity id")
		return
	}
	var in activityInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	a := &domain.Activity{
		Title:        in.Title,
		Description:  in.Description,
		CategoryID:   in.CategoryID,
		ActivityDate: in.ActivityDate,
		StartTime:    in.StartTime,
		EndTime:      in.EndTime,
		DurationMin:  in.DurationMin,
		Status:       in.Status,
		Priority:     in.Priority,
		Notes:        in.Notes,
	}
	item, err := h.service.Update(c.Request.Context(), id, a)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, item)
}

func (h *ActivityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid activity id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, map[string]bool{"deleted": true})
}

func (h *ActivityHandler) Stats(c *gin.Context) {
	q := c.Request.URL.Query()
	start := strings.TrimSpace(q.Get("start"))
	end := strings.TrimSpace(q.Get("end"))
	if d := strings.TrimSpace(q.Get("date")); d != "" {
		start = d
		end = d
	}
	for _, d := range []string{start, end} {
		if d != "" {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				response.Error(c, http.StatusBadRequest, "invalid date, expected YYYY-MM-DD")
				return
			}
		}
	}
	stats, err := h.service.Stats(c.Request.Context(), start, end)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, stats)
}