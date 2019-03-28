package main

import (

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/epicmr/auto_release/routers"
)

func main() {
	beego.BConfig.WebConfig.ViewsPath = "dist"
    logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	beego.Run()
}
