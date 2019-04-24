package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	ms "github.com/epicmr/auto_release/models/mysql"
	_ "github.com/epicmr/auto_release/routers"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	db, _ := ms.InitDb()
	defer db.Close()
	db.AutoMigrate(&ms.Host{}, &ms.Serv{}, &ms.ServEnv{})
	// orm.RegisterDriver("mysql", orm.DRMySQL)
	// orm.RegisterModel(&ms.Host{})
	//orm.RegisterDataBase("default", "mysql", "root:root@/orm_test?charset=utf8")
	// orm.RegisterDataBase("default", "mysql", "auto_release:auto_release@tcp(120.25.154.225:3309)/dev_release?charset=utf8")
}

func main() {
	beego.BConfig.WebConfig.ViewsPath = "dist"
	logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	beego.Run()
}
