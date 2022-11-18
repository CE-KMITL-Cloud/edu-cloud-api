package libvirt

import (
	"libvirt.org/go/libvirt"
)

// ArchCapabilities -
type ArchCapabilities struct {
	WordSize string
	Emulator string
	Domains  []string
	Machines []MachineDetail
	Features []string
	OsType   string
}

// DomCapabilities -
type DomCapabilities struct {
	Path                   string
	Domain                 string
	Machine                string
	VcpuMax                string
	IoThreads              string
	OsSupport              string
	LoaderSupport          string
	Loader                 []string
	LoaderEnums            []OsLoaderEnum
	CPUModes               []string
	CPUCustomModels        []string
	DiskSupport            string
	DiskDevices            []string
	DiskBus                []string
	GraphicsSupport        string
	GraphicsTypes          []string
	VideoSupport           string
	VideoTypes             []string
	HostDevSupport         string
	HostDevTypes           []string
	HostDevStartupPolicies []string
	HostDevSubSysTypes     []string
	FeaturesGicSupport     string
	FeatureGenIDSupport    string
	FeatureVMCoreInfo      string
	FeatureSevSupport      string
}

// MachineDetail -
type MachineDetail struct {
	Machine   string
	MaxCPU    string
	Canonical string
}

// HostInstance -
type HostInstance struct {
	Name        string
	Status      libvirt.DomainState
	UUID        string
	VCPU        string
	Memory      int
	Title       string
	Description string
}

// UserInstance -
type UserInstance struct {
	Name        string
	Status      libvirt.DomainState
	UUID        string
	VCPU        string
	Memory      int
	Title       string
	Description string
}

// HostEmulator -
type HostEmulator struct {
	Arch     string
	Emulator string
}

// HostDomainType -
type HostDomainType struct {
	Arch       string
	DomainType string
}

// HostMachine -
type HostMachine struct {
	Arch    string
	Machine string
}

// OsLoaderEnum -
type OsLoaderEnum struct {
	Enum  string
	Value string
}

// ArchUEFI -
type ArchUEFI struct {
	Arch string
	UEFI []string
}

// CacheMode -
type CacheMode struct {
	Default      string
	None         string
	WriteThrough string
	WriteBack    string
	DirectSync   string
	Unsafe       string
}

// IOMode -
type IOMode struct {
	Default string
	Native  string
	Threads string
}

// DiscardMode -
type DiscardMode struct {
	Default string
	Ignore  string
	Unmap   string
}

// DetectZeroMode -
type DetectZeroMode struct {
	Default string
	On      string
	Off     string
	Unmap   string
}

// NetworkModel -
type NetworkModel struct {
	Default string
	E1000   string
	Virtio  string
}

// ImageFormat -
type ImageFormat struct {
	Raw   string
	Qcow  string
	Qcow2 string
}

// FileExtension -
type FileExtension struct {
	Img   string
	Qcow  string
	Qcow2 string
}
