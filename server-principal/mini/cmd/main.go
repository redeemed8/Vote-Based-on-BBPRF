package main

import (
	"github.com/spf13/viper"
	"mini/app"
	_ "mini/config"
)

var Port = ":" + viper.GetString("server.port")

func main() {
	r := app.AppsRouter()
	_ = r.Run(Port)
}
