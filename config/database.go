package config

import (
	"log"
	"ticert/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {
	cfg := GetConfig()

	dsn := cfg.GetDatabaseDSN()

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}

	log.Printf("Connected to MySQL database successfully at %s:%s", cfg.DBHost, cfg.DBPort)

	if cfg.DBPassword == "" {
		log.Println("Warning: DB_PASSWORD is not set. Using default password.")
	}

	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	log.Println("Database migration completed")
}
