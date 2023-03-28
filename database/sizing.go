// Package database - database's functions
package database

import (
	"errors"
	"log"

	"github.com/edu-cloud-api/model"
)

// GetAllTemplates - getting all instance templates
func GetAllTemplates() ([]model.Sizing, error) {
	var templates []model.Sizing
	db := ConnectDb()
	db.Table("sizing").Find(&templates)
	if len(templates) == 0 {
		log.Println("Error: Could not get instance templates list")
		return templates, errors.New("error: unable to list instance templates")
	}
	return templates, nil
}

// GetAllTemplatesID - getting all instance templates's ID
func GetAllTemplatesID() ([]string, error) {
	var templates []string
	db := ConnectDb()
	db.Table("sizing").Select("vmid").Find(&templates)
	if len(templates) == 0 {
		log.Println("Error: Could not get instance templates's ID list")
		return templates, errors.New("error: unable to list instances templates's ID")
	}
	return templates, nil
}

// IsSizingTemplate - check vm's ID that are in templates preset
func IsSizingTemplate(vmid string) (bool, error) {
	_, getTemplateErr := GetTemplate(vmid)
	if getTemplateErr != nil {
		return false, getTemplateErr
	}
	return true, nil
}
