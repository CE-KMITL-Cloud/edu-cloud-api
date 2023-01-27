// Package internal - internal functions
package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// use as regex to filter only worker node
const workerNodeExp = `work-[-]?\d[\d,]*[\.]?[\d{2}]*`

// AllocateNode - allocate which node is the best choice to have interaction with (e.g. cloning, creating)
// GET /api2/json/cluster/resources
func AllocateNode(cookies model.Cookies) (string, error) {
	// Getting all nodes in cluster
	log.Println("Getting cluster's resources ...")

	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = "/api2/json/cluster/resources"
	urlStr := u.String()

	// check status of the vm
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		log.Println(err)
	}

	// Getting cookie
	req.AddCookie(&cookies.Cookie)
	req.Header.Add("CSRFPreventionToken", cookies.CSRFPreventionToken.Value)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// If not 200 OK then log error
	if resp.StatusCode != 200 {
		log.Println("error: with status", resp.Status)
		return "", errors.New(resp.Status)
	}

	// Parsing response
	node := model.Node{}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}
	// log.Println(string(body))

	// Unmarshal body to struct
	if marshalErr := json.Unmarshal(body, &node); marshalErr != nil {
		log.Println(marshalErr)
	}

	var nodeList []model.Resources

	// Regex and return only worker nodes
	for i := 0; i < len(node.Resources); i++ {
		r, _ := regexp.Compile(workerNodeExp) // match node which start with work-{number}
		if node.Resources[i].Type == "node" && r.MatchString(node.Resources[i].Node) {
			nodeList = append(nodeList, node.Resources[i])
		}
	}
	// log.Println(nodeList)

	// Compare all of them which node is the best {mem, cpu}
	var selectedNode model.Resources
	var maxFreeMemoryPercent, maxFreeCPUPercent float64 = -math.MaxFloat64, -math.MaxFloat64
	for _, node := range nodeList {
		freeMemoryPercent := (float64(node.MaxMem-node.Mem) / float64(node.MaxMem)) * 100
		freeCPUPercent := (float64(node.MaxCPU-node.CPU) / float64(node.MaxCPU)) * 100
		log.Printf("node: %s, free mem: %f, free cpu: %f", node.Node, freeMemoryPercent, freeCPUPercent)
		if freeMemoryPercent > maxFreeMemoryPercent || (freeMemoryPercent == maxFreeMemoryPercent && freeCPUPercent > maxFreeCPUPercent) {
			maxFreeMemoryPercent = freeMemoryPercent
			maxFreeCPUPercent = freeCPUPercent
			selectedNode = node
		}
	}

	// return the best node's name
	log.Printf("Return selected node : %s", selectedNode.Node)
	return selectedNode.Node, nil
}
