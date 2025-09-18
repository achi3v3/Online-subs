package models

import (
	"time"
)

// rubles — алиас для валюты
type rubles uint

// UserSubs — структура подписки пользователя на сервис
type UserSubs struct {
	ID          uint      `json:"id" gorm:"primaryKey; column:id"`
	ServiceName string    `json:"service_name" gorm:"not null; column:service_name"`
	Price       rubles    `json:"price" gorm:"not null; column:price"`
	UserID      string    `json:"user_id" gorm:"not null; column:user_id"`
	StartDate   time.Time `json:"start_date" gorm:"not null; column:start_date"`
	EndDate     time.Time `json:"end_date" gorm:"not null; column:end_date"`
}
