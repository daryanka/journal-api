package clients

import "github.com/daryanka/myorm"

var ClientOrm *myorm.MyOrm

func InitializeOrm() {
	ClientOrm = myorm.DBInit(myorm.DBConnection{
		DBUsername:     "root",
		DBPassword:     "",
		DBProtocol:     "tcp",
		DBAddress:      "127.0.0.1",
		DBName:         "journal",
		DBDriver:       "mysql",
		ConnectionName: "default",
	})
}