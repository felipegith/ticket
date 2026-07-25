package config

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
type Config struct {
	PGCONECTIONSTRING string
	REDISCONFIG       RedisConfig
	RABBITMQURL       string
	ELASTICSEARCHURL  string
}

func Load() Config {
	return Config{
		PGCONECTIONSTRING: getEnviromentVariable("PGCONECTIONSTRING", "host=localhost port=5499 user=root password=123456 dbname=ticketing sslmode=disable"),
		REDISCONFIG: RedisConfig{
			Host:     getEnviromentVariable("REDIS_HOST", "localhost"),
			Port:     getEnviromentVariableInt("REDIS_PORT", 6379),
			Password: getEnviromentVariable("REDIS_PASSWORD", ""),
			DB:       getEnviromentVariableInt("REDIS_DB", 0),
		},
		RABBITMQURL:      getEnviromentVariable("RABBITMQ_URL", "amqp://root:123456@localhost:5677/"),
		ELASTICSEARCHURL: getEnviromentVariable("ELASTICSEARCH_URL", "http://localhost:9299"),
	}
}

func getEnviromentVariable(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnviromentVariableInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return fallback
}
