package logger

import (
	"SnickersShopPet1.0/internal/config"
	"go.uber.org/zap"
)

func LoadLogger(cfg *config.Config) (*zap.Logger, error) {

	switch cfg.Env {
	case "develop":
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		return cfg.Build()
	case "production":
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		return cfg.Build()
	default: //local
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		return cfg.Build()
	}

}
