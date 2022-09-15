package config

// Config holds the application configuration values
type Config struct {
	HTTPServerConfig      *HTTPServerConfig
	MySQLRepositoryConfig *MySQLRepositoryConfig
}

// HTTPServerConfig holds the http server configuration values
type HTTPServerConfig struct {
	HTTPServerHost string
	WaitToShutdown uint
}

// MySQLRepositoryConfig holds the MySQL configuration values
type MySQLRepositoryConfig struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBType     string
}
