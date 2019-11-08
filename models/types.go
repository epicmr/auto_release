package models

import ms "auto_release/models/mysql"

type ServConf struct {
	ms.Serv
	ServMd5 string `json:"serv_md5"`
}

type ServEnvWithMd5 struct {
	ms.Serv
	ServEnvMap map[string]ms.ServEnv `json:"servenv_list"`
}
