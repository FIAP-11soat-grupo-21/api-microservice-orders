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

	MessageBroker struct {
		Type string // "sqs" ou "rabbitmq"
		
		// SQS
		SQS struct {
			PaymentQueueURL string 
			KitchenQueueURL string 
			AWSRegion       string
		}
		
		// RabbitMQ
		RabbitMQ struct {
			URL          string // (ex: amqp://user:pass@host:port/)
			PaymentQueue string
			KitchenQueue string
		}
	}
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
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

	// Message Broker Configuration
	c.MessageBroker.Type = getEnv("MESSAGE_BROKER_TYPE", "sqs")
	
	// SQS
	c.MessageBroker.SQS.PaymentQueueURL = getEnv("SQS_PAYMENT_QUEUE_URL", "")
	c.MessageBroker.SQS.KitchenQueueURL = getEnv("SQS_KITCHEN_QUEUE_URL", "")
	c.MessageBroker.SQS.AWSRegion = getEnv("AWS_REGION", "us-east-2")
	
	// RabbitMQ
	c.MessageBroker.RabbitMQ.URL = getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	c.MessageBroker.RabbitMQ.PaymentQueue = getEnv("RABBITMQ_PAYMENT_QUEUE", "payment.confirmation")
	c.MessageBroker.RabbitMQ.KitchenQueue = getEnv("RABBITMQ_KITCHEN_QUEUE", "kitchen.orders")

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
