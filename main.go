package main

import (
	"strings"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	ms "auto_release/models/mysql"
	_ "auto_release/routers"
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
	db.AutoMigrate(&ms.Env{}, &ms.Host{}, &ms.Serv{}, &ms.ServEnv{}, &ms.RouteItem{}, &ms.User{}, &ms.UserConf{})
	db.Model(&ms.ServEnv{}).AddUniqueIndex("idx_serv_env", "serv_id", "env_id")
	db.Model(&ms.RouteItem{}).AddUniqueIndex("idx_parentid_name", "parent_id", "name")
	//a := gormadapter.NewAdapterByDB(db)
	//ctl.E = casbin.NewEnforcer("conf/rbac.conf", a)
}

func main() {
	//beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
	//beego.BConfig.WebConfig.ViewsPath = "dist"
	logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
	beego.Run()
	//e.SavePolicy()
}
