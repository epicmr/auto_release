package routers

import (
	"auto_release/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.DelStaticPath("/static")
	beego.SetStaticPath("//", "web/dist")

	// /api namespace
	apiNS := beego.NewNamespace("/api",

		// Handle user requests
		beego.NSRouter("servs", &controllers.ApiController{}, "get:GetServs"),
		beego.NSRouter("hosts", &controllers.ApiController{}, "get:GetHosts"),
		beego.NSRouter("confs", &controllers.ApiController{}, "get:GetConfs"),
		beego.NSRouter("refresh", &controllers.ApiController{}, "get:GetConfsWithMd5"),
		beego.NSRouter("conf", &controllers.ApiController{}, "post:UpdateServsConf"),
		//		beego.NSRouter("/users/:id([0-9]+)", &controllers.ApiController{}, "get:GetUser"),
		//		beego.NSRouter("/users", &controllers.ApiController{}, "post:AddUser"),
		//		beego.NSRouter("/users", &controllers.ApiController{}, "put:UpdateUser"),
		//		beego.NSRouter("/users/:id([0-9]+)", &controllers.ApiController{}, "delete:DeleteUser"),
	)

	beego.AddNamespace(apiNS)

	releaseNS := beego.NewNamespace("/release",

		// Handle user requests
		beego.NSRouter("pack", &controllers.ReleaseController{}, "post:Pack"),
		beego.NSRouter("trans", &controllers.ReleaseController{}, "post:Trans"),
		beego.NSRouter("post", &controllers.ReleaseController{}, "post:Post"),
	)

	beego.AddNamespace(releaseNS)
}
