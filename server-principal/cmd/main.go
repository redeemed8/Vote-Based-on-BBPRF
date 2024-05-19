package main

import (
	"github.com/spf13/viper"
	"mini_pkcs/app"
	_ "mini_pkcs/config"
)

var Port = ":" + viper.GetString("server.port")

func main() {
	r := app.AppsRouter()
	_ = r.Run(Port)
}
