package apiserver

type Config struct {
	BindAddr    string
	LogLevel    string
	DatabaseURL string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":9000",
		LogLevel:    "debug",
		DatabaseURL: "host=localhost dbname=avito-job sslmode=disable port=5432 password=1234 user=d",
	}
}
