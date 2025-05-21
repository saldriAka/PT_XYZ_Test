package connection

import (
	"fmt"
	"os"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/util"
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabase(conf config.Database) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")

	gormConfig := &gorm.Config{
		// Optional: set logger level, misalnya silent supaya tidak terlalu verbose
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	util.PanicIfError(err)

	sqlDB, err := db.DB()
	util.PanicIfError(err)

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	fmt.Println("✅ Successfully connected to the database using GORM")

	return db
}
