package rest

import (
	"context"
	"database/sql"
	"developerApplication/internal/adapters/repository/postgres/developer_application"
	"developerApplication/internal/core/service"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	UploadBuild(ctx context.Context, params service.UploadBuildParams) (developer_application.DeveloperApplication, error)
	GetApplication(ctx context.Context, id string) (developer_application.DeveloperApplication, error)
	ListApplications(ctx context.Context, developerID string, page, size int) ([]developer_application.DeveloperApplication, error)
	StartVerification(ctx context.Context, id string) (developer_application.DeveloperApplication, error)
	GetVerificationStatus(ctx context.Context, id string) (string, error)
	PublishApplication(ctx context.Context, id string) error
}

type Handler struct {
	service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{service: svc}
}

// UploadBuild godoc
// @Summary Загрузить сборку
// @Description Загружает APK-файл и создаёт черновик приложения.
// @Tags applications
// @Accept multipart/form-data
// @Produce json
// @Param developer_id          formData string true  "ID разработчика"
// @Param code_name             formData string true  "Кодовое имя"
// @Param category_id           formData int    false "ID категории (по умолчанию 0)"
// @Param android_package_name  formData string true  "Android package name (напр. com.TPU.DIVE_E9)"
// @Param default_locale        formData string true  "Локаль по умолчанию (напр. ru)"
// @Param name                  formData string true  "Название приложения"
// @Param short_title           formData string true  "Краткий заголовок"
// @Param version               formData string true  "Версия (напр. 1.0.0)"
// @Param description           formData string false "Описание"
// @Param goals                 formData string false "Цели"
// @Param tasks                 formData string false "Задачи"
// @Param results               formData string false "Результаты"
// @Param challenges            formData string false "Челленджи"
// @Param location              formData string false "Местоположение"
// @Param video_cover           formData string false "Обложка видео"
// @Param safety                formData string false "Техника безопасности"
// @Param web_video             formData string false "Видео для web"
// @Param inapp_video           formData string false "Видео для in-app"
// @Param web_background_image  formData string false "Фон для web"
// @Param inapp_background_image formData string false "Фон для in-app (launcher)"
// @Param file                  formData file   true  "APK-файл"
// @Success 201 {object} ApplicationOutDTO
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

	result, err := h.service.UploadBuild(c.Request.Context(), service.UploadBuildParams{
		DeveloperID:          form.DeveloperID,
		CodeName:             form.CodeName,
		CategoryID:           form.CategoryID,
		AndroidPackageName:   form.AndroidPackageName,
		DefaultLocale:        form.DefaultLocale,
		Name:                 form.Name,
		ShortTitle:           form.ShortTitle,
		Version:              form.Version,
		Description:          form.Description,
		Goals:                form.Goals,
		Tasks:                form.Tasks,
		Results:              form.Results,
		Challenges:           form.Challenges,
		Location:             form.Location,
		VideoCover:           form.VideoCover,
		Safety:               form.Safety,
		WebVideo:             form.WebVideo,
		InappVideo:           form.InappVideo,
		WebBackgroundImage:   form.WebBackgroundImage,
		InappBackgroundImage: form.InappBackgroundImage,
		FileName:             header.Filename,
		Reader:               file,
		Size:                 header.Size,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toApplicationOutDTO(result))
}

// GetApplication godoc
// @Summary Получить приложение
// @Tags applications
// @Produce json
// @Param id path string true "ID приложения"
// @Success 200 {object} ApplicationOutDTO
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

	c.JSON(http.StatusOK, toApplicationOutDTO(app))
}

// ListApplications godoc
// @Summary Список приложений разработчика
// @Tags applications
// @Produce json
// @Param developer_id query string true  "ID разработчика"
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

	dtos := make([]ApplicationOutDTO, len(apps))
	for i, a := range apps {
		dtos[i] = toApplicationOutDTO(a)
	}
	c.JSON(http.StatusOK, ListApplicationsOutDTO{Items: dtos, Page: page, Size: size})
}

// StartVerification godoc
// @Summary Запустить верификацию
// @Tags applications
// @Produce json
// @Param id path string true "ID приложения"
// @Success 200 {object} ApplicationOutDTO
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

	c.JSON(http.StatusOK, toApplicationOutDTO(app))
}

// GetVerificationStatus godoc
// @Summary Статус верификации
// @Tags applications
// @Produce json
// @Param id path string true "ID приложения"
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
// @Produce json
// @Param id path string true "ID приложения"
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

	if err := h.service.PublishApplication(c.Request.Context(), id); err != nil {
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