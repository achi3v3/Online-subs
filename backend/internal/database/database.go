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

var db *gorm.DB

func Init() *gorm.DB {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Info("database.Init: No .env file found")
	}

	host, exists := os.LookupEnv("DB_HOST")
	if exists {
		logrus.Infof("database.Init: DB_HOST=%s", host)
	}

	port, exists := os.LookupEnv("DB_PORT")
	if exists {
		logrus.Infof("database.Init: DB_PORT=%s", port)
	}

	user, exists := os.LookupEnv("DB_USER")
	if exists {
		logrus.Infof("database.Init: DB_USER=%s", user)
	}

	password, exists := os.LookupEnv("DB_PASSWORD")
	if exists {
		logrus.Infof("database.Init: DB_PASSWORD=***") // Don't log actual password
	}

	dbname, exists := os.LookupEnv("DB_NAME")
	if exists {
		logrus.Infof("database.Init: DB_NAME=%s", dbname)
	}

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		logrus.Fatal("database.Init: Не все переменные окружения для подключения к БД установлены")
	}
	// Формируем строку подключения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	//dsn := "host=localhost user=postgres password=discolover dbname=botans port=8001 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("database.Init: Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&models.UserSubs{})
	return db
}

func Get() *gorm.DB {
	if db == nil {
		db = Init()
	}
	return db
}
