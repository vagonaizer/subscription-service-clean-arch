package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/vagonaizer/effective-mobile/subscription-service/internal/config"
)

const (
	defaultConfigPath    = "configs/config.yaml"
	defaultMigrationsDir = "file://internal/infrastructure/database/postgres/migrations"
)

func main() {
	var (
		configPath    = flag.String("config", defaultConfigPath, "path to configuration file")
		migrationsDir = flag.String("migrations-dir", defaultMigrationsDir, "path to migrations directory")
		action        = flag.String("action", "up", "migration action: up, down, version, force")
		steps         = flag.Int("steps", 0, "number of steps for up/down migration")
		version       = flag.Int("version", 0, "target version for migration")
	)
	flag.Parse()

	if envConfigPath := os.Getenv("CONFIG_PATH"); envConfigPath != "" {
		*configPath = envConfigPath
	}

	cfg := config.NewConfig()
	if err := cfg.Load(*configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := cfg.Database.DSN()
	log.Printf("connecting to database: %s", hidePassword(dsn))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		*migrationsDir,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "version":
		var currentVersion uint
		var dirty bool
		currentVersion, dirty, err = m.Version()
		if err != nil {
			log.Fatalf("failed to get current version: %v", err)
		}
		log.Printf("current version: %d, dirty: %t", currentVersion, dirty)
		return
	case "force":
		if *version < 0 {
			log.Fatal("version must be specified (>= 0) for force action")
		}
		err = m.Force(*version)
	default:
		log.Fatalf("unknown action: %s", *action)
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no migrations to apply")
			return
		}
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migration completed successfully")
}

func hidePassword(dsn string) string {
	return "postgres://user:***@host:port/dbname?sslmode=disable"
}
