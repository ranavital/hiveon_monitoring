package logger

import (
	"fmt"
	"hiveon_monitoring/config"
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

	Logging.DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime)
	Logging.InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	Logging.WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	Logging.ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	Logging.file = file
	return nil
}

func Close() {
	Logging.file.Close()
}

func (l *Log) Debug(format string, params ...interface{}) {
	l.DebugLogger.Printf(format, params...)
	fmt.Printf(format, params...)
	fmt.Println()
}

func (l *Log) Info(format string, params ...interface{}) {
	l.InfoLogger.Printf(format, params...)
	fmt.Printf(format, params...)
	fmt.Println()
}

func (l *Log) Warning(format string, params ...interface{}) {
	l.WarningLogger.Printf(format, params...)
	fmt.Printf(format, params...)
	fmt.Println()
}

func (l *Log) Error(format string, params ...interface{}) {
	l.ErrorLogger.Printf(format, params...)
	fmt.Printf(format, params...)
	fmt.Println()
}
