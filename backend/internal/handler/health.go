package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"personal-erp-backend/pkg/response"
)

type HealthHandler struct{ db *gorm.DB }

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		response.Error(c, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	response.Data(c, http.StatusOK, map[string]string{"status": "ok"})
}