package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MamangRust/payment-gateway-monolith-grpc/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewClient(logger logger.LoggerInterface) (*sql.DB, error) {
	dbDriver := viper.GetString("DB_DRIVER")

	var connStr string
	switch dbDriver {
	case "postgres":
		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			viper.GetString("DB_HOST"),
			viper.GetString("DB_PORT"),
			viper.GetString("DB_USERNAME"),
			viper.GetString("DB_NAME"),
			viper.GetString("DB_PASSWORD"),
		)
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			viper.GetString("DB_USERNAME"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_HOST"),
			viper.GetString("DB_PORT"),
			viper.GetString("DB_NAME"),
		)
	default:
		logger.Error("Unsupported database driver", zap.String("DB_DRIVER", dbDriver))
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	con, err := sql.Open(dbDriver, connStr)

	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := con.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	maxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}
	con.SetMaxOpenConns(maxOpenConns)

	maxIdleConns := viper.GetInt("DB_MAX_IDLE_CONNS")
	if maxIdleConns == 0 {
		maxIdleConns = 5
	}
	con.SetMaxIdleConns(maxIdleConns)

	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
	}
	con.SetConnMaxLifetime(connMaxLifetime)

	logger.Debug("Database connection established successfully with connection pool settings",
		zap.String("DB_DRIVER", dbDriver),
		zap.Int("MaxOpenConns", maxOpenConns),
		zap.Int("MaxIdleConns", maxIdleConns),
		zap.Duration("ConnMaxLifetime", connMaxLifetime),
	)
	return con, nil
}
