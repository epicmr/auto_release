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
	db.AutoMigrate(&ms.Env{}, &ms.Host{}, &ms.Serv{}, &ms.ServEnv{})
	db.Model(&ms.ServEnv{}).AddUniqueIndex("idx_serv_env", "serv_id", "env_id")

    //env := ms.Env{
    //    Name:     "stg",
    //    ServType: 14}

    //db.Create(&env)

    //serv := ms.Serv{
    //    ServName:  "deal",
    //    ServType:  1,
    //    LocalPath: "/data/upgrade/cgi"}
    //db.Create(&serv)

    //if env.ID > 0 {
    //logs.Info(env)
    //host := ms.Host{
    //    EnvID:    env.ID,
    //    Name:     "STG-ALL-83",
    //    ServType: 14}

    //db.Create(&host)
    //}

    //if serv.ID > 0 && env.ID > 0 {
    //servEnv := ms.ServEnv{
    //    ServID:     serv.ID,
    //    ServName:   serv.ServName,
    //    EnvID:      env.ID,
    //    Env:        env.Name,
    //    RemotePath: "/var/www/cgi-bin/deal/"}

    //db.Create(&servEnv)
    //}

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
