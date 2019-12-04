package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	ms "auto_release/models/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/casbin/casbin"
)

var (
	baseObjGroup string
	baseSubGroup string
	E            *casbin.Enforcer
)

func init() {
	baseObjGroup = "root"
	baseSubGroup = "base"
}

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
	logs.Info("Status:[%d] Message:[%s]", r.status, r.message)
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

//GetEnvs return envs
func (c *APIController) GetEnvs() {
	_info_type, _ := c.GetInt("info_type")
	db, _ := ms.InitDb()

	var envs []ms.Env
	if _info_type == 1 {
		db = db.Preload("Hosts")
	}
	db.Find(&envs)

	c.setData(envs)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//UpdateEnv update env
func (c *APIController) UpdateEnv() {
	var env ms.Env
	json.Unmarshal(c.Ctx.Input.RequestBody, &env)
	logs.Info(string(c.Ctx.Input.RequestBody))

	db, _ := ms.InitDb()
	if env.ID > 0 {
		db.Save(&env)
	} else {
		db.Create(&env)
	}

	c.setData(env)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetHosts return hosts
func (c *APIController) GetHosts() {
	db, _ := ms.InitDb()

	var hosts []ms.Host
	db.Find(&hosts)

	c.setData(hosts)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//UpdateHost update host
func (c *APIController) UpdateHost() {
	var host ms.Host
	json.Unmarshal(c.Ctx.Input.RequestBody, &host)
	logs.Info(string(c.Ctx.Input.RequestBody))

	db, _ := ms.InitDb()
	if host.ID > 0 {
		db.Save(&host)
	} else {
		db.Create(&host)
	}

	c.setData(host)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetConfs returns confs
func (c *APIController) GetConfs() {
	db, _ := ms.InitDb()

	var envs []ms.Env
	var servs []ms.Serv
	db.Find(&envs)
	db.Debug().Preload("ServEnvs").Find(&servs).GetErrors()

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

	//修改localpath显示为用户自己设定的本地目录
	c.ModifyLocalPath(servs)
	for i, _ := range servs {
		logs.Info("ID[%v], Name[%v],localPath [%v]", servs[i].ID, servs[i].ServName, servs[i].LocalPath)
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

	db.Preload("Hosts").Find(&envs).GetErrors()
	db.Preload("ServEnvs").Where("serv_name = ?", servName).First(&serv).GetErrors()

	mapServEnv := make(map[string]*ms.ServEnv)
	for i, servEnv := range serv.ServEnvs {
		mapServEnv[servEnv.Env] = &serv.ServEnvs[i]
	}

	for _, env := range envs {
		servType1 := serv.ServType
		var host ms.Host
		for _, h := range env.Hosts {
			if (1<<uint8(servType1))&h.ServType > 0 {
				host = h
			}
		}
		servType2 := host.ServType
		if (1<<uint8(servType1))&servType2 > 0 {
			if servEnv, ok := mapServEnv[env.Name]; ok {
				var stderr, stdout bytes.Buffer
				s := fmt.Sprintf("md5sum %s/%s", servEnv.RemotePath, serv.ServName)
				logs.Info(host.Name)
				cmd := exec.Command("ssh", host.Name, s)
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				if err != nil {
					//c.setError(1, fmt.Sprintf("HostName:[%s] exec[%s] failed. ", host.Name, s))
					logs.Info("HostName:[%s] exec[%s] failed. Error:[%s]", host.Name, s, stderr.String())
					continue
					//goto end
				}
				vecList := strings.Split(stdout.String(), " ")
				logs.Info(vecList)
				if len(vecList) > 0 {
					servEnv.ServMd5 = vecList[0]
				}
			} //else {
			//		servEnv := ms.ServEnv{
			//			ServName: serv.ServName,
			//			Env:      env.Name}
			//		serv.ServEnvs = append(serv.ServEnvs, servEnv)
			//	}
		}
	}

	//获取本地md5
	/*{
		var stderr, stdout bytes.Buffer
		s := fmt.Sprintf("md5sum %s/%s", serv.LocalPath, serv.ServName)
		cmd := exec.Command("/bin/bash", "-c", s)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			//c.setError(1, fmt.Sprintf("local exec[%s] failed. ", s))
			logs.Info("exec[%s] failed. Error:[%s]", s, stderr.String())
			goto step
		}

		vecList := strings.Split(stdout.String(), " ")
		logs.Info(vecList)
		if len(vecList) > 0 {
			serv.ServMd5 = vecList[0]
		}
	}*/
	err := c.GetLocalPathMd5(&serv)
	if err != nil {
		goto step
	}

step:
	c.setData(serv)

	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//UpdateServsConf update conf
func (c *APIController) UpdateServsConf() {
	var serv ms.Serv
	var userconf ms.UserConf
	var user ms.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &serv)
	logs.Info(string(c.Ctx.Input.RequestBody))
	for i, _ := range serv.ServEnvs {
		serv.ServEnvs[i].ServName = serv.ServName
	}
	db, _ := ms.InitDb()
	if serv.ID > 0 {
		db.Debug().Save(&serv)

	} else {
		db.Debug().Create(&serv)
	}

	//备份更新的路径，防止后面查询等操作会丢失。
	LocalPathBak := serv.LocalPath

	//更新用户对应的user_conf表
	phone := c.GetSession("current_user")
	db.Where("phone = ?", phone).Find(&user)
	//新增服务传进来的json没有serv_id,需要插入后数据库再分配，所以通过查询后重新更新一下serv
	db.Debug().Where("serv_name = ?", serv.ServName).First(&serv)

	if db.Where("serv_id = ? AND user_id = ?", serv.ID, user.UserID).First(&userconf).RecordNotFound() { //创建
		userconf.LocalPath = LocalPathBak
		userconf.UserID = user.UserID
		userconf.ServID = serv.ID
		db.Debug().Create(&userconf)
	} else { //更新
		db.Debug().Model(&userconf).Where("serv_id = ? AND user_id = ?", serv.ID, user.UserID).Update("local_path", LocalPathBak)
	}

	c.setData(serv)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func (c *APIController) GetServs() {
	_env := c.GetString("env")
	db, _ := ms.InitDb()
	var servs, servs1 []ms.Serv
	var env ms.Env
	db.Preload("Hosts").Where("name = ?", _env).Find(&env)
	db.Debug().Preload("ServEnvs").Find(&servs).GetErrors()

	validServType := 0
	for _, host := range env.Hosts {
		validServType |= host.ServType
	}

	for _, serv := range servs {
		if (1<<uint8(serv.ServType))&validServType > 0 {
			servs1 = append(servs1, serv)
		}
	}

	c.setData(servs1)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func (c *APIController) CheckMD5() {
	var ob ms.ServFlt
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	logs.Info(ob)
	db, _ := ms.InitDb()

	var env ms.Env
	var serv ms.Serv

	db.Preload("Hosts").Where("name = ?", ob.Env).First(&env).GetErrors()
	db.Preload("ServEnvs").Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()

	mapServEnv := make(map[string]*ms.ServEnv)
	for i, servEnv := range serv.ServEnvs {
		mapServEnv[servEnv.Env] = &serv.ServEnvs[i]
	}

	var host ms.Host
	servType1 := serv.ServType
	for _, h := range env.Hosts {
		if (1<<uint8(servType1))&h.ServType > 0 {
			host = h
		}
	}

	servType2 := host.ServType

	//获取本地md5
	/*{
		var stderr, stdout bytes.Buffer
		s := fmt.Sprintf("md5sum %s/%s", serv.LocalPath, serv.ServName)
		cmd := exec.Command("/bin/bash", "-c", s)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			c.setError(1, fmt.Sprintf("本地[%s]MD5SUM失败. ", serv.ServName))
			logs.Info("exec[%s] failed. Error:[%s]", s, stderr.String())
			goto end
		}

		vecList := strings.Split(stdout.String(), " ")
		logs.Info(vecList)
		if len(vecList) > 0 {
			serv.ServMd5 = vecList[0]
		}
	}*/
	err := c.GetLocalPathMd5(&serv)
	if err != nil {
		c.setError(1, fmt.Sprintf("本地[%s]MD5SUM失败. ", serv.ServName))
		goto end
	}

	if (1<<uint8(servType1))&servType2 > 0 {
		if servEnv, ok := mapServEnv[env.Name]; ok {
			var stderr, stdout bytes.Buffer
			s := fmt.Sprintf("md5sum %s/%s", servEnv.RemotePath, serv.ServName)
			logs.Info(host.Name)
			cmd := exec.Command("ssh", host.Name, s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				c.setError(1, fmt.Sprintf("主机[%s]服务[%s]MD5sum失败. ", host.Name, serv.ServName))
				logs.Info("HostName:[%s] exec[%s] failed. Error:[%s]", host.Name, s, stderr.String())
				goto end
			}
			vecList := strings.Split(stdout.String(), " ")
			logs.Info(vecList)
			if len(vecList) > 0 {
				servEnv.ServMd5 = vecList[0]
				if serv.ServMd5 != servEnv.ServMd5 {
					c.setError(1, fmt.Sprintf("服务:[%s]检查MD5失败. ", ob.ServName))
					logs.Error("Serv:[%s] Env:[%s] Check md5sum failed. Serv:[%s] ServEnv:[%s]",
						ob.ServName, ob.Env, serv.ServMd5, servEnv.ServMd5)
					goto end
				}
			}
		}
	}

	c.setData(serv)

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func (c *APIController) CheckTime() {
	var ob ms.ServFlt
	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
	db, _ := ms.InitDb()

	var env ms.Env
	var serv ms.Serv
	var err error
	db.Preload("Hosts").Where("name = ?", ob.Env).First(&env).GetErrors()
	db.Preload("ServEnvs").Where("serv_name = ?", ob.ServName).First(&serv).GetErrors()
	mapServEnv := make(map[string]*ms.ServEnv)
	for i, servEnv := range serv.ServEnvs {
		mapServEnv[servEnv.Env] = &serv.ServEnvs[i]
	}

	var host ms.Host
	servType1 := serv.ServType
	servType2 := 0
	for _, h := range env.Hosts {
		if (1<<uint8(servType1))&h.ServType > 0 {
			host = h
			if ((1 << uint8(servType1)) & 12) > 0 {
				var stderr, stdout bytes.Buffer
				s := fmt.Sprintf("ps -eo lstart,cmd |grep -v grep |grep -wv vi|grep -wv vim|grep -v tail|grep -w %s", strings.TrimSuffix(ob.ServName, ".so"))
				cmd := exec.Command("ssh", host.Name, s)
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				err := cmd.Run()
				if err != nil {
					c.setError(1, fmt.Sprintf("[%s] exec[%s] failed. ", host.Name, s))
					logs.Error("HostName:[%s] exec[%s] failed. Error:[%s]", host.Name, s, stderr.String())
					goto end
				}

				vecList := strings.Split(stdout.String(), "\n")
				var runMMax int64
				for _, detail := range vecList {
					vecDetail := strings.Split(strings.TrimSpace(detail), " ./")
					if len(vecDetail) > 1 {
						sec, _ := time.Parse("Mon Jan 2 15:04:05 2006", vecDetail[0])
						runM := (time.Now().Unix() - (sec.Unix() - 8*3600)) / 60
						if runM > runMMax {
							runMMax = runM
						}
						if runM > 5 {
							c.setError(1, fmt.Sprintf("服务[%s]启动时间[%d]分", ob.ServName, runM))
							logs.Error("[%s] run time [%d] m", ob.ServName, runM)
							goto end
						}
					}
				}
				logs.Info("服务[%s]启动时间[%d]分", ob.ServName, runMMax)
			}
		}
	}

	servType2 = host.ServType
	//获取本地md5
	/*{
		var stderr, stdout bytes.Buffer
		s := fmt.Sprintf("md5sum %s/%s", serv.LocalPath, serv.ServName)
		cmd := exec.Command("/bin/bash", "-c", s)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			c.setError(1, fmt.Sprintf("本地[%s]MD5SUM失败. ", serv.ServName))
			logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
			goto end
		}

		vecList := strings.Split(stdout.String(), " ")
		if len(vecList) > 0 {
			serv.ServMd5 = vecList[0]
		}
	}*/
	err = c.GetLocalPathMd5(&serv)
	if err != nil {
		c.setError(1, fmt.Sprintf("本地[%s]MD5SUM失败. ", serv.ServName))
		goto end
	}

	if (1<<uint8(servType1))&servType2 > 0 {
		if servEnv, ok := mapServEnv[env.Name]; ok {
			var stderr, stdout bytes.Buffer
			s := fmt.Sprintf("md5sum %s/%s", servEnv.RemotePath, serv.ServName)
			cmd := exec.Command("ssh", host.Name, s)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				c.setError(1, fmt.Sprintf("主机[%s]服务[%s]MD5sum失败. ", host.Name, serv.ServName))
				logs.Error("HostName:[%s] exec[%s] failed. Error:[%s]", host.Name, s, stderr.String())
				goto end
			}
			vecList := strings.Split(stdout.String(), " ")
			if len(vecList) > 0 {
				servEnv.ServMd5 = vecList[0]
				if serv.ServMd5 != servEnv.ServMd5 {
					c.setError(1, fmt.Sprintf("服务[%s]检查MD5失败. ", ob.ServName))
					logs.Error("Serv:[%s] Env:[%s] Check md5sum failed. Serv:[%s] ServEnv:[%s]",
						ob.ServName, ob.Env, serv.ServMd5, servEnv.ServMd5)
					goto end
				}
			}
		}
	}

	c.setData(serv)

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//UpdateItem update item
func (c *APIController) UpdateItem() {
	var item ms.RouteItem
	json.Unmarshal(c.Ctx.Input.RequestBody, &item)
	logs.Info(string(c.Ctx.Input.RequestBody))
	item.RouteItems = item.RouteItems[0:0]

	db, _ := ms.InitDb()
	if item.ID > 0 {
		if item.ID == item.ParentID {
			c.setError(1, "不能设置自己为父节点 ")
			logs.Error("cant set itself as parent")
			goto end
		}
		db.Save(&item)
	} else {
		db.Create(&item)
	}

	c.setData(item)
end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func Fill(parentID uint64, m map[uint64][]ms.RouteItem) []ms.RouteItem {
	var items []ms.RouteItem
	var ok bool
	if items, ok = m[parentID]; ok {
		for i, item := range items {
			items[i].RouteItems = Fill(item.ID, m)
		}
	}

	return items
}

//GetItems return items
func (c *APIController) GetItemsTree() {
	db, _ := ms.InitDb()

	var items []ms.RouteItem
	accesslevel := c.GetAccessLevel()
	db.Where("parent_id = ?", "0").Or("id IN (?)", accesslevel).Find(&items)
	m := make(map[uint64][]ms.RouteItem)
	for _, item := range items {
		m[item.ParentID] = append(m[item.ParentID], item)
	}

	s := Fill(0, m)

	c.setData(s)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetItems return items
func (c *APIController) GetAllItems() {
	db, _ := ms.InitDb()

	var items []ms.RouteItem
	db.Find(&items)
	c.setData(items)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetItems return items
func (c *APIController) UpdateUserGroup() {
	var groups []ms.UserGroup
	json.Unmarshal(c.Ctx.Input.RequestBody, &groups)
	logs.Info(string(c.Ctx.Input.RequestBody))

	for _, group := range groups {
		if group.Name == "" {
			group.Name = baseObjGroup
		}
		if group.Group == "" {
			group.Group = baseSubGroup
		}
		E.AddNamedGroupingPolicy(group.Type, group.Name, group.Group)
	}
	E.SavePolicy()

	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//GetItems return items
func (c *APIController) UpdateGroup() {
	var groups []ms.UserGroup
	json.Unmarshal(c.Ctx.Input.RequestBody, &groups)
	logs.Info(string(c.Ctx.Input.RequestBody))

	for _, group := range groups {
		if group.Name == "" {
			group.Name = baseObjGroup
		}
		if group.Group == "" {
			group.Group = baseSubGroup
		}
		E.AddNamedGroupingPolicy(group.Type, group.Name, group.Group)
	}
	E.SavePolicy()

	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func (c *APIController) ModifyLocalPath(servs []ms.Serv) {
	var user ms.User
	var userConfs []ms.UserConf
	db, _ := ms.InitDb()
	phone := c.GetSession("current_user")
	logs.Info("==============current_user_phone: [%v]", phone)
	db.Where("phone = ?", phone).First(&user)
	logs.Info("userID: [%v]", user.UserID)
	db.Where("user_id = ?", user.UserID).Find(&userConfs)
	for i, serv := range servs {
		for j, userConf := range userConfs {
			if serv.ID == userConf.ServID {
				servs[i].LocalPath = userConfs[j].LocalPath
				logs.Info("ID[%v], Name[%v], 匹配成功，修改localPath [%v]", servs[i].ID, servs[i].ServName, servs[i].LocalPath)
			}
		}
	}
}

//Grant permission to user
func (c *APIController) Grant() {
	var user ms.User
	db, _ := ms.InitDb()
	json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	//先记录更新的权限，否则下面查询用户操作会覆盖user.AccessLevel
	accessLevel := user.AccessLevel
	if user.AccessLevel == "" || (user.UserID == 0 && user.Phone == "") {
		c.setError(2, "授权失败，缺少参数")
		logs.Error("授权失败，缺少参数")
		goto end
	}
	if db.Debug().Where("user_id = ? OR phone = ? ", user.UserID, user.Phone).First(&user).RecordNotFound() {
		c.setError(2, "授权失败！请先注册，再授权")
		logs.Error("授权失败，用户未注册")
		goto end
	}
	//开始授权
	db.Debug().Model(&user).Update("access_level", accessLevel)
	user.Password = ""
	c.setData(user)
end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//获取用户accesslevel字段
func (c *APIController) GetAccessLevel() []string {
	db, _ := ms.InitDb()

	var user ms.User
	phone := c.GetSession("current_user")
	db.Where("phone = ?", phone).Find(&user)
	logs.Info("user:[%v]", user)
	var accesslevel []string
	accesslevel = strings.Split(user.AccessLevel, ";")
	logs.Info("===========accesslevel:[%v]", accesslevel)
	return accesslevel
}

func (c *APIController) GetUser() {
	phone := c.Ctx.Input.Param(":phone")
	db, _ := ms.InitDb()

	var user ms.User
	db.Where("phone = ?", phone).First(&user)

	c.setData(user)
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

//获取本地md5
func (c *APIController) GetLocalPathMd5(serv *ms.Serv) error {
	var stderr, stdout bytes.Buffer
	var user ms.User
	var userConf ms.UserConf
	db, _ := ms.InitDb()

	phone := c.GetSession("current_user")
	db.Where("phone = ?", phone).Find(&user)
	db.Where("user_id = ? AND serv_id = ?", user.UserID, serv.ID).Find(&userConf)
	serv.LocalPath = userConf.LocalPath
	s := fmt.Sprintf("md5sum %s/%s", userConf.LocalPath, serv.ServName)
	cmd := exec.Command("/bin/bash", "-c", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logs.Info("exec[%s] failed. Error:[%s]", s, stderr.String())
		return err
	}

	vecList := strings.Split(stdout.String(), " ")
	logs.Info(vecList)
	if len(vecList) > 0 {
		serv.ServMd5 = vecList[0]
	}
	return nil
}
