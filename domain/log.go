package domain

type (
	// By implementing all interface below, you can control how your log works
	// e.g: if you use logger from third party, maybe you can use logger.Info
	// or logger.Debug ( provided in that 3rd party ).
	// But if you plan to change that 3rd party, the second 3rd party may have
	// different func (e.g: don't have `Debug` feature)
	// But if you implement all interface below in your logger, you don't need
	// to change how you write your log in your code.
	// In `./config/logger.go`, I have implement 2 logger, `logrus` and `stdLib log`
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
