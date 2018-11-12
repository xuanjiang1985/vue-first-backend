package controllers

import (
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/gin-gonic/gin.v1"
)

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

	type msgData struct {
		Content string `form:"content" json:"content" binding:"required"`
	}
	var msgdata msgData
	if c.Bind(&msgdata) != nil {
		c.JSON(200, gin.H{
			"msg":  "内容不能为空",
			"code": 403,
		})
		return
	}
	_, err = db.Exec(`INSERT INTO messages (content) VALUES (?)`, msgdata.Content)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  err.Error(),
			"code": 500,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "提交成功",
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
			"msg":  err.Error(),
			"code": 403,
		})
		return
	}
	seelog.Info("获取成功")
	c.JSON(200, gin.H{
		"code": 200,
		"data": msg,
		"msg":  "操作成功",
	})
}
