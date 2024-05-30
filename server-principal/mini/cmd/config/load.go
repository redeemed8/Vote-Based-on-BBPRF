package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var RDB *redis.Client

func init() {
	viper.SetConfigName("application")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s \n", err))
	}
	fmt.Println("config loaded successfully ...")
}

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //	慢 SQL阈值
			LogLevel:      logger.Info, //	级别
			Colorful:      true,        //	彩色
		})
	dsn := viper.GetString("mysql.dns") +
		"/" + viper.GetString("mysql.basename") +
		"?" + viper.GetString("mysql.others")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1) //	退出
	}
	fmt.Println("mysql loaded successfully ...")
}

