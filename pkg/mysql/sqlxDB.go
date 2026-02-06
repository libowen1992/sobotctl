package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	DefaultMaxOpenConns = 5
	DefaultMaxIdleConns = 2
)
//利用sqlx库链接数据库
func NewSqx(host, user, password, dbname, charset string, port, maxOpenConns, maxIdleConns int) (db *sqlx.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, port, dbname, charset)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	if maxOpenConns <= 0 {
		maxOpenConns = DefaultMaxOpenConns
	}
	if maxIdleConns <= 0 {
		maxIdleConns = DefaultMaxIdleConns
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	return db, nil
}

func Close(db *sqlx.DB) {
	_ = db.Close()
}
