package main

import (
	"app/internal/database"
	"app/internal/subs"
	"net/http"
	"os"

	_ "app/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Swagger Users Subscribtions
// @version 1.3
// @description REST-сервис для агрегации данных об онлайн-подписках пользователей
// @contact.url	https://github.com/achi3v3
// @contact.email aamir-tutaev@mail.ru
// @host localhost:8080
// @BasePath /
func main() {
	// Настройка логирования
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Установка уровня логирования из переменной окружения
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // значение по умолчанию
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warnf("Invalid log level '%s', using 'info' as default", logLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)
	// Инициализация базы данных
	database.Init()
	logger.Infof("Database initialized successfully")

	// Создание экземпляров репозитория, сервиса и обработчиков
	repo := subs.NewRepository(logger)
	service := subs.NewService(repo, logger)
	handlers := subs.NewHandlers(service, logger)

	// Создание роутера
	router := gin.Default()

	// Регистрация обработчиков
	subsGroup := router.Group("/subs")
	{
		subsGroup.POST("", handlers.CreateSub)
		subsGroup.GET("/:id", handlers.GetSubByID)
		subsGroup.PUT("/:id", handlers.UpdateSub)
		subsGroup.DELETE("/:id", handlers.DeleteSub)
		subsGroup.GET("", handlers.ListSubs)
		subsGroup.GET("/total", handlers.GetTotalPriceForPeriod)
	}

	// Базовый эндпоинт
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello by Effective Mobile")
	})

	// Эндпоинт документации Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Infof("Documentation API: http://localhost:8080/swagger/index.html#/")

	// Запуск сервера
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // значение по умолчанию
	}
	logger.Infof("Server starting on port %s", port)

	if err := router.Run(":" + port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
