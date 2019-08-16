package logger

type (
	// By implementing all interface below, you can control how your log works
	// e.g: if you use logger from third party, maybe you can use logger.Info
	// or logger.Debug ( provided in that 3rd party ).
	// But if you plan to change that 3rd party, the second 3rd party may have
	// different func (e.g: don't have `Debug` feature)
	// But if you implement all interface below in your logger, you don't need
	// to change how you write your log in your code.
	// In `./logger`, I have implement 2 logger, `logrus` and `stdLib log`
	// and you can choose which log engine do you want to use by overriding
	// `APP_LOG_ENGINE` env var
	AppLoggerInterface interface {
		SetLevel(l int)

		WithField(key string, value interface{}) AppLoggerInterface
		WithError(err error) AppLoggerInterface

		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
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
	appLogger = NewStdLogger()
)

func SetLogger(l AppLoggerInterface) {
	appLogger = l
}

func SetLevel(l int) {
	appLogger.SetLevel(l)
}

func WithField(key string, value interface{}) AppLoggerInterface {
	return appLogger.WithField(key, value)
}

func WithError(err error) AppLoggerInterface {
	return appLogger.WithError(err)
}

func Debug(args ...interface{}) {
	appLogger.Debug(args...)
}

func Info(args ...interface{}) {
	appLogger.Info(args...)
}

func Warn(args ...interface{}) {
	appLogger.Warn(args...)
}

func Error(args ...interface{}) {
	appLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	appLogger.Fatal(args...)
}
