// Package model - structs
package model

// NodeResource - struct of node's resources
type NodeResource struct {
	Nodes []Node `json:"data"`
}

// Node - struct of node's resources detail
type Node struct {
	ID      string  `json:"id"`
	Node    string  `json:"node"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
	MaxDisk uint64  `json:"maxdisk"`
	MaxCPU  float64 `json:"maxcpu"`
	MaxMem  uint64  `json:"maxmem"`
	Disk    uint64  `json:"disk"`
	CPU     float64 `json:"cpu"`
	Mem     uint64  `json:"mem"`
	UpTime  uint64  `json:"uptime"`
}

// VMSpec - struct of VM's specification
type VMSpec struct {
	Memory uint64
	CPU    float64
	Disk   uint64
}

// StorageResource - struct of storage's resources
type StorageResource struct {
	Storages []Storage `json:"data"`
}

// Storage - struct of storage's resources detail
type Storage struct {
	ID         string `json:"id"`
	Node       string `json:"node"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	MaxDisk    uint64 `json:"maxdisk"`
	Disk       uint64 `json:"disk"`
	Storage    string `json:"storage"`
	Content    string `json:"content"`
	Shared     uint8  `json:"shared"`
	PluginType string `json:"plugintype"`
}
