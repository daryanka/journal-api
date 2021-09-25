package clients

import (
	"github.com/daryanka/myorm"
	"os"
)

var ClientOrm *myorm.MyOrm

func InitializeOrm() {
	ClientOrm = myorm.DBInit(myorm.DBConnection{
		DBUsername:     os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASS"),
		DBProtocol:     "tcp",
		DBAddress:      "127.0.0.1",
		DBName:         "journal",
		DBDriver:       "mysql",
		ConnectionName: "defaultCon",
	})

	db, _ := ClientOrm.GetDBConnection("defaultCon")
	err := db.Ping()
	if err != nil {
		panic("Error connecting to database " + err.Error())
	}
}