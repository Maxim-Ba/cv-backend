package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testDB хранит соединение с тестовой БД для всех тестов
var testDB *sql.DB

// pgContainer хранит ссылку на контейнер PostgreSQL
var pgContainer *postgres.PostgresContainer

// TestMain настраивает тестовое окружение с PostgreSQL контейнером
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Запускаем PostgreSQL контейнер
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	pgContainer = container

	// Получаем строку подключения
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	// Подключаемся к БД
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Проверяем соединение
	if err := testDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Применяем миграции
	if err := runMigrations(testDB); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Запускаем тесты
	code := m.Run()

	// Очистка
	if err := testDB.Close(); err != nil {
		log.Printf("failed to close database: %v", err)
	}

	if err := container.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}

// runMigrations применяет миграции к тестовой БД
func runMigrations(db *sql.DB) error {
	// Создаем таблицы из миграций
	migrations := []string{
		// 0001_init.up.sql
		`CREATE TABLE IF NOT EXISTS tag (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			hex_color TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS education (
			id BIGSERIAL PRIMARY KEY,
			name TEXT,
			year INT NOT NULL,
			course TEXT NOT NULL,
			organization TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS technology (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL UNIQUE,
			description TEXT,
			logo_url TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS technologies_tag (
			tag_id BIGINT NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
			technology_id BIGINT NOT NULL REFERENCES technology (id) ON DELETE CASCADE,
			PRIMARY KEY (tag_id, technology_id)
		)`,
		`CREATE TABLE IF NOT EXISTS work_history (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			about TEXT NOT NULL,
			logo_url BYTEA,
			period_start DATE NOT NULL,
			period_end DATE,
			what_i_did TEXT[],
			projects TEXT[]
		)`,
		`CREATE TABLE IF NOT EXISTS work_history_technology (
			work_history_id BIGINT NOT NULL REFERENCES work_history (id) ON DELETE CASCADE,
			technology_id BIGINT NOT NULL REFERENCES technology (id) ON DELETE CASCADE,
			PRIMARY KEY (work_history_id, technology_id)
		)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w\nQuery: %s", err, migration)
		}
	}

	return nil
}

// cleanupTable очищает указанную таблицу и сбрасывает счетчик ID
func cleanupTable(t *testing.T, tableName string) {
	t.Helper()

	_, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
	if err != nil {
		t.Fatalf("failed to cleanup table %s: %v", tableName, err)
	}

	// Сбрасываем sequence если она существует
	_, _ = testDB.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", tableName))
}

// cleanupAllTables очищает все таблицы перед тестом
func cleanupAllTables(t *testing.T) {
	t.Helper()

	tables := []string{
		"work_history_technology",
		"technologies_tag",
		"work_history",
		"technology",
		"education",
		"tag",
	}

	for _, table := range tables {
		_, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("failed to cleanup table %s: %v", table, err)
		}
	}

	// Сбрасываем sequences
	sequences := []string{"tag_id_seq", "education_id_seq", "technology_id_seq", "work_history_id_seq"}
	for _, seq := range sequences {
		_, _ = testDB.Exec(fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq))
	}
}
