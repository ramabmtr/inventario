package config

type env struct {
	App struct {
		Address   string
		Debug     bool
		Name      string
		LogEngine string
	}
	Database struct {
		Engine string
		URL    string
	}
}

var Env env

func InitEnvVar() {
	Env.App.Address = GetEnv("APP_ADDRESS").SetDefault(":8080").ToString()
	Env.App.Debug = GetEnv("APP_DEBUG").SetDefault("1").ToBool()
	Env.App.Name = GetEnv("APP_NAME").SetDefault("Inventario").ToString()
	Env.App.LogEngine = GetEnv("APP_LOG_ENGINE").OneOf(LogEngineLogrus, LogEngineStdlib).SetDefault(LogEngineLogrus).ToString()

	Env.Database.Engine = GetEnv("DATABASE_ENGINE").OneOf(DatabaseEngineSqlite3).SetDefault(DatabaseEngineSqlite3).ToString()
	Env.Database.URL = GetEnv("DATABASE_URL").SetDefault("inventario.db").ToString()
}
