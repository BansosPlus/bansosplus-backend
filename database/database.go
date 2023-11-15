package database

import (
	"database/sql"
	"fmt"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/BansosPlus/bansosplus-backend.git/utility"
)

func Open(conf utility.Configuration) (db *sql.DB, gormDB *gorm.DB, err error) {
	defer utility.RecoverError()

	strConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", conf.Database.Username, conf.Database.Password, conf.Database.Hostname, conf.Database.Port, conf.Database.DatabaseName)
	db, err = sql.Open("mysql", strConnectionString)
	if err != nil {
		return
	}

	gormDB, err = gorm.Open(gormMySQL.New(gormMySQL.Config{
		Conn: db,
	}), &gorm.Config{})

	return
}