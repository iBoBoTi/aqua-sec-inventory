package main

import (
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/iBoBoTi/aqua-sec-inventory/cmd"
	"github.com/iBoBoTi/aqua-sec-inventory/config"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/repository"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/service"
	grpc2 "github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/transport/grpc"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/transport/rest"
	"github.com/iBoBoTi/aqua-sec-inventory/internal/notification-service/usecase"
	"github.com/iBoBoTi/aqua-sec-inventory/pkg/db"
	pb "github.com/iBoBoTi/aqua-sec-inventory/proto/notification"
)

var serverCmd = &cobra.Command{
	Use:   "notification-server",
	Short: "Start the Aqua Security Cloud Resource Inventory Notification Server",
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
		notificationRepo := repository.NewNotificationRepository(pgDB)

		// Init Usecases
		notificationUC := usecase.NewNotificationUsecase(notificationRepo)

		// Initialize RabbitMQ for notifications
		log.Println("RabbitMQ URL: ", cfg.RabbitMQ.URL)
		notifier, err := service.NewRabbitMQNotifier(cfg.RabbitMQ.URL, notificationRepo)
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
		router := rest.NewRouter(notificationUC, notifier)

		// Start Rest HTTP server in a goroutine
		go func() {
			log.Printf("Notification Server is running on port %s", cfg.Server.Port)
			if err := router.Run(":" + cfg.Server.Port); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}()

		// Setup and start gRPC server
		grpcServer := grpc.NewServer()
		grpcNotificationService := grpc2.NewNotificationGRPCService(notificationUC)
		pb.RegisterNotificationServiceServer(grpcServer, grpcNotificationService)

		reflection.Register(grpcServer)

		// Listen on gRPC port
		grpcAddr := cfg.GRPCServer.Port
		listener, err := net.Listen("tcp", ":"+grpcAddr)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port %s: %v", grpcAddr, err)
		}

		log.Printf("Notification gRPC server is running on port %s", grpcAddr)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC Server error: %v", err)
		}
	},
}

func main() {
	root := &cobra.Command{Use: "aqua-sec-cloud-inventory-notification"}
	root.AddCommand(serverCmd)
	// Attach other subcommands from cmd package
	root.AddCommand(cmd.RootCmd.Commands()...)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
