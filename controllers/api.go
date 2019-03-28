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
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/epicmr/auto_release/models"
)

type JsonRetMessage struct {
	status  int
	message string
}

func (r *JsonRetMessage) setError(_status int, _message string) {
	r.status = _status
	r.message = _message
}

func (r *JsonRetMessage) genRetJson() map[string]string {
	var retMap map[string]string
	retMap = make(map[string]string)

	logs.Info("Status:[%d] Message:[%s]", r.status, r.message)
	retMap["status"] = strconv.Itoa(r.status)
	retMap["message"] = r.message

	return retMap
}

type ApiController struct {
	beego.Controller
	JsonRetMessage
}

func (this *ApiController) GetHosts() {
	db := models.InitDb()
	ctx := context.Background()
	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}
	models.CloseDb(db)

	env_map := make(map[string][]models.Host)
	for _, host := range host_list {
		env_map[host.Env] = append(env_map[host.Env], host)
	}

	this.Data["json"] = &env_map
	this.ServeJSON()
}

func (this *ApiController) GetConfs() {
	db := models.InitDb()
	ctx := context.Background()

	err, serv_list := models.BatchQueryServ(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	err, servenv_list := models.BatchQueryServEnv(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	servenvlist_map := make(map[string]models.ServEnv)
	for _, servenv := range servenv_list {
		servenvlist_map[servenv.ServName+servenv.Env] = servenv
	}

	models.CloseDb(db)

	servconf_list := make([]models.ServConf, 0)
	for _, serv := range serv_list {
		servenv_map := make(map[string]models.ServEnv)
		for _, host := range host_list {
			serv_type1, _ := strconv.Atoi(serv.ServType)
			serv_type2, _ := strconv.Atoi(host.ServType)
			if (1<<uint8(serv_type1))&serv_type2 > 0 {
				servenv_map[host.Env] = models.ServEnv{
					ServName:   serv.ServName,
					Env:        host.Env,
					RemotePath: servenvlist_map[serv.ServName+host.Env].RemotePath,
					ServMd5:    ""}
			}
		}
		servconf := models.ServConf{
			Serv:       serv,
			ServEnvMap: servenv_map}
		servconf_list = append(servconf_list, servconf)
	}

	this.Data["json"] = &servconf_list
	this.ServeJSON()
}

func (this *ApiController) GetConfsWithMd5() {
	serv_name := this.GetString("serv_name")
	db := models.InitDb()
	ctx := context.Background()

	err, serv_list := models.BatchQueryServ(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	err, servenv_list := models.BatchQueryServEnv(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	servenvlist_map := make(map[string]models.ServEnv)
	for _, servenv := range servenv_list {
		servenvlist_map[servenv.ServName+servenv.Env] = servenv
	}

	models.CloseDb(db)

	servconf_list := make([]models.ServConf, 0)
	for _, serv := range serv_list {
		servenv_map := make(map[string]models.ServEnv)
		for _, host := range host_list {
			serv_type1, _ := strconv.Atoi(serv.ServType)
			serv_type2, _ := strconv.Atoi(host.ServType)
			if (1<<uint8(serv_type1))&serv_type2 > 0 {

				serv_md5 := ""
				//获取远程md5
				if serv.ServName == serv_name {
					remote := host.HostName
					s := fmt.Sprintf("ssh %s \"md5sum %s/%s\"", remote, servenvlist_map[serv.ServName+host.Env].RemotePath, serv.ServName)
					logs.Info(s)

					var stderr, stdout bytes.Buffer
					cmd := exec.Command("/bin/sh", "-c", s)
					cmd.Stdout = &stdout
					cmd.Stderr = &stderr
					err = cmd.Run()
					if err != nil {
						this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
						logs.Error(stdout.String())
						logs.Error(stderr.String())
						continue
					}

					vec_list := strings.Split(stdout.String(), " ")
					logs.Info(vec_list)
					if len(vec_list) > 0 {
						serv_md5 = vec_list[0]
					}
				}

				servenv_map[host.Env] = models.ServEnv{
					ServName:   serv.ServName,
					Env:        host.Env,
					RemotePath: servenvlist_map[serv.ServName+host.Env].RemotePath,
					ServMd5:    serv_md5}
			}
		}

		//获取本地md5
		if serv.ServName == serv_name {
			s := fmt.Sprintf("md5sum %s/%s", serv.LocalPath, serv.ServName)
			logs.Info(s)

			var stderr, stdout bytes.Buffer
			cmd := exec.Command("/bin/sh", "-c", s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
				logs.Error(stdout.String())
				logs.Error(stderr.String())
				continue
			}

			vec_list := strings.Split(stdout.String(), " ")
			logs.Info(vec_list)
			if len(vec_list) > 0 {
				serv.ServMd5 = vec_list[0]
			}
		}

		servconf := models.ServConf{
			Serv:       serv,
			ServEnvMap: servenv_map}
		servconf_list = append(servconf_list, servconf)
	}

	this.Data["json"] = &servconf_list
	this.ServeJSON()
}

func (this *ApiController) UpdateServsConf() {
	var servconf models.ServConf
	json.Unmarshal(this.Ctx.Input.RequestBody, &servconf)
	timestr := time.Now().Format("2006-01-02 15:04:05")
	logs.Info(servconf)

	db := models.InitDb()
	ctx := context.Background()
	serv_list_old, err := models.QueryServByName(ctx, nil, db, servconf.Serv.ServName)
	if err != nil {
		logs.Info(err)
	}

	if len(serv_list_old) > 0 {
		_servconf := serv_list_old[0]
		_servconf.LocalPath = servconf.Serv.LocalPath
		_servconf.LastUpdateTime = timestr
		logs.Info("UpdateServ", _servconf)
		err = models.UpdateServ(ctx, nil, db, &_servconf)
		if err != nil {
			logs.Info(err)
		}
	} else {
		_servconf := servconf.Serv
		_servconf.CreateTime = timestr
		_servconf.LastUpdateTime = timestr
		logs.Info("InsertServ", _servconf)
		err = models.InsertServ(ctx, nil, db, _servconf)
		if err != nil {
			logs.Info(err)
		}
	}

	err, servenv_list := models.BatchQueryServEnv(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	servenvlist_map := make(map[string]models.ServEnv)
	for _, servenv := range servenv_list {
		servenvlist_map[servenv.ServName+servenv.Env] = servenv
	}

	for _, servenv := range servconf.ServEnvMap {
		servenv_old, ok := servenvlist_map[servenv.ServName+servenv.Env]
		if ok {
			_servenv := servenv_old
			_servenv.RemotePath = servenv.RemotePath
			_servenv.LastUpdateTime = timestr
			logs.Info("UpdateServEnv", _servenv)
			err = models.UpdateServEnv(ctx, nil, db, &_servenv)
			if err != nil {
				logs.Info(err)
			}
		} else {
			_servenv := servenv
			_servenv.CreateTime = timestr
			_servenv.LastUpdateTime = timestr
			logs.Info("InsertServEnv", _servenv)
			err = models.InsertServEnv(ctx, nil, db, _servenv)
			if err != nil {
				logs.Info(err)
			}
		}
	}

	models.CloseDb(db)

	this.Data["json"] = this.genRetJson()
	this.ServeJSON()
}

//// GetUser returns a user by id
//func (c *ApiController) GetUser() {
//	ID := c.Ctx.Input.Param(":id")
//
//	userID, err := strconv.Atoi(ID)
//
//	if err != nil {
//		logs.Info("UserID error")
//	}
//
//	c.Data["json"] = models.GetUser(userID)
//	c.ServeJSON()
//}
//

//func (this *ApiController) GetConfs() {
//	db := models.InitDb()
//	servs := models.BatchQueryServ(db)
//	envs := models.BatchQueryEnv(db)
//	//serv2hosts := models.BatchQueryServHost(db)
//	models.CloseDb(db)
//
//	//for _, serv := range servs {
//	//	serv2host_list := make([]models.Serv2Host, 0)
//	//	for _, env := range envs {
//	//		for _, serv2env := range serv2envs {
//	//			serv2host := models.Serv2Host{
//	//				HostName:   env["host_name"],
//	//				RemotePath: serv2env["remote_path"]}
//	//		}
//	//	}
//	//}
//
//	for _, serv := range servs {
//		//serv2host_list := make([]map[string]string, 0)
//		for _, env := range envs {
//			for _, host := range env {
//				if serv["serv_type"] == host["serv_type"] {
//					serv2host := make(map[string]string)
//					serv2host["host_name"] = host["host_name"]
//					//serv2host["remote_path"] = host["host_name"]
//				}
//			}
//		}
//	}
//
//	this.Data["json"] = &servs
//	this.ServeJSON()
//
//}
//
func (this *ApiController) GetServs() {
	env := this.GetString("env")
	db := models.InitDb()
	ctx := context.Background()
	err, serv_list_old := models.BatchQueryServ(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	err, host_list := models.BatchQueryHost(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	err, servenv_list := models.BatchQueryServEnv(ctx, nil, db)
	if err != nil {
		logs.Info(err)
	}

	serv_type2 := 0
	for _, host := range host_list {
		if host.Env == env {
			serv_type, _ := strconv.Atoi(host.ServType)
			serv_type2 |= serv_type
		}
	}

	servenvlist_map := make(map[string]models.ServEnv)
	for _, servenv := range servenv_list {
		servenvlist_map[servenv.ServName+servenv.Env] = servenv
	}

	models.CloseDb(db)

	serv_map := make(map[string]models.Serv)
	for _, serv := range serv_list_old {
		serv_type1, _ := strconv.Atoi(serv.ServType)
		if (1<<uint8(serv_type1))&serv_type2 > 0 {
			servenv_old, ok := servenvlist_map[serv.ServName+env]
			if ok && len(servenv_old.RemotePath) > 0 {
				serv_map[serv.ServName] = serv
			}
		}
	}

	for _, host := range host_list {
		serv_type2, _ := strconv.Atoi(host.ServType)
		if host.Env == env && (serv_type2&12) > 0 {
			remote := host.HostName
			logs.Info("HostName", remote)

			var stderr, stdout bytes.Buffer
			var mapTime map[string]string
			mapTime = make(map[string]string)
			var vec_detail, vec_name, vec_list []string
			s := fmt.Sprintf("ssh %s \"ps -eo etime,cmd |grep -v grep |grep cont_server\"", remote)

			cmd := exec.Command("/bin/sh", "-c", s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				this.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
				logs.Error(stdout.String())
				logs.Error(stderr.String())
				goto end
			}
			//logs.Info(stdout.String())

			vec_list = strings.Split(stdout.String(), "\n")
			for _, detail := range vec_list {
				vec_detail = strings.Split(strings.TrimSpace(detail), " ")
				//logs.Info(vec_detail)
				if len(vec_detail) > 1 && strings.HasSuffix(vec_detail[1], "cont_server") {
					vec_name = strings.Split(vec_detail[1], "/")
					if len(vec_name) > 1 {
						mapTime[vec_name[1]] = vec_detail[0]
					}
				}
			}

			for name, serv := range serv_map {
				serv_state := models.ServState{
					HostName: host.HostName,
					ServTime: mapTime[strings.TrimSuffix(serv.ServName, ".so")]}
				logs.Info(serv.ServName, serv_state)
				serv.ServState = append(serv.ServState, serv_state)
				serv_map[name] = serv
			}
			logs.Info(len(mapTime))
		}
	}

end:
	this.Data["json"] = &serv_map
	this.ServeJSON()
}

//// AddUser adds a new user
//func (c *ApiController) AddUser() {
//	var user models.User
//	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
//	c.Data["json"] = models.AddUser(&user)
//	c.ServeJSON()
//}
//
//// UpdateUser updates existing user by id
//func (c *ApiController) UpdateUser() {
//	var user models.User
//	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
//	c.Data["json"] = models.UpdateUser(&user)
//	c.ServeJSON()
//}
//
//// DeleteUser deletes existing user by id
//func (c *ApiController) DeleteUser() {
//	ID := c.Ctx.Input.Param(":id")
//
//	userID, err := strconv.Atoi(ID)
//
//	if err != nil {
//		logs.Info("UserID error")
//	}
//
//	if models.DeleteUser(userID) {
//		c.Abort("204")
//	}
//
//	c.Abort("404")
//}
