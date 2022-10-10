package main

import (
	"fmt"
	"log"

	virt "github.com/edu-cloud-api/libvirt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// Get any item from .env
func getFromENV(item string) string {
	value, ok := viper.Get(item).(string)
	if !ok {
		log.Fatalf("Error while getting item : %s", value)
	}
	return value
}

func main() {
	// Setup .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// Collecting variables from .env
	AUTHNAME := getFromENV("AUTHNAME")
	PASSPHASE := getFromENV("PASSPHASE")

	// test object
	objConnection := virt.Connection{
		Host: "captain-2.ce.kmitl.cloud",
		// Host: "node-02.ce.kmitl.cloud", // 10.55.0.12
		// Host: "10.20.20.100",
		Username: AUTHNAME,
		// Username:  "ce",
		Passwd:    PASSPHASE,
		Conn_type: "tls",
	}

	// Start connection with libvirt
	conn := virt.CreateCompute(objConnection)
	log.Printf("connection pointer : %s", conn)

	// Need to close connection after process done
	defer conn.Close()

	// log.Println(virt.Get_secrets(conn))
	// log.Println(conn.GetURI())
	// log.Println(conn.ListInterfaces())
	// log.Println(conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE))
	// log.Println(conn.ListDomains())
	// log.Println(conn.LookupDomainByName("ce-cloud-freeipa-2"))
	// log.Println(conn.ListNetworks())
	// log.Println(conn.ListStoragePools())
	// log.Println(conn.ListNWFilters())
	// log.Println(conn.IsAlive())
	log.Println(virt.Get_cap_xml(conn))

	app := fiber.New()

	// Get /health-check
	app.Get("/", func(c *fiber.Ctx) error {
		msg := "✋ Healthy"
		return c.SendString(msg)
	})

	api := app.Group("/api")

	// GET /api/register
	api.Get("/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	log.Fatal(app.Listen(":3001"))
}
