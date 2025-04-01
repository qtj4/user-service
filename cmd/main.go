package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"

	"github.com/yourusername/user-service/config"
	"github.com/yourusername/user-service/internal/event"
	"github.com/yourusername/user-service/internal/handler"
	"github.com/yourusername/user-service/internal/repository"
	"github.com/yourusername/user-service/internal/service"
	"github.com/yourusername/user-service/proto"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// PostgreSQL connection
	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// RabbitMQ connection
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	eventPublisher := event.NewEventPublisher(ch)
	userService := service.NewUserService(userRepo, eventPublisher)
	restHandler := handler.NewUserRESTHandler(userService)
	grpcHandler := handler.NewUserGRPCServer(userService)

	// Start REST server
	go func() {
		r := gin.Default()
		r.GET("/users/:id", restHandler.GetUser)
		r.GET("/users/:id/orders", restHandler.GetUserOrders)
		r.POST("/users", restHandler.CreateUser)
		r.PUT("/users/:id", restHandler.UpdateUser)
		r.DELETE("/users/:id", restHandler.DeleteUser)
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, grpcHandler)
	log.Printf("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}