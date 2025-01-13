package cmd

import (
	"database/sql"
	"log"

	"github.com/iBoBoTi/aqua-sec-inventory/config"
	"github.com/iBoBoTi/aqua-sec-inventory/pkg/db"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed predefined cloud resources",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		conn, err := db.NewPostgresDB(cfg.DB)
		if err != nil {
			log.Fatalf("Could not connect to Postgres: %v", err)
		}
		defer conn.Close()

		if err := seedResources(conn); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Println("Seeding successful!")
	},
}

func init() {
	RootCmd.AddCommand(seedCmd)
}

func seedResources(db *sql.DB) error {
	resources := []struct {
		Name   string
		Type   string
		Region string
	}{
		{"aws_vpc_main", "VPC", "us-east-1"},
		{"gcp_vm_instance", "Compute", "us-central1"},
		{"azure_sql_db", "Database", "eastus"},
		// Add more as needed...
	}

	for _, r := range resources {
		_, err := db.Exec(`
            INSERT INTO resources (name, type, region, created_at, updated_at)
            VALUES ($1, $2, $3, NOW(), NOW())
            ON CONFLICT (name) DO NOTHING
        `, r.Name, r.Type, r.Region)
		if err != nil {
			return err
		}
	}
	return nil
}
