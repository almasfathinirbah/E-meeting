package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"e_meeting/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Service represents a service that interacts with a database.
// It provides methods for health checking, connection management,
// and access to both standard SQL and GORM database instances.
type Service interface {
	// Health returns a map of health status information including:
	// - connection status (up/down)
	// - error message if any
	// - number of open connections
	// - number of in-use connections
	// - number of idle connections
	Health() map[string]string

	// Close terminates the database connection and releases all resources.
	Close() error

	// DB returns the underlying sql.DB instance for raw SQL operations.
	DB() *sql.DB

	// GormDB returns the GORM database instance for ORM operations.
	GormDB() *gorm.DB
}

// service implements the Service interface and holds the database connections
type service struct {
	db     *sql.DB  // Standard SQL database connection
	gormDB *gorm.DB // GORM ORM database connection
}

var (
	dbInstance *service
	once       sync.Once
)

// New creates a new database service instance with the provided configuration.
// It ensures only one database connection is created using the singleton pattern.
// Parameters:
//   - cfg: Configuration containing database connection parameters
//
// Returns:
//   - A Service interface implementation for database operations
func New(cfg *config.Config) Service {
	once.Do(func() {
		// Create GORM connection
		gormDB, err := NewPostgresDB(
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			fmt.Sprintf("%d", cfg.DBPort),
		)
		if err != nil {
			log.Fatalf("Failed to create GORM connection: %v", err)
		}

		// Get underlying sql.DB
		sqlDB, err := gormDB.DB()
		if err != nil {
			log.Fatalf("Failed to get sql.DB from GORM: %v", err)
		}

		// Set connection pool settings
		sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConnections)
		sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConnections)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(30 * time.Minute)

		dbInstance = &service{
			db:     sqlDB,
			gormDB: gormDB,
		}
	})

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", s.gormDB.Migrator().CurrentDatabase())
	return s.db.Close()
}

func (s *service) DB() *sql.DB {
	return s.db
}

func (s *service) GormDB() *gorm.DB {
	return s.gormDB
}

// NewPostgresDB establishes a new connection to a PostgreSQL database using GORM.
// Parameters:
//   - host: Database server hostname
//   - user: Database username
//   - password: Database password
//   - dbname: Name of the database to connect to
//   - port: Database server port
//
// Returns:
//   - A configured GORM DB instance and any error that occurred during connection
func NewPostgresDB(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(zerolog.NewConsoleWriter(), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}
