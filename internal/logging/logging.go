package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hvs-fasya/micro/internal/configure"
)

//SetLoggers set global production logger + special logger for requests histoy
func SetLoggers() *zap.Logger {
	var logger, _ = zap.NewDevelopment()
	if configure.Cfg.Env == "prod" {
		globalLogCfg := zap.NewProductionConfig()
		globalLogCfg.DisableStacktrace = true
		logger, _ = globalLogCfg.Build()
	}
	zap.ReplaceGlobals(logger) //logger shared global logger for common purposes - stderr

	requestsLogCfg := &zap.Config{ // special logger for requests history - logs to file (for study purposes)
		Encoding:          "console",
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:       []string{configure.Cfg.Server.RequestsLog},
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	if configure.Cfg.Env != "prod" {
		requestsLogCfg.OutputPaths = append(requestsLogCfg.OutputPaths, "stdout") //log also to stdout if not prod
	}
	requestsLogger, _ := requestsLogCfg.Build()
	return requestsLogger
}
