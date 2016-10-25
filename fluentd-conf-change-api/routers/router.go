package routers

import (
	"github.com/astaxie/beego"
	"t-bee/change-conf/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/r", &controllers.OwnController{}, `post:ChangeRegexp`)
}
