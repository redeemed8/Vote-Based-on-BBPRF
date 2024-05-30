package config

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

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

var OssClient *oss.Client
var IVBucket *oss.Bucket

const (
	Endpoint        = "http://oss-cn-beijing.aliyuncs.com"
	AccessKeyId     = "xxxx"
	AccessKeySecret = "xxxx"
	IVBucketName    = "xxxx"
)

func init() {
	var err1 error
	OssClient, err1 = oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err1 != nil {
		fmt.Println("ossClient初始化错误: " + err1.Error())
		os.Exit(1)
	}
	var err2 error
	IVBucket, err2 = OssClient.Bucket(IVBucketName)
	if err2 != nil {
		fmt.Println("ossBucket初始化错误: " + err2.Error())
		os.Exit(1)
	}
	fmt.Println("oss初始化完成....")
}
