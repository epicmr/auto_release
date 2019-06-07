package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "github.com/casbin/casbin"
    "github.com/casbin/gorm-adapter"
	ms "github.com/epicmr/auto_release/models/mysql"
	_ "github.com/epicmr/auto_release/routers"
	_ "github.com/go-sql-driver/mysql"
)

var (
    e *casbin.Enforcer
)

func init() {
	db, _ := ms.InitDb()
	db.AutoMigrate(&ms.Env{}, &ms.Host{}, &ms.Serv{}, &ms.ServEnv{})
	db.Model(&ms.ServEnv{}).AddUniqueIndex("idx_serv_env", "serv_id", "env_id")

    a := gormadapter.NewAdapterByDB(db)
    e = casbin.NewEnforcer("conf/rbac.conf", a)
    e.AddPolicy("alice", "data", "read")
    e.AddNamedGroupingPolicy("g", "alice", "group")
    e.AddNamedGroupingPolicy("g2", "data1", "data")
}

func main() {
	beego.BConfig.WebConfig.ViewsPath = "dist"
	logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
	beego.Run()
    //e.SavePolicy()
}
