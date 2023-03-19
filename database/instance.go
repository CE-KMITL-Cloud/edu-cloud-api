// Package database - database's functions
package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// GetAllInstancesByOwner - getting all instances
func GetAllInstancesByOwner(ownerid string) ([]model.Instance, error) {
	var instances []model.Instance
	db := ConnectDb()
	db.Table("instance").Where("ownerid = ?", ownerid).Find(&instances)
	if len(instances) == 0 {
		log.Println("Error: Could not get instance list")
		return instances, errors.New("error: unable to list instances")
	}
	return instances, nil
}

// GetAllInstancesIDByOwner - getting all instances's ID by given owner ID
func GetAllInstancesIDByOwner(ownerid string) ([]string, error) {
	var instances []string
	db := ConnectDb()
	db.Table("instance").Select("vmid").Where("ownerid = ?", ownerid).Find(&instances)
	if len(instances) == 0 {
		log.Println("Error: Could not get instance's ID list from given owner ID")
		return instances, errors.New("error: unable to list instances's ID from given owner ID")
	}
	return instances, nil
}

// GetInstance - getting instance from given vmid
func GetInstance(vmid string) (model.Instance, error) {
	var instance model.Instance
	db := ConnectDb()
	db.Table("instance").Where("vmid = ?", vmid).Find(&instance)
	if instance == (model.Instance{}) {
		log.Println("Error: Could not get instance id :", vmid)
		return instance, fmt.Errorf("error: unable to get instance id : %s", vmid)
	}
	return instance, nil
}

// GetInstanceTemplate - getting instance template from given vmid
func GetInstanceTemplate(vmid string) (model.Instance, error) {
	var instance model.Instance
	db := ConnectDb()
	db.Table("instance").Where("vmid = ? AND is_template = ?", vmid, true).Find(&instance)
	if instance == (model.Instance{}) {
		log.Println("Error: Could not get instance template id :", vmid)
		return instance, fmt.Errorf("error: unable to get instance template id : %s", vmid)
	}
	return instance, nil
}

// GetAllInstanceTemplatesIDByOwner - getting all instance templates's ID from given ownerid
func GetAllInstanceTemplatesIDByOwner(ownerid string) []string {
	var instances []string
	db := ConnectDb()
	db.Table("instance").Select("vmid").Where("ownerid = ? AND is_template = ?", ownerid, true).Find(&instances)
	return instances
}

// GetAllInstanceTemplatesByOwner - getting all instance templates from given ownerid
func GetAllInstanceTemplatesByOwner(ownerid string) []model.Instance {
	var instances []model.Instance
	db := ConnectDb()
	db.Table("instance").Where("ownerid = ? AND is_template = ?", ownerid, true).Find(&instances)
	return instances
}

// GetTemplate - getting instance template from given vmid
func GetTemplate(vmid string) (model.Sizing, error) {
	var template model.Sizing
	db := ConnectDb()
	db.Table("template").Where("vmid = ?", vmid).Find(&template)
	if template == (model.Sizing{}) {
		log.Println("Error: Could not get instance template id :", vmid)
		return template, fmt.Errorf("error: unable to get instance template id : %s", vmid)
	}
	return template, nil
}

// CreateInstance - creating new instance
func CreateInstance(vmid, ownerid, node, name string, spec model.VMSpec) (model.Instance, error) {
	db := ConnectDb()
	newInstance := model.Instance{
		VMID:       vmid,
		OwnerID:    ownerid,
		Node:       node,
		IsTemplate: false,
		Name:       name,
		MaxCPU:     spec.CPU,
		MaxRAM:     config.BytetoGB(spec.Memory),
		MaxDisk:    config.BytetoGB(spec.Disk),
		CreateTime: time.Now().UTC().Format("2006-01-02"),
		ExpireTime: time.Now().UTC().AddDate(0, 4, 0).Format("2006-01-02"),
	}
	checked, err := CheckInstanceLimit(ownerid, spec)
	if err != nil {
		return model.Instance{}, fmt.Errorf("error: could not create instance due to %s", err)
	}
	if checked {
		if createErr := db.Table("instance").Create(&newInstance).Error; createErr != nil {
			log.Println("Error: Could not create instance due to", createErr)
			return newInstance, fmt.Errorf("error: could not create instance due to %s", createErr)
		}
		return newInstance, nil
	}
	return model.Instance{}, errors.New("error: could not create instance due to user's instance limit has reached")
}

