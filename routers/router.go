package routers

import (
	"github.com/epicmr/auto_release/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.DelStaticPath("/static")
	beego.SetStaticPath("//", "web/dist")

	// /api namespace
	apiNS := beego.NewNamespace("/api",

		// Handle user requests
		beego.NSRouter("servs", &controllers.APIController{}, "get:GetServs"),
		beego.NSRouter("envs", &controllers.APIController{}, "get:GetEnvs"),
		beego.NSRouter("env", &controllers.APIController{}, "post:UpdateEnv"),
		beego.NSRouter("itemstree", &controllers.APIController{}, "get:GetItemsTree"),
		beego.NSRouter("items", &controllers.APIController{}, "get:GetAllItems"),
		beego.NSRouter("item", &controllers.APIController{}, "post:UpdateItem"),
		beego.NSRouter("hosts", &controllers.APIController{}, "get:GetHosts"),
		beego.NSRouter("host", &controllers.APIController{}, "post:UpdateHost"),
		beego.NSRouter("confs", &controllers.APIController{}, "get:GetConfs"),
		beego.NSRouter("refresh", &controllers.APIController{}, "get:GetConfsWithMd5"),
		beego.NSRouter("conf", &controllers.APIController{}, "post:UpdateServsConf"),
        beego.NSRouter("checkmd5", &controllers.APIController{}, "post:CheckMD5"),
        beego.NSRouter("checktime", &controllers.APIController{}, "post:CheckTime"),
        beego.NSRouter("usergroup", &controllers.APIController{}, "post:UpdateUserGroup"),
		//		beego.NSRouter("/users/:id([0-9]+)", &controllers.APIController{}, "get:GetUser"),
		//		beego.NSRouter("/users", &controllers.APIController{}, "post:AddUser"),
		//		beego.NSRouter("/users", &controllers.APIController{}, "put:UpdateUser"),
		//		beego.NSRouter("/users/:id([0-9]+)", &controllers.APIController{}, "delete:DeleteUser"),
	)

	beego.AddNamespace(apiNS)

    releaseNS := beego.NewNamespace("/release", // Handle user requests
    beego.NSRouter("pack", &controllers.ReleaseController{}, "post:Pack"),
    beego.NSRouter("trans", &controllers.ReleaseController{}, "post:Trans"),
    beego.NSRouter("post", &controllers.ReleaseController{}, "post:Post"),
)

	beego.AddNamespace(releaseNS)

    releaseThirdNS := beego.NewNamespace("/third", // Handle user requests
    beego.NSRouter("deleteUser", &controllers.ThirdController{}, "get:DeleteUser"),
    )

	beego.AddNamespace(releaseThirdNS)
}
