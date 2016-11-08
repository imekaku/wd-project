package routers

import (
	"github.com/astaxie/beego"
	"github.com/fluentd_conf/controllers"
)

func init() {
	//	beego.Router("/", &controllers.MainController{})
	//	beego.Router("/service/get", &controllers.ServiceController{}, `get:GetServiceList`)
	//	beego.Router("/service/add", &controllers.ServiceController{}, `post:ChangeRegexp`)
	//	beego.Router("/service/change", &controllers.ServiceController{}, `post:ChangeRegexp`)
	//	beego.Router("/service/delete", &controllers.ServiceController{}, `delete:DeleteService`)
	//	beego.Router("/deploy", &controllers.ServiceController{}, `post:Deploy`)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/service", &controllers.ServiceController{}, `get:GetServicesList;post:AddServiceRegexp`)
	beego.Router("/service/:service", &controllers.ServiceController{},
		`get:GetServiceRegexp;put:ChangeServiceRegexp;delete:DeleteService`)
}
