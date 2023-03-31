// Package model - structs
package model

// CloneBody - struct for request Cloning VM
type CloneBody struct {
	Name    string `json:"name"`
	Storage string `json:"storage"` // Storage name - {"ceph-vm, ceph-vm2 ..."}
	CIUser  string `json:"ciuser"`
	CIPass  string `json:"cipassword"`
}

// CreateBody - struct for request Creating VM
type CreateBody struct {
	Name    string  `json:"name"`
	Memory  uint64  `json:"memory"`
	Sockets uint64  `json:"sockets"`
	Cores   float64 `json:"cores"`
	Onboot  uint8   `json:"onboot"`  // {0, 1}
	Storage string  `json:"storage"` // Storage name - {"ceph-vm, ceph-vm2 ..."}
	Disk    string  `json:"disk"`    // Amount of Disk in GiB
	CDROM   string  `json:"cdrom"`
	Net0    string  `json:"net0"`
	SCSIHW  string  `json:"scsihw"`
}

// TemplateBody - struct for request Templating VM
type TemplateBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// DeleteBody - struct for request Deleting VM
type DeleteBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// EditBody - struct for request Editing VM configuration
type EditBody struct {
	Memory uint64  `json:"memory"`
	Cores  float64 `json:"cores"`
	Disk   uint64  `json:"disk"`
}

// StartBody - struct for request Starting VM
type StartBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// StopBody - struct for request Stopping VM
type StopBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// ShutdownBody - struct for request Shutting down VM
type ShutdownBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// SuspendBody - struct for request Suspending VM
type SuspendBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// ResumeBody - struct for request Resuming VM
type ResumeBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// ResetBody - struct for request Resetting VM
type ResetBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// RebootBody - struct for request Rebooting VM
type RebootBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}

// VncProxyBody - struct for request Get VNC Ticket
type VncProxyBody struct {
	VMID uint64 `json:"vmid"`
	Node string `json:"node"`
}
