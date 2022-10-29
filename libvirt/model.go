package libvirt

import (
	"libvirt.org/go/libvirt"
)

type HostInstance struct {
	Name        string
	Status      libvirt.DomainState
	UUID        string
	VCPU        string
	Memory      int
	Title       string
	Description string
}

type UserInstance struct {
	Name        string
	Status      libvirt.DomainState
	UUID        string
	VCPU        string
	Memory      int
	Title       string
	Description string
}

type HostEmulator struct {
	Arch     string
	Emulator string
}

type HostDomainType struct {
	Arch       string
	DomainType string
}

type HostMachine struct {
	Arch    string
	Machine string
}

type OsLoaderEnum struct {
	Enum  string
	Value string
}

type ArchUEFI struct {
	Arch string
	UEFI []string
}

type CacheMode struct {
	Default      string
	None         string
	WriteThrough string
	WriteBack    string
	DirectSync   string
	Unsafe       string
}

type IOMode struct {
	Default string
	Native  string
	Threads string
}

type DiscardMode struct {
	Default string
	Ignore  string
	Unmap   string
}

type DetectZeroMode struct {
	Default string
	On      string
	Off     string
	Unmap   string
}

type NetworkModel struct {
	Default string
	E1000   string
	Virtio  string
}

type ImageFormat struct {
	Raw   string
	Qcow  string
	Qcow2 string
}

type FileExtension struct {
	Img   string
	Qcow  string
	Qcow2 string
}
