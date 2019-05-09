package controllers

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	ms "github.com/epicmr/auto_release/models/mysql"
)

// JSONRetMsg represent return data
type JSONRetMsg struct {
	status  int
	message string
	data    interface{}
	m       map[string]interface{}
}

//GenRetJSON m to jsonM
func (r *JSONRetMsg) GenRetJSON() map[string]interface{} {
	if r.m == nil {
		r.m = make(map[string]interface{})
	}
	logs.Debug("Status:[%d] Message:[%s]", r.status, r.message)
	r.m["status"] = r.status
	r.m["message"] = r.message

	if r.status == 0 {
		r.m["data"] = r.data
	}

	return r.m
}

func (r *JSONRetMsg) setError(_status int, _message string) {
	r.status = _status
	r.message = _message
}

func (r *JSONRetMsg) setData(i interface{}) {
	r.data = i
}

//APIController for backstage api
type APIController struct {
	beego.Controller
	JSONRetMsg
}

//GetHosts return envs
func (c *APIController) GetHosts() {
	db, _ := ms.InitDb()

	var envs []ms.Env
	db.Find(&envs)

	envMap := make(map[string][]ms.Env)
	for _, env := range envs {
		envMap[env.Name] = append(envMap[env.Name], env)
	}

	c.setData(envMap)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetConfs returns confs
func (c *APIController) GetConfs() {
	db, _ := ms.InitDb()

	var envs []ms.Env
	var servs []ms.Serv

	db.Find(&envs)
	db.Debug().Find(&envs).GetErrors()
	db.Debug().Preload("ServEnvs").Find(&servs).GetErrors()
	logs.Info(envs)

	mapServEnv := make(map[string]ms.ServEnv)
	for _, serv := range servs {
		for _, servEnv := range serv.ServEnvs {
			mapServEnv[serv.ServName+servEnv.Env] = servEnv
		}
	}

	for i, serv := range servs {
		for _, env := range envs {
			servType1 := serv.ServType
			servType2 := env.ServType
			if (1<<uint8(servType1))&servType2 > 0 {
				if _, ok := mapServEnv[serv.ServName+env.Name]; !ok {
					servEnv := ms.ServEnv{
						ServName:   serv.ServName,
						EnvID:      env.ID,
						Env:        env.Name,
						RemotePath: ""}
					servs[i].ServEnvs = append(servs[i].ServEnvs, servEnv)
				}
			}
		}
	}

	c.setData(servs)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetConfsWithMd5 returns with md5
func (c *APIController) GetConfsWithMd5() {
	servName := c.GetString("serv_name")
	db, _ := ms.InitDb()

	var envs []ms.Env
	var serv ms.Serv

	db.Debug().Preload("Hosts").Find(&envs).GetErrors()
	db.Debug().Preload("ServEnvs").Where("serv_name = ?", servName).First(&serv).GetErrors()

	db.Debug().Save(&serv)

	mapServEnv := make(map[string]*ms.ServEnv)
	for i, servEnv := range serv.ServEnvs {
		mapServEnv[servEnv.Env] = &serv.ServEnvs[i]
	}

	for _, env := range envs {
		host := env.Hosts[0]
		servType1 := serv.ServType
		servType2 := env.ServType
		if (1<<uint8(servType1))&servType2 > 0 {
			if servEnv, ok := mapServEnv[env.Name]; ok {
				var stderr, stdout bytes.Buffer
				s := fmt.Sprintf("md5sum %s/%s", servEnv.RemotePath, serv.ServName)
				cmd := exec.Command("ssh", host.Name, s)
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				if err != nil {
					c.setError(1, fmt.Sprintf("HostName:[%s] exec[%s] failed. ", host.Name, s))
					logs.Error("HostName:[%s] exec[%s] failed. Error:[%s]", host.Name, s, stderr.String())
					goto end
				}
				vecList := strings.Split(stdout.String(), " ")
				logs.Debug(vecList)
				if len(vecList) > 0 {
					servEnv.ServMd5 = vecList[0]
				}
			} else {
				servEnv := ms.ServEnv{
					ServName: serv.ServName,
					Env:      env.Name}
				serv.ServEnvs = append(serv.ServEnvs, servEnv)
			}
		}
	}

	//获取本地md5
	{
		var stderr, stdout bytes.Buffer
		s := fmt.Sprintf("md5sum %s/%s", serv.LocalPath, serv.ServName)
		cmd := exec.Command("/bin/bash", "-c", s)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			c.setError(1, fmt.Sprintf("local exec[%s] failed. ", s))
			logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
			goto end
		}

		vecList := strings.Split(stdout.String(), " ")
		logs.Debug(vecList)
		if len(vecList) > 0 {
			serv.ServMd5 = vecList[0]
		}
	}

	c.setData(serv)

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//UpdateServsConf update conf
func (c *APIController) UpdateServsConf() {
    var serv ms.Serv
    json.Unmarshal(c.Ctx.Input.RequestBody, &serv)
    logs.Info(string(c.Ctx.Input.RequestBody))
    logs.Debug(serv.ServType)

    for i, _ := range serv.ServEnvs {
        serv.ServEnvs[i].ServName = serv.ServName
    }

	db, _ := ms.InitDb()
    if serv.ID > 0 {
        db.Debug().Save(&serv)
    }else {
        db.Debug().Create(&serv)
    }

    c.setData(serv)
    c.Data["json"] = c.GenRetJSON()
    c.ServeJSON()
}

func (c *APIController) GetServs() {
	db, _ := ms.InitDb()
	var servs []ms.Serv

	db.Debug().Preload("ServEnvs").Find(&servs).GetErrors()

    c.setData(servs)
    c.Data["json"] = c.GenRetJSON()
    c.ServeJSON()
}

// 	var servconf models.ServConf
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &servconf)
// 	curTime := time.Now().Unix()
// 	logs.Debug(servconf)

// 	db, _ := ms.InitDb()
// 	ctx := context.Background()
// 	servListOld, err := ms.QueryServByName(ctx, nil, db, servconf.Serv.ServName)
// 	if err != nil {
// 		logs.Debug(err)
// 	}

// 	if len(servListOld) > 0 {
// 		_servconf := servListOld[0]
// 		_servconf.LocalPath = servconf.Serv.LocalPath
// 		_servconf.LastUpdateTime = timestr
// 		logs.Debug("UpdateServ", _servconf)
// 		err = ms.UpdateServ(ctx, nil, db, &_servconf)
// 		if err != nil {
// 			logs.Debug(err)
// 		}
// 	} else {
// 		_servconf := servconf.Serv
// 		_servconf.CreateTime = curTime
// 		_servconf.UpdateTime = curTime
// 		logs.Debug("InsertServ", _servconf)
// 		err = ms.InsertServ(ctx, nil, db, _servconf)
// 		if err != nil {
// 			logs.Debug(err)
// 		}
// 	}

// 	servEnvList, err := ms.BatchQueryServEnv(ctx, nil, db)
// 	if err != nil {
// 		logs.Debug(err)
// 	}

// 	servEnvListMap := make(map[string]ms.ServEnv)
// 	for _, servEnv := range servEnvList {
// 		servEnvListMap[servEnv.ServName+servEnv.Env] = servEnv
// 	}

// 	for _, servEnv := range servconf.ServEnvMap {
// 		servEnvOld, ok := servEnvListMap[servEnv.ServName+servEnv.Env]
// 		if ok {
// 			_servenv := servEnvOld
// 			_servenv.RemotePath = servEnv.RemotePath
// 			_servenv.UpdateTime = curTime
// 			logs.Debug("UpdateServEnv", _servenv)
// 			err = ms.UpdateServEnv(ctx, nil, db, &_servenv)
// 			if err != nil {
// 				logs.Debug(err)
// 			}
// 		} else {
// 			_servenv := servEnv
// 			_servenv.CreateTime = curTime
// 			_servenv.UpdateTime = curTime
// 			logs.Debug("InsertServEnv", _servenv)
// 			err = ms.InsertServEnv(ctx, nil, db, _servenv)
// 			if err != nil {
// 				logs.Debug(err)
// 			}
// 		}
// 	}

// 	ms.CloseDb(db)

// 	c.Data["json"] = c.GenRetJSON()
// 	c.ServeJSON()
// }

// //GetServs return servs
// 	env := c.GetString("env")
// 	db, _ := ms.InitDb()
// 	ctx := context.Background()
// 	servListOld, err := ms.BatchQueryServ(ctx, nil, db)
// 	if err != nil {
// 		logs.Debug(err)
// 	}

// 	hostList, err := ms.BatchQueryHost(ctx, nil, db)
// 	if err != nil {
// 		logs.Debug(err)
// 	}

// 	servEnvList, err := ms.BatchQueryServEnv(ctx, nil, db)
// 	if err != nil {
// 		logs.Debug(err)
// 	}

// 	servType2 := 0
// 	for _, env := range hostList {
// 		if env.Env == env {
// 			//servType, _ := strconv.Atoi(env.ServType)
// 			servType := env.ServType
// 			servType2 |= servType
// 		}
// 	}

// 	servEnvListMap := make(map[string]ms.ServEnv)
// 	for _, servEnv := range servEnvList {
// 		servEnvListMap[servEnv.ServName+servEnv.Env] = servEnv
// 	}

// 	ms.CloseDb(db)

// 	servMap := make(map[string]ms.Serv)
// 	for _, serv := range servListOld {
// 		servType1 := serv.ServType
// 		if (1<<uint8(servType1))&servType2 > 0 {
// 			servEnvOld, ok := servEnvListMap[serv.ServName+env]
// 			if ok && len(servEnvOld.RemotePath) > 0 {
// 				servMap[serv.ServName] = serv
// 			}
// 		}
// 	}

// 	for _, env := range hostList {
// 		servType2 := env.ServType
// 		if env.Env == env && (servType2&12) > 0 {
// 			remote := env.HostName
// 			logs.Debug("HostName", remote)

// 			var stderr, stdout bytes.Buffer
// 			var mapTime map[string]string
// 			mapTime = make(map[string]string)
// 			var vecDetail, vecName, vecList []string
// 			s := fmt.Sprintf("ssh %s \"ps -eo etime,cmd |grep -v grep |grep cont_server\"", remote)

// 			cmd := exec.Command("/bin/sh", "-c", s)
// 			cmd.Stdout = &stdout
// 			cmd.Stderr = &stderr
// 			err = cmd.Run()
// 			if err != nil {
// 				//c.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
// 				logs.Error(stdout.String())
// 				logs.Error(stderr.String())
// 				goto end
// 			}

// 			vecList = strings.Split(stdout.String(), "\n")
// 			for _, detail := range vecList {
// 				vecDetail = strings.Split(strings.TrimSpace(detail), " ")
// 				if len(vecDetail) > 1 && strings.HasSuffix(vecDetail[1], "cont_server") {
// 					vecName = strings.Split(vecDetail[1], "/")
// 					if len(vecName) > 1 {
// 						mapTime[vecName[1]] = vecDetail[0]
// 					}
// 				}
// 			}

// 			for name, serv := range servMap {
// 				servState := ms.ServState{
// 					HostName: env.HostName,
// 					ServTime: mapTime[strings.TrimSuffix(serv.ServName, ".so")]}
// 				logs.Debug(serv.ServName, servState)
// 				//serv.ServState = append(serv.ServState, servState)
// 				servMap[name] = serv
// 			}
// 			logs.Debug(len(mapTime))
// 		}
// 	}
// 	logs.Debug("GetServs")

// end:
// 	c.Data["json"] = &servMap
// 	c.ServeJSON()
// }
