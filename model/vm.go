// Package model - structs
package model

// VMs - struct for retriving which node we want to list VMs
type VMs struct {
	Node string `form:"node"`
}

// VM - struct for retriving which VM & node we want to see info
type VM struct {
	Node string `form:"node"`
	Qemu string `form:"qemu"`
}
