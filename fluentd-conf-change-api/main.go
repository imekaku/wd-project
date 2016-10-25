package main

import (
	"github.com/astaxie/beego"
	_ "t-bee/change-conf/routers"
)

func main() {
	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.BeeLogger.DelLogger("console")
	beego.SetLevel(beego.LevelError)
	beego.Run()
}
