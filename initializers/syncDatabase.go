package initializers

import (
	"github.com/yhanli/go-jwt-asymmetric/models"
)

func SyncDatabase() {
	DB.AutoMigrate((&models.Post{}))
	DB.AutoMigrate((&models.User{}))
}
