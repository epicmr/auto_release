package models

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	ms "auto_release/models/mysql"
)

func GenSpec(serv *ms.Serv, servenv ms.ServEnv) {
	rpmPath := "/root/rpmbuild/SPECS/"
	srcName := fmt.Sprintf("%s/template/TEMPLATE_%d.spec", rpmPath, serv.ServType)
	destName := fmt.Sprintf("%s/%s", rpmPath, serv.ServName+".spec")
	suffix := path.Ext(serv.ServName)
	restartArg := strings.TrimSuffix(serv.ServName, suffix)
	CopyFile(srcName, destName)

	local_path := strings.Replace(serv.LocalPath, "/", "\\/", -1)
	remote_path := strings.Replace(servenv.RemotePath, "/", "\\/", -1)

	time := time.Now()

	v := "060102"
	v_str := time.Format(v)
	r := "150405"
	r_str := time.Format(r)

	var sed string
	sed += fmt.Sprintf("sed -i 's/VERSION/%s/g' %s;", v_str, destName)
	sed += fmt.Sprintf("sed -i 's/RELEASE/%s/g' %s;", r_str, destName)
	sed += fmt.Sprintf("sed -i 's/FILENAME/%s/g' %s;", serv.ServName+"-"+strings.ToLower(servenv.Env), destName)
	sed += fmt.Sprintf("sed -i 's/SERVNAME/%s/g' %s;", serv.ServName, destName)
	sed += fmt.Sprintf("sed -i 's/LOCAL_PATH/%s/g' %s;", local_path, destName)
	sed += fmt.Sprintf("sed -i 's/REMOTE_PATH/%s/g' %s;", remote_path, destName)

	sed += fmt.Sprintf("sed -i 's/RESTART_ARG/%s/g' %s;", restartArg, destName)
	logs.Info(sed)

	cmd := exec.Command("/bin/sh", "-c", sed)

	var stderr, stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if nil != err {
		return
	}
}
