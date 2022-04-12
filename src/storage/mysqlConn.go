package storage

import (
	conf "HCCTV/conf"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	path := conf.GetEnv(`MYSQL_USER`) + ":" + conf.GetEnv(`MYSQL_PASSWORD`) + "@tcp(dev_db)/" + conf.GetEnv(`MYSQL_DATABASE`)
	DB, err := sql.Open("mysql", path)
	if err != nil {
		log.Println("🚫 DB Con failed..")
		panic(err)
	}
	db = DB
	log.Println("✅ Load successfully")

}

func DB() *sql.DB {
	return db
}
