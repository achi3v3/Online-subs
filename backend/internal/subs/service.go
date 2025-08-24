package subs

import (
	"app/internal/models"
	"errors"
	"log"
	"time"

	"github.com/sirupsen/logrus"
)

// Service — контракт для работы с service
type Service interface {
	CreateSub(sub *models.UserSubs) error
	GetSubByID(id uint) (*models.UserSubs, error)
	UpdateSub(sub *models.UserSubs) error
	DeleteSub(id uint) error
	ListSubs() ([]models.UserSubs, error)
	GetTotalPriceForPeriod(startDate, endDate time.Time, userID, serviceName string) (uint, error)
}

// service  — структура, реализующая интерфейс Service
type service struct {
	repo   Repository
	logger *logrus.Logger
}

// NewService — конструктор servoce
func NewService(repo Repository, logger *logrus.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// CreateSub создает новую подписку с валидацией
func (s *service) CreateSub(sub *models.UserSubs) error {
	s.logger.Infof("service.CreateSub: Creating subscription")
	// Валидация обязательных полей
	if sub.ServiceName == "" {
		s.logger.Warnf("service.CreateSub: service_name is required")
		return errors.New("service_name is required")
	}
	if sub.Price <= 0 {
		s.logger.Warnf("service.CreateSub: price must be greater than 0")
		return errors.New("price must be greater than 0")
	}
	if sub.UserID == "" {
		s.logger.Warnf("service.CreateSub: user_id is required")
		return errors.New("user_id is required")
	}
	if sub.StartDate.IsZero() {
		s.logger.Warnf("service.CreateSub: start_date is required")
		return errors.New("start_date is required")
	}
	if sub.EndDate.IsZero() {
		s.logger.Warnf("service.CreateSub: end_date is required")
		return errors.New("end_date is required")
	}
	if sub.EndDate.Before(sub.StartDate) {
		s.logger.Warnf("service.CreateSub: end_date must be after start_date")
		return errors.New("end_date must be after start_date")
	}

	return s.repo.Create(sub)
}

// GetSubByID возвращает подписк по ID
func (s *service) GetSubByID(id uint) (*models.UserSubs, error) {
	log.Printf("service.GetSubByID: Fetching subscription with ID %d", id)
	return s.repo.GetByID(id)
}

// UpdateSub обновляет существующую подписку с валидацией
func (s *service) UpdateSub(sub *models.UserSubs) error {
	log.Printf("service.UpdateSub: Updating subscription with ID %d", sub.ID)
	// Валидация обязательных полей
	if sub.ID == 0 {
		s.logger.Warnf("service.UpdateSub: id is required for update")
		return errors.New("id is required for update")
	}
	if sub.ServiceName == "" {
		s.logger.Warnf("service.UpdateSub: service_name is required")
		return errors.New("service_name is required")
	}
	if sub.Price <= 0 {
		s.logger.Warnf("service.UpdateSub: price must be greater than 0")
		return errors.New("price must be greater than 0")
	}
	if sub.UserID == "" {
		s.logger.Warnf("service.UpdateSub: user_id is required")
		return errors.New("user_id is required")
	}
	if sub.StartDate.IsZero() {
		s.logger.Warnf("service.UpdateSub: start_date is required")
		return errors.New("start_date is required")
	}
	if sub.EndDate.IsZero() {
		s.logger.Warnf("service.UpdateSub: end_date is required")
		return errors.New("end_date is required")
	}
	if sub.EndDate.Before(sub.StartDate) {
		s.logger.Warnf("service.UpdateSub: end_date must be after start_date")
		return errors.New("end_date must be after start_date")
	}

	return s.repo.Update(sub)
}

// DeleteSub удаляет подписку по ID
func (s *service) DeleteSub(id uint) error {
	s.logger.Infof("service.DeleteSub: Deleting subscription with ID %d", id)
	return s.repo.Delete(id)
}

// ListSubs возвращает список всех подписок
func (s *service) ListSubs() ([]models.UserSubs, error) {
	s.logger.Infof("service.ListSubs: Fetching list of all subscriptions")
	return s.repo.List()
}

// GetTotalPriceForPeriod подсчитывает суммарную стоимостт подписок за период
func (s *service) GetTotalPriceForPeriod(startDate, endDate time.Time, userID, serviceName string) (uint, error) {
	s.logger.Infof("service.GetTotalPriceForPeriod: Calculating total price for period %s to %s, userID: %s, serviceName: %s", startDate, endDate, userID, serviceName)
	// Валидация обязательных полей
	if startDate.IsZero() {
		s.logger.Warnf("service.GetTotalPriceForPeriod: start_date is required")
		return 0, errors.New("start_date is required")
	}
	if endDate.IsZero() {
		s.logger.Warnf("service.GetTotalPriceForPeriod: end_date is required")
		return 0, errors.New("end_date is required")
	}
	if endDate.Before(startDate) {
		s.logger.Warnf("service.GetTotalPriceForPeriod: end_date must be after start_date")
		return 0, errors.New("end_date must be after start_date")
	}

	return s.repo.GetTotalPriceForPeriod(startDate, endDate, userID, serviceName)
}
