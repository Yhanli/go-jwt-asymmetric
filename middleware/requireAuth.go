package middleware

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
)

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

func RequireAuth(c *gin.Context) {
	// command to generate the rsa key pair
	// openssl genrsa -out cert/id_rsa 4096
	// openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub

	// get cookies
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	//decode and validate
	data, err := ioutil.ReadFile("./id_rsa.pub")
	check(err)
	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	check(err)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// }

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		// return []byte(os.Getenv("SECRET")), nil
		return key, nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)

		}
		// find the user
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)

		}

		// attach to req
		c.Set("user", user)
		c.Set("claim", claims)

		// continue
		c.Next()

	} else {
		fmt.Println("token not valid")
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
