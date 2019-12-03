package mysql

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

var (
	db   *gorm.DB
	Conn string
)

func init() {
	//Conn = "auto_release:auto_release@tcp(localhost:3309)/dev_release?charset=utf8&parseTime=true"
	Conn = "auto_release:auto_release@tcp(localhost:3309)/run_release?charset=utf8&parseTime=true"
}

type Base struct {
	ID          uint64     `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"update_at"`
	DeletedAt   *time.Time `sql:"index" json:"-"`
	DataVersion int        `gorm:"not null;default:0;comment:'数据版本'" json:"-"`
}

type Env struct {
	Base
	Name     string `gorm:"unique;size:16;not null;default:'';comment:'HOST'" json:"name"`
	ServType int    `gorm:"not null;default:0;comment:'服务类型'" json:"serv_type"`
	Hosts    []Host `gorm:"ForeignKey:EnvId;AssociationForeignKey:ID" json:"hosts"`
}

type Host struct {
	Base
	EnvID    uint64 `gorm:"not null;default:0;comment:'EnvID'" json:"env_id"`
	Name     string `gorm:"unique;size:16;not null;default:'';comment:'HOST'" json:"name"`
	ServType int    `gorm:"not null;default:0;comment:'服务类型'" json:"serv_type"`
}

type Serv struct {
	Base
	ServName  string    `gorm:"unique;size:32;not null;default:'';comment:'服务名称'" json:"serv_name"`
	ServType  int       `gorm:"not null;default:0;comment:'服务类型'" json:"serv_type"`
	LocalPath string    `gorm:"-" json:"local_path"`
	ServMd5   string    `gorm:"-" json:"serv_md5"`
	ServEnvs  []ServEnv `gorm:"ForeignKey:ServID;AssociationForeignKey:ID" json:"serv_envs"`
}

type ServEnv struct {
	Base
	ServID     uint64 `gorm:"not null;default:0;;comment:'服务ID'" json:"serv_id"`
	ServName   string `gorm:"size:32;not null;default:'';comment:'服务名称'" json:"serv_name"`
	EnvID      uint64 `gorm:"not null;default:0;;comment:'环境ID'" json:"env_id"`
	Env        string `gorm:"size:32;not null;default:'';comment:'发布环境'" json:"env"`
	RemotePath string `gorm:"size:256;not null;default:'';comment:'安装路径'" json:"remote_path"`
	ServMd5    string `gorm:"-" json:"serv_md5"`
}

type RouteItem struct {
	Base
	ParentID   uint64      `gorm:"not null;default:0;;comment:'组ID'" json:"parent_id"`
	Name       string      `gorm:"size:32;not null;default:'';comment:'名称'" json:"name"`
	Idx        string      `gorm:"size:32;not null;default:'';comment:'index'" json:"index"`
	Class      string      `gorm:"size:32;not null;default:'';comment:'图标'" json:"class"`
	RouteItems []RouteItem `gorm:"ForeignKey:ID;AssociationForeignKey:ParentID" json:"items"`
}

type User struct {
	Base
	UserID      uint64 `gorm:"not null;default:0;comment:'用户ID'" json:"user_id"`
	Phone       string `gorm:"size:32;not null;default:'';comment:'登录手机号'" json:"phone"`
	Password    string `gorm:"size:32;not null;default:'';comment:'登录密码'" json:"password"`
	AccessLevel string `gorm:"size:32;not null;default:'';comment:'权限等级'" json:"access_level"`
}

type UserConf struct {
	Base
	ServID    uint64 `gorm:"not null;default:0;comment:'ServID'" json:"serv_id"`
	UserID    uint64 `gorm:"not null;default:0;comment:'用户ID'" json:"user_id"`
	LocalPath string `gorm:"size:256;not null;default:'';comment:'本地路径'" json:"local_path"`
}

type ServConf struct {
	Serv       Serv               `json:"serv"`
	ServEnvMap map[string]ServEnv `json:"servenv_list"`
}

type UserGroup struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type ServState struct {
	HostName string `json:"host_name" sql:"-"`
	ServTime string `json:"serv_time" sql:"-"`
}

type ServFlt struct {
	Env      string `json:"env"`
	ServName string `json:"serv_name"`
}

type LoginRequest struct {
	Phone  string `json:"phone"`
	Passwd string `json:"passwd"`
}

func InitDb() (*gorm.DB, error) {
	if db != nil {
		logs.Info("%#v", db.DB().Stats())
		return db, nil
	}
	var err error
	db, err = gorm.Open("mysql", Conn)
	if err != nil {
		logs.Error("Mysql::Open failed. " + err.Error())
		return nil, err
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 auto_increment=1")
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(100)
	return db.Debug(), nil
}
