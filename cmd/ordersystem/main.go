package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thiagomaganha/go-clean-architecture/configs"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/database"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/web"
	"github.com/thiagomaganha/go-clean-architecture/internal/infra/web/webserver"
	"github.com/thiagomaganha/go-clean-architecture/internal/usecase"
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

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(listOrdersUseCase)
	webserver.AddHandler("/order", webOrderHandler.ListOrders)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start()
}

func runMigrations(db *sql.DB, dbName string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("could not start sql migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://sql/migrations",
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
