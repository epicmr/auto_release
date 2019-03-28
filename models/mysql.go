package models

import (
	"context"
	"database/sql"

	"github.com/arthas29/sqlmapper"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

type Serv struct {
	RecordId       string      `json:"record_id" sql:"Frecord_id"`
	ServName       string      `json:"serv_name" sql:"Fserv_name"`
	ServType       string      `json:"serv_type" sql:"Fserv_type"`
	LocalPath      string      `json:"local_path" sql:"Flocal_path"`
	CreateTime     string      `json:"-" sql:"Fcreate_time"`
	LastUpdateTime string      `json:"-" sql:"Flast_update_time"`
	DataVersion    string      `json:"-" sql:"Fdata_version"`
	ServState      []ServState `json:"serv_state" sql:"-"`
	ServMd5        string      `json:"serv_md5" sql:"-"`
}

type Host struct {
	RecordId       string `json:"record_id" sql:"Frecord_id"`
	HostName       string `json:"host_name" sql:"Fhost_name"`
	ServType       string `json:"serv_type" sql:"Fserv_type"`
	Env            string `json:"env" sql:"Fenv"`
	CreateTime     string `json:"-" sql:"Fcreate_time"`
	LastUpdateTime string `json:"-" sql:"Flast_update_time"`
	DataVersion    string `json:"-" sql:"Fdata_version"`
}

type ServEnv struct {
	RecordId       string `json:"record_id" sql:"Frecord_id"`
	ServName       string `json:"serv_name" sql:"Fserv_name"`
	Env            string `json:"env" sql:"Fenv"`
	RemotePath     string `json:"remote_path" sql:"Fremote_path"`
	CreateTime     string `json:"-" sql:"Fcreate_time"`
	LastUpdateTime string `json:"-" sql:"Flast_update_time"`
	DataVersion    string `json:"-" sql:"Fdata_version"`
	ServMd5        string `json:"serv_md5" sql:"-"`
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

func InitDb() *sql.DB {
	//db, err := sql.Open("mysql", "auto_release:auto_release@tcp(localhost:3309)/dev_release?charset=utf8")
	db, err := sql.Open("mysql", "auto_release:auto_release@tcp(localhost:3309)/run_release?charset=utf8")
	if err != nil {
		beego.Error("Mysql::Open failed. ")
		return nil
	}
	return db
}

func BatchQueryHost(ctx context.Context, tx *sql.Tx, db *sql.DB) (error, []Host) {
	var row Host
	fm, err := sqlmapper.NewFieldsMap("t_host", &row)
	if err != nil {
		return err, nil
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return err, nil
	}

	var objs []Host
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*Host))
	}

	return nil, objs
}

func BatchQueryServ(ctx context.Context, tx *sql.Tx, db *sql.DB) (error, []Serv) {
	var row Serv
	fm, err := sqlmapper.NewFieldsMap("t_serv", &row)
	if err != nil {
		return err, nil
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return err, nil
	}

	var objs []Serv
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*Serv))
	}

	return nil, objs
}

func BatchQueryServEnv(ctx context.Context, tx *sql.Tx, db *sql.DB) (error, []ServEnv) {
	var row ServEnv
	fm, err := sqlmapper.NewFieldsMap("t_serv_env", &row)
	if err != nil {
		return err, nil
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return err, nil
	}

	var objs []ServEnv
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*ServEnv))
	}

	return nil, objs
}

func QueryServByName(ctx context.Context, tx *sql.Tx, db *sql.DB, serv_name string) (
	[]Serv, error) {

	var row Serv
	row.ServName = serv_name
	fm, err := sqlmapper.NewFieldsMap("t_serv", &row)
	if err != nil {
		return nil, err
	}

	objptrs, err := fm.SQLSelectRowsByFieldNameInDB(ctx, tx, db, "Fserv_name")
	if err != nil {
		return nil, err
	}

	var objs []Serv
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*Serv))
	}

	return objs, nil
}

