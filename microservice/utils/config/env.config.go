package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv   string
	APIPort string
	APIHost string

	Database struct {
		RunMigrations bool
		Host          string
		Name          string
		Port          string
		Username      string
		Password      string
	}

	RabbitMQ struct {
		Host         string
		Port         string
		User         string
		Password     string
		PaymentQueue string
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func LoadConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.Load()
	})
	return instance
}

func (c *Config) Load() *Config {
	dotEnvPath := ".env"
	_, err := os.Stat(dotEnvPath)

	if err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	c.GoEnv = getEnv("GO_ENV")
	c.APIPort = getEnv("API_PORT")
	c.APIHost = getEnv("API_HOST")

	c.Database.RunMigrations = getEnv("DB_RUN_MIGRATIONS") == "true"
	c.Database.Host = getEnv("DB_HOST")
	c.Database.Name = getEnv("DB_NAME")
	c.Database.Port = getEnv("DB_PORT")
	c.Database.Username = getEnv("DB_USERNAME")
	c.Database.Password = getEnv("DB_PASSWORD")

	c.RabbitMQ.Host = getEnv("RABBITMQ_HOST")
	c.RabbitMQ.Port = getEnv("RABBITMQ_PORT")
	c.RabbitMQ.User = getEnv("RABBITMQ_USER")
	c.RabbitMQ.Password = getEnv("RABBITMQ_PASSWORD")
	c.RabbitMQ.PaymentQueue = getEnv("RABBITMQ_PAYMENT_QUEUE")

	return c
}

func (c *Config) IsProduction() bool {
	return c.GoEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.GoEnv == "development"
}

var (
	instance *Config
	once     sync.Once
)
