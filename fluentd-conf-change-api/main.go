package main

import (
	"github.com/astaxie/beego"
	"github.com/wd-project/fluentd-conf-change-api/conf"
	_ "github.com/wd-project/fluentd-conf-change-api/routers"
)

func main() {
	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.BeeLogger.DelLogger("console")
	beego.SetLevel(beego.LevelEmergency)
	conf.ParseConfig("./conf/cfg.json")
	beego.Run()
}
