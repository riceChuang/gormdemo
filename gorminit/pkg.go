package gorminit

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	ConnGorm *gorm.DB
	ConnSqlx *sqlx.DB
)

func InitializeGorm() (db *gorm.DB, err error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", Username, Password, Port, DBname)

	db, err = gorm.Open("mysql", connectionString)

	if err != nil {
		log.Println(err)
		log.Println("Connection Failed to Open")
	}
	log.Println("Connection Established")
	ConnGorm = db
	return
}

func InitializeSqlx() (db *sqlx.DB, err error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", Username, Password, Port, DBname)
	db, err = sqlx.Open("mysql", connectionString)

	if err != nil {
		log.Println(err)
		log.Println("Connection Failed to Open")
	}
	log.Println("Connection Established")
	ConnSqlx = db
	return
}
