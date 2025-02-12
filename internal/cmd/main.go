package main

import (
	"avito-internship-2025/internal/config"
	"avito-internship-2025/internal/logger"
	"avito-internship-2025/internal/migrations"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
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
		migrationsPath := path.Join(os.Getenv("GOPATH"), "/migrations")
		migrations.RunMigrations(databaseURL, migrationsPath)
		os.Exit(0)
	}

	logg.Info("Запуск сервера")

	// TODO: app init next steps
}
