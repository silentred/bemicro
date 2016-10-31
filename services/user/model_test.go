package user

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initORM() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "jason:jason@/fame?charset=utf8")
	orm.SetMaxOpenConns("default", 3)

	//orm.Debug = true
	//orm.DebugLog = orm.NewLog(os.Stdout)

	orm.RegisterModel(new(User))
}

func TestCRUD(t *testing.T) {
	initORM()

	o := orm.NewOrm()

	user := User{UID: 10000001}
	o.Read(&user)
	fmt.Println(user)

	user.FollowerCount = 100
	user.Username = "test2"
	o.Update(&user, "FollowerCount", "Username")
	fmt.Println(user)

	f := map[string]string{"NoExists": "hehe", "FollowerCount": "123", "FollowingCount": "321", "Coin": "999"}
	changeField(&user, f)
	fmt.Println(user)

}

func changeField(user *User, fields map[string]string) {
	userV := reflect.ValueOf(user).Elem()
	// if userV.Kind() == reflect.Ptr {
	// 	userV = userV.Elem()
	// }

	for key, value := range fields {
		keyV := userV.FieldByName(key)
		if keyV.IsValid() {
			switch keyV.Kind() {
			case reflect.Uint64, reflect.Uint32:
				if s, err := strconv.Atoi(value); err == nil {
					keyV.SetUint(uint64(s))
				}
			case reflect.Int, reflect.Int64:
				if s, err := strconv.Atoi(value); err == nil {
					keyV.SetInt(int64(s))
				}
			}

		} else {
			fmt.Println(key, "not exists")
		}
	}

}
