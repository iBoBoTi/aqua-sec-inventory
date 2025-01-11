package main

import (
    "log"
    "os"

    "github.com/spf13/cobra"

    "github.com/iBoBoTi/aqua-sec-inventory/cmd"
    "github.com/iBoBoTi/aqua-sec-inventory/config"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/repository"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/service"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/transport/rest"
    "github.com/iBoBoTi/aqua-sec-inventory/internal/usecase"
    "github.com/iBoBoTi/aqua-sec-inventory/pkg/db"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Start the Aqua Security Cloud Resource Inventory server",
    Run: func(c *cobra.Command, args []string) {
        // Load config
        cfg := config.LoadConfig()

        // Init DB
        pgDB, err := db.NewPostgresDB(cfg.DB)
        if err != nil {
            log.Fatalf("Could not connect to Postgres: %v", err)
        }
        defer pgDB.Close()

        // Init Repositories
        customerRepo := repository.NewCustomerRepository(pgDB)
        resourceRepo := repository.NewResourceRepository(pgDB)
        notificationRepo := repository.NewNotificationRepository(pgDB)

        // Init Usecases
        customerUC := usecase.NewCustomerUsecase(customerRepo)
        resourceUC := usecase.NewResourceUsecase(resourceRepo, customerRepo)
        notificationUC := usecase.NewNotificationUsecase(notificationRepo)

        // Initialize RabbitMQ (or any MQ) for notifications
		log.Println("RabbitMQ URL: ", cfg.RabbitMQ.URL)
        notifier, err := service.NewRabbitMQNotifier(cfg.RabbitMQ.URL)
        if err != nil {
            log.Fatalf("Could not connect to RabbitMQ: %v", err)
        }
        defer notifier.Close()

        // Start listening for notifications in a separate goroutine
        go func() {
            if err := notifier.Listen(); err != nil {
                log.Printf("[WARNING] Notification listener stopped: %v\n", err)
            }
        }()

        // Setup Gin Router
        router := rest.NewRouter(customerUC, resourceUC, notificationUC, notifier)

        // Start HTTP server
        log.Printf("Server is running on port %s", cfg.Server.Port)
        if err := router.Run(":" + cfg.Server.Port); err != nil {
            log.Fatalf("Server error: %v", err)
        }
    },
}

func main() {
    root := &cobra.Command{Use: "aqua-sec-cloud-inventory"}
    root.AddCommand(serverCmd)
    // Attach other subcommands from cmd package
    root.AddCommand(cmd.RootCmd.Commands()...)

    if err := root.Execute(); err != nil {
        os.Exit(1)
    }
}
