// Package model - structs
package model

// Node - struct of node's resources
type Node struct {
	Resources []Resources `json:"data"`
}

// Resources - struct of node's resources detail
type Resources struct {
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
