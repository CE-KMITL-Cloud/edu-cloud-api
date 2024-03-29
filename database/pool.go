// Package database - database's functions
package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
	"github.com/lib/pq"
)

// GetAllPools - getting all pools
func GetAllPools() ([]model.Pool, error) {
	var pools []model.Pool
	DB.Table("pool").Find(&pools)
	if len(pools) == 0 {
		log.Println("Error: Could not get pools")
		return pools, errors.New("error: unable to list pools")
	}
	return pools, nil
}

// GetPoolsByOwner - getting all pools by given owner
func GetPoolsByOwner(owner string) ([]model.Pool, error) {
	var pools []model.Pool
	DB.Table("pool").Where("owner = ?", owner).Find(&pools)
	if len(pools) == 0 {
		log.Printf("Error: Could not get pools by given owner : %s", owner)
		return pools, fmt.Errorf("error: unable to list pools from given owner : %s", owner)
	}
	return pools, nil
}

// GetPoolsByVMID - getting all pools by given VMID
func GetPoolsByVMID(vmid string) ([]model.Pool, error) {
	var pools []model.Pool
	DB.Table("pool").Where("vmid @> ARRAY[?]::text[]", vmid).Find(&pools)
	if len(pools) == 0 {
		log.Printf("Error: Could not get pools by given vmid : %s", vmid)
		return pools, fmt.Errorf("error: unable to list pools from given vmid : %s", vmid)
	}
	return pools, nil
}

// GetPoolByCode - getting pool by given course code, owner
func GetPoolByCode(code, owner string) (model.Pool, error) {
	var pool model.Pool
	if err := DB.Table("pool").Where("owner = ? AND code = ?", owner, code).Find(&pool).Error; err != nil || pool.ID == 0 {
		log.Printf("Error: Could not get pool by given owner : %s, code : %s", owner, code)
		return pool, fmt.Errorf("error: unable to list pool from given owner : %s, code : %s", owner, code)
	}
	return pool, nil
}

// GetAllPoolsByMember - getting all pools that user is member by given username
func GetAllPoolsByMember(member string) ([]model.Pool, error) {
	var pools []model.Pool
	query := fmt.Sprintf("'%s' = ANY (\"member\")", member)
	if err := DB.Table("pool").Where(query).Find(&pools).Error; err != nil {
		log.Printf("Error: Could not get pools by given member's username : %s", member)
		return pools, fmt.Errorf("error: unable to list pools from given member's username : %s", member)
	}
	return pools, nil
}

// CreatePool - creating pool
func CreatePool(body *model.CreatePoolBody) (model.Pool, error) {
	newPool := model.Pool{
		Owner:      body.Owner,
		Code:       body.Code,
		Name:       body.Name,
		VMID:       []string{},
		Member:     []string{},
		Status:     true,
		CreateTime: time.Now().UTC().Format(config.TIME_FORMAT),
		ExpireTime: time.Now().UTC().AddDate(0, 4, 0).Format(config.TIME_FORMAT),
	}
	if createErr := DB.Table("pool").Create(&newPool).Error; createErr != nil {
		log.Println("Error: Could not create pool due to", createErr)
		return model.Pool{}, fmt.Errorf("error: could not create pool due to %s", createErr)
	}
	return newPool, nil
}

// DeletePool - delete pool from given code, owner
func DeletePool(code, owner string) error {
	if err := DB.Table("pool").Where("code = ? AND owner = ?", code, owner).Delete(&model.Pool{}).Error; err != nil {
		log.Println("Error: Could not delete pool due to", err)
		return fmt.Errorf("error: could not delete pool due to %s", err)
	}
	return nil
}

// MarkPoolExpired - mark pool as expired by given ID
func MarkPoolExpired(id uint64) error {
	if err := DB.Model(&model.Pool{}).Table("pool").Where("id = ?", id).UpdateColumn("status", false).Error; err != nil {
		log.Println("Error: Could not mark pool as expired ID :", id)
		return fmt.Errorf("error: unable to mark pool as expired ID : %d", id)
	}
	return nil
}

// AddPoolMembers - edit pool by given code, owner
func AddPoolMembers(code, owner string, members pq.StringArray) error {
	if err := DB.Model(&model.Instance{}).Table("pool").Where("code = ? AND owner = ?", code, owner).UpdateColumn("member", members).Error; err != nil {
		log.Printf("Error: Could not update pool code : %s, owner : %s", code, owner)
		return fmt.Errorf("error: unable to update pool code : %s, owner : %s", code, owner)
	}
	return nil
}

// AddPoolInstances - edit pool by given code, owner
func AddPoolInstances(code, owner string, instances pq.StringArray) error {
	if err := DB.Model(&model.Instance{}).Table("pool").Where("code = ? AND owner = ?", code, owner).UpdateColumn("vmid", instances).Error; err != nil {
		log.Printf("Error: Could not update pool code : %s, owner : %s", code, owner)
		return fmt.Errorf("error: unable to update pool code : %s, owner : %s", code, owner)
	}
	return nil
}

// IsPoolMember - check is given username a one of pool's member
func IsPoolMember(code, owner, username string) bool {
	pool, getPoolErr := GetPoolByCode(code, owner)
	if getPoolErr != nil {
		return false
	}
	if config.Contains(pool.Member, username) {
		log.Printf("Found user : %s in pool which owner : %s, code : %s", username, owner, code)
		return true
	}
	log.Printf("Not found user : %s in pool which owner : %s, code : %s", username, owner, code)
	return false
}

// IsPoolOwner - check is given username a one of pool's owner
func IsPoolOwner(code, owner, username, group string) bool {
	pool, getPoolErr := GetPoolByCode(code, owner)
	if getPoolErr != nil {
		return false
	}
	if pool.Owner == username || group == config.ADMIN {
		log.Printf("Found user : %s is owner of pool which owner : %s, code : %s", username, owner, code)
		return true
	}
	log.Printf("Not found user : %s is owner of pool which owner : %s, code : %s", username, owner, code)
	return false
}

// PoolInstanceDuplicate - check given vmid is exist in specific pool
func PoolInstanceDuplicate(code, owner, vmid string) (bool, error) {
	pool, getPoolErr := GetPoolByCode(code, owner)
	if getPoolErr != nil {
		return false, fmt.Errorf("error: unable to list pool from given owner : %s, code : %s", owner, code)
	}
	if config.Contains(pool.VMID, vmid) {
		log.Printf("Found vmid : %s in pool which owner : %s, code : %s", vmid, owner, code)
		return true, nil
	}
	log.Printf("Not found vmid : %s in pool which owner : %s, code : %s", vmid, owner, code)
	return false, nil
}
