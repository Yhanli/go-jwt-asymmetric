package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yhanli/go-jwt-asymmetric/initializers"
	"github.com/yhanli/go-jwt-asymmetric/models"
)

func PostCreate(c *gin.Context) {

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"post": post,
	})

}

func PostIndex(c *gin.Context) {
	// Get the post
	var posts []models.Post
	initializers.DB.Find(&posts)

	// respond
	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func PostShow(c *gin.Context) {
	// Get id of post
	id := c.Param("id")
	// Get the post
	var post models.Post
	initializers.DB.First(&post, id)

	// respond
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {
	// Get id of post
	id := c.Param("id")
	// Get data off request body

	var body struct {
		Body  string
		Title string
	}
	c.Bind(&body)
	// find the post
	fmt.Println(body)
	var post models.Post
	initializers.DB.First(&post, id)

	// update it
	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title, Body: body.Body,
	})

	// post.Title = body.Title
	// post.Body = body.Body

	// initializers.DB.Save(&post)

	// respond
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostDelete(c *gin.Context) {
	// Get id of post
	id := c.Param("id")
	// Delete the post

	initializers.DB.Delete(&models.Post{}, id)

	// respond
	c.Status(200)
}
