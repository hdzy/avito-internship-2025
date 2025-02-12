package migrations

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // импорт драйвера для PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // импорт источника файлов миграций
)

// RunMigrations принимает адрес базы данных и путь к файлам миграции
// и применяет миграции к базе данных
func RunMigrations(databaseURL, migrationsPath string) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Ошибка создания миграций: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
	log.Println("Миграции применены успешно")
}
