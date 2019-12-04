package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	models "auto_release/models"
	ms "auto_release/models/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//ReleaseController struct
type ReleaseController struct {
	beego.Controller
	JSONRetMsg
}

// Pack returns a list of users
func (c *ReleaseController) Pack() {
	logs.Info("Pack() begin")
	var ob ms.ServFlt
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info("解析json参数: serv.ServName[%v], serv.Env[%v]", ob.ServName, ob.Env)
	db, _ := ms.InitDb()

	//更新spec版本号
	var serv ms.Serv
	db.Debug().Preload("ServEnvs").Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()
	if ob.Env == "local" {
		c.ModifyLocalPath(&serv)
		logs.Info("local环境，修改localPath为用户当前配置的路径。[%v]", serv.LocalPath)
	} else {
		if serv.ServType == 1 {
			serv.LocalPath = CgiPath
		} else if serv.ServType == 2 {
			serv.LocalPath = AppPath
		} else if serv.ServType == 3 {
			serv.LocalPath = GoPath
		}
		logs.Info("其他环境，localPath.[%v], servType.[%v]", serv.LocalPath, serv.ServType)
	}

	var _servenv ms.ServEnv
	for _, servenv := range serv.ServEnvs {
		if servenv.ServName == ob.ServName &&
			strings.ToLower(servenv.Env) == strings.ToLower(ob.Env) {
			_servenv = servenv
		}
	}

	logs.Info("%#v", _servenv)
	models.GenSpec(&serv, _servenv)

	//打包
	var stderr, stdout bytes.Buffer
	s := "cd /root/rpmbuild;rpmbuild -bb SPECS/" + ob.ServName + ".spec"
	logs.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if nil != err {
		c.setError(1, fmt.Sprintf("local exec[%s] failed. ", s))
		logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
		goto end
	}

	logs.Info(ob.ServName + " package passed")

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

// //Trans to server
func (c *ReleaseController) Trans() {
	var ob ms.ServFlt
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)

	db, _ := ms.InitDb()

	var env ms.Env
	var serv ms.Serv
	db.Preload("Hosts").Where("name = ?", ob.Env).First(&env).GetErrors()
	db.Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()

	var installRpm string
	var stderr, stdout bytes.Buffer
	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
	logs.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		c.setError(1, fmt.Sprintf("local exec[%s] failed. ", s))
		logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
		goto end
	}

	//包名
	installRpm = "/root/rpmbuild/RPMS/x86_64/" + strings.Trim(stdout.String(), "\r\n")

	for _, host := range env.Hosts {
		servType1 := serv.ServType
		servType2 := host.ServType
		if (1<<uint8(servType1))&servType2 > 0 {
			logs.Info("scp [%s] [%s]", installRpm, host.Name)
			cmd = exec.Command("scp", installRpm, host.Name+":/data/upgrade/")

			err = cmd.Run()
			if err != nil {
				c.setError(1, fmt.Sprintf("scp [%s] [%s] failed. ", installRpm, host.Name))
				logs.Error("scp [%s] [%s] failed. Error:[%s]", installRpm, host.Name, stderr.String())
				goto end
			}
		}
	}

	logs.Info(ob.ServName + " translate passed")

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

// //Post replace exe and restart
func (c *ReleaseController) Post() {
	var ob ms.ServFlt
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)

	db, _ := ms.InitDb()

	var env ms.Env
	var serv ms.Serv
	db.Preload("Hosts").Where("name = ?", ob.Env).First(&env).GetErrors()
	db.Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()

	var installRpm string
	var stderr, stdout bytes.Buffer
	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
	logs.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		c.setError(1, fmt.Sprintf("local exec[%s] failed. ", s))
		logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
		goto end
	}

	//包名
	installRpm = "/data/upgrade/" + strings.Trim(stdout.String(), "\r\n")
	logs.Info("Install RPM :", installRpm)

	for _, host := range env.Hosts {
		servType1 := serv.ServType
		servType2 := host.ServType
		if (1<<uint8(servType1))&servType2 > 0 {
			s = fmt.Sprintf("rpm -U --force %s", installRpm)
			logs.Info(host.Name, s)
			cmd = exec.Command("ssh", host.Name, s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				c.setError(1, fmt.Sprintf("install [%s] [%s] failed. ", installRpm, host.Name))
				logs.Error("install [%s] [%s] failed. Error:[%s]", installRpm, host.Name, stderr.String())
				goto end
			}
		}
	}

	logs.Info(ob.ServName + " install passed")

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func (c *ReleaseController) ModifyLocalPath(serv *ms.Serv) {
	var user ms.User
	var userConf ms.UserConf
	db, _ := ms.InitDb()
	phone := c.GetSession("current_user")
	logs.Info("==============current_user_phone: [%v]", phone)
	db.Where("phone = ?", phone).First(&user)
	logs.Info("userID: [%v]", user.UserID)
	db.Where("user_id = ? AND serv_id = ?", user.UserID, serv.ID).First(&userConf)
	serv.LocalPath = userConf.LocalPath
}
