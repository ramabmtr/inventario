package logger

import (
	"fmt"
	"log"
	"os"
)

type (
	stdLoggerRepository struct {
		param  map[string]interface{}
		logger *log.Logger
		level  int
	}
)

// NewStdLogger create appLogger based on golang std lib logger
func NewStdLogger() AppLoggerInterface {
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

func (c *stdLoggerRepository) WithField(key string, val interface{}) AppLoggerInterface {
	if c.param == nil {
		c.param = make(map[string]interface{})
	}
	c.param[key] = val
	return c
}

func (c *stdLoggerRepository) WithError(err error) AppLoggerInterface {
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
