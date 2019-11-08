package controllers

import (
	ms "auto_release/models/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//MainController main controller
type MainController struct {
	beego.Controller
}

func (c *MainController) Prepare() {
	//c.Layout = "layout/app.tpl"
}

// Home goes to index.html file from ./../dist
func (c *MainController) Home() {
	//c.TplName = "main/index.tpl"
	c.TplName = "main/index.html"
}

func (c *MainController) Login() {
	c.TplName = "session/new.tpl"
}

func (c *MainController) Create() {
	phone := c.GetString("phone")
	password := c.GetString("password")

	db, _ := ms.InitDb()
	var user ms.User
	db.Where("phone = ?", phone).Find(&user)

	if user.Password == password {
		logs.Info("验证成功，开始跳转")
		c.SetSession("current_user", phone)
		c.Redirect("/", 302)
		c.StopRun()
	}
	logs.Info("验证失败，回到登录界面")
	c.Redirect("/session/login", 302)
}

func (c *MainController) Logout() {
	c.DelSession("current_user")
	c.DestroySession()
	c.Redirect("/session/login", 302)
}
