// go run migration.go up 1(version)
// go run migration.go  1(version)

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type customLogger struct{}

func (l *customLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *customLogger) Verbose() bool {
	return true
}

func main() {
	fmt.Println("--Migrate Start--")
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbURL := os.Getenv("POSTGRES_URI")
	m, err := migrate.New(
		"file://../migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	m.Log = &customLogger{}

	args := os.Args
	if len(args) < 3 {
		log.Fatalf("Usage: go run main.go [up|down] [target_version]")
	}

	command := args[1]
	targetVersion, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("Invalid target version: %v", err)
	}

	switch command {
	case "up":
		log.Printf("Migrating up to version %d...", targetVersion)
		err = m.Migrate(uint(targetVersion))
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate up to version %d: %v", targetVersion, err)
		}
		log.Printf("Successfully migrated up to version %d", targetVersion)

	case "down":
		if targetVersion == 0 {
			log.Fatalf("Cannot rollback to version 0. Use 'm.Down()' to reset all migrations.")
		}

		log.Printf("Rolling back down to version %d...", targetVersion)
		err = m.Migrate(uint(targetVersion))
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback down to version %d: %v", targetVersion, err)
		}
		log.Printf("Successfully rolled back to version %d", targetVersion)

	default:
		log.Fatalf("Invalid command: %s. Use 'up' or 'down'", command)
	}

	log.Println("--Migrate Complete--")
}
