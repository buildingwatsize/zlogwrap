package zlogwrap

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

const ( // Todo: It's should be a config
	ServiceNameTag   = "service"
	TransactionIDTag = "transaction-id"
	URLTag           = "url"
)

type zerologWrapper interface {
	SetField(key string, anything interface{}) zerologWrapper // Set field in logs
	Debug(anything ...interface{})                            // level 0
	Info(anything ...interface{})                             // level 1
	Warn(anything ...interface{})                             // level 2
	Error(anything ...interface{})                            // level 3
	Fatal(anything ...interface{})                            // level 4
	Panic(anything ...interface{})                            // level 5
	GetLogEvent(zerolog.Level) *zerolog.Event                 // With Caller (file:line)
}

func (c Config) SetField(key string, anything interface{}) zerologWrapper {
	if key == "" {
		return &c
	}
	switch v := anything.(type) {
	case bool:
		c.Logger = c.Logger.With().Bool(key, v).Logger()
	case int:
		c.Logger = c.Logger.With().Int(key, v).Logger()
	case int64:
		c.Logger = c.Logger.With().Int64(key, v).Logger()
	case float64:
		c.Logger = c.Logger.With().Float64(key, v).Logger()
	case []byte:
		c.Logger = c.Logger.With().RawJSON(key, v).Logger()
	case []int:
		c.Logger = c.Logger.With().Ints(key, v).Logger()
	case []int64:
		c.Logger = c.Logger.With().Ints64(key, v).Logger()
	case []float64:
		c.Logger = c.Logger.With().Floats64(key, v).Logger()
	case []string:
		c.Logger = c.Logger.With().Strs(key, v).Logger()
	default:
		c.Logger = c.Logger.With().Str(key, fmt.Sprintf("%v", anything)).Logger()
	}
	return &c
}

func (c *Config) createLogTemplate(zLevel zerolog.Level) *zerolog.Event {
	var logTemplate *zerolog.Event

	logTemplate = c.Logger.WithLevel(zLevel)

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

	logTemplate := c.createLogTemplate(zerolog.DebugLevel)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Info(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(zerolog.InfoLevel)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Warn(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(zerolog.WarnLevel)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Error(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(zerolog.ErrorLevel)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Fatal(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(zerolog.FatalLevel)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) Panic(anything ...interface{}) {
	if c.Hidden {
		return
	}

	logString := toString(anything...)

	logTemplate := c.createLogTemplate(zerolog.PanicLevel)
	logTemplate.Msgf("%v", logString)
}

func (c *Config) GetLogEvent(level zerolog.Level) *zerolog.Event {
	return c.Logger.WithLevel(level)
}
