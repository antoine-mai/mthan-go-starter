package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// DatabaseService handles the connection pool of the application database.
type DatabaseService struct {
	db     *sql.DB
	logger *LoggerService
	cfg    *Config
	mu     sync.Mutex
}

// NewDatabaseService initializes a DatabaseService.
func NewDatabaseService(cfg *Config, logger *LoggerService) *DatabaseService {
	return &DatabaseService{
		cfg:    cfg,
		logger: logger,
	}
}

// Connect opens and validates the database connection.
func (s *DatabaseService) Connect() (*sql.DB, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return s.db, nil
	}

	var driverName string
	switch s.cfg.DbDriver {
	case "postgres":
		driverName = "postgres"
		if s.cfg.DbUrl == "" {
			return nil, fmt.Errorf("postgres database driver requires a database URL to be configured in config.yaml (database.url)")
		}
		s.logger.Info("Connecting to Postgres database", "url", s.cfg.DbUrl)
	case "sqlite":
		fallthrough
	default:
		driverName = "sqlite"
		// Ensure the parent directory of the SQLite database file exists.
		dir := filepath.Dir(s.cfg.DbUrl)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for database file %s: %w", dir, err)
		}
		s.logger.Info("Connecting to SQLite database", "path", s.cfg.DbUrl)
	}

	db, err := sql.Open(driverName, s.cfg.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pooling limits
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Validate connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	s.logger.Info("Database connection established successfully")
	s.db = db
	return db, nil
}

// Close terminates the active database connection.
func (s *DatabaseService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db == nil {
		return nil
	}

	s.logger.Info("Closing database connection")
	err := s.db.Close()
	s.db = nil
	return err
}

// GetDB returns the active database connection pool.
func (s *DatabaseService) GetDB() *sql.DB {
	return s.db
}
