package cmd

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"

	"github.com/iBoBoTi/aqua-sec-inventory/config"
	"github.com/iBoBoTi/aqua-sec-inventory/pkg/db"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		conn, err := db.NewPostgresDB(cfg.DB)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer conn.Close()

		runMigrations(conn, cfg.DB.MigrationsPath)
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}

func runMigrations(db *sql.DB, migrationsPath string) {
	goose.SetBaseFS(embedMigrations)
	if migrationsPath == "" {
		absPath, err := filepath.Abs("migrations")
		if err != nil {
			log.Fatalf("Failed to resolve migrations path: %v", err)
		}
		migrationsPath = absPath
	}

	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
}
