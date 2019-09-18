package logger

import (
	"os"

	"github.com/ztrue/tracerr"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func SetupLogger() {
	logger, err := zap.NewProduction()

	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		os.Exit(127)
	}
	defer logger.Sync()

	Logger = logger.Sugar()
}
