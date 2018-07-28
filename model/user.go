package model

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

// User correspond user table in mysql
type User struct {
	ID       int    `xorm:"pk autoincr 'id'"`
	Username string `xorm:"notnull unique"`
	Email    string `xorm:"notnull unique"`
	Password string `xorm:"notnull"`
}

// TableName defines table name
func (u User) TableName() string {
	return "user"
}

// Insert a user to the table
func (u *User) Insert(session *xorm.Session) (int, error) {
	affected, err := session.InsertOne(u)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(affected), nil
}

// GetWithEmail get user with given email
func (u *User) GetWithEmail(session *xorm.Session) (bool, error) {
	email := u.Email
	return session.Where("email = ?", email).Get(u)
}

// GetWithUsername get user with given username
func (u *User) GetWithUsername(session *xorm.Session) (bool, error) {
	username := u.Username
	return session.Where("username = ?", username).Get(u)
}