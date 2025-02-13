package main

import (
	"avito-internship-2025/internal/config"
	"avito-internship-2025/internal/logger"
	"avito-internship-2025/internal/migrations"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
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

	logg.Info("Запуск сервера")
}
