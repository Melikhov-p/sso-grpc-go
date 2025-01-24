package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	lvlInfo  = "INFO"
	lvlDebug = "DEBUG"

	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// BuildLogger строит zap.Logger с необходимым уровнем логирования.
func SetupLogger(env string) (*zap.Logger, error) {
	// выбираем уровень логирования в зависимости от окружения
	var logLevel string
	switch env {
	case envLocal:
		logLevel = lvlDebug
	case envDev:
		logLevel = lvlDebug
	case envProd:
		logLevel = lvlInfo
	default:
		logLevel = lvlDebug
	}
	lvl, err := zap.ParseAtomicLevel(logLevel)

	if err != nil {
		return nil, fmt.Errorf("failed to parse atomic level: %w", err)
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()

	// Настройка формата времени
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	// Отключение stacktrace
	cfg.DisableStacktrace = true

	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build config %w", err)
	}
	return zl, nil
}
