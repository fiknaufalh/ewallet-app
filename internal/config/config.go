package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	App         AppConfig
	Transaction TransactionConfig
	Security    SecurityConfig
}

type ServerConfig struct {
	Port string
	Host string
	Mode string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

type TransactionConfig struct {
	MaxWithdrawalAmount float64
	MinWithdrawalAmount float64
	MaxTopupAmount      float64
	MinTopupAmount      float64
}

type SecurityConfig struct {
	IdempotencyKeyExpiration time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
		fmt.Println("No config file found, using environment variables")
	}

	setDefaults()

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
			Host: viper.GetString("SERVER_HOST"),
			Mode: viper.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetString("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Name:            viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSL_MODE"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		App: AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Version:     viper.GetString("APP_VERSION"),
			Environment: viper.GetString("APP_ENVIRONMENT"),
		},
		Transaction: TransactionConfig{
			MaxWithdrawalAmount: viper.GetFloat64("MAX_WITHDRAWAL_AMOUNT"),
			MinWithdrawalAmount: viper.GetFloat64("MIN_WITHDRAWAL_AMOUNT"),
			MaxTopupAmount:      viper.GetFloat64("MAX_TOPUP_AMOUNT"),
			MinTopupAmount:      viper.GetFloat64("MIN_TOPUP_AMOUNT"),
		},
		Security: SecurityConfig{
			IdempotencyKeyExpiration: viper.GetDuration("IDEMPOTENCY_KEY_EXPIRATION"),
		},
	}

	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("GIN_MODE", "debug")

	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "ewallet")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	// Transaction defaults
	viper.SetDefault("MAX_WITHDRAWAL_AMOUNT", 10000000)
	viper.SetDefault("MIN_WITHDRAWAL_AMOUNT", 10000)
	viper.SetDefault("MAX_TOPUP_AMOUNT", 50000000)
	viper.SetDefault("MIN_TOPUP_AMOUNT", 10000)

	// Security defaults
	viper.SetDefault("IDEMPOTENCY_KEY_EXPIRATION", "24h")
}