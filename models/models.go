package models

import (
	"WowjoyProject/FileServer/pkg/loggin"
	"WowjoyProject/FileServer/pkg/setting"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type DownData struct {
	instance_key                           sql.NullInt64
	jpgfile, jpgremote, dcmremote, dcmfile sql.NullString
}

func init() {
	db, _ = sql.Open("mysql", setting.DBConn)
	// 数据库最大连接数
	db.SetMaxOpenConns(setting.MaxConn)
	db.SetMaxIdleConns(setting.MaxConn)
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
	loggin.Debug("数据库连接成功...")
}

func CloseDB() {
	defer db.Close()
}
