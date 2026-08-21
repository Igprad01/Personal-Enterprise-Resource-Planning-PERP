package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"personal-erp-backend/internal/application"
	"personal-erp-backend/pkg/response"
)

type CategoryHandler struct {
	service *application.CategoryService
}

func NewCategoryHandler(service *application.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

type categoryInput struct {
	Name  string  `json:"name"`
	Color string  `json:"color"`
	Icon  *string `json:"icon"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, items)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var in categoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Color = strings.TrimSpace(in.Color)

	item, err := h.service.Create(c.Request.Context(), in.Name, in.Color, in.Icon)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusCreated, item)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid category id")
		return
	}
	var in categoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Color = strings.TrimSpace(in.Color)

	item, err := h.service.Update(c.Request.Context(), id, in.Name, in.Color, in.Icon)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, item)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.WriteAppError(c, err)
		return
	}
	response.Data(c, http.StatusOK, map[string]bool{"deleted": true})
}