package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
	"ios-go/conf"
	"ios-go/controllers"
	"strings"
)

var hmacSampleSecret = []byte(conf.JwtKey)

func main() {
	r := gin.Default()
	r.Use(CORS())
	r.Static("/public", "./public")
	r.POST("/msg", controllers.PostMsg)
	r.GET("/msg", controllers.GetMsg)
	r.POST("/login", controllers.PostLogin)
	r.POST("/register", controllers.PostRegister)
	//jwt auth group
	v1 := r.Group("/v1")
	v1.Use(JWTMiddleware())
	{
		v1.GET("/change-name", controllers.GetChangeName)
	}
	r.Run(":8080")
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//if has token
		authString := c.Request.Header.Get("Authorization")
		if len(authString) == 0 {
			c.JSON(401, gin.H{"code": 401, "msg": "invalid token"})
			c.AbortWithStatus(401)
			return
		}
		//confirm token
		tokenString := strings.Replace(authString, "Bearer ", "", -1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecret, nil
		})
		if token.Valid {
			c.Next()
		} else {
			c.JSON(401, gin.H{"code": 401, "msg": err.Error()})
			c.AbortWithStatus(401)
		}
	}
}
