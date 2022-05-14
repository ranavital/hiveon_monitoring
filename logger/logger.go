package logger

import (
	"hiveon_monitoring/config"
	"io"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type Log struct {
	file          *os.File
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
}

var Logging Log

func Init() error {
	file, err := os.OpenFile(config.AppConfig.LoggerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, file)
	Logging.DebugLogger = log.New(mw, "DEBUG: ", log.Ldate|log.Ltime)
	Logging.InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime)
	Logging.WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime)
	Logging.ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime)
	Logging.file = file
	return nil
}

func Close() {
	Logging.file.Close()
}

func (l *Log) Debug(format string, params ...interface{}) {
	l.DebugLogger.Printf(format, params...)

}

func (l *Log) Info(format string, params ...interface{}) {
	l.InfoLogger.Printf(format, params...)

}

func (l *Log) Warning(format string, params ...interface{}) {
	l.WarningLogger.Printf(format, params...)
}

func (l *Log) Error(format string, params ...interface{}) {
	l.ErrorLogger.Printf(format, params...)
}
