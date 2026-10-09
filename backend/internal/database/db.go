package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"lastmile-backend/internal/config"
)

type DB struct {
	*sql.DB
}

func Connect(cfg *config.Config) (*DB, error) {
	conn, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database at %s:%s: %w", cfg.DBHost, cfg.DBPort, err)
	}

	log.Printf("Connected successfully to PostgreSQL [%s:%s/%s]", cfg.DBHost, cfg.DBPort, cfg.DBName)
	db := &DB{conn}

	// Ejecutar migraciones iniciales si corresponde
	if err := db.RunMigrations(); err != nil {
		log.Printf("Warning during migrations: %v", err)
	}

	return db, nil
}

func (db *DB) RunMigrations() error {
	migrationFiles := []string{
		"migrations/001_init.sql",
		"migrations/002_seed.sql",
		"migrations/003_users.sql",
		"migrations/004_tenants.sql",
		"migrations/005_tenant_fleet.sql",
		"migrations/006_tenant_orders_routes.sql",
		"migrations/007_user_driver_link.sql",
		"backend/migrations/001_init.sql",
		"backend/migrations/002_seed.sql",
		"backend/migrations/003_users.sql",
		"backend/migrations/004_tenants.sql",
		"backend/migrations/005_tenant_fleet.sql",
		"backend/migrations/006_tenant_orders_routes.sql",
		"backend/migrations/007_user_driver_link.sql",
		"../migrations/001_init.sql",
		"../migrations/002_seed.sql",
		"../migrations/003_users.sql",
		"../migrations/004_tenants.sql",
		"../migrations/005_tenant_fleet.sql",
		"../migrations/006_tenant_orders_routes.sql",
		"../migrations/007_user_driver_link.sql",
	}

	executed := make(map[string]bool)

	for _, file := range migrationFiles {
		cleanName := filepath.Base(file)
		if executed[cleanName] {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		log.Printf("Executing migration: %s", cleanName)
		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			trimmed := strings.TrimSpace(q)
			if trimmed == "" {
				continue
			}
			if _, err := db.Exec(trimmed); err != nil {
				// Log y continuar si ya existía
				log.Printf("Migration step notice in %s: %v", cleanName, err)
			}
		}
		executed[cleanName] = true
	}

	return nil
}
