package controllers

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/yhanli/go-jwt-asymmetric/initializers"
	"github.com/yhanli/go-jwt-asymmetric/models"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get email and pass off req body
	var body struct {
		Email    string
		Password string
		Name     string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	// create the user
	user := models.User{Email: body.Email, Password: string(hash), Name: body.Name}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadKey(pemData []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemData)

	if block == nil {
		return nil, fmt.Errorf("unable to load key")
	}

	if block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("wrong type of key - %s", block.Type)
	}

	return x509.ParseECPrivateKey(block.Bytes)
}

func Login(c *gin.Context) {
	// get email and password
	var body struct {
		Email    string
		Password string
		Name     string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// look up user

	var user models.User
	initializers.DB.First(&user, "email=?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalide user or password",
		})
		return
	}

	// compare send pass hash and saved hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalide user or password",
		})
		return
	}

	// generate jwt token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": user.ID,
	// 	"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	// 	"kid": "{json file}", // kid claim is required to the middleware, but its also optional
	// })
	// tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	// if err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Failed to make token",
	// 	})
	// 	return
	// }

	// with private key
	// command to generate the rsa key pair
	// openssl genrsa -out cert/id_rsa 4096
	// openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub
	data, err := ioutil.ReadFile("./id_rsa")
	check(err)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	check(err)
	// privateKey := x509.MarshalPKCS1PrivateKey(key)
	claim := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"kid": "id_rsa.pub", // kid claim is required to the middleware, but its also optional
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	tokenString, err := t.SignedString(key)
	check(err)
	// send back
	// send as cookies
	c.SetSameSite(http.SameSiteLaxMode)
	// c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	// send as json
	c.JSON(http.StatusOK, gin.H{
		// "token": tokenString,
	})
}

func Validate(c *gin.Context) {
	_user, _ := c.Get("user")
	user, ok := _user.(models.User)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalide user ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "I'm logged in",
		"user":    user,
	})
}
