// Package main ...
package main

import (
	"log"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/router"
	"github.com/edu-cloud-api/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

func main() {
	database.Initialize()
	app := fiber.New()
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.GetFromENV("ENCRYPT_KEY"),
	}))
	router.SetupRoutes(app)

	// Set up a new cron job scheduler
	cron := schedule.SetupCron()

	// Schedule job - VM
	markExpireVMErr := schedule.CronJob(cron, schedule.MarkExpireVM, "* * * ? * *")
	if markExpireVMErr != nil {
		log.Println("Error adding Mark Expire VM job scheduled job:", markExpireVMErr)
	}
	// expireVMErr := schedule.CronJob(cron, schedule.ExpireVM, "* * * ? * *")
	// if expireVMErr != nil {
	// 	log.Println("Error adding Expire VM job scheduled job:", expireVMErr)
	// }

	// Schedule job - User
	markExpireUserErr := schedule.CronJob(cron, schedule.MarkExpireUser, "* * * ? * *")
	if markExpireUserErr != nil {
		log.Println("Error adding Mark Expire VM job scheduled job:", markExpireUserErr)
	}

	// Schedule job - Pool
	MarkExpirePoolErr := schedule.CronJob(cron, schedule.MarkExpirePool, "* * * ? * *")
	if MarkExpirePoolErr != nil {
		log.Println("Error adding Mark Expire VM job scheduled job:", MarkExpirePoolErr)
	}

	// Start cron job scheduler
	cron.Start()

	log.Fatal(app.Listen(":3001"))
}
