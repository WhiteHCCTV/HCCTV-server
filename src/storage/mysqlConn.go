package storage

import (
	conf "HCCTV/conf"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	path := conf.GetEnv(`MYSQL_USER`) + ":" + conf.GetEnv(`MYSQL_PASSWORD`) + "@tcp(dev_db)/" + conf.GetEnv(`MYSQL_DATABASE`)
	DB, err := sql.Open("mysql", path)
	if err != nil {
		panic(err)
	}
	db = DB
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to:", version)
}

func DB() *sql.DB {

	return db
}
