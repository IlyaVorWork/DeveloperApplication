package rest

import (
	"context"
	"database/sql"
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	"developerApplication/internal/core/service"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	UploadBuild(ctx context.Context, developerID, version string, fileName string, reader io.Reader, size int64) (developer_application.DeveloperApplication, error)
	GetApplication(ctx context.Context, id string) (developer_application.DeveloperApplication, error)
	ListApplications(ctx context.Context, developerID string, page, size int) ([]developer_application.DeveloperApplication, error)
	StartVerification(ctx context.Context, id string) (developer_application.DeveloperApplication, error)
	GetVerificationStatus(ctx context.Context, id string) (string, error)
	PublishApplication(ctx context.Context, id, developerID, codeName string, categoryID int64) error
}

type Handler struct {
	service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{service: svc}
}

// UploadBuild godoc
// @Summary Загрузить сборку
// @Description Загружает APK-файл в хранилище и создаёт запись о сборке.
// @Tags applications
// @Accept multipart/form-data
// @Produce json
// @Param developer_id formData string true "ID разработчика"
// @Param version      formData string true "Версия сборки"
// @Param file         formData file   true "APK-файл"
// @Success 201 {object} developer_application.DeveloperApplication
// @Failure 400 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /applications [post]
func (h *Handler) UploadBuild(c *gin.Context) {
	var form UploadBuildInDTO
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	app, err := h.service.UploadBuild(c.Request.Context(), form.DeveloperID, form.Version, header.Filename, file, header.Size)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, app)
}

// GetApplication godoc
// @Summary Получить сборку
// @Tags applications
// @Produce json
// @Param id path string true "ID сборки"
// @Success 200 {object} developer_application.DeveloperApplication
// @Failure 404 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /applications/{id} [get]
func (h *Handler) GetApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	app, err := h.service.GetApplication(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, app)
}

// ListApplications godoc
// @Summary Список сборок разработчика
// @Tags applications
// @Produce json
// @Param developer_id query string true "ID разработчика"
// @Param page         query int    false "Номер страницы"
// @Param size         query int    false "Размер страницы"
// @Success 200 {object} ListApplicationsOutDTO
// @Failure 400 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /applications [get]
func (h *Handler) ListApplications(c *gin.Context) {
	developerID := strings.TrimSpace(c.Query("developer_id"))
	if developerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "developer_id is required"})
		return
	}

	page, err := parsePositiveIntQuery(c, "page", 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	size, err := parsePositiveIntQuery(c, "size", 50)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apps, err := h.service.ListApplications(c.Request.Context(), developerID, page, size)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, ListApplicationsOutDTO{Items: apps, Page: page, Size: size})
}

// StartVerification godoc
// @Summary Запустить верификацию
// @Tags applications
// @Produce json
// @Param id path string true "ID сборки"
// @Success 200 {object} developer_application.DeveloperApplication
// @Failure 404 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /applications/{id}/verify [post]
func (h *Handler) StartVerification(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	app, err := h.service.StartVerification(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, app)
}

// GetVerificationStatus godoc
// @Summary Статус верификации
// @Tags applications
// @Produce json
// @Param id path string true "ID сборки"
// @Success 200 {object} VerificationStatusOutDTO
// @Failure 404 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /applications/{id}/verify [get]
func (h *Handler) GetVerificationStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	status, err := h.service.GetVerificationStatus(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, VerificationStatusOutDTO{Status: status})
}

// PublishApplication godoc
// @Summary Опубликовать приложение
// @Tags applications
// @Accept json
// @Produce json
// @Param id   path string               true "ID сборки"
// @Param body body PublishApplicationInDTO true "Данные для публикации"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponseDTO
// @Failure 404 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /applications/{id}/publish [post]
func (h *Handler) PublishApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var body PublishApplicationInDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PublishApplication(c.Request.Context(), id, body.DeveloperID, body.CodeName, body.CategoryID); err != nil {
		writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parsePositiveIntQuery(c *gin.Context, key string, fallback int) (int, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, errors.New("invalid " + key)
	}
	return value, nil
}

func writeError(c *gin.Context, err error) {
	if errors.Is(err, context.Canceled) {
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "request canceled"})
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if errors.Is(err, service.ErrNotVerified) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
