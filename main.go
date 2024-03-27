package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sfjgoapic/app"
	"sfjgoapic/handlers"
	"sfjgoapic/server"
	"syscall"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

// Config represents the structure of the YAML configuration file
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

func main() {
	// Load configuration from config.yaml file
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Create connection string for PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Ping the database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	// Create a scheduler instance
	scheduler := app.NewSJFScheduler(db)

	// Create a handler instance using NewHandler function from handlers package
	handler := handlers.NewHandler(scheduler)

	router := server.NewRouter(handler, 8099)

	s := server.NewServer(":8099", router.Router(), ":8098")

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v. Shutting down gracefully...", sig)
		cancel()
	}()

	// Start the server
	go func() {
		log.Println("Starting server on :8099")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for the context to be cancelled (graceful shutdown)
	<-ctx.Done()
	log.Println("Shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server stopped gracefully.")
}
