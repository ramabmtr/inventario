package logger

import (
	"github.com/sirupsen/logrus"
)

type (
	logrusLoggerRepository struct {
		param  map[string]interface{}
		logger *logrus.Entry
	}
)

// NewLogrusLogger create appLogger based on 3rd party logger (sirupsen/logrus)
func NewLogrusLogger() AppLoggerInterface {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	customFormatter.FullTimestamp = true
	std := logrus.StandardLogger()
	std.SetFormatter(customFormatter)

	return &logrusLoggerRepository{
		param:  make(map[string]interface{}),
		logger: logrus.NewEntry(std),
	}
}

func (c *logrusLoggerRepository) processParam() *logrus.Entry {
	l := c.logger
	for k, v := range c.param {
		l = l.WithField(k, v)
	}

	c.param = nil
	return l
}

func (c *logrusLoggerRepository) SetLevel(l int) {
	lvl := logrus.InfoLevel
	switch l {
	case DebugLevel:
		lvl = logrus.DebugLevel
	case FatalLevel:
		lvl = logrus.FatalLevel
	case ErrorLevel:
		lvl = logrus.ErrorLevel
	case WarnLevel:
		lvl = logrus.WarnLevel
	case InfoLevel:
		fallthrough
	default:
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)
}

func (c *logrusLoggerRepository) WithField(key string, val interface{}) AppLoggerInterface {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param[key] = val
	return c
}

func (c *logrusLoggerRepository) WithError(err error) AppLoggerInterface {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param["err_message"] = err.Error()
	return c
}

func (c *logrusLoggerRepository) Debug(args ...interface{}) {
	l := c.processParam()
	l.Debug(args...)
}

func (c *logrusLoggerRepository) Info(args ...interface{}) {
	l := c.processParam()
	l.Info(args...)
}

func (c *logrusLoggerRepository) Warn(args ...interface{}) {
	l := c.processParam()
	l.Warn(args...)
}

func (c *logrusLoggerRepository) Error(args ...interface{}) {
	l := c.processParam()
	l.Error(args...)
}

func (c *logrusLoggerRepository) Fatal(args ...interface{}) {
	l := c.processParam()
	l.Fatal(args...)
}
