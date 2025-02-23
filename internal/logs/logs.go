package logs

import (
	"atomwoz.com/remitly_task/internal/config"
	"github.com/charmbracelet/log"
)

//These functions are wrappers for the underlying loggin lib used to log messages to the console.

func Log(format string, args ...interface{}) {
	if config.GetDebugMode() {
		log.Infof(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Warn(format string, args ...interface{}) {
	if config.GetDebugMode() {
		log.Warnf(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
