package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gitlab.paradise-soft.com.tw/gormdemo/gorm/gorminit"
)

type User struct {
	Uid       int        `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string     `gorm:"type:varchar(20);not null;index:ip_idx"`
	Account   string     `gorm:"type:varchar(30);not null"`
	Password  string     `gorm:"type:varchar(30);not null"`
	CreatedAt *time.Time `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"not null"`
}

type Order struct {
	ID        int        `gorm:"primary_key;AUTO_INCREMENT"`
	UserID    int        `gorm:"not null"`
	Name      string     `gorm:"type:varchar(20);not null"`
	Price     int        `gorm:"not null"`
	CreatedAt *time.Time `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"not null"`
}

type UserOrder struct {
	UserName  string `gorm:"user_name" db:"user_name"`
	OrderName string `gorm:"order_name" db:"order_name"`
}

type SqlxOpt struct {
	ID int `json:"id" db:"id"`
}

func Test_GormCreate(t *testing.T) {
	db, err := gorminit.InitializeGorm()
	if err != nil {
		fmt.Println(err)
	}
	//db.SingularTable(true)
	db.Debug().DropTableIfExists(&Order{}, &User{})
	if err := db.Debug().AutoMigrate(&User{}, &Order{}).Error; err != nil {
		fmt.Println(err)
	}
	db.Debug().Model(&Order{}).AddForeignKey("user_id", "users(uid)", "CASCADE", "CASCADE")
	db.Debug().Create(&User{Name: "我叫測試", Account: "tt", Password: "PP"})
	db.Debug().Create(&User{Name: "我是測試2號", Account: "aa", Password: "bb"})

	db.Debug().Create(&Order{UserID: 1, Name: "訂單一號", Price: 999})
	db.Debug().Create(&Order{UserID: 2, Name: "訂單二號", Price: 333})

}

func Test_GormSelect(t *testing.T) {

	db, err := gorminit.InitializeGorm()
	if err != nil {
		fmt.Println(err)
	}

	user := []*User{}
	//取得第一筆資料
	db.Debug().First(&user, 1)
	//Last(&user), Find(&user), db.First(&user, 1)
	fmt.Println(*user[0])
	dbName := "我是測試2號"

	//---------------------- select ALL ------------------//
	//FIND name = ?
	db.Debug().Where("name = ?", dbName).Find(&user)
	fmt.Printf("%s", prettyprint(user))
	//IN ID = ?
	db.Debug().Where("uid in (?)", []int{1, 2}).Find(&user)
	fmt.Printf("%s", prettyprint(user))
	//AND
	tx := db.Begin()
	tx = tx.Where("uid = ?", 2)
	tx = tx.Where("name = ?", dbName)
	tx.Debug().Find(&user)
	fmt.Printf("%s", prettyprint(user))
	tx.Commit()
	//Struct
	db.Debug().Where(&User{Uid: 1, Account: "tt"}).Find(&user)
	fmt.Printf("%s", prettyprint(user))
	//Or
	db.Debug().Where("uid = ?", 1).Or("uid = ?", 2).Find(&user)
	fmt.Printf("%s", prettyprint(user))

	//---------------------- select ---------------------//
	db.Debug().Select("uid, name, created_at, updated_at").Find(&user)
	fmt.Printf("%s", prettyprint(user))

	//-------------------FirstOrCreate-------------------//
	utest := &User{}
	db.Debug().FirstOrCreate(&utest, User{Name: "我是測試3號"})
	fmt.Printf("%s", prettyprint(utest))

	//----------------------- Join ----------------------//
	userOrders := []*UserOrder{}
	db.Debug().Table("users").Select("users.name AS user_name, orders.name AS order_name").Joins("left join orders on users.uid = orders.user_id").Scan(&userOrders)
	fmt.Printf("%s", prettyprint(userOrders))

	//------------------- exec SQL ------------------//
	execUser := []*User{}
	db.Debug().Raw("SELECT uid, name FROM users WHERE uid < ?", 10).Scan(&execUser)
	fmt.Printf("%s", prettyprint(execUser))

}

func Test_GormUpdate(t *testing.T) {
	db, err := gorminit.InitializeGorm()
	if err != nil {
		fmt.Println(err)
	}
	user := &User{Uid: 1}
	db.Debug().Model(&user).Update("name", "改變名子")

	db.Debug().Exec("UPDATE users SET name=? WHERE id = ?", "改回來惹", 1)

	db.Debug().Model(&user).Select("name").Updates(map[string]interface{}{"name": "賈霸天溝", "account": "abcde"})

	db.Debug().Model(&user).Omit("name").Updates(map[string]interface{}{"name": "不要問你會怕", "account": "abcde"})

}

func Test_GormDelete(t *testing.T) {

	db, err := gorminit.InitializeGorm()
	if err != nil {
		fmt.Println(err)
	}
	user := &User{Uid: 1}
	db.Debug().Delete(&user)
}

func Test_GormTransation(t *testing.T) {
	db, err := gorminit.InitializeGorm()
	if err != nil {
		fmt.Println(err)
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := []*User{}
	user = append(user, &User{Name: "我不會超過", Account: "tt", Password: "PP"})
	user = append(user, &User{Name: "我會超過字喔我會超過字喔我會超過字喔我會超過字喔我會超過字喔我會超過字喔我會超過字喔", Account: "tt", Password: "PP"})

	for _, value := range user {
		if err := db.Debug().Create(value).Error; err != nil {
			tx.Rollback()
		}
	}

	fmt.Println(tx.Commit().Error)
}

func prettyprint(data interface{}) []byte {
	b, _ := json.Marshal(data)
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "  ")
	return out.Bytes()
}
