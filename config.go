package zlogwrap

import (
	"strings"

	"github.com/rs/zerolog"
)

const ( // Todo: It's should be a config
	ServiceNameTag   = "service"
	TransactionIDTag = "transaction-id"
	URLTag           = "url"
)

const (
	logDebugType = "debug"
	logInfoType  = "info"
	logWarnType  = "warn"
	logErrorType = "error"
	logFatalType = "fatal"
	logPanicType = "panic"
)

type zerologWrapper interface {
	Debug(anything ...interface{}) // level 0
	Info(anything ...interface{})  // level 1
	Warn(anything ...interface{})  // level 2
	Error(anything ...interface{}) // level 3
	Fatal(anything ...interface{}) // level 4
	Panic(anything ...interface{}) // level 5
}

func (c *Config) createLogTemplate(typeLog string) *zerolog.Event {
	var logTemplate *zerolog.Event
	switch typeLog {
	case logDebugType:
		logTemplate = c.Logger.Debug()
	case logInfoType:
		logTemplate = c.Logger.Info()
	case logWarnType:
		logTemplate = c.Logger.Warn()
	case logErrorType:
		logTemplate = c.Logger.Error()
	case logFatalType:
		logTemplate = c.Logger.Fatal()
	case logPanicType:
		logTemplate = c.Logger.Panic()
	default:
		logTemplate = c.Logger.Log()
	}

	if c.ServiceName != "" {
		logTemplate = logTemplate.Str(ServiceNameTag, c.ServiceName)
	}

	if c.Context != nil {
		if txID := string(c.Context.Response().Header.Peek(TransactionIDTag)); txID != "" {
			logTemplate = logTemplate.Str(strings.ReplaceAll(TransactionIDTag, "-", "_"), txID)
		}

		if url := c.Context.OriginalURL(); url != "" {
			logTemplate = logTemplate.Str(URLTag, url)
		}
	}

	return logTemplate
}

func (c *Config) Debug(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(logDebugType)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Info(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(logInfoType)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Warn(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(logWarnType)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Error(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(logErrorType)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Fatal(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(logFatalType)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Panic(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(logPanicType)
	logTemplate.Msgf("%v", logString)
}
