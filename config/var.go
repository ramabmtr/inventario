package config

const (
	LogEngineStdlib = "stdlib"
	LogEngineLogrus = "logrus"

	DatabaseEngineSqlite3 = "sqlite3"

	DefaultLimit  = 10
	DefaultOffset = 0

	IncomingTransactionType = "IN"
	OutgoingTransactionType = "OUT"

	QueryDateFormatLayout = "2006-01-02"
	ISO8601Format         = "2006-01-02T15:04:05.000Z"
)
