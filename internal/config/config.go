package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host         string        `env:"SERVER_HOST" envDefault:"localhost"`
	Port         int           `env:"SERVER_PORT" envDefault:"8080"`
	Version      string        `env:"SERVER_VERSION" envDefault:"1.0.0"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	BodyLimit    int           `env:"SERVER_BODY_LIMIT" envDefault:"4194304"` // 4MB
}

type DatabaseConfig struct {
	Host            string        `env:"DB_HOST" envDefault:"localhost" validate:"required"`
	Port            int           `env:"DB_PORT" envDefault:"5432" validate:"min=1,max=65535"`
	User            string        `env:"DB_USER" envDefault:"postgres" validate:"required"`
	Password        string        `env:"DB_PASSWORD" envDefault:"" validate:"required"`
	Name            string        `env:"DB_NAME" envDefault:"webhook_api" validate:"required"`
	SSLMode         string        `env:"DB_SSL_MODE" envDefault:"disable" validate:"oneof=disable require verify-ca verify-full"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25" validate:"min=1"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5" validate:"min=1"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m"`
}

type RedisConfig struct {
	Host         string        `env:"REDIS_HOST" envDefault:"localhost" validate:"required"`
	Port         int           `env:"REDIS_PORT" envDefault:"6379" validate:"min=1,max=65535"`
	Password     string        `env:"REDIS_PASSWORD" envDefault:""`
	DB           int           `env:"REDIS_DB" envDefault:"0" validate:"min=0,max=15"`
	PoolSize     int           `env:"REDIS_POOL_SIZE" envDefault:"10" validate:"min=1"`
	MinIdleConns int           `env:"REDIS_MIN_IDLE_CONNS" envDefault:"5" validate:"min=1"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
}

type JWTConfig struct {
	SecretKey            string        `env:"JWT_SECRET_KEY" validate:"required,min=32"`
	AccessTokenDuration  time.Duration `env:"JWT_ACCESS_TOKEN_DURATION" envDefault:"24h"`
	RefreshTokenDuration time.Duration `env:"JWT_REFRESH_TOKEN_DURATION" envDefault:"168h"` // 7 days
	Issuer               string        `env:"JWT_ISSUER" envDefault:"webhook-api"`
	CachePrefix          string        `env:"JWT_CACHE_PREFIX" envDefault:"jwt:"`
}

type LoggerConfig struct {
	Level       string `env:"LOG_LEVEL" envDefault:"info" validate:"oneof=debug info warn error"`
	Format      string `env:"LOG_FORMAT" envDefault:"json" validate:"oneof=json console"`
	Development bool   `env:"LOG_DEVELOPMENT" envDefault:"false"`
	Caller      bool   `env:"LOG_CALLER" envDefault:"true"`
	Stacktrace  bool   `env:"LOG_STACKTRACE" envDefault:"false"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *DatabaseConfig) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
