package main

import (
	"fmt"
	seelog "github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/gin-gonic/gin.v1"
	"ios-go/conf"
	"time"
)

var sqlconn string = conf.Conn
var logger = conf.Logger
var hmacSampleSecret = []byte("wang ba da 211")

type Messages struct {
	Id         int    `json:"id"`
	Content    string `json:"content"`
	Created_at string `json:"created_at"`
}

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
	r.POST("/msg", PostMsg)
	r.GET("/msg", GetMsg)
	r.GET("/token", GetToken)
	v1 := r.Group("/v1")
	v1.Use(JWTMiddleware())
	{
		v1.GET("/test", JustTest)
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
		tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTgxMDM2MTUsIm1haWwiOiJzb2dvbmd5dUAxNjMuY29tIiwibmFtZSI6Inpob3VnYW5nIn0.LxdLyTSlxZpVGT5ManFrurFw9keeeWZwfv3JyJS0NXk"
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
			c.JSON(200, gin.H{"code": 401, "msg": err.Error()})
			c.AbortWithStatus(200)
		}
	}
}

func PostMsg(c *gin.Context) {
	//开启日志
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()
	//数据库连接
	db, err := sqlx.Connect("mysql", sqlconn)
	if err != nil {
		seelog.Error("can't connect db ", err)
		return
	}
	defer db.Close()

	type Data struct {
		Content string `form:"content" json:"content" binding:"required"`
	}
	var data Data
	if c.Bind(&data) != nil {
		c.JSON(200, gin.H{
			"error":  "内容不能为空",
			"status": 0,
		})
		return
	}
	_, err = db.Exec(`INSERT INTO messages (content) VALUES (?)`, data.Content)
	if err != nil {
		c.JSON(200, gin.H{
			"error":  err.Error(),
			"status": 0,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
	})
}

func GetMsg(c *gin.Context) {
	//开启日志
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()
	//数据库连接
	db, err := sqlx.Connect("mysql", sqlconn)
	if err != nil {
		seelog.Error("can't connect db ", err)
		return
	}
	defer db.Close()
	var msg []Messages
	skip := c.Query("skip")
	err = db.Select(&msg, "SELECT * FROM messages ORDER BY id DESC LIMIT ?,10", skip)
	if err != nil {
		seelog.Error("can't read db ", err)
		c.JSON(200, gin.H{
			"error":  err.Error(),
			"status": 0,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 1,
		"data":   msg,
	})
}

func GetToken(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": "zhougang",
		"mail": "sogongyu@163.com",
		"exp":  time.Now().Add(time.Second * 20).Unix(),
	})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": "Server error!"})
	}
	c.JSON(200, gin.H{"code": 200, "msg": "ok", "jwt": tokenString})
}

func JustTest(c *gin.Context) {
	c.JSON(200, gin.H{"code": 200, "msg": "ok, I am passed via jwt token"})
}
