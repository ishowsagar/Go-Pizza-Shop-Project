package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// types declartion
type User struct {
	ID uint  `gorm:"primaryKey"`
	Username string `gorm:"unqiueIndex;not null"`
	Password string `gorm:"not null"`
}

// model that stores db
type UserModel struct {
	Db *gorm.DB
}

// method that belongs to the UserModel{gormDb} -> gets user from db query --> compare pass hash and then call next method
func(u *UserModel) AuthenticateUser(usernamePayload,passwordPayload string) (*User,error) {
	var user User

	//! get user from db 👇 --> we use "?" as sql injection prevention placeholder for query args
	query := "username=?"
	err := u.Db.Where(query,usernamePayload).First(&user).Error // query gives req --> populate first res into user var
	if err != nil {
		return nil,fmt.Errorf("Invalid credentials : user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(passwordPayload))
	if err != nil {	
		return nil,fmt.Errorf("Invalid credentials : user not found")
	}
	// otherwise return that user
	return &user,nil
}

// method that gets user from db by passin in id of that user
func(u *UserModel) GetUserByID(id string) (*User,error) {
	var user User

	// get user from id via db call
	query := "id=?"
	// todo : used where instead of first --> might need to change if gives err
	err := u.Db.Where(query,id).First(&user).Error // get res from db "orders/users" where this is query condition n args --> get the first res into this var to store this res
	if err == gorm.ErrRecordNotFound {
		return nil,err
	}
	return &user,nil
}
