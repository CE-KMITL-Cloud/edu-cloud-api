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

// CreateUser - creating new user
// todo : test
func CreateUser(username, group, password, name string) (model.User, error) {
	db := ConnectDb()
	newUser := model.User{
		Username:   username,
		Password:   password, // need to see best's approach to encrypt password
		Name:       name,
		Status:     true,
		CreateTime: time.Now().UTC().Format("2006-01-02"),
		ExpireTime: time.Now().UTC().AddDate(4, 0, 0).Format("2006-01-02"),
	}
	if createErr := db.Table(group).Create(&newUser).Error; createErr != nil {
		log.Println("Error: Could not create user due to", createErr)
		return newUser, fmt.Errorf("error: could not create user due to %s", createErr)
	}
	return newUser, nil
}

// DeleteUser - delete user and user's instance limit by given username
func DeleteUser(username, group string) error {
	db := ConnectDb()
	if err := db.Table(group).Where("username = ?", username).Delete(&model.User{}).Error; err != nil {
		log.Println("Error: Could not delete user due to", err)
		return fmt.Errorf("error: could not delete user due to %s", err)
	}
	return nil
}

// EditUser - edit user by given username, group
// todo : test
func EditUser(username, group string, modifiedUser model.User) error {
	db := ConnectDb()
	if err := db.Model(&model.User{}).Table(group).Where("username = ?", username).Updates(&modifiedUser).Error; err != nil {
		log.Println("Error: Could not update username :", username)
		return fmt.Errorf("error: unable to update username : %s", username)
	}
	return nil
}

// EditInstanceLimit - edit user's instance limit by given username
// todo : test
func EditInstanceLimit(limit model.InstanceLimit) error {
	db := ConnectDb()
	if err := db.Model(&model.InstanceLimit{}).Table("instance_limit").Where("username = ?", limit.Username).Updates(&limit).Error; err != nil {
		log.Println("Error: Could not update instance limit of username :", limit.Username)
		return fmt.Errorf("error: unable to update instance limit of username : %s", limit.Username)
	}
	return nil
}

// DeleteInstanceLimit - delete user and user's instance limit by given username
// todo : test
func DeleteInstanceLimit(username string) error {
	db := ConnectDb()
	if err := db.Table("instance_limit").Where("username = ?", username).Delete(&model.InstanceLimit{}).Error; err != nil {
		log.Println("Error: Could not delete user's instance limit due to", err)
		return fmt.Errorf("error: could not delete user's instance limit due to %s", err)
	}
	return nil
}

// GetInstanceLimit - getting user's instance limit from given username
func GetInstanceLimit(username string) (model.InstanceLimit, error) {
	var limit model.InstanceLimit
	db := ConnectDb()
	db.Table("instance_limit").Where("username = ?", username).Find(&limit)
	if limit == (model.InstanceLimit{}) {
		log.Println("Error: Could not get instance limit of username :", username)
		return limit, fmt.Errorf("error: unable to get instance limit of username : %s", username)
	}
	log.Println("Got instance limit from db :", limit)
	return limit, nil
}
