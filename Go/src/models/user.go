package models

import (
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/sirupsen/logrus"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	UserId   int    `json:"userId"`
}

func GetUser(username, password string) (*User, error) {

	initDbConnection()
	defer closeDbConnection()
	user := &User{}
	results, err := db.Query("SELECT password, username, userId FROM Users where username = ? and password = ?", username, password)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer results.Close()
	if results.Next() {
		err = results.Scan(&user.Password, &user.UserName, &user.UserId)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return user, err
}

func GetUserById(id string) (*User, error) {

	initDbConnection()
	defer closeDbConnection()
	user := &User{}
	results, err := db.Query("SELECT password, username, userId FROM Users where userId = ?", id)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer results.Close()
	if results.Next() {
		err = results.Scan(&user.Password, &user.UserName, &user.UserId)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return user, err
}
