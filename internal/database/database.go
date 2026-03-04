package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/tritrongnguyen/repo-reviewer.git/pkg/logger"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	reporeviewer "github.com/tritrongnguyen/repo-reviewer.git"
)

type Service interface {

	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	GetDB() *pgxpool.Pool
}

type service struct {
	pool *pgxpool.Pool
}

var (
	database   = os.Getenv("DB_DATABASE")
	username   = os.Getenv("DB_USERNAME")
	password   = os.Getenv("DB_PASSWORD")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Log.Fatal("Unable to create connection pool", zap.Error(err))
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Log.Fatal("Unable to ping database", zap.Error(err))
	}

	dbInstance = &service{
		pool: pool,
	}

	dbInstance.runMigrations()

	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)
	err := s.pool.Ping(ctx)

	if err != nil {
		stats["status"] = "down"
		return stats
	}

	sdt := s.pool.Stat()
	stats["status"] = "up"
	stats["total_conns"] = strconv.Itoa(int(sdt.TotalConns()))
	stats["idle_conns"] = strconv.Itoa(int(sdt.IdleConns()))
	stats["acquired_conns"] = strconv.Itoa(int(sdt.AcquiredConns()))
	stats["max_conns"] = strconv.Itoa(int(sdt.MaxConns()))

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)

	s.pool.Close()

	return nil
}

func (s *service) GetDB() *pgxpool.Pool {
	return s.pool
}

func (s *service) runMigrations() {
	d, err := iofs.New(reporeviewer.MigrationFiles, "migrations")
	if err != nil {
		logger.Log.Fatal("Error when init iofs migration", zap.Error(err))
	}

	sqlDB := stdlib.OpenDBFromPool(s.pool)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal("Error when init migrate driver", zap.Error(err))
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		logger.Log.Fatal("Error when init migrate instance", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal("Error when running migration", zap.Error(err))
	}

	logger.Log.Info("Database migration done successfully")
}
