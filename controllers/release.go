package controllers

import (
	"github.com/astaxie/beego"
)

//ReleaseController struct
type ReleaseController struct {
	beego.Controller
	JSONRetMsg
}

// // Pack returns a list of users
// func (c *ReleaseController) Pack() {
// 	var ob ms.ServFlt
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
// 	logs.Info(ob)

// 	//更新spec版本号
// 	db, _ := ms.InitDb()
// 	ctx := context.Background()

// 	servList, err := ms.QueryServByName(ctx, nil, db, ob.ServName)
// 	if err != nil {
// 		logs.Info(err)
// 	}

// 	servEnvList, err := ms.BatchQueryServEnv(ctx, nil, db)
// 	if err != nil {
// 		logs.Info(err)
// 	}

// 	ms.CloseDb(db)

// 	var _servenv ms.ServEnv
// 	for _, servenv := range servEnvList {
// 		if servenv.ServName == servList[0].ServName &&
// 			servenv.Env == ob.Env {
// 			_servenv = servenv
// 		}
// 	}

// 	models.GenSpec(&servList[0], _servenv)

// 	//打包
// 	s := "cd /root/rpmbuild;rpmbuild -bb SPECS/" + ob.ServName + ".spec"
// 	logs.Info(s)
// 	cmd := exec.Command("/bin/sh", "-c", s)

// 	var stderr, stdout bytes.Buffer
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr

// 	err = cmd.Run()
// 	if nil != err {
// 		logs.Error(err)
// 		c.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "pack failed. ")
// 		logs.Error(stdout.String())
// 		logs.Error(stderr.String())
// 		goto end
// 	}

// 	logs.Info(ob.ServName + " package passed")

// end:
// 	c.Data["json"] = c.GenRetJSON()
// 	c.ServeJSON()
// }

// //Trans to server
// func (c *ReleaseController) Trans() {
// 	var ob ms.ServFlt
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
// 	logs.Info(ob)

// 	db, _ := ms.InitDb()
// 	ctx := context.Background()

// 	servList, err := ms.QueryServByName(ctx, nil, db, ob.ServName)
// 	if err != nil {
// 		logs.Info(err)
// 	}

// 	hostList, err := ms.BatchQueryHost(ctx, nil, db)
// 	if err != nil {
// 		logs.Info(err)
// 	}

// 	ms.CloseDb(db)

// 	var installRpm string
// 	var stderr, stdout bytes.Buffer

// 	//最新版本
// 	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
// 	logs.Info(s)
// 	cmd := exec.Command("/bin/sh", "-c", s)

// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr

// 	err = cmd.Run()
// 	if err != nil {
// 		c.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "trans list failed. ")
// 		logs.Error(stdout.String())
// 		logs.Error(stderr.String())
// 		goto end
// 	}

// 	//包名
// 	installRpm = strings.Trim(stdout.String(), "\r\n")
// 	logs.Info("Install RPM :", installRpm)

// 	for _, host := range hostList {
// 		servType1 := servList[0].ServType
// 		//servType2, _ := strconv.Atoi(host.ServType)
// 		servType2 := host.ServType
// 		if host.Env == ob.Env && (1<<uint8(servType1))&servType2 > 0 {
// 			remote := host.HostName
// 			logs.Info("HostName", remote)

// 			//传输
// 			s = fmt.Sprintf("cd /root/rpmbuild/RPMS/x86_64;scp %s %s:/data/upgrade", installRpm, remote)
// 			logs.Info(s)
// 			cmd = exec.Command("/bin/sh", "-c", s)

// 			err = cmd.Run()
// 			if err != nil {
// 				c.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "trans failed. ")
// 				logs.Error(stdout.String())
// 				logs.Error(stderr.String())
// 				goto end
// 			}
// 		}
// 	}

// 	logs.Info(ob.ServName + " translate passed")

// end:
// 	c.Data["json"] = c.GenRetJSON()
// 	c.ServeJSON()
// }

// //Post replace exe and restart
// func (c *ReleaseController) Post() {
// 	var ob ms.ServFlt
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &ob)
// 	logs.Info(ob)

// 	db, _ := ms.InitDb()
// 	ctx := context.Background()

// 	servList, err := ms.QueryServByName(ctx, nil, db, ob.ServName)
// 	if err != nil {
// 		logs.Info(err)
// 	}

// 	hostList, err := ms.BatchQueryHost(ctx, nil, db)
// 	if err != nil {
// 		logs.Info(err)
// 	}

// 	ms.CloseDb(db)

// 	var installRpm string
// 	var stderr, stdout bytes.Buffer

// 	//最新版本
// 	s := fmt.Sprintf("ls -lt /root/rpmbuild/RPMS/x86_64/ |grep -w %s |grep -w %s | awk -F' ' '{print $9}' |head -n 1", ob.Env, ob.ServName)
// 	logs.Info(s)
// 	cmd := exec.Command("/bin/sh", "-c", s)
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr

// 	err = cmd.Run()
// 	if err != nil {
// 		c.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install list failed. ")
// 		logs.Error(stdout.String())
// 		logs.Error(stderr.String())
// 		goto end
// 	}

// 	//包名
// 	installRpm = strings.Trim(stdout.String(), "\r\n")
// 	logs.Info("Install RPM :", installRpm)

// 	for _, host := range hostList {
// 		servType1 := servList[0].ServType
// 		//servType2, _ := strconv.Atoi(host.ServType)
// 		servType2 := host.ServType
// 		if host.Env == ob.Env && (1<<uint8(servType1))&servType2 > 0 {
// 			remote := host.HostName
// 			logs.Info("HostName", remote)

// 			s = fmt.Sprintf("ssh %s 'rpm -U --force /data/upgrade/%s'", remote, installRpm)
// 			logs.Info(s)
// 			cmd = exec.Command("/bin/sh", "-c", s)
// 			cmd.Stdout = &stdout
// 			cmd.Stderr = &stderr
// 			err = cmd.Run()
// 			if err != nil {
// 				logs.Error(err)
// 				c.setError(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus(), "install failed. ")
// 				logs.Error(stdout.String())
// 				logs.Error(stderr.String())
// 				goto end
// 			}
// 		}
// 	}

// 	logs.Info(ob.ServName + " install passed")

// end:
// 	c.Data["json"] = c.GenRetJSON()
// 	c.ServeJSON()
// }
