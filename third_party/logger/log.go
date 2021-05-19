package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Logger logrus.FieldLogger

var Log logrus.FieldLogger

func BuildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}

type Util struct {
	// Stdout is true if the output needs to goto standard out
	Stdout bool `yaml:"stdout"`
	// Level is the desired log level
	Level string `yaml:"level"`
	// OutputFile is the path to the log output file
	OutputFile string `yaml:"outputFile"`
}

const fileMode = os.FileMode(0644)

// NewLogger newLogger builds and returns a new bark
// logger for this logging configuration
func (cfg *Util) NewLogger() *logrus.Logger {

	logger := logrus.New()
	logger.Out = ioutil.Discard
	logger.Level = parseLogrusLevel(cfg.Level)
	logger.Formatter = getFormatter()

	if cfg.Stdout {
		logger.Out = os.Stdout
	}

	if len(cfg.OutputFile) > 0 {
		outFile := createLogFile(cfg.OutputFile)
		logger.Out = outFile
		if cfg.Stdout {
			logger.Out = io.MultiWriter(os.Stdout, outFile)
		}
	}

	return logger
}

func getFormatter() logrus.Formatter {
	formatter := &logrus.TextFormatter{}
	formatter.FullTimestamp = true
	return formatter
}

func createLogFile(path string) *os.File {
	dir := filepath.Dir(path)
	if len(dir) > 0 && dir != "." {
		if err := os.MkdirAll(dir, fileMode); err != nil {
			log.Fatalf("error creating log directory %v, err=%v", dir, err)
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, fileMode)
	if err != nil {
		log.Fatalf("error creating log file %v, err=%v", path, err)
	}
	return file
}

// parseLogrusLevel converts the string log
// level into a logrus level
func parseLogrusLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
