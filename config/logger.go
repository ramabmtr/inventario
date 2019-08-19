package config

import "github.com/ramabmtr/inventario/repository/logger"

func InitLogger() {
	// setup logger with engine defined in `APP_LOG_ENGINE` env var
	switch Env.App.LogEngine {
	case LogEngineLogrus:
		logger.SetLogger(logger.NewLogrusLogger())
	case LogEngineStdlib:
		fallthrough
	default:
		logger.SetLogger(logger.NewLogrusLogger())
	}

	if Env.App.Debug {
		logger.SetLevel(logger.DebugLevel)
	}
}
