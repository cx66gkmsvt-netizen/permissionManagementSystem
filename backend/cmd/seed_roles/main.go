package main

import (
	"log"
	"os"
	"strings"

	"user-center/internal/config"
	"user-center/internal/repository"
)

func main() {
	// 1. Load Config
	// Warning: config.Load() might rely on current working directory to find config.yaml
	cfg := config.Load()

	// 2. Init DB
	if err := repository.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	db := repository.GetDB()

	// 3. Read SQL File
	// Assuming execution from backend directory
	sqlPath := "migrations/seed_roles_users.sql"
	content, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatalf("Failed to read SQL file %s: %v", sqlPath, err)
	}

	// 4. Split and Execute
	// Basic split by semicolon. Might break if semicolons are in strings,
	// but our generated SQL is simple statements.
	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			// Log error but continue (e.g. for duplicates)
			log.Printf("Warning executing statement: %v", err)
		}
	}

	log.Println("Seeding completed successfully.")
}