func InsertServ(ctx context.Context, tx *sql.Tx, db *sql.DB, rows ...Serv) error {
	for i, tlen := 0, len(rows); i < tlen; i++ {
		fm, err := sqlmapper.NewFieldsMap("t_serv", &rows[i])
		if err != nil {
			return err
		}

		err = fm.SQLInsert(ctx, tx, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateServ(ctx context.Context, tx *sql.Tx, db *sql.DB, row *Serv) error {
	fm, err := sqlmapper.NewFieldsMap("t_serv", row)
	if err != nil {
		return err
	}

	err = fm.SQLUpdateByPriKey(ctx, tx, db)
	if err != nil {
		return err
	}

	return nil
}

func InsertServEnv(ctx context.Context, tx *sql.Tx, db *sql.DB, rows ...ServEnv) error {
	for i, tlen := 0, len(rows); i < tlen; i++ {
		fm, err := sqlmapper.NewFieldsMap("t_serv_env", &rows[i])
		if err != nil {
			return err
		}

		err = fm.SQLInsert(ctx, tx, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateServEnv(ctx context.Context, tx *sql.Tx, db *sql.DB, row *ServEnv) error {
	fm, err := sqlmapper.NewFieldsMap("t_serv_env", row)
	if err != nil {
		return err
	}

	err = fm.SQLUpdateByPriKey(ctx, tx, db)
	if err != nil {
		return err
	}

	return nil
}

//func InsertData(db *sql.DB, serv *Serv) {
//	t := "2006-01-02 15:04:05"
//	time := time.Now()
//	timestr := time.Format(t)
//	beego.Info(timestr)
//	stmt, err := db.Prepare("insert into t_release (Fmd5sum,Fenv,Fserv_name,Fserv_type,Flocal_path,Fremote_path,Flast_update_time,Fcreate_time) values (?,?,?,?,?,?,?,?)")
//	if err != nil {
//		beego.Error("Mysql::Prepare failed. ", err)
//		return
//	}
//
//	md5sum := fmt.Sprintf("%x", md5.Sum([]byte(serv.Env+strings.Trim(serv.RemotePath, "/"))))
//	res, err := stmt.Exec(md5sum, serv.Env, serv.ServName, serv.ServType, serv.LocalPath, serv.RemotePath, timestr, timestr)
//	if err != nil {
//		beego.Error("Mysql::Exec failed. ", err)
//		return
//	}
//
//	id, err := res.LastInsertId()
//	if err != nil {
//		beego.Error("Mysql::LastInsertId failed. ", err)
//		return
//	}
//	beego.Info(id)
//}
//
//func UpdateData(db *sql.DB, serv *Serv) {
//	t := "2006-01-02 15:04:05"
//	time := time.Now()
//	timestr := time.Format(t)
//	beego.Info(timestr)
//	stmt, err := db.Prepare(`UPDATE t_release SET Fmd5sum = ?,Fenv = ?,Fserv_type=?,Flocal_path=?,Fremote_path=?,Flast_update_time=? WHERE Fserv_name=?`)
//	if err != nil {
//		beego.Error("Mysql::Prepare failed. ", err)
//		return
//	}
//
//	md5sum := fmt.Sprintf("%x", md5.Sum([]byte(serv.Env+strings.Trim(serv.RemotePath, "/"))))
//	res, err := stmt.Exec(md5sum, serv.Env, serv.ServType, serv.LocalPath, serv.RemotePath, timestr, serv.ServName)
//	beego.Info(serv)
//	if err != nil {
//		beego.Error("Mysql::Exec failed. ", err)
//		return
//	}
//	num, err := res.RowsAffected()
//	beego.Info(num)
//
//}
//
//func QueryData(db *sql.DB, flt *ServFlt) map[string]string {
//	rows, err := db.Query("SELECT Fmd5sum, Fenv,Fserv_name,Fserv_type,Flocal_path,Fremote_path,Flast_update_time,Fcreate_time from t_release where Fserv_name = ? and Fenv = ?", flt.ServName, flt.Env)
//	if err != nil {
//		beego.Error("Mysql::Query failed. ", err)
//		return nil
//	}
//
//	rows.Columns()
//	columns, _ := rows.Columns()
//	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]string, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//	beego.Info(len(columns))
//
//	record := make(map[string]string)
//	for rows.Next() {
//		//将行数据保存到record字典
//		err = rows.Scan(scanArgs...)
//		for i, col := range values {
//			record[columns[i][1:]] = col
//		}
//		beego.Info(record)
//	}
//	return record
//}
//
//func QueryServsByEnv(db *sql.DB, flt *ServFlt) []Serv {
//	rows, err := db.Query("SELECT Fserv_name,Fenv,Fremote_path,Flast_update_time,Fcreate_time from t_serv_env where Fenv = ?", flt.Env)
//	if err != nil {
//		beego.Error("Mysql::Query failed. ", err)
//		return nil
//	}
//
//	rows.Columns()
//	columns, _ := rows.Columns()
//	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]string, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//
//	for rows.Next() {
//		//将行数据保存到record字典
//		record := make(map[string]string)
//		err = rows.Scan(scanArgs...)
//		for i, col := range values {
//			record[columns[i][1:]] = col
//		}
//	}
//	return nil
//}
//
//func QueryServsByEnv1(db *sql.DB, flt *ServFlt) map[string]map[string]string {
//	rows, err := db.Query("SELECT Fserv_name,Fenv,Fremote_path,Flast_update_time,Fcreate_time from t_serv_env where Fenv = ?", flt.Env)
//	if err != nil {
//		beego.Error("Mysql::Query failed. ", err)
//		return nil
//	}
//
//	rows.Columns()
//	columns, _ := rows.Columns()
//	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]string, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//
//	specFileName := make(map[string]map[string]string)
//	for rows.Next() {
//		//将行数据保存到record字典
//		record := make(map[string]string)
//		err = rows.Scan(scanArgs...)
//		for i, col := range values {
//			record[columns[i][1:]] = col
//		}
//		specFileName[record["serv_name"]] = record
//	}
//	return specFileName
//}
//
////func BatchQueryServ(db *sql.DB) map[string]map[string]string {
////	rows, err := db.Query("SELECT Fserv_name,Fserv_type,Flocal_path,Flast_update_time,Fcreate_time from t_serv")
////	if err != nil {
////		beego.Error("Mysql::Query failed.", err)
////		return nil
////	}
////
////	rows.Columns()
////	columns, _ := rows.Columns()
////	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
////	scanArgs := make([]interface{}, len(columns))
////	values := make([]string, len(columns))
////	for i := range values {
////		scanArgs[i] = &values[i]
////	}
////
////	specFileName := make(map[string]map[string]string)
////	for rows.Next() {
////		//将行数据保存到record字典
////		record := make(map[string]string)
////		err = rows.Scan(scanArgs...)
////		for i, col := range values {
////			record[columns[i][1:]] = col
////		}
////		specFileName[record["serv_name"]] = record
////	}
////	return specFileName
////}
//
//func BatchQueryServHost(db *sql.DB) map[string][]map[string]string {
//	rows, err := db.Query("SELECT Fenv,Fhost_name,Flast_update_time,Fcreate_time from t_host")
//	if err != nil {
//		beego.Error("Mysql::Query failed.", err)
//		return nil
//	}
//
//	rows.Columns()
//	columns, _ := rows.Columns()
//	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]string, len(columns))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//	beego.Info(len(columns))
//
//	specFileName := make(map[string][]map[string]string)
//	for rows.Next() {
//		//将行数据保存到record字典
//		record := make(map[string]string)
//		err = rows.Scan(scanArgs...)
//		for i, col := range values {
//			record[columns[i][1:]] = col
//		}
//		specFileName[record["env"]] = append(specFileName[record["env"]], record)
//	}
//	return specFileName
//}

func CloseDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		beego.Error("Mysql::Close failed. ")
		return
	}
	return
}
