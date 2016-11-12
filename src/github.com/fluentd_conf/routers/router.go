package routers

import (
	"github.com/astaxie/beego"
	"github.com/fluentd_conf/controllers"
)

func init() {
	beego.Router("/services", &controllers.ServiceController{}, `get:GetServicesList;post:AddServiceRegexp`)
	beego.Router("/services/:service", &controllers.ServiceController{},
		`get:GetServiceRegexp;put:ChangeServiceRegexp;delete:DeleteService`)
}
