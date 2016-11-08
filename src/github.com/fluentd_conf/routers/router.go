package routers

import (
	"github.com/astaxie/beego"
	"github.com/fluentd_conf/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/service", &controllers.ServiceController{}, `get:GetServicesList;post:AddServiceRegexp`)
	beego.Router("/service/:service", &controllers.ServiceController{},
		`get:GetServiceRegexp;put:ChangeServiceRegexp;delete:DeleteService`)
}
