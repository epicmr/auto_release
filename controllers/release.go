package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/epicmr/auto_release/models"
)

type ReleaseController struct {
	beego.Controller
	JsonRetMessage
}

// GetUsers returns a list of users
func (this *ReleaseController) Pack() {
	var ob models.ServFlt
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)

	//更新spec版本号
	db := models.InitDb()
	ctx := context.Background()

	serv_list, err := models.QueryServByName(ctx, nil, db, ob.ServName)
	if err != nil {
		logs.Info(err)
	}

	err, servenv_list := models.BatchQueryServEnv(ctx, nil, db)
	if err != nil {
		logs.Info(err)
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
	logs.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)

	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if nil != err {
		logs.Error(err)
		this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "pack failed. ")
		logs.Error(stdout.String())
		logs.Error(stderr.String())
		goto end
	}

	logs.Info(ob.ServName + " package passed")

end:
	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}

func (this *ReleaseController) Trans() {
	var ob models.ServFlt
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)

	db := models.InitDb()
	ctx := context.Background()

	serv_list, err := models.QueryServByName(ctx, nil, db, ob.ServName)
	if err != nil {
		logs.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	models.CloseDb(db)

	var install_rpm string
	var stderr, stdout bytes.Buffer

	//最新版本
	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
	logs.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "trans list failed. ")
		logs.Error(stdout.String())
		logs.Error(stderr.String())
		goto end
	}

	//包名
	install_rpm = strings.Trim(stdout.String(), "\r\n")
	logs.Info("Install RPM :", install_rpm)

	for _, host := range host_list {
		serv_type1, _ := strconv.Atoi(serv_list[0].ServType)
		serv_type2, _ := strconv.Atoi(host.ServType)
		if host.Env == ob.Env && (1<<uint8(serv_type1))&serv_type2 > 0 {
			remote := host.HostName
			logs.Info("HostName", remote)

			//传输
			s = fmt.Sprintf("cd /root/rpmbuild/RPMS/x86_64;scp %s %s:/data/upgrade", install_rpm, remote)
			logs.Info(s)
			cmd = exec.Command("/bin/sh", "-c", s)

			err = cmd.Run()
			if err != nil {
				this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "trans failed. ")
				logs.Error(stdout.String())
				logs.Error(stderr.String())
				goto end
			}
		}
	}

	logs.Info(ob.ServName + " translate passed")

end:
	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}

func (this *ReleaseController) Post() {
	var ob models.ServFlt
	json.Unmarshal(this.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)

	db := models.InitDb()
	ctx := context.Background()

	serv_list, err := models.QueryServByName(ctx, nil, db, ob.ServName)
	if err != nil {
		logs.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	models.CloseDb(db)

	var install_rpm string
	var stderr, stdout bytes.Buffer

	//最新版本
	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
	logs.Info(s)
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install list failed. ")
		logs.Error(stdout.String())
		logs.Error(stderr.String())
		goto end
	}

	//包名
	install_rpm = strings.Trim(stdout.String(), "\r\n")
	logs.Info("Install RPM :", install_rpm)

	for _, host := range host_list {
		serv_type1, _ := strconv.Atoi(serv_list[0].ServType)
		serv_type2, _ := strconv.Atoi(host.ServType)
		if host.Env == ob.Env && (1<<uint8(serv_type1))&serv_type2 > 0 {
			remote := host.HostName
			logs.Info("HostName", remote)

			s = fmt.Sprintf("ssh %s 'rpm -U --force /data/upgrade/%s'", remote, install_rpm)
			logs.Info(s)
			cmd = exec.Command("/bin/sh", "-c", s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				logs.Error(err)
				this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
				logs.Error(stdout.String())
				logs.Error(stderr.String())
				goto end
			}
		}
	}

	logs.Info(ob.ServName + " install passed")

end:
	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}
