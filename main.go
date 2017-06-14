package main

import (
	"gopkg.in/gin-gonic/gin.v1"
)

type BookStruct struct {
	Num   int     `json:"num"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Profile struct {
	Wife string     `json:"wife"`
	Age  int        `json:"age"`
	Sex  int        `json:"sex"`
	Book BookStruct `json:"book"`
}

func main() {
	r := gin.Default()
	r.Use(CORS())
	r.Static("/public", "./public")
	r.GET("/test", func(c *gin.Context) {
		profile := Profile{"wenjuan", 28, 0, BookStruct{10, "Plan B", 34.21}}
		c.JSON(200, profile)
	})
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
