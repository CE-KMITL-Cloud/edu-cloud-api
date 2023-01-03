// Package internal - internal functions
package internal

import (
	"io/ioutil"
	"log"
	"net/http"
)

// GetVM - GET /api2/json/nodes/{node}/qemu/{vmid}
func GetVM(hostURL string) {
	// Return objects

	resp, err := http.Get(hostURL)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)

}

// GetVMs - GET /api2/json/nodes/{node}/qemu
func GetVMs(hostURL string) {
	// Return objects

	resp, err := http.Get(hostURL)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}

// CreateVM - POST /api2/json/nodes/{node}/qemu
func CreateVM()

// DeleteVM - DELETE /api2/json/nodes/{node}/qemu/{vmid}
func DeleteVM()

// CloneVM - POST /api2/json/nodes/{node}/qemu/{vmid}/clone
func CloneVM()

// CreateTemplate - POST /api2/json/nodes/{node}/qemu/{vmid}/template
func CreateTemplate()
