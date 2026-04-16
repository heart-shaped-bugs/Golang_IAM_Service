package main

import (
	"database/sql"
	"embed"
	"errors"
	route "iam-service/api/http"
	"iam-service/api/http/handlers"
	pg_repo "iam-service/internal/repositories/postgres"
	auth "iam-service/internal/usecases"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

func main() {
	godotenv.Load(".env.example")

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	//MIGRATIONS
	// Создаём драйвер источника миграций
	d, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		log.Fatal("migration source:", err)
	}

	// Создаём драйвер БД
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("migration driver:", err)
	}

	// Запускаем миграции
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatal("migrate init:", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("migrate up:", err)
	}
	//END OF MIGRATIONS

	repo := pg_repo.NewUserRepository(db)

	authService := auth.NewAuthService(repo, jwtSecret)

	authHandler := handlers.NewAuthHandler(authService)

	router := route.NewRouter(authHandler)

	log.Println("Server starting on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
