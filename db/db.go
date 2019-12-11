package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	err error
)

const (
	MaxConns int = 100
	MixConns int = 2
)

func init(){
	db, err = sql.Open("mysql", "poem:poem@tcp(127.0.0.1:3306)/shici?charset=utf8&parseTime=true")
	if err != nil  {
		panic(err)
	}
	db.SetMaxIdleConns(MaxConns)
	db.SetMaxOpenConns(MixConns)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func checkError(err error ) bool {
	if err != nil {
		return true
	}
	return false
}
