package conf

import (
	"flag"
	"gopkg.in/ini.v1"
	"log"
	"net"
	"strconv"
	"time"
)

type Config struct {
	Gin   Gin
	Mysql Mysql
	Redis Redis
	Email Email
}

type Gin struct {
	Addr    string
	Port    int64
	Release bool
	Https   bool
	Host    string
}

type Mysql struct {
	Addr     string
	Port     int64
	UserName string
	Password string
	DB       string
}

type Redis struct {
	Addr     string
	Port     int64
	UserName string
	Password string
	Prefix   string
	DB       int
	Expire   time.Duration
}

type Email struct {
	User     string
	Password string
	Host     string
}

var (
	Cfg Config
)

func init() {
	filePath := flag.String("conf", "", "your config file path, required!!!")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("You must run it like 'go-sso --conf your-config-path'")
	}

	iniFile, err := ini.Load(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	parseINI(iniFile)

}

func parseINI(f *ini.File) {
	// gin
	g := f.Section("gin")

	s := g.Key("addr").String()
	ip := net.ParseIP(s)
	if ip == nil {
		log.Fatal("gin.addr parse error, check your config file")
	}
	Cfg.Gin.Addr = ip.To4().String()

	i, err := g.Key("port").Int64()
	if err != nil || i < 1 || i > 65535 {
		log.Fatal("gin.port parse error, check your config file")
	}
	Cfg.Gin.Port = i

	b, err := g.Key("release").Bool()
	if err != nil {
		log.Fatal("gin.release parse error, check your config file")
	}
	Cfg.Gin.Release = b

	b, err = g.Key("https").Bool()
	if err != nil {
		log.Fatal("gin.https parse error, check your config file")
	}
	Cfg.Gin.Https = b

	s = g.Key("host").String()
	if s == "" {
		log.Fatal("gin.host parse error, check your config file")
	}
	Cfg.Gin.Host = s

	// mysql
	m := f.Section("mysql")

	s = m.Key("addr").String()
	ip = net.ParseIP(s)
	if ip == nil {
		log.Fatal("mysql.addr parse error, check your config file")
	}
	Cfg.Mysql.Addr = ip.To4().String()

	i, err = m.Key("port").Int64()
	if err != nil || i < 1 || i > 65535 {
		log.Fatal("mysql.port parse error, check your config file")
	}
	Cfg.Mysql.Port = i

	s = m.Key("user").String()
	if s == "" {
		log.Fatal("mysql.user parse error, check your config file")
	}
	Cfg.Mysql.UserName = s

	Cfg.Mysql.Password = m.Key("password").String()

	s = m.Key("db").String()
	if s == "" {
		log.Fatal("mysql.db parse error, check your config file")
	}
	Cfg.Mysql.DB = s

	// redis
	r := f.Section("redis")

	s = r.Key("addr").String()
	ip = net.ParseIP(s)
	if ip == nil {
		log.Fatal("redis.addr parse error, check your config file")
	}
	Cfg.Redis.Addr = ip.To4().String()

	i, err = r.Key("port").Int64()
	if err != nil || i < 1 || i > 65535 {
		log.Fatal("redis.port parse error, check your config file")
	}
	Cfg.Redis.Port = i

	Cfg.Redis.Prefix = r.Key("prefix").String()
	Cfg.Redis.DB, _ = r.Key("db").Int()
	Cfg.Redis.Password = r.Key("password").String()
	Cfg.Redis.UserName = r.Key("user").String()
	Cfg.Redis.Expire = r.Key("expire").MustDuration()

	e := f.Section("email")
	user := e.Key("user").String()
	if user == "" {
		log.Fatal("input email user")
	}

	pd := e.Key("password").String()
	if pd == "" {
		log.Fatal("input email password")
	}

	host := e.Key("host").String()
	if host == "" {
		log.Fatal("input email host")
	}

	Cfg.Email.User = user
	Cfg.Email.Password = pd
	Cfg.Email.Host = host
}

func GetListenAddr() string {
	return Cfg.Gin.Addr + ":" + strconv.FormatInt(Cfg.Gin.Port, 10)
}
