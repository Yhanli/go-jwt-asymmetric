package main

import (
	"github.com/yhanli/go-jwt-asymmetric/initializers"
	"github.com/yhanli/go-jwt-asymmetric/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate((&models.Post{}))
	initializers.DB.AutoMigrate((&models.User{}))
}
