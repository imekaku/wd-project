package main

import (
	"github.com/astaxie/beego"
	"github.com/fluentd_conf/conf"
	_ "github.com/fluentd_conf/routers"
)

func main() {
	beego.SetLogger("file", `{"filename":"./logs/change-fluent-conf.log"}`)
	beego.BeeLogger.DelLogger("console")
	beego.SetLevel(beego.LevelEmergency)
	beego.SetLogFuncCall(true)
	conf.ParseConfig("./conf/cfg.json")
	beego.Run()
}
