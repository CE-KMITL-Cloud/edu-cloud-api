// Package database - database's functions
package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// CreateInstanceLimit - create user's instance limit by given username, group
func CreateInstanceLimit(username, group string) error {
	db := ConnectDb()
	var (
		maxCPU, maxRAM, maxDisk float64
		maxInstance             uint64
	)
	switch group {
	case config.STUDENT:
		maxCPU, maxRAM, maxDisk, maxInstance = 4, 4, 40, 1
	case config.FACULTY:
		maxCPU, maxRAM, maxDisk, maxInstance = 12, 12, 120, 3
	case config.ADMIN:
		maxCPU, maxRAM, maxDisk, maxInstance = 120, 120, 1200, 30
	default:
		log.Printf("Error: Could not create instance limit of username %s due to group %s is invalid", username, group)
		return fmt.Errorf("error: unable to create instance limit of username %s due to group invalid", username)
	}
	limit := model.InstanceLimit{
		Username:    username,
		MaxCPU:      maxCPU,
		MaxRAM:      maxRAM,
		MaxDisk:     maxDisk,
		MaxInstance: maxInstance,
	}
	if err := db.Model(&model.InstanceLimit{}).Table("instance_limit").Create(&limit).Error; err != nil {
		log.Println("Error: Could not create instance limit of username :", limit.Username)
		return fmt.Errorf("error: unable to create instance limit of username : %s", limit.Username)
	}
	return nil
}

// EditInstanceLimit - edit user's instance limit by given username
func EditInstanceLimit(username string, body *model.EditInstanceLimit) error {
	db := ConnectDb()
	if body.MaxCPU > 0 && body.MaxRAM > 0 && body.MaxDisk > 0 && body.MaxInstance > 0 {
		limit := model.InstanceLimit{
			Username:    username,
			MaxCPU:      body.MaxCPU,
			MaxRAM:      body.MaxRAM,
			MaxDisk:     body.MaxDisk,
			MaxInstance: body.MaxInstance,
		}
		if err := db.Model(&model.InstanceLimit{}).Table("instance_limit").Where("username = ?", username).Updates(&limit).Error; err != nil {
			log.Println("Error: Could not update instance limit of username :", username)
			return fmt.Errorf("error: unable to update instance limit of username : %s", username)
		}
		return nil
	}
	return fmt.Errorf("error: unable to update instance limit of username : %s due to invalid limitation", username)
}

// DeleteInstanceLimit - delete user and user's instance limit by given username
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

// CheckInstanceLimit - check has instance limit reached already? and return boolean
func CheckInstanceLimit(username string, vmSpec model.VMSpec) (bool, error) {
	db := ConnectDb()
	var (
		sumCPU, remainCPU                                           float64
		limitRAM, limitDisk, sumRAM, sumDisk, remainRAM, remainDisk uint64
	)
	limit, err := GetInstanceLimit(username)
	if err != nil {
		return false, fmt.Errorf("error: could not get instance limit due to %s", err)
	}
	limitRAM, limitDisk = config.GBtoByteFloat(limit.MaxRAM), config.GBtoByteFloat(limit.MaxDisk)
	// check instance count
	var instanceCount int64
	if countErr := db.Table("instance").Where("ownerid = ?", username).Count(&instanceCount).Error; countErr != nil {
		log.Println("Error: Could not count instance due to", countErr)
		return false, fmt.Errorf("error: could not count instance due to %s", countErr)
	}
	log.Println("limit instance count :", limit.MaxInstance)
	log.Println("own instance count :", instanceCount)
	if instanceCount >= int64(limit.MaxInstance) {
		log.Println("Error: Maximum instance has reached")
		return false, errors.New("error: maximum instance has reached")
	}

	// sum spec from all instance user have and store as sumCPU, sumRAM, sumDisk then compare with limit
	if instanceCount > 0 {
		instances, instancesErr := GetAllInstancesByOwner(username)
		log.Println("all instances :", instances)
		if instancesErr != nil {
			log.Println("Error: Could not get instances due to", instancesErr)
			return false, fmt.Errorf("error: could not get instances due to %s", instancesErr)
		}
		for _, instance := range instances {
			sumCPU += instance.MaxCPU
			sumRAM += config.GBtoByteFloat(instance.MaxRAM)
			sumDisk += config.GBtoByteFloat(instance.MaxDisk)
		}
	}
	log.Printf("limit = cpu : %f, ram : %d, disk : %d", limit.MaxCPU, limitRAM, limitDisk)
	log.Printf("own spec = cpu : %f, ram : %d, disk : %d", sumCPU, sumRAM, sumDisk)
	if limitRAM > sumRAM && limit.MaxCPU > sumCPU && limitDisk > sumDisk {
		log.Println("have sufficient spec for creating VM, check spec of request's vm and remaining limit")
		remainCPU, remainRAM, remainDisk = limit.MaxCPU-sumCPU, limitRAM-sumRAM, limitDisk-sumDisk
		log.Printf("remaining limit = cpu : %f, ram : %d, disk : %d", remainCPU, remainRAM, remainDisk)
		log.Printf("VM spec = cpu : %f, ram : %d, disk : %d", vmSpec.CPU, vmSpec.Memory, vmSpec.Disk)
		if remainCPU > vmSpec.CPU && remainRAM > vmSpec.Memory && remainDisk > vmSpec.Disk {
			log.Println("able to create VM :D")
			return true, nil
		}
	}
	return false, errors.New("error: maximum instance limit has reached")
}
