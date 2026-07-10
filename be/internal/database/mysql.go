package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig holds MySQL configuration.
type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewMySQL creates a new MySQL database connection.
func NewMySQL(cfg MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=false",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	log.Println("MySQL connection established")

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
		log.Printf("MySQL connection pool: max_open=%d, max_idle=%d", maxOpen, maxIdle)
	}

	return db, nil
}
