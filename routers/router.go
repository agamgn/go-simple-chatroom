package routers
import (
	"../controllers"
	"github.com/astaxie/beego"
	)
/**
@author: agamgn
@date:	2019-04-12
 */

func init()  {
	beego.Router("/",&controllers.HomeContraller{})
	beego.Router("/Room",&controllers.ServersController{})
	beego.Router("/Room/WsRoom",&controllers.ServersController{},"get:WsRoom")


}

