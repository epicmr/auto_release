package main

import (

	"github.com/astaxie/beego"
	_ "github.com/epicmr/auto_release/routers"
)

func main() {
	beego.BConfig.WebConfig.ViewsPath = "dist"
    beego.SetLogger("file", `{"filename":"logs/log.log"}`)
	beego.Run()
}
