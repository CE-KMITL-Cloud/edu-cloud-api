// Package cluster - Cluster functions
package cluster

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/model"
)

// AllocateNode - allocate which node is the best choice to have interaction with (e.g. cloning, creating)
// GET /api2/json/cluster/resources
func AllocateNode(spec model.VMSpec, storage string, cookies model.Cookies) ([]model.Node, string, error) {
	log.Println("Getting nodes from cluster's resources ...")
	url := config.GetURL("/api2/json/cluster/resources")
	nodeResource := model.NodeResource{}
	storageResource := model.StorageResource{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []model.Node{}, "", err
	}
	if marshalErr := json.Unmarshal(body, &nodeResource); marshalErr != nil {
		return []model.Node{}, "", marshalErr
	}
	if marshalStorageErr := json.Unmarshal(body, &storageResource); marshalStorageErr != nil {
		return []model.Node{}, "", marshalStorageErr
	}
	var selectedStorage model.Storage
	for i := 0; i < len(storageResource.Storages); i++ {
		if storageResource.Storages[i].Type == "storage" && storageResource.Storages[i].Storage == storage && storageResource.Storages[i].PluginType == "rbd" {
			selectedStorage = storageResource.Storages[i]
		}
	}
	maxFreeDisk := selectedStorage.MaxDisk - selectedStorage.Disk
	log.Printf("storage: %s, free disk: %d", selectedStorage.Storage, maxFreeDisk)
	// Regex and return only worker nodes
	var nodeList []model.Node
	for j := 0; j < len(nodeResource.Nodes); j++ {
		r, _ := regexp.Compile(config.WorkerNode) // match node which start with work-{number}
		if nodeResource.Nodes[j].Type == "node" && r.MatchString(nodeResource.Nodes[j].Node) && nodeResource.Nodes[j].Status != "offline" {
			nodeList = append(nodeList, nodeResource.Nodes[j])
		}
	}
	// Compare all of them which node is the best {mem, cpu}
	var selectedNode model.Node
	var maxFreeMemory, maxFreeMemoryPercent uint64
	var maxFreeCPU, maxFreeCPUPercent float64 = -math.MaxFloat64, -math.MaxFloat64
	for _, node := range nodeList {
		freeMemory, freeCPU := node.MaxMem-node.Mem, node.MaxCPU-node.CPU
		freeMemoryPercent := (freeMemory / node.MaxMem) * 100
		freeCPUPercent := (freeCPU / node.MaxCPU) * 100
		// log.Printf("node: %s, free mem: %d, free cpu: %f", node.Node, freeMemory, freeCPU)
		if freeMemoryPercent > maxFreeMemoryPercent || (freeMemoryPercent == maxFreeMemoryPercent && freeCPUPercent > maxFreeCPUPercent) {
			maxFreeMemoryPercent, maxFreeCPUPercent = freeMemoryPercent, freeCPUPercent
			maxFreeMemory, maxFreeCPU = freeMemory, freeCPU
			selectedNode = node
		}
	}
	log.Printf("selected node: %s, free mem: %d, free cpu: %f, free disk: %d", selectedNode.Node, maxFreeMemory/config.Gigabyte, maxFreeCPU, maxFreeDisk/config.Gigabyte)
	log.Printf("vm spec mem: %d, cpu: %f, disk: %d", spec.Memory/config.Gigabyte, spec.CPU, spec.Disk/config.Gigabyte)
	if maxFreeMemory > spec.Memory && maxFreeCPU > spec.CPU && maxFreeDisk > spec.Disk {
		log.Printf("Return selected node : %s", selectedNode.Node)
		return nodeList, selectedNode.Node, nil
	}
	log.Printf("Selected node : %s have no enough free space", selectedNode.Node)
	return nodeList, selectedNode.Node, errors.New("error: Node have no enough free space")
}

// GetStorageList - Getting RBD storage list
// GET /api2/json/cluster/resources
func GetStorageList(cookies model.Cookies) ([]string, error) {
	log.Println("Getting storages from cluster's resources ...")
	url := config.GetURL("/api2/json/cluster/resources")
	storageResource := model.StorageResource{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []string{}, err
	}
	if marshalErr := json.Unmarshal(body, &storageResource); marshalErr != nil {
		return []string{}, marshalErr
	}
	// Filter recources to get only RBD storage
	var storages []string
	for i := 0; i < len(storageResource.Storages); i++ {
		if storageResource.Storages[i].Type == "storage" && storageResource.Storages[i].PluginType == "rbd" {
			if !config.Contains(storages, storageResource.Storages[i].Storage) {
				storages = append(storages, storageResource.Storages[i].Storage)
			}
		}
	}
	log.Println(storages)
	return storages, nil
}

// GetISOList - Getting ISO file list
// GET /api2/json/nodes/{node}/storage/{storage}/content
func GetISOList(cookies model.Cookies) ([]string, error) {
	log.Println("Getting ISO file list from cluster's resources ...")
	url := config.GetURL("/api2/json/nodes/ops1/storage/cephfs/content")
	storageContent := model.ISOList{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []string{}, err
	}
	if marshalErr := json.Unmarshal(body, &storageContent); marshalErr != nil {
		return []string{}, marshalErr
	}
	var ISOList []string
	for _, iso := range storageContent.ISOList {
		ISOList = append(ISOList, strings.TrimPrefix(iso.Volid, config.ISO))
	}
	log.Println(ISOList)
	return ISOList, nil
}

// GetNodes - Getting nodes
// GET /api2/json/cluster/resources
func GetNodes(cookies model.Cookies) ([]model.Node, error) {
	log.Println("Getting node information from given node ...")
	url := config.GetURL("/api2/json/cluster/resources")
	nodeResource := model.NodeResource{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return []model.Node{}, err
	}
	if marshalErr := json.Unmarshal(body, &nodeResource); marshalErr != nil {
		return []model.Node{}, marshalErr
	}
	// Regex and return only worker nodes
	var nodeList []model.Node
	for i := 0; i < len(nodeResource.Nodes); i++ {
		r, _ := regexp.Compile(config.WorkerNode) // match node which start with work-{number}
		if nodeResource.Nodes[i].Type == "node" && r.MatchString(nodeResource.Nodes[i].Node) {
			nodeList = append(nodeList, nodeResource.Nodes[i])
		}
	}
	return nodeList, nil
}

// GetNode - Getting node information from given name
// GET /api2/json/cluster/resources
func GetNode(name string, cookies model.Cookies) (model.Node, error) {
	log.Println("Getting node information from given node ...")
	url := config.GetURL("/api2/json/cluster/resources")
	nodeResource := model.NodeResource{}
	body, err := config.SendRequestWithErr(http.MethodGet, url, nil, cookies)
	if err != nil {
		return model.Node{}, err
	}
	if marshalErr := json.Unmarshal(body, &nodeResource); marshalErr != nil {
		return model.Node{}, marshalErr
	}
	// Regex and return only worker nodes
	var nodeList []model.Node
	for i := 0; i < len(nodeResource.Nodes); i++ {
		r, _ := regexp.Compile(config.WorkerNode) // match node which start with work-{number}
		if nodeResource.Nodes[i].Type == "node" && r.MatchString(nodeResource.Nodes[i].Node) {
			nodeList = append(nodeList, nodeResource.Nodes[i])
		}
	}
	matchNode := model.Node{}
	for _, node := range nodeList {
		if node.Node == name {
			matchNode = node
		}
	}
	return matchNode, nil
}
