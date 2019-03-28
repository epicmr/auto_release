package main

import (
	_ "auto_release/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.WebConfig.ViewsPath = "dist"
    beego.SetLogger("file", `{"filename":"logs/log.log"}`)
	beego.Run()
}
