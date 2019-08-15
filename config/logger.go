package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ramabmtr/inventario/domain"
	"github.com/sirupsen/logrus"
)

type (
	logrusLoggerRepository struct {
		param  map[string]interface{}
		logger *logrus.Entry
	}

	stdLoggerRepository struct {
		param  map[string]interface{}
		logger *log.Logger
		level  int
	}
)

const (
	FatalLevel int = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var (
	// set app log to use stdlib by default
	// can be overridden by ENV Var, by calling `InitLogger()`
	AppLogger = NewStdLogger()
)

func InitLogger() {
	switch Env.App.LogEngine {
	case LogEngineLogrus:
		AppLogger = NewLogrusLogger()
	case LogEngineStdlib:
		fallthrough
	default:
		AppLogger = NewStdLogger()
	}
}

// NewLogrusLogger create appLogger based on 3rd party logger (sirupsen/logrus)
func NewLogrusLogger() domain.AppLoggerInterface {
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

func (c *logrusLoggerRepository) WithField(key string, val interface{}) domain.AppLoggerInterface {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param[key] = val
	return c
}

func (c *logrusLoggerRepository) WithError(err error) domain.AppLoggerInterface {
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

// NewStdLogger create appLogger based on golang std lib logger
func NewStdLogger() domain.AppLoggerInterface {
	return &stdLoggerRepository{
		param:  make(map[string]interface{}),
		logger: log.New(os.Stderr, "", log.LstdFlags),
		level:  InfoLevel,
	}
}

func (c *stdLoggerRepository) processParam() []interface{} {
	args := make([]interface{}, 0)
	for k, v := range c.param {
		args = append(args, fmt.Sprintf("[%s: %v", k, v))
	}
	c.param = nil
	return args
}

func (c *stdLoggerRepository) isLevelEnabled(i int) bool {
	return c.level >= i
}

func (c *stdLoggerRepository) SetLevel(l int) {
	c.level = l
}

func (c *stdLoggerRepository) WithField(key string, val interface{}) domain.AppLoggerInterface {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param[key] = val
	return c
}

func (c *stdLoggerRepository) WithError(err error) domain.AppLoggerInterface {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param["err_message"] = err.Error()
	return c
}

func (c *stdLoggerRepository) Debug(args ...interface{}) {
	if c.isLevelEnabled(DebugLevel) {
		level := []interface{}{"[DEBUG]"}
		param := c.processParam()
		args = append(level, args...)
		args = append(args, param...)
		c.logger.Println(args...)
	}
}

func (c *stdLoggerRepository) Info(args ...interface{}) {
	if c.isLevelEnabled(InfoLevel) {
		level := []interface{}{"[INFO ]"}
		param := c.processParam()
		args = append(level, args...)
		args = append(args, param...)
		c.logger.Println(args...)
	}
}

func (c *stdLoggerRepository) Warn(args ...interface{}) {
	if c.isLevelEnabled(WarnLevel) {
		level := []interface{}{"[WARN ]"}
		param := c.processParam()
		args = append(level, args...)
		args = append(args, param...)
		c.logger.Println(args...)
	}
}

func (c *stdLoggerRepository) Error(args ...interface{}) {
	if c.isLevelEnabled(ErrorLevel) {
		level := []interface{}{"[ERROR]"}
		param := c.processParam()
		args = append(level, args...)
		args = append(args, param...)
		c.logger.Println(args...)
	}
}

func (c *stdLoggerRepository) Fatal(args ...interface{}) {
	if c.isLevelEnabled(FatalLevel) {
		level := []interface{}{"[FATAL]"}
		param := c.processParam()
		args = append(level, args...)
		args = append(args, param...)
		c.logger.Fatal(args...)
	}
}
