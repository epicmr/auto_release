package mysql

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

type Base struct {
	ID          int64 `gorm:"primary_key;not null;comment:'ID'" json:"id"`
	CreateTime  int64 `gorm:"not null;default:0;comment:'创建时间'" json:"create_time"`
	UpdateTime  int64 `gorm:"not null;default:0;comment:'更新时间'" json:"update_time"`
	DataVersion int   `gorm:"not null;default:0;comment:'数据版本'" json:"-"`
}

type Env struct {
	Base
	Name     string `gorm:"size:16;not null;default:'';comment:'HOST'" json:"host_name"`
	ServType int    `gorm:"not null;default:0;comment:'服务类型'" json:"serv_type"`
	Hosts    []Host `gorm:"ForeignKey:EnvId;AssociationForeignKey:ID" json:"hosts"`
}

type Host struct {
	Base
	EnvID    int64  `gorm:"not null;default:0;comment:'EnvID'" json:"env_id"`
	HostName string `gorm:"size:16;not null;default:'';comment:'HOST'" json:"host_name"`
	ServType int    `gorm:"not null;default:0;comment:'服务类型'" json:"serv_type"`
	User     string `gorm:"size:32;not null;default:'';comment:'user'" json:"user"`
	IP       string `gorm:"size:16;not null;default:'';comment:'IP'" json:"ip"`
	Port     int    `gorm:"not null;default:0;comment:'Port'" json:"port"`
	KeyFile  string `gorm:"size:64;not null;default:'';comment:'KeyFile'" json:"key_file"`
}

type Serv struct {
	Base
	ServName  string    `gorm:"size:32;not null;default:'';comment:'服务名称'" json:"serv_name"`
	ServType  int       `gorm:"not null;default:0;comment:'服务类型'" json:"serv_type"`
	LocalPath string    `gorm:"size:256;not null;default:'';comment:'本地路径'" json:"local_path"`
	ServMd5   string    `gorm:"-" json:"serv_md5"`
	ServEnvs  []ServEnv `gorm:"ForeignKey:ServName;AssociationForeignKey:ServName" json:"serv_envs"`
}

type ServEnv struct {
	Base
	ServName   string `gorm:"size:32;not null;default:'';comment:'服务名称'" json:"serv_name"`
	Env        string `gorm:"size:32;not null;default:'';comment:'发布环境'" json:"env"`
	RemotePath string `gorm:"size:256;not null;default:'';comment:'安装路径'" json:"remote_path"`
	ServMd5    string `gorm:"-" json:"serv_md5"`
}

type ServConf struct {
	Serv       Serv               `json:"serv"`
	ServEnvMap map[string]ServEnv `json:"servenv_list"`
}

type ServState struct {
	HostName string `json:"host_name" sql:"-"`
	ServTime string `json:"serv_time" sql:"-"`
}

type ServFlt struct {
	Env      string `json:"env"`
	ServName string `json:"serv_name"`
}

func InitDb() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	db, err := gorm.Open("mysql", "auto_release:auto_release@tcp(120.25.154.225:3309)/release?charset=utf8&parseTime=true")
	if err != nil {
		logs.Error("Mysql::Open failed. ")
		return nil, nil
	}
	return db, nil
}
