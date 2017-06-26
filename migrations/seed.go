package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"ios-go/conf"
	"log"
	"time"
)

var sqlconn string = conf.Conn

func main() {
	//开启日志
	//数据库连接
	db, err := sqlx.Connect("mysql", sqlconn)
	if err != nil {
		log.Println("can't connect db ", err)
		return
	}
	defer db.Close()
	password := []byte("62553380")
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	unix_time := time.Now().Unix()
	result, err := db.Exec(`INSERT INTO users (phone,password,created_at,updated_at) VALUES (?,?,?,?)`, "18922860697", hashedPassword, unix_time, unix_time)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// if err != nil {
	// 	c.JSON(200, gin.H{
	// 		"code": 400,
	// 		"msg":  "手机号已存在",
	// 	})
	// 	return
	// }
	// c.JSON(200, gin.H{
	// 	"code": 200,
	// 	"msg":  "登录成功",
	// })
	log.Println("数据填入成功。", result)
}
