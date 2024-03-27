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
	"github.com/rs/cors" // Import the cors package
	"gopkg.in/yaml.v2"
)

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
	// Read config file and set up database connection
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	// Create an instance of the scheduler and handler
	scheduler := app.NewSJFScheduler(db)
	handler := handlers.NewHandler(scheduler)

	// Create a router and wrap it with CORS middleware
	router := server.NewRouter(handler, 8099)
	c := cors.Default().Handler(router.Router())

	// Create and start the HTTP server
	s := server.NewServer(":8099", c, ":8098")

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
