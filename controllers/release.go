package controllers

import (
	"auto_release/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/astaxie/beego"
)

type ReleaseController struct {
	beego.Controller
	JsonRetMessage
}

// GetUsers returns a list of users
func (this *ReleaseController) Pack() {
	var ob models.ServFlt
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	beego.Info(ob)

	//更新spec版本号
	db := models.InitDb()
	ctx := context.Background()

	serv_list, err := models.QueryServByName(ctx, nil, db, ob.ServName)
	if err != nil {
		beego.Info(err)
	}

	err, servenv_list := models.BatchQueryServEnv(ctx, nil, db)
	if err != nil {
		beego.Info(err)
	}

	models.CloseDb(db)

	var _servenv models.ServEnv
	for _, servenv := range servenv_list {
		if servenv.ServName == serv_list[0].ServName &&
			servenv.Env == ob.Env {
			_servenv = servenv
		}
	}

	models.GenSpec(&serv_list[0], _servenv)

	//打包
	s := "cd /root/rpmbuild;rpmbuild -bb SPECS/" + ob.ServName + ".spec"
	beego.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)

	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if nil != err {
		beego.Error(err)
		this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "pack failed. ")
		beego.Error(stdout.String())
		beego.Error(stderr.String())
		goto end
	}

	beego.Info(ob.ServName + " package passed")

end:
	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}

func (this *ReleaseController) Trans() {
	var ob models.ServFlt
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	beego.Info(ob)

	db := models.InitDb()
	ctx := context.Background()

	serv_list, err := models.QueryServByName(ctx, nil, db, ob.ServName)
	if err != nil {
		beego.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		beego.Info(err)
	}

	models.CloseDb(db)

	var install_rpm string
	var stderr, stdout bytes.Buffer

	//最新版本
	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
	beego.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "trans list failed. ")
		beego.Error(stdout.String())
		beego.Error(stderr.String())
		goto end
	}

	//包名
	install_rpm = strings.Trim(stdout.String(), "\r\n")
	beego.Info("Install RPM :", install_rpm)

	for _, host := range host_list {
		serv_type1, _ := strconv.Atoi(serv_list[0].ServType)
		serv_type2, _ := strconv.Atoi(host.ServType)
		if host.Env == ob.Env && (1<<uint8(serv_type1))&serv_type2 > 0 {
			remote := host.HostName
			beego.Info("HostName", remote)

			//传输
			s = fmt.Sprintf("cd /root/rpmbuild/RPMS/x86_64;scp %s %s:/data/upgrade", install_rpm, remote)
			beego.Info(s)
			cmd = exec.Command("/bin/sh", "-c", s)

			err = cmd.Run()
			if err != nil {
				this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "trans failed. ")
				beego.Error(stdout.String())
				beego.Error(stderr.String())
				goto end
			}
		}
	}

	beego.Info(ob.ServName + " translate passed")

end:
	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}

func (this *ReleaseController) Post() {
	var ob models.ServFlt
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	beego.Info(ob)

	db := models.InitDb()
	ctx := context.Background()

	serv_list, err := models.QueryServByName(ctx, nil, db, ob.ServName)
	if err != nil {
		beego.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		beego.Info(err)
	}

	models.CloseDb(db)

	var install_rpm string
	var stderr, stdout bytes.Buffer

	//最新版本
	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
	beego.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install list failed. ")
		beego.Error(stdout.String())
		beego.Error(stderr.String())
		goto end
	}

	//包名
	install_rpm = strings.Trim(stdout.String(), "\r\n")
	beego.Info("Install RPM :", install_rpm)

	for _, host := range host_list {
		serv_type1, _ := strconv.Atoi(serv_list[0].ServType)
		serv_type2, _ := strconv.Atoi(host.ServType)
		if host.Env == ob.Env && (1<<uint8(serv_type1))&serv_type2 > 0 {
			remote := host.HostName
			beego.Info("HostName", remote)

			s = fmt.Sprintf("ssh %s 'rpm -U --force /data/upgrade/%s'", remote, install_rpm)
			beego.Info(s)
			cmd = exec.Command("/bin/sh", "-c", s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				beego.Error(err)
				this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
				beego.Error(stdout.String())
				beego.Error(stderr.String())
				goto end
			}
		}
	}

	beego.Info(ob.ServName + " install passed")

end:
	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}
