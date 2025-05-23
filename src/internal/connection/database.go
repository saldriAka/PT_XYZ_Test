package connection

import (
	"fmt"
	"os"
	"time"

	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabase(conf config.Database) *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("DSN:", dsn)

	var db *gorm.DB
	var err error

	// Retry logic: coba koneksi ke DB max 10x
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				// Set koneksi pool
				sqlDB.SetMaxIdleConns(5)
				sqlDB.SetMaxOpenConns(20)
				sqlDB.SetConnMaxLifetime(60 * time.Minute)
				sqlDB.SetConnMaxIdleTime(10 * time.Minute)
				fmt.Println("✅ Successfully connected to the database using GORM")
				return db
			}
		}

		fmt.Printf("⏳ Retry %d: waiting for database connection...\n", i+1)
		time.Sleep(3 * time.Second)
	}

	// Jika masih error setelah 10x, panic
	util.PanicIfError(err)
	return nil
}