// DeleteInstance - delete instance & decrease instance count in instance_limit by given vmid
func DeleteInstance(vmid string) error {
	db := ConnectDb()
	if err := db.Table("instance").Where("vmid = ?", vmid).Delete(&model.Instance{}).Error; err != nil {
		log.Println("Error: Could not delete instance due to", err)
		return fmt.Errorf("error: could not delete instance due to %s", err)
	}
	return nil
}

// EditInstance - edit instance by given vmid
func EditInstance(username string, modifiedInstance model.Instance) error {
	db := ConnectDb()
	if err := db.Model(&model.Instance{}).Table("instance").Where("vmid = ?", modifiedInstance.VMID).Updates(&modifiedInstance).Error; err != nil {
		log.Println("Error: Could not update instance :", modifiedInstance.VMID)
		return fmt.Errorf("error: unable to update instance : %s", modifiedInstance.VMID)
	}
	return nil
}

// TemplateInstance - update column `is_template` to true
func TemplateInstance(vmid string) error {
	db := ConnectDb()
	if err := db.Model(&model.Instance{}).Table("instance").Where("vmid = ?", vmid).UpdateColumn("is_template", true).Error; err != nil {
		log.Println("Error: Could not template instance :", vmid)
		return fmt.Errorf("error: unable to template instance : %s", vmid)
	}
	return nil
}

// ResizeDisk - update column `max_disk` according to template's max_disk in GiB
func ResizeDisk(vmid string, maxDisk float64) error {
	db := ConnectDb()
	if err := db.Model(&model.Instance{}).Table("instance").Where("vmid = ?", vmid).UpdateColumn("max_disk", maxDisk).Error; err != nil {
		log.Println("Error: Could not resize disk of instance :", vmid)
		return fmt.Errorf("error: unable to resize disk of instance : %s", vmid)
	}
	return nil
}

// CheckInstanceOwner - check owner of the given VMID
func CheckInstanceOwner(username, vmid string) (bool, error) {
	instance, getInstanceErr := GetInstance(vmid)
	if getInstanceErr != nil {
		log.Printf("Error: Getting instance ID : %s from DB due to %s", vmid, getInstanceErr)
		return false, getInstanceErr
	}
	group, getGroupErr := GetUserGroup(username)
	if getGroupErr != nil {
		log.Printf("Error: Getting user's group name : %s from DB due to %s", username, getGroupErr)
		return false, getGroupErr
	}
	if instance.OwnerID != username && group != config.ADMIN {
		log.Printf("Error: user is not owner of VM : %s", vmid)
		return false, fmt.Errorf("user is not owner of the given VM : %s", vmid)
	}
	return true, nil
}

// CheckInstanceTemplateOwner - check vm's or template's owner of the given VMID
func CheckInstanceTemplateOwner(username, vmid string) (bool, error) {
	template, getTemplateErr := GetInstanceTemplate(vmid)
	if getTemplateErr != nil {
		return false, getTemplateErr
	}
	group, getGroupErr := GetUserGroup(username)
	if getGroupErr != nil {
		log.Printf("Error: Getting user's group name : %s from DB due to %s", username, getGroupErr)
		return false, getGroupErr
	}
	if template.OwnerID != username && group != config.ADMIN {
		log.Printf("Error: user is not owner of VM : %s", vmid)
		return false, fmt.Errorf("user is not owner of the given VM : %s", vmid)
	}
	return true, nil
}
