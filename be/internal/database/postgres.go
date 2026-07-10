package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQLConfig holds PostgreSQL configuration.
type PostgreSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewPostgreSQL creates a new PostgreSQL database connection.
func NewPostgreSQL(cfg PostgreSQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Println("PostgreSQL connection established")

	sqlDB, err := db.DB()
	if err == nil {
		maxOpenStr := os.Getenv("DB_MAX_OPEN_CONNS")
		if maxOpenStr == "" {
			maxOpenStr = "25"
		}
		maxIdleStr := os.Getenv("DB_MAX_IDLE_CONNS")
		if maxIdleStr == "" {
			maxIdleStr = "10"
		}
		maxOpen, _ := strconv.Atoi(maxOpenStr)
		maxIdle, _ := strconv.Atoi(maxIdleStr)
		sqlDB.SetMaxOpenConns(maxOpen)
		sqlDB.SetMaxIdleConns(maxIdle)
		log.Printf("PostgreSQL connection pool: max_open=%d, max_idle=%d", maxOpen, maxIdle)
	}

	return db, nil
}
