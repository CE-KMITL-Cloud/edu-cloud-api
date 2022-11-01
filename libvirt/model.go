package libvirt

import (
	"libvirt.org/go/libvirt"
)

type ArchCapabilities struct {
	WordSize string
	Emulator string
	Domains  []string
	Machines []MachineDetail
	Features []string
	OsType   string
}

// ! Not finalize yet
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
	CpuModes               []string
	CpuCustomModels        []string
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
	FeatureGenIdSupport    string
	FeatureVMCoreInfo      string
	FeatureSevSupport      string
}

type MachineDetail struct {
	Machine   string
	MaxCPU    string
	Canonical string
}

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
