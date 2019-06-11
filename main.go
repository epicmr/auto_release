package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "github.com/casbin/casbin"
    "github.com/casbin/gorm-adapter"
	ms "github.com/epicmr/auto_release/models/mysql"
	ctl "github.com/epicmr/auto_release/controllers"
	_ "github.com/epicmr/auto_release/routers"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db, _ := ms.InitDb()
	db.AutoMigrate(&ms.Env{}, &ms.Host{}, &ms.Serv{}, &ms.ServEnv{}, &ms.RouteItem{})
	db.Model(&ms.ServEnv{}).AddUniqueIndex("idx_serv_env", "serv_id", "env_id")
	db.Model(&ms.RouteItem{}).AddUniqueIndex("idx_parentid_name", "parent_id", "name")

    a := gormadapter.NewAdapterByDB(db)
    ctl.E = casbin.NewEnforcer("conf/rbac.conf", a)
}

func main() {
	beego.BConfig.WebConfig.ViewsPath = "dist"
	logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
	beego.Run()
    //e.SavePolicy()
}
