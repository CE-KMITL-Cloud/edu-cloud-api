// Package model - structs
package model

// VM - struct for VM
type VM struct {
	Info VMInfo `json:"data"`
}

// VMInfo - struct for VM's info
type VMInfo struct {
	CPU            float64 `json:"cpu"`
	NetOut         uint64  `json:"netout"`
	DiskWrite      uint64  `json:"diskwrite"`
	UpTime         uint64  `json:"uptime"`
	MaxMem         uint64  `json:"maxmem"`
	PID            uint32  `json:"pid"`
	Name           string  `json:"name"`
	VMID           uint16  `json:"vmid"`
	MaxDisk        uint64  `json:"maxdisk"`
	Status         string  `json:"status"`
	Disk           uint16  `json:"disk"`
	DiskRead       uint64  `json:"diskread"`
	CPUs           uint16  `json:"cpus"`
	Mem            uint64  `json:"mem"`
	NetIn          uint64  `json:"netin"`
	Lock           string  `json:"lock"`
	QmpStatus      string  `json:"qmpstatus"`
	HA             HA      `json:"ha"`
	FreeMem        uint64  `json:"freemem"`
	Balloon        uint64  `json:"balloon"`
	RunningQEMU    string  `json:"running-qemu"`
	RunningMachine string  `json:"running-machine"`
}

// HA - struct for HA object use in VMInfo
type HA struct {
	Manage int `json:"managed"`
}

// VMList - struct for VM List
type VMList struct {
	Info []VMListInfo `json:"data"`
}

// VMListInfo - struct for VM List's info
type VMListInfo struct {
	CPU       float64 `json:"cpu"`
	NetOut    uint64  `json:"netout"`
	DiskWrite uint64  `json:"diskwrite"`
	UpTime    uint64  `json:"uptime"`
	MaxMem    uint64  `json:"maxmem"`
	PID       uint32  `json:"pid"`
	Name      string  `json:"name"`
	VMID      uint16  `json:"vmid"`
	MaxDisk   uint64  `json:"maxdisk"`
	Status    string  `json:"status"`
	Disk      uint16  `json:"disk"`
	DiskRead  uint64  `json:"diskread"`
	CPUs      uint16  `json:"cpus"`
	Mem       uint64  `json:"mem"`
	NetIn     uint64  `json:"netin"`
}

// VMResponse - struct for VM's response
type VMResponse struct {
	Info string `json:"data"`
}

// CloneBody - struct for request Cloning VM
type CloneBody struct {
	NewID uint32 `form:"newid"`
	Name  string `form:"name"`
}

// CreateBody - struct for request Creating VM
type CreateBody struct {
	VMID    uint32 `form:"vmid"`
	Name    string `form:"name"`
	Memory  uint32 `form:"memory"`
	Sockets uint32 `form:"sockets"`
	Cores   uint32 `form:"cores"`
	Onboot  uint8  `form:"onboot"` // {0, 1}
	SCSI0   string `form:"scsi0"`
	CDROM   string `form:"cdrom"`
	Net0    string `form:"net0"`
	SCSIHW  string `form:"scsihw"`
}

// VMTemplate - struct for VM Template
type VMTemplate struct {
	Info VMTemplateInfo `json:"data"`
}

// VMTemplateInfo - VM Template detail info
type VMTemplateInfo struct {
	Template  uint8   `json:"template"` // {0, 1}
	CPU       float64 `json:"cpu"`
	NetOut    uint64  `json:"netout"`
	DiskWrite uint64  `json:"diskwrite"`
	UpTime    uint64  `json:"uptime"`
	MaxMem    uint64  `json:"maxmem"`
	Name      string  `json:"name"`
	VMID      uint16  `json:"vmid"`
	MaxDisk   uint64  `json:"maxdisk"`
	Status    string  `json:"status"`
	Disk      uint16  `json:"disk"`
	DiskRead  uint64  `json:"diskread"`
	CPUs      uint16  `json:"cpus"`
	Mem       uint64  `json:"mem"`
	NetIn     uint64  `json:"netin"`
	QmpStatus string  `json:"qmpstatus"`
	HA        HA      `json:"ha"`
}
