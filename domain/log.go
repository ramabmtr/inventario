package domain

type (
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
