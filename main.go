// Package main -
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
		Passwd:   PASSPHASE,
		ConnType: "tls",
	}

	// Start connection with libvirt
	conn := virt.CreateCompute(objConnection)
	log.Println("connection pointer :", conn)

	// log.Println("connection pointer :", conn)
	// log.Println(virt.GetEmulator(conn, "i686"))
	// log.Println(virt.GetInstances(conn))
	// log.Println(virt.GetSnapshots(conn))
	// log.Println(virt.GetHostInstances(conn))
	// log.Println(virt.GetUserInstances(conn, "ce-cloud-freeipa-2"))
	// log.Println(virt.GetNetDevices(conn))
	// log.Println(virt.GetMachineTypes(conn, "i686"))
	// log.Println(virt.GetEmulators(conn))
	// log.Println(virt.GetHypervisorsDomainType(conn))
	// log.Println(virt.GetHypervisorsMachines(conn))
	// log.Println(virt.GetDomCapXML(conn, "i686", "pc-i440fx-4.2"))
	// log.Println(virt.GetCapabilities(conn, "i686"))
	// log.Println(virt.GetDomainCapabilities(conn, "i686", "pc-i440fx-focal"))
	// log.Println(virt.GetDomainCapabilities(conn, "x86_64", ""))
	log.Println(virt.CreateInstance(conn, "test-instance", "2", "2", "no-mode", "uuid", "i686", "pc", "BIOS", "", "", "", "", "", "", "", "", "", "", "", virt.OsLoaderEnum{}))
	// log.Println(virt.GetVersion(conn))
	// log.Println(virt.GetLibVersion(conn))
	// log.Println(virt.GetOsLoaders(conn, "i686", ""))
	// log.Println(virt.GetOsLoaderEnums(conn, "i686", "pc"))
	// log.Println(virt.GetDiskBusTypes(conn, "i686", "pc"))
	// log.Println(virt.GetDiskDeviceTypes(conn, "i686", "pc"))
	// log.Println(virt.GetGraphicTypes(conn, "i686", "pc"))
	// log.Println(virt.GetCPUModes(conn, "i686", "pc"))
	// log.Println(virt.GetVideoModels(conn, "i686", "pc"))
	// log.Println(virt.GetCPUCustomTypes(conn, "i686", "pc"))
	// log.Println(virt.GetVideoModels(conn, "i686", "pc-i440fx-focal"))
	// log.Println(virt.FindUEFIPathForArch(conn, "x86_64", "pc-i440fx-focal"))
	// log.Println(virt.LabelForFirmwarePath(conn, "x86_64", "/usr/share/OVMF/OVMF_CODE.fd"))
	// log.Println(virt.GetStorages(conn, false))
	// log.Println(virt.GetInterfaces(conn))
	// log.Println(virt.GetNetworks(conn))

	// Need to close connection after process done
	defer conn.Close()

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
