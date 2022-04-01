package logger

import "log"

// By default, log everything
const LOG_LEVEL = 2

// Usage:
// Logger is a wrapper struct that allows utility in what is being
// logged for a given level or task
// 0 -> Don't log anything
// 1 -> Log the API details
// 2 -> Log API details and file matches (all the things)

type Logger struct {
	LoggingLevel int
}

func Log(loggingStr string, loggingLevel int) {
	logger := &Logger{
		LoggingLevel: LOG_LEVEL,
	}
	switch loggingLevel {
	case 1:
		logger.Info(loggingStr)
	case 2:
		logger.Warn(loggingStr)
	}
}

func (l *Logger) Info(logStr string) {
	if l.LoggingLevel >= 1 {
		log.Println(logStr)
	}
}

func (l *Logger) Warn(logStr string) {
	if l.LoggingLevel >= 2 {
		log.Println(logStr)
	}
}
