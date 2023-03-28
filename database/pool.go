// Package database - database's functions
package database

import (
	"fmt"
	"log"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// GetPoolsByOwner - getting all pools by given owner
func GetPoolsByOwner(owner string) ([]model.Pool, error) {
	var pools []model.Pool
	db := ConnectDb()
	db.Table("pool").Where("owner = ?", owner).Find(&pools)
	if len(pools) == 0 {
		log.Printf("Error: Could not get pools by given owner : %s", owner)
		return pools, fmt.Errorf("error: unable to list pools from given owner : %s", owner)
	}
	return pools, nil
}

// GetPoolByCode - getting pool by given course code, owner
func GetPoolByCode(code, owner string) (model.Pool, error) {
	var pool model.Pool
	db := ConnectDb()
	if err := db.Table("pool").Where("owner = ? AND code = ?", owner, code).Find(&pool).Error; err != nil || pool.ID == 0 {
		log.Printf("Error: Could not get pool by given owner : %s, code : %s", owner, code)
		return pool, fmt.Errorf("error: unable to list pool from given owner : %s, code : %s", owner, code)
	}
	return pool, nil
}

// CreatePool - creating pool
func CreatePool(body *model.CreatePoolBody) (model.Pool, error) {
	db := ConnectDb()
	newPool := model.Pool{
		Owner:      body.Owner,
		Code:       body.Code,
		Name:       body.Name,
		VMID:       []string{},
		Member:     []string{},
		CreateTime: time.Now().UTC().Format("2006-01-02"),
		ExpireTime: time.Now().UTC().AddDate(0, 4, 0).Format("2006-01-02"),
	}
	if createErr := db.Table("pool").Create(&newPool).Error; createErr != nil {
		log.Println("Error: Could not create pool due to", createErr)
		return model.Pool{}, fmt.Errorf("error: could not create pool due to %s", createErr)
	}
	return newPool, nil
}

// DeletePool - delete pool from given code, owner
func DeletePool(code, owner string) error {
	db := ConnectDb()
	if err := db.Table("pool").Where("code = ? AND owner = ?", code, owner).Delete(&model.Pool{}).Error; err != nil {
		log.Println("Error: Could not delete pool due to", err)
		return fmt.Errorf("error: could not delete pool due to %s", err)
	}
	return nil
}

// EditPool - edit pool by given code, owner
func EditPool(username string, modifiedPool model.Pool) error {
	db := ConnectDb()
	if err := db.Model(&model.Instance{}).Table("pool").Where("code = ? AND owner = ?", modifiedPool.Code, modifiedPool.Owner).Updates(&modifiedPool).Error; err != nil {
		log.Printf("Error: Could not update pool code : %s, owner : %s", modifiedPool.Code, modifiedPool.Owner)
		return fmt.Errorf("error: unable to update pool code : %s, owner : %s", modifiedPool.Code, modifiedPool.Owner)
	}
	return nil
}

// IsPoolMember - check is given username a one of pool's member
func IsPoolMember(code, owner, username string) (bool, error) {
	pool, getPoolErr := GetPoolByCode(code, owner)
	if getPoolErr != nil {
		return false, fmt.Errorf("error: unable to list pool from given owner : %s, code : %s", owner, code)
	}
	if config.Contains(pool.Member, username) {
		log.Printf("Found user : %s in pool which owner : %s, code : %s", username, owner, code)
		return true, nil
	}
	log.Printf("Not found user : %s in pool which owner : %s, code : %s", username, owner, code)
	return false, nil
}

// IsPoolOwner - check is given username a one of pool's owner
func IsPoolOwner(code, owner, username string) (bool, error) {
	pool, getPoolErr := GetPoolByCode(code, owner)
	if getPoolErr != nil {
		return false, fmt.Errorf("error: unable to list pool from given owner : %s, code : %s", owner, code)
	}
	if pool.Owner == username {
		log.Printf("Found user : %s is owner of pool which owner : %s, code : %s", username, owner, code)
		return true, nil
	}
	log.Printf("Not found user : %s is owner of pool which owner : %s, code : %s", username, owner, code)
	return false, nil
}
