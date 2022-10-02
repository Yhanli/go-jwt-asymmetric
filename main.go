package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yhanli/go-jwt-asymmetric/controllers"
	"github.com/yhanli/go-jwt-asymmetric/initializers"
	"github.com/yhanli/go-jwt-asymmetric/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	// Get data off req body

	// Create a quest

	//

	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.POST("/post", controllers.PostCreate)
	r.PUT("/update_post/:id", controllers.PostUpdate)
	r.GET("/post", controllers.PostIndex)
	r.GET("/post/:id", controllers.PostShow)
	r.DELETE("/post/:id", controllers.PostDelete)

	r.Run()
}
