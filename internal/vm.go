// Package internal - internal functions
package internal

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/model"
)

// GetVM - GET /api2/json/nodes/{node}/qemu/{vmid}
// ! need to re-think about how we can get specific vm's info
// func GetVM(hostURL string, user model.User) {
// 	// Set Header
// 	log.Println(user.CSRFPreventionToken)

// 	// Construct request
// 	resp, err := http.Get(hostURL)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// We Read the response body on the line below.
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// Convert the body to type string
// 	sb := string(body)
// 	log.Print(sb)
// }

// GetVMs - GET /api2/json/nodes/{node}/qemu
func GetVMs(hostURL, username string, cookies model.Cookies) {
	//
	user := model.User{}
	database.DB.Db.Find(&user, "username = ?", username)
	log.Println(user)

	// Construct new request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, hostURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Getting cookie
	log.Println(cookies)
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

	// GET request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Fatalln("error: with status", resp.Status)
	}

	// We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Convert the body to type string
	sb := string(body)
	log.Print(sb)
	// // Unmarshal body to struct
	// if marshalErr := json.Unmarshal(body, &ticket); marshalErr != nil {
	// 	return ticket, marshalErr
	// }
}

// // CreateVM - POST /api2/json/nodes/{node}/qemu
// func CreateVM()

// // DeleteVM - DELETE /api2/json/nodes/{node}/qemu/{vmid}
// func DeleteVM()

// // CloneVM - POST /api2/json/nodes/{node}/qemu/{vmid}/clone
// func CloneVM()

// // CreateTemplate - POST /api2/json/nodes/{node}/qemu/{vmid}/template
// func CreateTemplate()
