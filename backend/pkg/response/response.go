package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"personal-erp-backend/internal/domain"
)

func Data(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func WriteAppError(c *gin.Context, err error) {
	var ve *domain.ValidationError
	if errors.As(err, &ve) {
		Error(c, http.StatusBadRequest, ve.Message)
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		Error(c, http.StatusNotFound, "resource not found")
		return
	}
	Error(c, http.StatusInternalServerError, "internal server error")
}