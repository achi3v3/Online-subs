package subs

import (
	"app/internal/models"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Handlers — контракт для HTTP-обработчиков
type Handlers interface {
	CreateSub(c *gin.Context)
	GetSubByID(c *gin.Context)
	UpdateSub(c *gin.Context)
	DeleteSub(c *gin.Context)
	ListSubs(c *gin.Context)
	GetTotalPriceForPeriod(c *gin.Context)
}

// handlers  — структура, реализующая интерфейс Handlers
type handlers struct {
	service Service
	logger  *logrus.Logger
}

// NewHandlers — конструктор handlers
func NewHandlers(service Service, logger *logrus.Logger) Handlers {
	return &handlers{
		service: service,
		logger:  logger,
	}
}

// CreateSub godoc
// @Summary Создаёт запись подписки
// @Description Создает новую запись о подписке пользователя (формат времени: RFC3339 [2025-07-31T19:00:00Z])
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.UserSubs true "Данные подписки"
// @Success 201 {object} models.UserSubs
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs [post]
func (h *handlers) CreateSub(c *gin.Context) {
	h.logger.Info("handlers.CreateSub: Creating subscription")
	var sub models.UserSubs
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.logger.Errorf("handlers.CreateSub: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSub(&sub); err != nil {
		h.logger.Errorf("handlers.CreateSub: Failed to create subscription: %v", err)
		// Проверяем, является ли ошибка ошибкой БД
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "database") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	h.logger.Infof("handlers.CreateSub: Subscription created successfully with ID %d", sub.ID)
	c.JSON(http.StatusCreated, sub)
}

// GetSubByID godoc
// @Summary Получить подписку по ID
// @Description Возвращает информацию о подписке по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} models.UserSubs
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subs/{id} [get]
func (h *handlers) GetSubByID(c *gin.Context) {
	h.logger.Info("handlers.GetSubByID: Fetching subscription by ID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Errorf("handlers.GetSubByID: Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.service.GetSubByID(uint(id))
	if err != nil {
		h.logger.Warnf("handlers.GetSubByID: Subscription not found with ID %d", id)
		// Проверяем, является ли ошибка ошибкой БД
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "database") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		}
		return
	}

	h.logger.Infof("handlers.GetSubByID: Subscription with ID %d fetched successfully", id)
	c.JSON(http.StatusOK, sub)
}

// UpdateSub godoc
// @Summary Обновить подписку
// @Description Обновляет информацию о существующей подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body models.UserSubs true "Обновленные данные подписки"
// @Success 200 {object} models.UserSubs
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subs/{id} [put]
func (h *handlers) UpdateSub(c *gin.Context) {
	h.logger.Info("handlers.UpdateSub: Updating subscription")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Errorf("handlers.UpdateSub: Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var sub models.UserSubs
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.logger.Errorf("handlers.UpdateSub: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Убеждаемся, что ID из URL используется, а не из тела запроса
	sub.ID = uint(id)

	if err := h.service.UpdateSub(&sub); err != nil {
		h.logger.Errorf("handlers.UpdateSub: Failed to update subscription with ID %d: %v", id, err)
		// Проверяем, является ли ошибка ошибкой БД
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "database") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	h.logger.Infof("handlers.UpdateSub: Subscription with ID %d updated successfully", id)
	c.JSON(http.StatusOK, sub)
}

// DeleteSub godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subs/{id} [delete]
func (h *handlers) DeleteSub(c *gin.Context) {
	h.logger.Info("handlers.DeleteSub: Deleting subscription")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Errorf("handlers.DeleteSub: Invalid ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteSub(uint(id)); err != nil {
		h.logger.Warnf("handlers.DeleteSub: Subscription not found with ID %d", id)
		// Проверяем, является ли ошибка ошибкой БД
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "database") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		}
		return
	}

	h.logger.Infof("handlers.DeleteSub: Subscription with ID %d deleted successfully", id)
	c.Status(http.StatusNoContent)
}

// ListSubs godoc
// @Summary Список подписок
// @Description Возвращает список всех записей о подписках с пагинацией
// @Tags subscriptions
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество элементов на странице (по умолчанию 10, максимум 100)"
// @Success 200 {array} models.UserSubs
// @Router /subs [get]
func (h *handlers) ListSubs(c *gin.Context) {
	h.logger.Info("handlers.ListSubs: Fetching list of all subscriptions")

	// Получаем параметры пагинации
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Если указаны параметры пагинации, используем пагинацию
	if pageStr != "" || limitStr != "" {
		offset := (page - 1) * limit

		subs, total, err := h.service.ListSubsWithPagination(limit, offset)
		if err != nil {
			h.logger.Errorf("handlers.ListSubs: Failed to fetch list of subscriptions: %v", err)
			// Проверяем, является ли ошибка ошибкой БД
			if strings.Contains(err.Error(), "database") {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		h.logger.Infof("handlers.ListSubs: Fetched %d subscriptions (page %d, limit %d)", len(subs), page, limit)
		c.JSON(http.StatusOK, gin.H{
			"subscriptions": subs,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
		return
	}

	// Если параметры пагинации не указаны, возвращаем все подписки (обратная совместимость)
	subs, err := h.service.ListSubs()
	if err != nil {
		h.logger.Errorf("handlers.ListSubs: Failed to fetch list of subscriptions: %v", err)
		// Проверяем, является ли ошибка ошибкой БД
		if strings.Contains(err.Error(), "database") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	h.logger.Infof("handlers.ListSubs: Fetched %d subscriptions", len(subs))
	c.JSON(http.StatusOK, subs)
}

// GetTotalPriceForPeriod godoc
// @Summary Подсчитать сумму подписок за период
// @Description Подсчитывает суммарную стоимость всех подписок за выбранный период с фильтрацией по ID пользователя и названию подписки
// @Tags subscriptions
// @Produce json
// @Param start_date query string true "Дата начала периода (в формате RFC3339) [2025-06-30T19:00:00Z]"
// @Param end_date query string true "Дата окончания периода (в формате RFC3339) [2025-07-31T19:00:00Z]"
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} map[string]uint
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subs/total [get]
func (h *handlers) GetTotalPriceForPeriod(c *gin.Context) {
	h.logger.Info("handlers.GetTotalPriceForPeriod: Calculating total price for period")
	// Парсим параметрф запроса
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	// Проверяем необходмые параметры
	if startDateStr == "" {
		h.logger.Warn("handlers.GetTotalPriceForPeriod: start_date is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date is required"})
		return
	}
	if endDateStr == "" {
		h.logger.Warn("handlers.GetTotalPriceForPeriod: end_date is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date is required"})
		return
	}

	// Парсим даты
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		h.logger.Errorf("handlers.GetTotalPriceForPeriod: Invalid start_date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected RFC3339"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		h.logger.Errorf("handlers.GetTotalPriceForPeriod: Invalid end_date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected RFC3339"})
		return
	}

	total, err := h.service.GetTotalPriceForPeriod(startDate, endDate, userID, serviceName)
	if err != nil {
		h.logger.Errorf("handlers.GetTotalPriceForPeriod: Failed to calculate total price: %v", err)
		// Проверяем, является ли ошибка ошибкой БД
		if strings.Contains(err.Error(), "database") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	h.logger.Infof("handlers.GetTotalPriceForPeriod: Total price calculated: %d", total)
	c.JSON(http.StatusOK, gin.H{"total": total})
}
