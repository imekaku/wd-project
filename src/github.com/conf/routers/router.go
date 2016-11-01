package routers

import (
	"github.com/astaxie/beego"
	"github.com/conf/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/r", &controllers.OwnController{}, `post:ChangeRegexp`)
	beego.Router("/d", &controllers.OwnController{}, `delete:DeleteService`)
	beego.Router("/deploy", &controllers.OwnController{}, `post:Deploy`)
}
