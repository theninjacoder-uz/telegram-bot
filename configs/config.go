package configs

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	conf *Configuration
	once sync.Once
)

// Config loads configuration using atomic pattern
func Config() *Configuration {
	once.Do(func() {
		conf = load()
	})
	return conf
}

// Configuration ...
type Configuration struct {
	HTTPPort    string `json:"http_port"`
	LogLevel    string `json:"log_level"`
	Environment string `json:"environment"`

	ServerPort       string
	ServerHost       string
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	// context timeout in seconds
	CtxTimeout      int
	TelgramBotToken string
	TelgramBotURI   string
	RpcHost         string
	RpcAuthLogin    string
	RpcAuthPassword string
}

func load() *Configuration {

	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	var config Configuration

	v := viper.New()
	v.AutomaticEnv()

	config.Environment = v.GetString("ENVIRONMENT")
	config.HTTPPort = v.GetString("HTTP_PORT")
	config.LogLevel = v.GetString("LOG_LEVEL")
	config.CtxTimeout = v.GetInt("CONTEXT_TIMEOUT")
	config.TelgramBotURI = v.GetString("TELEGRAM_BOT_URI")
	config.TelgramBotToken = v.GetString("TELEGRAM_BOT_TOKEN")
	config.ServerHost = v.GetString("SERVER_HOST")
	config.ServerPort = v.GetString("SERVER_PORT")
	config.PostgresDatabase = v.GetString("POSTGRES_DB")
	config.PostgresUser = v.GetString("POSTGRES_USER")
	config.PostgresPassword = v.GetString("POSTGRES_PASSWORD")
	config.PostgresHost = v.GetString("POSTGRES_HOST")
	config.PostgresPort = v.GetInt("POSTGRES_PORT")
	config.RpcHost = v.GetString("RPC_HOST")
	config.RpcAuthLogin = v.GetString("RPC_AUTH_LOGIN")
	config.RpcAuthPassword = v.GetString("RPC_AUTH_PASSWORD")

	return &config
}
