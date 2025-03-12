package database

import (
	"context"
	"fmt"
	"log"
	"net"
	"starter/internal/config"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	databaseName := "database"
	databaseUsername := "username"
	databasePassword := "password"

	// Create a new pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// Check if Docker is running
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Create the Postgres container with options
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", databaseName),
			fmt.Sprintf("POSTGRES_USER=%s", databaseUsername),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", databasePassword),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// Auto-remove container when stopped
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Get host and port
	hostAndPort := resource.GetHostPort("5432/tcp")
	host, port, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		log.Fatalf("Could not split hostAndPort: %s", err)
	}

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// Create the config
	cfg = &config.Config{
		App: config.App{
			Port: "5432",
			Env:  "test",
		},
		Database: config.Database{
			Host:     host,
			Port:     port,
			Database: databaseName,
			Username: databaseUsername,
			Password: databasePassword,
			Schema:   "public",
		},
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database, cfg.Database.Schema)
	log.Println("Connecting to database on url: ", databaseURL)

	var db *pgxpool.Pool

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dbInstance = &service{
		db: &database{
			pool: db,
		},
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func TestNew(t *testing.T) {
	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("New() returned error: %s", err)
	}
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	ctx := context.Background()
	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("New() returned error: %s", err)
	}

	stats, err := srv.Health(ctx)
	if err != nil {
		t.Fatalf("expected Health to return nil, go %s", err)
	}

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "bing chilling" {
		t.Fatalf("expected message to be 'bing chilling', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv, err := New(cfg)
	if err != nil {
		t.Fatalf("New() returned error: %s", err)
	}

	srv.Close()
}
