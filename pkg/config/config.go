package config

import (
	"fmt"
	"log"
	"time"

	"github.com/YeHeng/go-web-api/pkg/color"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	// 登录验证 Token，Header 中传递的参数
	HeaderLoginToken = "Token"

	// 签名验证 Token，Header 中传递的参数
	HeaderSignToken = "Authorization"

	// 签名验证 Date，Header 中传递的参数
	HeaderSignTokenDate = "Authorization-Date"
)

type Configuration struct {
	AppName string `toml:"appName" default:"go-web-api"`
	Port    string `toml:"port" default:"9092"`
	Stack   bool   `toml:"stack" default:"false"`

	Secure struct {
		Salt string `toml:"salt" default:"CIjLJuMFQOS3Kd3r3T84P/EQy89NMcamvIA5iUh+Fr66PaHrlQvDSMASMj/mFzDl7t7fNGo587NhInEkcVNT9Tm0aZBK0sJHKjkSuyF6eHiIDQY5KewwpetfRnpigG5z57DC++Rflkuf+XkL8OZgvshfCVOKrAyBDFhlhl9T3Dk="`
	}

	Feature struct {
		DisablePProf   bool `toml:"disablePProf" disablePProf:"port" default:"false"`
		DisableSwagger bool `toml:"disableSwagger" default:"true"`
		PanicNotify    bool `toml:"panicNotify" default:"true"`
		RecordMetrics  bool `toml:"recordMetrics" default:"true"`
		EnableCors     bool `toml:"enableCors" default:"true"`
		EnableRate     bool `toml:"enableRate" default:"true"`
	} `toml:"feature"`

	Cors struct {
		AllowedOrigins   string `toml:"allowedOrigins" default:"*"`
		AllowedMethods   string `toml:"allowedMethods" default:"GET,POST,HEAD,PUT,PATCH,DELETE"`
		AllowedHeaders   string `toml:"allowedHeaders" default:"*"`
		AllowCredentials bool   `toml:"allowCredentials" default:"true"`
	} `toml:"cors"`

	Redis struct {
		Addr         string `toml:"addr" default:"127.0.0.1:3306"`
		Pass         string `toml:"pass" default:""`
		Db           int    `toml:"db" default:"0"`
		MaxRetries   int    `toml:"maxRetries" default:"3""`
		PoolSize     int    `toml:"poolSize" default:"10"`
		MinIdleConns int    `toml:"minIdleConns" default:"5"`
	} `toml:"redis"`

	Logger struct {
		Folder   string `toml:"folder" default:"./logs/"`
		Filename string `toml:"filename" default:"app.logger"`
		Level    string `toml:"level"  default:"info"`

		// MaxSize is the maximum size in megabytes of the logger file before it gets
		// rotated. It defaults to 100 megabytes.
		MaxSize int `toml:"maxsize"`

		// MaxAge is the maximum number of days to retain old logger files based on the
		// timestamp encoded in their filename.  Note that a day is defined as 24
		// hours and may not exactly correspond to calendar days due to daylight
		// savings, leap seconds, etc. The default is not to remove old logger files
		// based on age.
		MaxAge int `toml:"maxage"`

		// MaxBackups is the maximum number of old logger files to retain.  The default
		// is to retain all old logger files (though MaxAge may still cause them to get
		// deleted.)
		MaxBackups int `toml:"maxbackups"`

		// LocalTime determines if the time used for formatting the timestamps in
		// backup files is the computer's local time.  The default is to use UTC
		// time.
		LocalTime bool `toml:"localtime"`

		// Compress determines if the rotated logger files should be compressed
		// using gzip. The default is not to perform compression.
		Compress bool `toml:"compress"`
	} `toml:"logger"`

	Database struct {
		DbType          string        `toml:"dbType"`
		Addr            string        `toml:"addr"`
		Username        string        `toml:"username"`
		Password        string        `toml:"password"`
		DbName          string        `toml:"dbName"`
		SkipTransaction bool          `toml:"skipTransaction" default:"false"`
		MaxOpenConn     int           `toml:"maxOpenConn"`
		MaxIdleConn     int           `toml:"maxIdleConn"`
		ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
	} `toml:"database"`
}

var config = new(Configuration)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath("/etc/go-web-api")
	viper.AddConfigPath("$HOME/.go-web-api")
	viper.SetConfigType("toml")

	fmt.Println(color.Green("* [loading config.toml]"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(config); err != nil {
			panic(err)
		}
	})

}

func Get() Configuration {
	return *config
}
