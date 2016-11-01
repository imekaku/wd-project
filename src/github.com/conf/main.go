package main

import (
	"github.com/astaxie/beego"
	"github.com/conf/conf"
	_ "github.com/conf/routers"
)

func main() {
	beego.SetLogger("file", `{"filename":"./logs/change-fluent-conf.log"}`)
	beego.BeeLogger.DelLogger("console")
	beego.SetLevel(beego.LevelEmergency)
	conf.ParseConfig("./conf/cfg.json")
	beego.Run()
}
