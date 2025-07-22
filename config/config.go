package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerConfig ServerConfig
	LoggerConfig LoggerConfig
	DBConfig     DBConfig
	RedisConfig  RedisConfig
}

type ServerConfig struct {
	ServerPort    string `env:"SERVER_PORT"     envDefault:"8080"`
	JWTSecret     string `env:"JWT_SECRET"      envDefault:"secret"`
	RefreshSecret string `env:"REFRESH_SECRET"  envDefault:"refresh"`
	Environment   string `env:"ENVIRONMENT"     envDefault:"development"`
}

type LoggerConfig struct {
	Level      string `env:"LOG_LEVEL"       envDefault:"info"`
	Format     string `env:"LOG_FORMAT"      envDefault:"json"`
	OutputPath string `env:"LOG_OUTPUT_PATH" envDefault:"stdout"`
	Component  string `env:"LOG_COMPONENT"   envDefault:"default"`
}

type DBConfig struct {
	URL string `env:"DB_URL" envDefault:"postgres://user:password@localhost:5432/dbname?sslmode=disable"`
}

type RedisConfig struct {
	URL string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
