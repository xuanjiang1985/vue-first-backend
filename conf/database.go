package conf

import (
	"github.com/dlintw/goconf"
	"log"
	//"os"
)

var Conn string
var JwtKey string

//var basePath = os.Getenv("GOPATH")

//初始化数据库配置
func init() {
	conf, err := goconf.ReadConfigFile(basePath + "/src/vue-first-backend/.env")
	if err != nil {
		log.Println(err)
		return
	}
	user, _ := conf.GetString("mysql", "user")
	password, _ := conf.GetString("mysql", "password")
	host, _ := conf.GetString("mysql", "host")
	port, _ := conf.GetString("mysql", "port")
	db, _ := conf.GetString("mysql", "db")
	jwtkey, _ := conf.GetString("jwt", "jwtkey")
	Conn = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	JwtKey = jwtkey
}
