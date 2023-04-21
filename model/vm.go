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
	PID            uint64  `json:"pid"`
	Name           string  `json:"name"`
	VMID           uint64  `json:"vmid"`
	MaxDisk        uint64  `json:"maxdisk"`
	Status         string  `json:"status"`
	Disk           uint64  `json:"disk"`
	DiskRead       uint64  `json:"diskread"`
	CPUs           float64 `json:"cpus"`
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
	PID       uint64  `json:"pid"`
	Name      string  `json:"name"`
	VMID      uint64  `json:"vmid"`
	MaxDisk   uint64  `json:"maxdisk"`
	Status    string  `json:"status"`
	Disk      uint64  `json:"disk"`
	DiskRead  uint64  `json:"diskread"`
	CPUs      float64 `json:"cpus"`
	Mem       uint64  `json:"mem"`
	NetIn     uint64  `json:"netin"`
}

// VMResponse - struct for VM's response
type VMResponse struct {
	Info string `json:"data"`
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
	VMID      uint64  `json:"vmid"`
	MaxDisk   uint64  `json:"maxdisk"`
	Status    string  `json:"status"`
	Disk      uint64  `json:"disk"`
	DiskRead  uint64  `json:"diskread"`
	CPUs      float64 `json:"cpus"`
	Mem       uint64  `json:"mem"`
	NetIn     uint64  `json:"netin"`
	QmpStatus string  `json:"qmpstatus"`
	HA        HA      `json:"ha"`
}

// VMConfig - struct for VM config
type VMConfig struct {
	Info ConfigDetail `json:"data"`
}

// ConfigDetail - struct for VM config detail
type ConfigDetail struct {
	USB0         string `json:"usb0"`
	IPConfig0    string `json:"ipconfig0"`
	NameServer   string `json:"nameserver"`
	SCSIHW       string `json:"scsihw"`
	Meta         string `json:"meta"`
	SCSI0        string `json:"scsi0"`
	SearchDomain string `json:"searchdomain"`
	VMGenID      string `json:"vmgenid"`
	OSType       string `json:"ostype"`
	Tags         string `json:"tags"`
	BootDisk     string `json:"bootdisk"`
	VGA          string `json:"vga"`
	Net0         string `json:"net0"`
	Boot         string `json:"boot"`
	SMBIOS1      string `json:"smbios1"`
	Digest       string `json:"digest"`
	Numa         uint64 `json:"numa"`
	IDE2         string `json:"ide2"`
	Serial0      string `json:"serial0"`
	Agent        string `json:"agent"`
}

// VncProxy - struct for VNC Proxy
type VncProxy struct {
	Detail VncProxyResponse `json:"data"`
}

// VncProxyResponse - struct for VNC Proxy response
type VncProxyResponse struct {
	Ticket string `json:"ticket"`
	Port   string `json:"port"`
	Url    string `json:"url"`
}
