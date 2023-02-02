package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"

	Consumer  = "trigger_consumer"
	AllErrors = "consumer1.error" // errors
	AllInfo   = "consumer2.info"  // info
	AllDebug  = "consumer3.debug" // debug
	All       = "#"               // all

)

type Config struct {
	ServiceName string
	Environment string // debug, test, release
	Version     string

	DefaultOffset string
	DefaultLimit  string

	//PostgresMaxConnections int32
	RabbitMQHost     string
	RabbitMQPort     int
	RabbitMQUser     string
	RabbitMQPassword string

	RestServiceHost string
	RestServicePort int

	LogLevel string
}

// Load ...
func Load() Config {
	if err := godotenv.Load("/app/.env"); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.LogLevel = cast.ToString(getOrReturnDefaultValue("LOG_LEVEL", "debug"))

	config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "epa_go_api_gateway"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))

	//config.PostgresMaxConnections = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_MAX_CONNECTIONS", 30))
	config.RabbitMQUser = cast.ToString(getOrReturnDefaultValue("RABBIT_MQ_USER", "guest"))
	config.RabbitMQPassword = cast.ToString(getOrReturnDefaultValue("RABBIT_MQ_PASSWORD", "guest"))
	config.RabbitMQHost = cast.ToString(getOrReturnDefaultValue("RABBIT_MQ_HOST", "localhost"))
	config.RabbitMQPort = cast.ToInt(getOrReturnDefaultValue("RABBIT_MQ_PORT", 5672))

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10000000"))

	config.RestServiceHost = cast.ToString(getOrReturnDefaultValue("REST_SERVICE_HOST", "localhost"))
	config.RestServicePort = cast.ToInt(getOrReturnDefaultValue("REST_SERVICE_PORT", 8000))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
