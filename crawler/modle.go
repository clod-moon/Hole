package crawler

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"log"
)

var (
	Username string = "root"
	Password string = "root"
	DbName   string = "hole"
	Host     string = "192.168.1.224"
	Port     int    = 3306

	DBHd          *gorm.DB

)

func Init(){
	mysqlstr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", Username, Password, Host, Port, DbName)
	DB, err := gorm.Open("mysql", mysqlstr)
	if err != nil {
		log.Fatalf(" gorm.Open.err: %v", err)
	}
	DBHd = DB

	if !DBHd.HasTable(&Hole{}) {
		err := DBHd.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&Hole{}).Error
		if err != nil {
			panic(err)
		}
	}
}