package schedule

import (
	"fmt"
	"log"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/internal/qemu"
	"github.com/robfig/cron/v3"
)

// SetupCron - setting up cron job
func SetupCron() *cron.Cron {
	// Set up a new cron job scheduler
	c := cron.New(cron.WithSeconds())
	return c
}

// CronJob - function for retrieve any task as cron job
// 0 0 0 * * ?
func CronJob(cron *cron.Cron, task func() error, schedule string) error {
	_, err := cron.AddFunc(schedule, func() {
		if taskErr := task(); taskErr != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// ExpireVM - check expire date on instance table then mark it will be deleted
func ExpireVM() error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	instances := database.GetAllInstances()
	for _, instance := range instances {
		expireDate, _ := time.Parse(config.TIME_FORMAT, instance.ExpireTime)
		threeDaysAfter := expireDate.AddDate(0, 0, 2)
		if instance.WillBeExpire && instance.Expired && today.After(threeDaysAfter) {
			log.Printf("instance ID : %s, expire date : %s, today : %s", instance.VMID, instance.ExpireTime, today.Format(config.TIME_FORMAT))
			log.Printf("instance ID : %s was expired and will be deleted", instance.VMID)

			// Get VM's info
			vmStatusURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", instance.Node, instance.VMID))
			vm, err := qemu.GetVMUsingToken(vmStatusURL)
			if err != nil {
				log.Printf("Schedule job error : getting instance ID : %s due to %s", instance.VMID, err)
				return err
			}

			// If target VM's status is "running" then stop first
			if vm.Info.Status == "running" {
				// Stop VM
				vmStopURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/stop", instance.Node, instance.VMID))

				_, stopErr := qemu.PowerManagementUsingToken(vmStopURL, nil)
				if stopErr != nil {
					log.Printf("Error: Could not stop VMID : %s in %s : %s", instance.VMID, instance.Node, stopErr)
					return stopErr
				}

				// Waiting until stopping process has been completed
				stopped := qemu.CheckStatus(instance.Node, instance.VMID, []string{"stopped"}, false, (5 * time.Minute), time.Second)
				if stopped {
					log.Printf("Finished stopping VMID : %s in %s", instance.VMID, instance.Node)
				}
			}

			// Delete VM in Proxmox
			vmDeleteURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", instance.Node, instance.VMID))
			if deleteErr := qemu.DeleteVMUsingToken(vmDeleteURL); deleteErr != nil {
				log.Printf("Schedule job error : deleting instance ID : %s due to %s", instance.VMID, deleteErr)
				return deleteErr
			}
			deleted := qemu.DeleteCompletely(instance.Node, instance.VMID)
			if deleted {
				log.Printf("Finished deleting VMID : %s in %s", instance.VMID, instance.Node)

				// Delete VM in DB
				if deleteInstanceErr := database.DeleteInstance(instance.VMID); deleteInstanceErr != nil {
					log.Printf("Schedule job error : deleting instance ID : %s due to %s", instance.VMID, deleteInstanceErr)
					return deleteInstanceErr
				}
			}
			log.Printf("instance ID : %s was expired and completely deleted", instance.VMID)
		}
	}
	return nil
}

// MarkExpireVM - check expire date on instance table then mark it will be expired
func MarkExpireVM() error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	instances := database.GetAllInstances()
	for _, instance := range instances {
		expireDate, _ := time.Parse(config.TIME_FORMAT, instance.ExpireTime)
		oneWeekBefore := expireDate.AddDate(0, 0, -7)
		if !instance.WillBeExpire && !instance.Expired && today.After(oneWeekBefore) {
			log.Printf("instance ID : %s, expire date : %s, today : %s", instance.VMID, instance.ExpireTime, today.Format(config.TIME_FORMAT))
			log.Printf("instance ID : %s will be marked and will be expired within 7 days", instance.VMID)
			if err := database.MarkWillBeExpired(instance.VMID); err != nil {
				return err
			}
		}
		if today.Equal(expireDate) || today.After(expireDate) {
			if instance.WillBeExpire && !instance.Expired {
				log.Printf("instance ID : %s, expire date : %s, today : %s", instance.VMID, instance.ExpireTime, today.Format(config.TIME_FORMAT))
				log.Printf("instance ID : %s was expired and will be deleted within 3 days", instance.VMID)
				if err := database.MarkInstanceExpired(instance.VMID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// MarkExpireUser - check expire date on user table then mark it will be expired
func MarkExpireUser() error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	for _, group := range []string{config.STUDENT, config.FACULTY, config.ADMIN} {
		users, getUsersErr := database.GetAllUsersByGroup(group)
		if getUsersErr != nil {
			return getUsersErr
		}
		for _, user := range users {
			expireDate, _ := time.Parse(config.TIME_FORMAT, user.ExpireTime)
			oneMonthBefore := expireDate.AddDate(0, -1, 0)
			if user.Status && today.After(oneMonthBefore) {
				log.Printf("user ID : %s, expire date : %s, today : %s", user.Username, user.ExpireTime, today.Format(config.TIME_FORMAT))
				log.Printf("user ID : %s will be marked and will be expired within 30 days", user.Username)
				if err := database.MarkUserExpired(user.Username, group); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// MarkExpirePool - check expire date on pool table then mark it will be expired
func MarkExpirePool() error {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	pools, getPoolsErr := database.GetAllPools()
	if getPoolsErr != nil {
		return getPoolsErr
	}
	for _, pool := range pools {
		expireDate, _ := time.Parse(config.TIME_FORMAT, pool.ExpireTime)
		oneMonthBefore := expireDate.AddDate(0, -1, 0)
		// sevenDaysAfter := expireDate.AddDate(0, 0, 6)
		if pool.Status && today.After(oneMonthBefore) {
			log.Printf("pool ID : %d, expire date : %s, today : %s", pool.ID, pool.ExpireTime, today.Format(config.TIME_FORMAT))
			log.Printf("pool ID : %d will be marked and will be expired within 30 days", pool.ID)
			if err := database.MarkPoolExpired(pool.ID); err != nil {
				return err
			}
		}
		// if !pool.Status && today.After(sevenDaysAfter) {
		// 	log.Printf("pool ID : %d, expire date : %s, today : %s", pool.ID, pool.ExpireTime, today.Format(config.TIME_FORMAT))
		// 	log.Printf("pool ID : %d was expired, deleting ...", pool.ID)
		// 	if err := database.DeletePool(pool.Code, pool.Owner); err != nil {
		// 		return err
		// 	}
		// }
	}
	return nil
}
