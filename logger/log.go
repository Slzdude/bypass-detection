package logger

import (
	"github.com/apsdehal/go-logger"
	"os"
)

var Logger *logger.Logger

func InitLogger() (err error) {
	Logger, err = logger.New("MITM Server", 1, os.Stdout, logger.DebugLevel)
	if err != nil {
		return
	}
	return
}
