package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewClient creates a new database connection client.
//
// It takes a logger.LoggerInterface as argument and returns an *sql.DB
// and an error. The logger is used to log any errors that occur when
// connecting to the database.
//
// The function uses the viper package to retrieve the database connection
// settings from the configuration file. The connection settings used are:
//
// - DB_DRIVER: the database driver to use (postgres or mysql)
// - DB_HOST: the hostname of the database server
// - DB_PORT: the port number of the database server
// - DB_USERNAME: the username to use when connecting to the database
// - DB_PASSWORD: the password to use when connecting to the database
// - DB_NAME: the name of the database to connect to
// - DB_MAX_OPEN_CONNS: the maximum number of open connections to the database
// - DB_MAX_IDLE_CONNS: the maximum number of idle connections to the database
// - DB_CONN_MAX_LIFETIME: the maximum lifetime of a connection to the database
//
// The function will return an error if any of the connection settings are
// invalid or if the connection to the database fails.
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
