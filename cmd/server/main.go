package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"
	"promoCode/internal/config"
	"promoCode/internal/handlers"
	"promoCode/internal/service"
	"promoCode/internal/storage"
	"promoCode/internal/storage/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
	}
	slog.Info("Configuration loaded successfully")

	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error closing database connection", "error", err)
		}
	}()

	if err := migrations.RunMigrations(db.DB); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}
	slog.Info("Migrations applied successfully")

	promoService := service.NewPromoService(db)
	adminHandler := handlers.NewAdminHandler(promoService)
	apiHandler := handlers.NewAPIHandler(promoService)

	router := gin.Default()

	router.POST("/admin/create", adminHandler.CreatePromoCode)
	router.GET("/admin", adminHandler.AdminPage)
	router.POST("/api/apply", apiHandler.ApplyPromoCode)

	slog.Info("Server started", "address", ":8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
