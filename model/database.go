// Package model - structs
package model

// User - struct for user's info {student, faculty, admin}
type User struct {
	Username string
	Password string
	Name     string
	// TelMobile  string
	Status     bool
	CreateTime string
	ExpireTime string
}

// InstanceLimit - struct for instance limit
type InstanceLimit struct {
	Username    string
	MaxCPU      float64 // Amount of CPU limit
	MaxRAM      float64 // Amount of RAM limit in GiB
	MaxDisk     float64 // Amount of Disk limit in GiB
	MaxInstance uint64  // Amount of instance count limit
}

// Instance - struct for instance's info
type Instance struct {
	VMID       string `gorm:"column:vmid"`
	OwnerID    string `gorm:"column:ownerid"`
	Node       string
	Name       string
	IsTemplate bool
	MaxCPU     float64 // Amount of CPU limit
	MaxRAM     float64 // Amount of RAM limit in GiB
	MaxDisk    float64 // Amount of Disk limit in GiB
	CreateTime string
	ExpireTime string
}

// InstanceBody - struct for instance's request body
type InstanceBody struct {
	Name       string
	IsTemplate bool
	MaxCPU     float64 // Amount of CPU limit
	MaxRAM     float64 // Amount of RAM limit in GiB
	MaxDisk    float64 // Amount of Disk limit in GiB
	ExpireTime string
}

// Sizing - struct for instance's template
type Sizing struct {
	VMID       string `gorm:"column:vmid"`
	Node       string
	Name       string
	MaxCPU     float64 // Amount of CPU limit
	MaxRAM     float64 // Amount of RAM limit in GiB
	MaxDisk    float64 // Amount of Disk limit in GiB
	CreateTime string
}
