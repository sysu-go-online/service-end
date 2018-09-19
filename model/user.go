package model

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-xorm/xorm"
)

// User correspond user table in mysql
type User struct {
	ID         int        `xorm:"pk autoincr 'id'"`
	Username   string     `xorm:"notnull unique"`
	Email      string     `xorm:"notnull unique"`
	Password   string     `xorm:"notnull"`
	CreateTime *time.Time `xorm:"created"`
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

// AddUserHome create user home directory
func (u *User) AddUserHome() error {
	userHome := path.Join("/home", u.Username)
	userGit := path.Join(userHome, "git")
	err := os.MkdirAll(userHome, os.ModeDir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(userGit, os.ModeDir)
	if err != nil {
		return err
	}
	// add gitconfig and .gitconfig file
	f, err := os.OpenFile(path.Join(userGit, "gitconfig"), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.Close()
	f, err = os.OpenFile(path.Join(userGit, ".gitconfig"), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.Close()
	return nil
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

// GetWithUserID get user with given id
func (u *User) GetWithUserID(session *xorm.Session) (bool, error) {
	id := u.ID
	return session.Where("id = ?", id).Get(u)
}
