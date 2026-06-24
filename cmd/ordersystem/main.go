package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/thiagomaganha/go-clean-architecture/configs"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/database"
	grpchandler "github.com/thiagomaganha/go-clean-architecture/internal/infra/grpc/handler"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/grpc/pb"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/web"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/web/webserver"
	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
	embeddedSQL "github.com/thiagomaganha/go-clean-architecture/sql"
	"google.golang.org/grpc"
)

func main() {
	configs, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName)
	db, err := sql.Open(configs.DBDriver, dbConn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = runMigrations(db, configs.DBName); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	orderRepository := database.NewOrderRepository(db)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)

	// gRPC server
	grpcServer := grpc.NewServer()
	orderGrpcHandler := grpchandler.NewOrderGrpcHandler(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderGrpcHandler)

	go func() {
		lis, err := net.Listen("tcp", configs.GRPCServerPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port: %v", err)
		}
		fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// HTTP server
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(listOrdersUseCase, createOrderUseCase)
	webserver.AddHandler("GET", "/order", webOrderHandler.ListOrders)
	webserver.AddHandler("POST", "/order", webOrderHandler.CreateOrder)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start()
}

func runMigrations(db *sql.DB, dbName string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("could not start sql migration driver: %w", err)
	}

	src, err := iofs.New(embeddedSQL.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("could not open migrations source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		src,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not setup migrate instance: %w", err)
	}

	log.Println("Running database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}
