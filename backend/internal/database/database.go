package database

import (
	"app/internal/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init() (*gorm.DB, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Warnf("database.Init: Error loading .env file: %v", err)
		}
	}

	host := getEnvOrFail("DB_HOST")
	port := getEnvOrFail("DB_PORT")
	user := getEnvOrFail("DB_USER")
	password := getEnvOrFail("DB_PASSWORD")
	dbname := getEnvOrFail("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		logrus.Fatal("database.Init: Не все переменные окружения для подключения к БД установлены")
	}

	logrus.Infof("database.Init: Connecting to %s@%s:%s/%s", user, host, port, dbname)
	// Формируем строку подключения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	//dsn := "host=localhost user=postgres password=discolover dbname=botans port=8001 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database.Init: Failed to connect to database: %v", err)
		return nil, fmt.Errorf("database.Init: Failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.UserSubs{}); err != nil {
		return nil, fmt.Errorf("database.Init: Failed to auto-migrate: %w", err)
	}

	logrus.Info("database.Init: Successfully connected and migrated")
	return db, nil
}

func Get() *gorm.DB { //  нужно добавить Once из пакета sync, для потокобезопасности
	if db == nil {
		db, _ = Init()
	}
	return db
}

func getEnvOrFail(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		logrus.Fatalf("database.gerEnvOrFail: Environment variable %s is required", key)
	}
	return value
}
