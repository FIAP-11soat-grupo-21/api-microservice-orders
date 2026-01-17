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

	AWS struct {
		Endpoint        string
		Region          string
		AccessKeyID     string
		SecretAccessKey string
	}

	MessageBroker struct {
		Type string // "sqs" ou "rabbitmq"

		// SQS
		SQS struct {
			UpdateOrderStatusQueueURL string
			OrderErrorQueueURL        string
		}

		SNS struct {
			OrderErrorTopicARN   string
			OrderCreatedTopicARN string
		}

		// RabbitMQ
		RabbitMQ struct {
			URL         string // (ex: amqp://user:pass@host:port/)
			OrdersQueue string
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

	// AWS Configuration
	c.AWS.Endpoint = getEnv("AWS_ENDPOINT", "")
	c.AWS.Region = getEnv("AWS_REGION", "us-east-2")
	c.AWS.AccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "")
	c.AWS.SecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")

	// Message Broker Configuration
	c.MessageBroker.Type = getEnv("MESSAGE_BROKER_TYPE", "sqs")

	// SQS
	c.MessageBroker.SQS.UpdateOrderStatusQueueURL = getEnv("SQS_UPDATE_ORDER_STATUS_QUEUE_URL", "")
	c.MessageBroker.SQS.OrderErrorQueueURL = getEnv("SQS_ORDER_ERROR_QUEUE_URL", "")

	// SNS
	c.MessageBroker.SNS.OrderErrorTopicARN = getEnv("SNS_ORDER_ERROR_TOPIC_ARN", "")
	c.MessageBroker.SNS.OrderCreatedTopicARN = getEnv("SNS_ORDER_CREATED_TOPIC_ARN", "")

	// RabbitMQ
	c.MessageBroker.RabbitMQ.URL = getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	c.MessageBroker.RabbitMQ.OrdersQueue = getEnv("RABBITMQ_ORDERS_QUEUE", "orders.updates")

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
