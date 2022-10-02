package initializers

import (
	"log"
	"os"

	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}
