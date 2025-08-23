package main

import (
	"app/internal/database"
	"app/internal/subs"
	"net/http"

	_ "app/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Swagger Users Subscribtions
// @version 1.0
// @description REST-сервис для агрегации данных об онлайн-подписках пользователей
// @contact.url	https://github.com/achi3v3
// @contact.email aamir-tutaev@mail.ru
// @host localhost:8080
// @BasePath /
func main() {
	// Настройка логирования
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	// Инициализация базы данных
	database.Init()
	logrus.Infof("Database initialized successfully")

	// Создание экземпляров репозитория, сервиса и обработчиков
	repo := subs.NewRepository()
	service := subs.NewService(repo)
	handlers := subs.NewHandlers(service)

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
	logrus.Infof("Documentation API: http://localhost:8080/swagger/index.html#/")

	// Запуск сервера
	logrus.Infof("Server starting on :8080")

	if err := router.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
