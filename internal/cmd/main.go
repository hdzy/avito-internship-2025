package main

import (
	"avito-internship-2025/internal/config"
	"avito-internship-2025/internal/handlers"
	"avito-internship-2025/internal/logger"
	"avito-internship-2025/internal/migrations"
	"avito-internship-2025/internal/repository"
	"avito-internship-2025/internal/service"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Применить миграции и завершить работу")
	flag.Parse()

	cfg, err := config.LoadConfig("./config/")
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.Init(cfg.LogLevel, cfg.LogFormat)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("Логгер инициализирован", slog.String("env", cfg.Env))

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)

	if *migrateFlag {
		logg.Info("Запуск миграций", slog.String("migrationsPath", "./migrations"))
		migrationsPath := "./migrations"
		migrations.RunMigrations(databaseURL, migrationsPath)
		os.Exit(0)
	}

	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		logg.Error("Ошибка подключения к БД", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logg.Error("Ошибка подключения к БД (ping)", slog.Any("error", err))
		os.Exit(1)
	}
	logg.Info("Подключение к БД установлено")

	employeesRepo := repository.NewEmployeeRepository(db, logg)
	authService := service.NewAuthService(employeesRepo, cfg.JWTSecret, logg)
	authHandler := handlers.NewAuthHandler(authService, logg)

	http.HandleFunc("/api/auth", authHandler.Authenticate)

	logg.Info("Запуск сервера на порту", slog.Int("port", cfg.ServerPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), nil); err != nil {
		logg.Error("Ошибка запуска сервера", slog.Any("error", err))
		os.Exit(1)
	}
}
