#!/bin/sh

rm *.sql

CREATE_FILE_NAME=create_release
DROP_FILE_NAME=drop_release

CREATE=create
DROP=drop
UPDATE=update

if [ $# -lt 1 ]; then
echo "缺少参数"
echo "CREATE TABLE IF NOT EXISTS : $0 ${CREATE}"
echo ""
echo "================================================================================"
echo ""
echo "UPDATE TABLE : $0 ${UPDATE}"
echo ""
echo "================================================================================"
echo ""
echo "DROP TABLE : $0 ${DROP}"
echo ""
echo "================================================================================"
echo ""
exit
fi
    
echo "
CREATE TABLE IF NOT EXISTS \`t_release\` (
    Frecord_id              INT NOT NULL auto_increment                         COMMENT '记录id' ,
    Fmd5sum                 varchar(64) NOT NULL DEFAULT ''                     COMMENT '唯一MD5',
    Fenv                    varchar(32) NOT NULL DEFAULT ''                     COMMENT '发布环境',
    Fserv_name              varchar(32) NOT NULL DEFAULT ''                     COMMENT '服务名称',
    Fserv_type              INT NOT NULL DEFAULT 0                              COMMENT '服务类型' ,
    Flocal_path             varchar(256) NOT NULL DEFAULT ''                    COMMENT '本地路径',
    Fremote_path            varchar(256) NOT NULL DEFAULT ''                    COMMENT '安装路径',
    Fcreate_time            datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '创建时间',
    Flast_update_time       datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '最后更新时间',
    Fdata_version           int unsigned NOT NULL DEFAULT '0'                   COMMENT '数据版本',
	primary key(Frecord_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
" >> ${CREATE_FILE_NAME}.sql

echo "
CREATE TABLE IF NOT EXISTS \`t_serv\` (
    Frecord_id              INT NOT NULL auto_increment                         COMMENT '记录id' ,
    Fserv_name              varchar(32) NOT NULL DEFAULT ''                     COMMENT '服务名称',
    Fserv_type              INT NOT NULL DEFAULT 0                              COMMENT '服务类型' ,
    Flocal_path             varchar(256) NOT NULL DEFAULT ''                    COMMENT '本地路径',
    Fcreate_time            datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '创建时间',
    Flast_update_time       datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '最后更新时间',
    Fdata_version           int unsigned NOT NULL DEFAULT '0'                   COMMENT '数据版本',
	primary key(Frecord_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
" >> ${CREATE_FILE_NAME}.sql

echo "
CREATE TABLE IF NOT EXISTS \`t_env\` (
    Frecord_id              INT NOT NULL auto_increment                         COMMENT '记录id' ,
    Fenv                    varchar(32) NOT NULL DEFAULT ''                     COMMENT '发布环境',
    Fcreate_time            datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '创建时间',
    Flast_update_time       datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '最后更新时间',
    Fdata_version           int unsigned NOT NULL DEFAULT '0'                   COMMENT '数据版本',
	primary key(Frecord_id),
	unique key(Fenv)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
" >> ${CREATE_FILE_NAME}.sql

echo "
CREATE TABLE IF NOT EXISTS \`t_env_host\` (
    Frecord_id              INT NOT NULL auto_increment                         COMMENT '记录id' ,
    Fenv                    varchar(32) NOT NULL DEFAULT ''                     COMMENT '发布环境',
    Fhost_name              varchar(16) NOT NULL DEFAULT ''                     COMMENT 'HOST',
    Fserv_type              INT NOT NULL DEFAULT 0                              COMMENT '服务类型' ,
    Fcreate_time            datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '创建时间',
    Flast_update_time       datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '最后更新时间',
    Fdata_version           int unsigned NOT NULL DEFAULT '0'                   COMMENT '数据版本',
	primary key(Frecord_id),
	unique key(Fenv,Fhost_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
" >> ${CREATE_FILE_NAME}.sql

echo "
CREATE TABLE IF NOT EXISTS \`t_serv_env\` (
    Frecord_id              INT NOT NULL auto_increment                         COMMENT '记录id' ,
    Fserv_name              varchar(32) NOT NULL DEFAULT ''                     COMMENT '服务名称',
    Fenv                    varchar(32) NOT NULL DEFAULT ''                     COMMENT '发布环境',
    Fremote_path            varchar(256) NOT NULL DEFAULT ''                    COMMENT '安装路径',
    Fcreate_time            datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '创建时间',
    Flast_update_time       datetime NOT NULL DEFAULT '0000-00-00 00:00:00'     COMMENT '最后更新时间',
    Fdata_version           int unsigned NOT NULL DEFAULT '0'                   COMMENT '数据版本',
	primary key(Frecord_id),
	unique key(Fserv_name,Fenv)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
" >> ${CREATE_FILE_NAME}.sql

echo "
drop table envs;
drop table hosts;
drop table servs;
drop table serv_envs;
" >> ${DROP_FILE_NAME}.sql

db_connection="mysql -uauto_release -pauto_release -h119.23.163.203 -A  -s -N --default-character-set=utf8 -P3309"
 
[[ $1 == ${DROP} ]] && \
${db_connection} dev_release < ${DROP_FILE_NAME}.sql && \
echo "fihished drop talbe in file of ${DROP}.sql" && \
exit


[[ $1 == ${CREATE} ]] && \
${db_connection} dev_release < ${CREATE_FILE_NAME}.sql && \
exit

[[ $1 == ${UPDATE} ]] && \
${db_connection[1]} -f db_user < ${UPDATE}.sql && \
echo "finshed update table in file of ${UPDATE}.sql" && \
exit

echo "参数错误 : ${DROP} | ${CREATE} | ${UPDATE}" && \
exit
