package controllers

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//ThirdController for backstage api
type ThirdController struct {
	beego.Controller
	JSONRetMsg
}

//GetEnvs return envs
func (c *ThirdController) DeleteUser() {
    phone, _ := c.GetInt("phone")
	var stderr, stdout bytes.Buffer
	s := fmt.Sprintf("cd /home/tools/del_user_wechat/;./repair.py %d", phone)
	logs.Info(s)
	cmd := exec.Command("ssh", "STG-ALL-83", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

    err := cmd.Run()
	if nil != err {
        c.setError(1, fmt.Sprintf("STG-ALL-83 exec[%s] failed. ", s))
        logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
        goto end
	}

	logs.Info("DeleteUser %d passed", phone)

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}
