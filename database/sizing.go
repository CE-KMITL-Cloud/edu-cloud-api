// Package database - database's functions
package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/edu-cloud-api/model"
)

// GetAllTemplates - getting all instance templates
func GetAllTemplates() ([]model.Sizing, error) {
	var templates []model.Sizing
	DB.Table("sizing").Find(&templates)
	if len(templates) == 0 {
		log.Println("Error: Could not get instance templates list")
		return templates, errors.New("error: unable to list instance templates")
	}
	return templates, nil
}

// GetAllTemplatesID - getting all instance templates's ID
func GetAllTemplatesID() ([]string, error) {
	var templates []string
	DB.Table("sizing").Select("vmid").Find(&templates)
	if len(templates) == 0 {
		log.Println("Error: Could not get instance templates's ID list")
		return templates, errors.New("error: unable to list instances templates's ID")
	}
	return templates, nil
}

// GetTemplate - getting instance template from given vmid
func GetTemplate(vmid string) (model.Sizing, error) {
	var template model.Sizing
	DB.Table("sizing").Where("vmid = ?", vmid).Find(&template)
	if template == (model.Sizing{}) {
		log.Println("Error: Could not get instance template id :", vmid)
		return template, fmt.Errorf("error: unable to get instance template id : %s", vmid)
	}
	return template, nil
}

// IsSizingTemplate - check vm's ID that are in templates preset
func IsSizingTemplate(vmid string) (bool, error) {
	_, getTemplateErr := GetTemplate(vmid)
	if getTemplateErr != nil {
		return false, getTemplateErr
	}
	return true, nil
}
