package main

import (
	"strings"

	"github.com/astaxie/beego/context"

	ms "auto_release/models/mysql"
	_ "auto_release/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
)

var FilterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session("uid").(string)
	ok2 := strings.Contains(ctx.Request.RequestURI, "/login")
	if !ok && !ok2 {
		ctx.Redirect(302, "/release-stg")
	}
}

func init() {
	db, _ := ms.InitDb()
	db.AutoMigrate(&ms.Env{}, &ms.Host{}, &ms.ServEnv{}, &ms.RouteItem{}, &ms.User{}, &ms.UserConf{})
	db.Model(&ms.ServEnv{}).AddUniqueIndex("idx_serv_env", "serv_id", "env_id")
	db.Model(&ms.RouteItem{}).AddUniqueIndex("idx_parentid_name", "parent_id", "name")
	//a := gormadapter.NewAdapterByDB(db)
	//ctl.E = casbin.NewEnforcer("conf/rbac.conf", a)
}

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	//beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
	//beego.BConfig.WebConfig.ViewsPath = "dist"
	logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
	beego.Run()
	//e.SavePolicy()
}
