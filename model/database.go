// Package model - structs
package model

// User - struct for user's info {student, faculty, admin}
type User struct {
	Username   string
	Password   string
	Name       string
	Status     bool
	CreateTime string
	ExpireTime string
}

// CreateUserDB - create user in DB's body
type CreateUserDB struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Name     string `form:"name"`
	Group    string `form:"group"`
}

// EditUserDB - edit user in DB's body
type EditUserDB struct {
	Password   string `form:"password"`
	Name       string `form:"name"`
	Status     bool   `form:"status"`
	ExpireTime string `form:"expire_time"`
}

// InstanceLimit - struct for instance limit
type InstanceLimit struct {
	Username    string
	MaxCPU      float64 // Amount of CPU limit
	MaxRAM      float64 // Amount of RAM limit in GiB
	MaxDisk     float64 // Amount of Disk limit in GiB
	MaxInstance uint64  // Amount of instance count limit
}

// EditInstanceLimit - struct for edit instance limit
type EditInstanceLimit struct {
	MaxCPU      float64 `form:"max_cpu"`
	MaxRAM      float64 `form:"max_ram"`
	MaxDisk     float64 `form:"max_disk"`
	MaxInstance uint64  `form:"max_instance"`
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
