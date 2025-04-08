package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Arcanm/deliveryPlannerGolang/internal/application/services"
	"github.com/Arcanm/deliveryPlannerGolang/internal/domain/repositories"
	"github.com/Arcanm/deliveryPlannerGolang/internal/infrastructure/persistence/mongodb"
	grpcimpl "github.com/Arcanm/deliveryPlannerGolang/internal/interfaces/grpc"
	"github.com/Arcanm/deliveryPlannerGolang/internal/interfaces/http/handlers"
	"github.com/Arcanm/deliveryPlannerGolang/proto"
)

func main() {
	// Initialize MongoDB connection
	mongoClient, err := mongodb.NewClient()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Get database name from environment variable or use default
	dbName := "delivery_planner"
	if name := os.Getenv("MONGODB_DB"); name != "" {
		dbName = name
	}
	db := mongoClient.Database(dbName)

	// Initialize repositories
	driverRepo := repositories.NewDriverRepository(db)
	packageRepo := repositories.NewPackageRepository(db)
	routeRepo := repositories.NewRouteRepository(db)

	// Initialize services
	driverService := services.NewDriverService(driverRepo, routeRepo)
	packageService := services.NewPackageService(packageRepo)
	routeService := services.NewRouteService(routeRepo, driverRepo, packageRepo)

	// Initialize gRPC server
	grpcServer := grpcserver.NewServer()

	// Register gRPC services
	proto.RegisterDriverServiceServer(grpcServer, grpcimpl.NewDriverService(driverService))
	proto.RegisterPackageServiceServer(grpcServer, grpcimpl.NewPackageService(packageService, routeService))
	proto.RegisterRouteServiceServer(grpcServer, grpcimpl.NewRouteService(routeService))

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("Failed to listen for gRPC:", err)
	}

	// Initialize HTTP server
	router := gin.Default()

	// Initialize HTTP handlers
	driverHandler := handlers.NewDriverHandler(driverService)
	packageHandler := handlers.NewPackageHandler(packageService, routeService)
	routeHandler := handlers.NewRouteHandler(routeService)

	// Register HTTP routes
	driverHandler.RegisterRoutes(router)
	packageHandler.RegisterRoutes(router)
	routeHandler.RegisterRoutes(router)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Setup HTTP port
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Create channels for graceful shutdown
	grpcErr := make(chan error, 1)
	httpErr := make(chan error, 1)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			grpcErr <- err
		}
	}()

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := router.Run(":" + httpPort); err != nil {
			httpErr <- err
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either server error or interrupt signal
	select {
	case err := <-grpcErr:
		log.Printf("gRPC server error: %v", err)
	case err := <-httpErr:
		log.Printf("HTTP server error: %v", err)
	case <-quit:
		log.Println("Shutting down servers...")
		grpcServer.GracefulStop()
	}
}
