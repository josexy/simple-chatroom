package main

import (
	"flag"

	"github.com/josexy/gochatroom/global"
	"github.com/josexy/gochatroom/router"
)

var configPath string

func InitParse() {
	flag.StringVar(&configPath, "c", "conf/config.yaml", "config file path")
}

func main() {

	InitParse()
	flag.Parse()
	global.InitConfig(configPath)

	svr := NewServer(global.AppConfig.Server.WebPort, router.NewRouter())
	svr.Run()
}
