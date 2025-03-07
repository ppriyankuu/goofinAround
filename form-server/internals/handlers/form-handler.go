package handlers

import (
	"fmt"
	"form-server/internals/models"
	"form-server/internals/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	formService  services.FormService
	cacheService services.CacheService
}

func NewFormHandler(formService services.FormService, cacheService services.CacheService) *FormHandler {
	return &FormHandler{formService: formService, cacheService: cacheService}
}

func (h *FormHandler) SubmitForm(c *gin.Context) {
	var formData string
	if err := c.BindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := h.formService.SubmitForm(formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, form)
}

func (h *FormHandler) GetFormByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cachedFormData, err := h.cacheService.Get(fmt.Sprintf("form:%d", id))
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"id": id, "form_data": cachedFormData})
		return
	}

	form, err := h.formService.GetFormByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := h.cacheService.Set(fmt.Sprintf("form:%d", id), form.FormData, 10*time.Minute); err != nil {
		log.Printf("Failed to cache form: %v", err)
	}

	c.JSON(http.StatusOK, form)
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var form models.FormResponse
	if err := c.BindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form.ID = uint(id)

	if err := h.formService.UpdateForm(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.cacheService.Delete(fmt.Sprintf("form:%d", id)); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	c.JSON(http.StatusOK, form)
}

func (h *FormHandler) DeleteForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.formService.DeleteForm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.cacheService.Delete(fmt.Sprintf("form:%d", id)); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "form deleted successfully"})
}
