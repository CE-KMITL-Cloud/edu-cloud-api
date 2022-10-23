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

type UEFIArch struct {
	i686    []string
	x86_64  []string
	aarch64 []string
	armv7l  []string
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
