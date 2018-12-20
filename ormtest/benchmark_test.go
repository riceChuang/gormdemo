package test

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"gitlab.paradise-soft.com.tw/gormdemo/gorminit"
)

//example sql
//SELECT users.name AS user_name, orders.name AS order_name
//FROM users LEFT JOIN orders on users.uid = orders.user_id
//WHERE users.uid = :id

func BenchmarkGorm(t *testing.B) {
	gormdb, err := gorminit.InitializeGorm()
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < t.N; i++ {
		gormfun(gormdb)
	}
}

func BenchmarkSqlx(t *testing.B) {

	sqlxdb, err := gorminit.InitializeSqlx()
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < t.N; i++ {
		sqlxfun(sqlxdb)
	}

}

func gormfun(gormdb *gorm.DB) {
	userOrders := []*UserOrder{}
	gormdb.Table("users").Select("users.name AS user_name, orders.name AS order_name").Joins("left join orders on users.uid = orders.user_id").Where(&User{Uid: 1}).Scan(&userOrders)
}

func sqlxfun(sqlxdb *sqlx.DB) {
	perparsql := "SELECT users.name AS user_name, orders.name AS order_name FROM users LEFT JOIN orders on users.uid = orders.user_id WHERE users.uid = :id "
	namedStmt, err := sqlxdb.PrepareNamed(perparsql)
	if err != nil {
		fmt.Println(err)
	}
	defer namedStmt.Close()

	opts := SqlxOpt{ID: 1}
	models := make([]*UserOrder, 0)
	err = namedStmt.Select(&models, opts)
	if err != nil {
		fmt.Println(err)
	}
}
