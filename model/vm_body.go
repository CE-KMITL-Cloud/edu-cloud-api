// Package model - structs
package model

// CloneBody - struct for request Cloning VM
type CloneBody struct {
	Name    string `form:"name"`
	Storage string `form:"storage"` // Storage name - {"ceph-vm, ceph-vm2 ..."}
	CIUser  string `form:"ciuser"`
	CIPass  string `form:"cipassword"`
}

// CreateBody - struct for request Creating VM
type CreateBody struct {
	Name    string  `form:"name"`
	Memory  uint64  `form:"memory"`
	Sockets uint64  `form:"sockets"`
	Cores   float64 `form:"cores"`
	Onboot  uint8   `form:"onboot"`  // {0, 1}
	Storage string  `form:"storage"` // Storage name - {"ceph-vm, ceph-vm2 ..."}
	Disk    string  `form:"disk"`    // Amount of Disk in GiB
	CDROM   string  `form:"cdrom"`
	Net0    string  `form:"net0"`
	SCSIHW  string  `form:"scsihw"`
}

// TemplateBody - struct for request Templating VM
type TemplateBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// DeleteBody - struct for request Deleting VM
type DeleteBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// EditBody - struct for request Editing VM configuration
type EditBody struct {
	Memory uint64  `form:"memory"`
	Cores  float64 `form:"cores"`
	Disk   uint64  `form:"disk"`
}

// StartBody - struct for request Starting VM
type StartBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// StopBody - struct for request Stopping VM
type StopBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// ShutdownBody - struct for request Shutting down VM
type ShutdownBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// SuspendBody - struct for request Suspending VM
type SuspendBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// ResumeBody - struct for request Resuming VM
type ResumeBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// ResetBody - struct for request Resetting VM
type ResetBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// RebootBody - struct for request Rebooting VM
type RebootBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}

// VncProxyBody - struct for request Get VNC Ticket
type VncProxyBody struct {
	VMID uint64 `form:"vmid"`
	Node string `form:"node"`
}
