package models

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// @ type declration

// type that stores OrderModel type struct
type DBModel struct {
	OrderStore OrderModelStore //& as this stores all the methods belongs to the OrderModel
	UserStore UserModelStore //& this also stores db & methods belongs to it
	DbConn *gorm.DB
}


// initializes db connection and stores in DBModel type where ithas OrderModel that has gorm's type db connection, we are returning instance of the same
func ConnectionToDB(connectionString string) (*DBModel,error) {
	
	// opening sql connection using gorm.Open fnc
	db,err := gorm.Open(sqlite.Open(connectionString),&gorm.Config{
		// auto configuring db connection props on it
	})
	if err != nil {
		return nil,fmt.Errorf("failed to load connection due to wrong conn string -%v",err)
	}
	 
	// * for ex --> &User{} --> creates table with these fields of the struct literal
	err = db.AutoMigrate(&Order{},&OrderItem{},&User{}) //# pass in struct literal of type you'd want to create table for it
	if err != nil {
		return nil,fmt.Errorf("unexpected error occurred while doing auto migration-%v",err)
	}

	//returning db connection in instance of Type DbModel's nested interface stores
	dbModelInstance := &DBModel{
		OrderStore: NewOrderModel(db),
		UserStore: NewUserModel(db),
		DbConn : db,
	}
	return dbModelInstance,nil //! this stores sql Conn + OrderModel that stores methods on it, anf also we can later add models there to access structs from here
}

