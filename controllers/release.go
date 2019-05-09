package controllers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strings"
    "os/exec"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	ms "github.com/epicmr/auto_release/models/mysql"
	models "github.com/epicmr/auto_release/models"
)

//ReleaseController struct
type ReleaseController struct {
	beego.Controller
	JSONRetMsg
}

// Pack returns a list of users
func (c *ReleaseController) Pack() {
	var ob ms.ServFlt
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)

	db, _ := ms.InitDb()

	//更新spec版本号
	var serv ms.Serv
	db.Debug().Preload("ServEnvs").Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()

	var _servenv ms.ServEnv
	for _, servenv := range serv.ServEnvs {
		if servenv.ServName == ob.ServName &&
			servenv.Env == ob.Env {
			_servenv = servenv
		}
	}

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
	db.Debug().Preload("Hosts").First(&env).GetErrors()
	db.Debug().Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()

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
	db.Debug().Preload("Hosts").First(&env).GetErrors()
	db.Debug().Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()

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
