// Package main ...
package main

import (
	"log"

	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Initialize()
	app := fiber.New()

	// Configure CORS to allow credentials and set the allowed origin to your frontend URL
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:3000",
	}))

	router.SetupRoutes(app)

	// app.Get("/ws", fiberWs.New(func(c *fiberWs.Conn) {
	// 	vmConsoleWebsocketURL := c.Query("console")
	// 	u, _ := url.Parse(vmConsoleWebsocketURL)
	// 	dialer := gorillaWs.Dialer{}
	// 	vmConn, _, err := dialer.Dial(u.String(), nil)
	// 	if err != nil {
	// 		log.Println("dial:", err)
	// 	}

	// 	done := make(chan struct{})
	// 	go func() {
	// 		defer close(done)
	// 		for {
	// 			_, message, err := c.ReadMessage()
	// 			if err != nil {
	// 				log.Println("read:", err)
	// 				return
	// 			}
	// 			_ = vmConn.WriteMessage(gorillaWs.TextMessage, message)
	// 		}
	// 	}()

	// 	for {
	// 		_, message, err := vmConn.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			return
	// 		}
	// 		_ = c.WriteMessage(gorillaWs.TextMessage, message)
	// 	}
	// }))

	// Set up a new cron job scheduler
	// cron := schedule.SetupCron()

	// // Schedule job - VM
	// markExpireVMErr := schedule.CronJob(cron, schedule.MarkExpireVM, "* * * ? * *")
	// if markExpireVMErr != nil {
	// 	log.Println("Error adding Mark Expire VM job scheduled job:", markExpireVMErr)
	// }
	// // expireVMErr := schedule.CronJob(cron, schedule.ExpireVM, "* * * ? * *")
	// // if expireVMErr != nil {
	// // 	log.Println("Error adding Expire VM job scheduled job:", expireVMErr)
	// // }

	// // Schedule job - User
	// markExpireUserErr := schedule.CronJob(cron, schedule.MarkExpireUser, "* * * ? * *")
	// if markExpireUserErr != nil {
	// 	log.Println("Error adding Mark Expire VM job scheduled job:", markExpireUserErr)
	// }

	// // Schedule job - Pool
	// MarkExpirePoolErr := schedule.CronJob(cron, schedule.MarkExpirePool, "* * * ? * *")
	// if MarkExpirePoolErr != nil {
	// 	log.Println("Error adding Mark Expire VM job scheduled job:", MarkExpirePoolErr)
	// }

	// // Start cron job scheduler
	// cron.Start()

	log.Fatal(app.Listen(":3002"))
}
