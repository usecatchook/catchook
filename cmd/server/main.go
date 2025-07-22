package server

import (
	"log"

	"github.com/theotruvelot/catchook/config"
	"github.com/theotruvelot/catchook/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.InitLogger(logger.Config{
		Level:      cfg.LoggerConfig.Level,
		Format:     cfg.LoggerConfig.Format,
		OutputPath: cfg.LoggerConfig.OutputPath,
		Component:  cfg.LoggerConfig.Component,
	})
}
