package controllers

import (
	valid "github.com/asaskevich/govalidator"
	seelog "github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"ios-go/conf"
	"time"
)

var sqlconn = conf.Conn
var logger = conf.Logger
var hmacSampleSecret = []byte(conf.JwtKey)

func PostLogin(c *gin.Context) {
	//validate
	type Validator struct {
		Phone    string `valid:"required~手机不能为空,int~手机必须是数字,stringlength(11|11)~手机必须为11位"`
		Password string `valid:"required~密码不能为空,stringlength(6|60)~密码至少6位"`
	}
	data := &Validator{
		Phone:    c.PostForm("phone"),
		Password: c.PostForm("password"),
	}

	_, err := valid.ValidateStruct(data)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  err.Error(),
		})
		return
	}
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

	//get login user
	var authUser Users
	password := []byte(c.PostForm("password"))
	err = db.Get(&authUser, "SELECT * FROM users WHERE phone=?", c.PostForm("phone"))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  "账户或密码错误0",
		})
		return
	}
	errors := bcrypt.CompareHashAndPassword([]byte(authUser.Password), password)
	if errors != nil {
		//log.Println(user.Password)
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  "账户或密码错误1",
		})
		return
	}
	token := getToken(authUser)
	c.JSON(200, gin.H{
		"code":       200,
		"msg":        "登录成功",
		"token":      token,
		"userName":   authUser.Name,
		"userHeader": authUser.Header,
		"userId":     authUser.Id,
	})
}

func PostRegister(c *gin.Context) {
	//validate
	type Validator struct {
		Phone     string `valid:"required~手机不能为空,int~手机必须是数字,stringlength(11|11)~手机必须为11位"`
		Password  string `valid:"required~密码不能为空,stringlength(6|60)~密码至少6位"`
		Password2 string `valid:"required~密码确认不能为空,stringlength(6|60)~密码确认至少6位"`
	}
	data := &Validator{
		Phone:     c.PostForm("phone"),
		Password:  c.PostForm("password"),
		Password2: c.PostForm("password2"),
	}

	if c.PostForm("password") != c.PostForm("password2") {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  "两次输入的密码不相同;",
		})
		return
	}

	_, err := valid.ValidateStruct(data)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  err.Error(),
		})
		return
	}
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

	//register user
	password := []byte(c.PostForm("password"))
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	unix_time := time.Now().Unix()
	result, err := db.Exec(`INSERT INTO users (phone,password,created_at,updated_at) VALUES (?,?,?,?)`, c.PostForm("phone"), hashedPassword, unix_time, unix_time)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  "手机号已存在",
		})
		return
	}
	userId, _ := result.LastInsertId()
	token := getToken2(int(userId))
	c.JSON(200, gin.H{
		"code":       200,
		"msg":        "注册成功",
		"token":      token,
		"userName":   "匿名用户",
		"userId":     int(userId),
		"userHeader": "/public/images/header.jpg",
	})
}

func GetChangeName(c *gin.Context) {
	//validate
	type Validator struct {
		Name string `valid:"required~昵称不能为空,stringlength(1|11)~昵称不能超过11个字符"`
	}
	data := &Validator{
		Name: c.Query("name"),
	}

	_, err := valid.ValidateStruct(data)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  err.Error(),
		})
		return
	}
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
	unix_time := time.Now().Unix()
	//change user name
	_, err = db.Exec(`UPDATE users SET name=?, updated_at=? WHERE id=?`, c.Query("name"), unix_time, c.Query("userId"))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 403,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":     200,
		"msg":      "昵称修改成功",
		"userName": c.Query("name"),
	})
}

func getToken(user Users) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Second * 3600).Unix(),
	})
	tokenString, _ := token.SignedString(hmacSampleSecret)
	return tokenString
}

func getToken2(a int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  a,
		"exp": time.Now().Add(time.Second * 3600).Unix(),
	})
	tokenString, _ := token.SignedString(hmacSampleSecret)
	return tokenString
}
