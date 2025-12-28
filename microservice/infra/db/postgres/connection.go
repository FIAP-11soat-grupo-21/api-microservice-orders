package postgres

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"microservice/infra/db/postgres/models"
	"microservice/utils/config"
)

var (
	dbConnection *gorm.DB
	instance     *gorm.DB
	once         sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		instance = dbConnection
	})
	return instance
}

func Connect() {
	if dbConnection != nil {
		log.Println("Database connection already established")
		return
	}

	cfg := config.LoadConfig()

	dsn := "host=" + cfg.Database.Host +
		" user=" + cfg.Database.Username +
		" dbname=" + cfg.Database.Name +
		" password=" + cfg.Database.Password +
		" port=" + cfg.Database.Port

	queryLogLevel := logger.Info
	if cfg.IsProduction() {
		queryLogLevel = logger.Error
	}

	queryLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  queryLogLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 2 * time.Second

	for i := range maxRetries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: queryLogger,
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	dbConnection = db
}

func Close() {
	if dbConnection == nil {
		log.Println("Database connection already closed")
		return
	}

	sqlDriver, err := dbConnection.DB()
	if err != nil {
		log.Fatal("Failed to close database")
	}

	sqlDriver.Close()
}

func RunMigrations() {
	dbConnection.AutoMigrate(
		&models.OrderModel{},
		&models.OrderItemModel{},
		&models.OrderStatusModel{},
	)
}
