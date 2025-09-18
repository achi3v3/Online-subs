package subs

import (
	"app/internal/database"
	"app/internal/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// calculateMonths вычисляет количество месяцев между двумя датами
func calculateMonths(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())

	// Если день окончания меньше дня начала, вычитаем один месяц
	if end.Day() < start.Day() {
		months--
	}

	totalMonths := years*12 + months

	// Если totalMonths отрицательный, возвращаем 0
	if totalMonths < 0 {
		return 0
	}

	return totalMonths
}

// Repository — контракт для работы с подписками в бд
type Repository interface {
	Create(sub *models.UserSubs) error
	GetByID(id uint) (*models.UserSubs, error)
	Update(sub *models.UserSubs) error
	Delete(id uint) error
	List() ([]models.UserSubs, error)
	ListWithPagination(limit, offset int) ([]models.UserSubs, int64, error)
	GetTotalPriceForPeriod(startDate, endDate time.Time, userID, serviceName string) (uint, error)
}

// repository — структура, реализующая интерфейса Repository
type repository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewRepository — конструктор repository
func NewRepository(logger *logrus.Logger) Repository {
	return &repository{
		db:     database.Get(),
		logger: logger,
	}
}

// Create создает новую запись models.UserSubs в бд
func (r *repository) Create(sub *models.UserSubs) error {
	r.logger.Infof("repositor.Create: Creating subscription for user %s, service %s", sub.UserID, sub.ServiceName)
	err := r.db.Create(sub).Error
	if err != nil {
		r.logger.Fatalf("repositor.Create: Failed to create subscription: %v", err)
		return err
	}
	r.logger.Infof("repositor.Create: Subscription created successfully with ID %d", sub.ID)
	return nil
}

// GetByID возвращает подписку по ID
func (r *repository) GetByID(id uint) (*models.UserSubs, error) {
	r.logger.Infof("repository.GetByID: Fetching subscription with ID %d", id)
	var sub models.UserSubs
	err := r.db.First(&sub, id).Error
	if err != nil {
		r.logger.Warnf("repository.GetByID: Failed to fetch subscription with ID %d: %v", id, err)
		return nil, err // GORM возвращает gorm.ErrRecordNotFound если запись не найдена
	}
	r.logger.Infof("repository.GetByID: Subscription with ID %d fetched successfully", id)
	return &sub, nil
}

// Update обновляет существующую подписку
func (r *repository) Update(sub *models.UserSubs) error {
	r.logger.Infof("repository.Update: Updating subscription with ID %d", sub.ID)
	// Проверяем, существует ли подписка с таким ID
	var existingSub models.UserSubs
	err := r.db.First(&existingSub, sub.ID).Error
	if err != nil {
		r.logger.Warnf("repository.Update: Subscription with ID %d not found: %v", sub.ID, err)
		return gorm.ErrRecordNotFound
	}

	// Если подписка существует, обновляем её
	err = r.db.Save(sub).Error
	if err != nil {
		r.logger.Warnf("repository.Update: Failed to update subscription with ID %d: %v", sub.ID, err)
		return err
	}
	r.logger.Infof("repository.Update: Subscription with ID %d updated successfully", sub.ID)
	return nil
}

// Delete удаляет подписку по ID
func (r *repository) Delete(id uint) error {
	r.logger.Infof("repository.Delete: Deleting subscription with ID %d", id)
	// Проверяем, существует ли подписка с таким ID
	var existingSub models.UserSubs
	err := r.db.First(&existingSub, id).Error
	if err != nil {
		r.logger.Warnf("repository.Delete: Subscription with ID %d not found: %v", id, err)
		return gorm.ErrRecordNotFound
	}

	// Если подписка существует, удаляем её
	err = r.db.Delete(&models.UserSubs{}, id).Error
	if err != nil {
		r.logger.Errorf("repository.Delete: Failed to delete subscription with ID %d: %v", id, err)
		return err
	}
	r.logger.Infof("repository.Delete: Subscription with ID %d deleted successfully", id)
	return nil
}

// List возвращает список всех подписк
func (r *repository) List() ([]models.UserSubs, error) {
	r.logger.Info("repository.List: Fetching list of all subscriptions")
	var subs []models.UserSubs
	err := r.db.Find(&subs).Error
	if err != nil {
		r.logger.Errorf("repository.List: Failed to fetch list of subscriptions: %v", err)
		return nil, err
	}
	r.logger.Infof("repository.List: Fetched %d subscriptions", len(subs))
	return subs, nil
}

// ListWithPagination возвращает список подписок с пагинацией
func (r *repository) ListWithPagination(limit, offset int) ([]models.UserSubs, int64, error) {
	r.logger.Infof("repository.ListWithPagination: Fetching list of subscriptions with limit %d and offset %d", limit, offset)
	var subs []models.UserSubs
	var total int64

	// Получаем общее количество записей
	if err := r.db.Model(&models.UserSubs{}).Count(&total).Error; err != nil {
		r.logger.Errorf("repository.ListWithPagination: Failed to count subscriptions: %v", err)
		return nil, 0, err
	}

	// Получаем записи с пагинацией
	err := r.db.Limit(limit).Offset(offset).Find(&subs).Error
	if err != nil {
		r.logger.Errorf("repository.ListWithPagination: Failed to fetch list of subscriptions: %v", err)
		return nil, 0, err
	}

	r.logger.Infof("repository.ListWithPagination: Fetched %d subscriptions (limit %d, offset %d)", len(subs), limit, offset)
	return subs, total, nil
}

// GetTotalPriceForPeriod подсчитывает суммарную стоимость подписок за период
func (r *repository) GetTotalPriceForPeriod(startDate, endDate time.Time, userID, serviceName string) (uint, error) {
	r.logger.Infof("repository.GetTotalPriceForPeriod: Calculating total price for period %s to %s, userID: %s, serviceName: %s", startDate, endDate, userID, serviceName)
	query := r.db.Model(&models.UserSubs{}).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate) // колизии дат

	// Проверка полей
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// Получаем все подписки, которые пересекаются с заданным периодом
	var subs []models.UserSubs
	if err := query.Find(&subs).Error; err != nil {
		r.logger.Errorf("repository.GetTotalPriceForPeriod: Failed to fetch subscriptions: %v", err)
		return 0, err
	}

	// Подсчитываем сумму с учетом количества месяцев
	var total uint
	for _, sub := range subs {
		// Определяем период пересечения
		actualStart := sub.StartDate
		if startDate.After(actualStart) {
			actualStart = startDate
		}

		actualEnd := sub.EndDate
		if endDate.Before(actualEnd) {
			actualEnd = endDate
		}

		// Рассчитываем количество месяцев
		months := calculateMonths(actualStart, actualEnd)
		total += uint(months) * uint(sub.Price)
	}

	r.logger.Infof("repository.GetTotalPriceForPeriod: Total price calculated: %d", total)
	return total, nil
}
