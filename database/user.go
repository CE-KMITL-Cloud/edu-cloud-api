// Package database - database's functions
package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/edu-cloud-api/model"
)

// GetAllUsersByGroup - getting all users from given group
func GetAllUsersByGroup(group string) ([]model.User, error) {
	var users []model.User
	db := ConnectDb()
	db.Table(group).Find(&users)
	if len(users) == 0 {
		log.Printf("Error: Could not get %s list", group)
		return users, fmt.Errorf("error: unable to get %s list", group)
	}
	log.Println("Got user list from db :", users)
	return users, nil
}

// GetUser - getting user from given username, group
func GetUser(username, group string) (model.User, error) {
	var user model.User
	db := ConnectDb()
	db.Table(group).Where("username = ?", username).Find(&user)
	if user == (model.User{}) {
		log.Printf("Error: Could not get %s username : %s", group, username)
		return user, fmt.Errorf("error: unable to get %s username : %s", group, username)
	}
	log.Println("Got user from db :", user)
	return user, nil
}

// GetUsers - getting user's username list from all tables
func GetUsers() ([]string, error) {
	db := ConnectDb()
	var usernames []string
	err := db.Raw("SELECT username FROM student UNION SELECT username FROM faculty UNION SELECT username FROM admin").Pluck("username", &usernames).Error
	if err != nil {
		log.Printf("Error: Could not get users due to : %s", err)
		return usernames, fmt.Errorf("error: unable to get users due to : %s", err)
	}
	return usernames, nil
}

// GetUserGroup - getting user's group from {student, faculty, admin} tables
func GetUserGroup(username string) (string, error) {
	db := ConnectDb()
	var group string
	if err := db.Raw(`
    SELECT
        'student' AS group
    FROM
        student
    WHERE
        username = ?
    UNION
    SELECT
        'faculty' AS group
    FROM
        faculty
    WHERE
        username = ?
    UNION
    SELECT
        'admin' AS group
    FROM
        admin
    WHERE
        username = ?;
`, username, username, username).Scan(&group).Error; err != nil {
		return group, err
	}
	if group == "" {
		return group, errors.New("error: user not found")
	}
	return group, nil
}

// CreateUserDB - creating new user in DB
func CreateUserDB(body *model.CreateUserDB) (model.User, error) {
	db := ConnectDb()
	newUser := model.User{
		Username:   body.Username,
		Password:   body.Password, // need to see best's approach to encrypt password
		Name:       body.Name,
		Status:     true,
		CreateTime: time.Now().UTC().Format("2006-01-02"),
		ExpireTime: time.Now().UTC().AddDate(4, 0, 0).Format("2006-01-02"),
	}
	if createErr := db.Table(body.Group).Create(&newUser).Error; createErr != nil {
		log.Println("Error: Could not create user due to", createErr)
		return newUser, fmt.Errorf("error: could not create user due to %s", createErr)
	}
	return newUser, nil
}

// DeleteUserDB - delete user and user's instance limit by given username
func DeleteUserDB(username, group string) error {
	db := ConnectDb()
	if err := db.Table(group).Where("username = ?", username).Delete(&model.User{}).Error; err != nil {
		log.Println("Error: Could not delete user due to", err)
		return fmt.Errorf("error: could not delete user due to %s", err)
	}
	return nil
}

// EditUser - edit user by given username, group
func EditUser(username, group string, body *model.EditUserDB) error {
	db := ConnectDb()
	modifiedUser := model.User{
		Username:   username,
		Password:   body.Password, // need to see best's approach to encrypt password
		Name:       body.Name,
		Status:     true,
		CreateTime: time.Now().UTC().Format("2006-01-02"),
		ExpireTime: time.Now().UTC().AddDate(4, 0, 0).Format("2006-01-02"),
	}
	if err := db.Model(&model.User{}).Table(group).Where("username = ?", username).Updates(&modifiedUser).Error; err != nil {
		log.Println("Error: Could not update username :", username)
		return fmt.Errorf("error: unable to update username : %s", username)
	}
	return nil
}
