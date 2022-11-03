package libvirt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"libvirt.org/go/libvirt"
)

func GetDomCapXML(conn *libvirt.Connect, arch string, machine string) string {
	emulatorBin := GetEmulator(conn, arch)
	virtType := "qemu"
	hypervisorDomainTypes := GetHypervisorsDomainType(conn)
	for _, val := range hypervisorDomainTypes {
		if val.Arch == arch {
			if val.DomainType == "kvm" {
				virtType = "kvm"
				break
			}
		}
	}
	machineTypes := GetMachineTypes(conn, arch)
	if machine == "" || !contains(machineTypes, machine) {
		if contains(machineTypes, "pc") {
			machine = "pc"
		} else {
			machine = machineTypes[0]
		}
	}
	domCap, err := conn.GetDomainCapabilities(emulatorBin, arch, machine, virtType, 0)
	check(err)
	WriteStringtoFile(domCap, fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine))

	// return machine incase machine in other func is ""
	// don't need to return domCap
	return machine
}

func GetCapXML(conn *libvirt.Connect) string {
	cap, err := conn.GetCapabilities()
	check(err)
	WriteStringtoFile(cap, "xml/capabilities.xml")
	return cap
}

// Host Capabilities for specified architecture
func GetCapabilities(conn *libvirt.Connect, arch string) ArchCapabilities {
	GetCapXML(conn)
	archWordSize, archWordSizeErr := GetXPath("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/wordsize", arch))
	check(archWordSizeErr)
	archEmu, archEmuErr := GetXPath("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/emulator", arch))
	check(archEmuErr)
	archDomains, archDomainsErr := GetXPathsAttr("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/domain", arch), "type")
	check(archDomainsErr)
	archMachines, archMachinesErr := GetXPaths("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine", arch))
	check(archMachinesErr)
	maxCpu, maxCpuErr := GetChildElementsAttr("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine", arch), archMachines, "maxCpus")
	check(maxCpuErr)
	canonical, canonicalErr := GetChildElementsAttr("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine", arch), archMachines, "canonical")
	check(canonicalErr)
	var machineDetail []MachineDetail
	for i := range archMachines {
		machineDetail = append(machineDetail, MachineDetail{
			Machine:   archMachines[i],
			MaxCPU:    maxCpu[i],
			Canonical: canonical[i],
		})
	}
	archFeatures, archFeaturesErr := GetParentTags("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']", arch), "features")
	check(archFeaturesErr)
	archOsType, archOsTypeErr := GetParentText("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']", arch), "os_type")
	check(archOsTypeErr)
	return ArchCapabilities{
		WordSize: archWordSize,
		Emulator: archEmu,
		Domains:  archDomains,
		Machines: machineDetail,
		Features: archFeatures,
		OsType:   archOsType,
	}
}

func GetDomainCapabilities(conn *libvirt.Connect, arch string, machine string) DomCapabilities {
	machine = GetDomCapXML(conn, arch, machine)
	var (
		loaders                []string
		loaderEnums            []OsLoaderEnum
		cpuCustomModel         []string
		diskDevice             []string
		diskBus                []string
		graphicsTypes          []string
		videoTypes             []string
		hostDevTypes           []string
		hostDevStartupPolicies []string
		hostDevSubSysTypes     []string
	)
	path, pathErr := GetXPath(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/path")
	check(pathErr)
	domain, domainErr := GetXPath(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/domain")
	check(domainErr)
	domMachine, domMachineErr := GetXPath(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/machine")
	check(domMachineErr)
	vcpu, vcpuErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/vcpu", "max")
	check(vcpuErr)
	ioThreads, ioThreadsErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/iothreads", "supported")
	check(ioThreadsErr)
	osSupport, osSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/os", "supported")
	check(osSupportErr)
	loaderSupport, loaderSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/os/loader", "supported")
	check(loaderSupportErr)
	if loaderSupport == "yes" {
		loaders = GetOsLoaders(conn, arch, machine)
		loaderEnums = GetOsLoaderEnums(conn, arch, machine)
	}
	cpuModes := GetCPUModes(conn, arch, machine)
	if contains(cpuModes, "custom") {
		// supported and unknown cpu models
		cpuCustomModel = GetCPUCustomTypes(conn, arch, machine)
	}
	diskSupport, diskSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/devices/disk", "supported")
	check(diskSupportErr)
	if diskSupport == "yes" {
		diskDevice = GetDiskDeviceTypes(conn, arch, machine)
		diskBus = GetDiskBusTypes(conn, arch, machine)
	}
	graphicsSupport, graphicsSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/devices/graphics", "supported")
	check(graphicsSupportErr)
	if graphicsSupport == "yes" {
		graphicsTypes = GetGraphicTypes(conn, arch, machine)
	}
	videoSupport, videoSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/devices/video", "supported")
	check(videoSupportErr)
	if videoSupport == "yes" {
		videoTypes = GetVideoModels(conn, arch, machine)
	}
	hostDevSupport, hostDevSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/devices/hostdev", "supported")
	check(hostDevSupportErr)
	if hostDevSupport == "yes" {
		hostDevTypes = GetHostDevModes(conn, arch, machine)
		hostDevStartupPolicies = GetHostDevStartupPolicies(conn, arch, machine)
		hostDevSubSysTypes = GetHostDevSubSysTypes(conn, arch, machine)
	}
	gicSupport, gicSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/features/gic", "supported")
	check(gicSupportErr)
	genIdSupport, genIdSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/features/genid", "supported")
	check(genIdSupportErr)
	vmCoreInfoSupport, vmCoreInfoSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/features/vmcoreinfo", "supported")
	check(vmCoreInfoSupportErr)
	sevSupport, sevSupportErr := GetXPathAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/features/sev", "supported")
	check(sevSupportErr)
	return DomCapabilities{
		Path:                   path,
		Domain:                 domain,
		Machine:                domMachine,
		VcpuMax:                vcpu,
		IoThreads:              ioThreads,
		OsSupport:              osSupport,
		LoaderSupport:          loaderSupport,
		Loader:                 loaders,
		LoaderEnums:            loaderEnums,
		CpuModes:               cpuModes,
		CpuCustomModels:        cpuCustomModel,
		DiskSupport:            diskSupport,
		DiskDevices:            diskDevice,
		DiskBus:                diskBus,
		GraphicsSupport:        graphicsSupport,
		GraphicsTypes:          graphicsTypes,
		VideoSupport:           videoSupport,
		VideoTypes:             videoTypes,
		HostDevSupport:         hostDevSupport,
		HostDevTypes:           hostDevTypes,
		HostDevStartupPolicies: hostDevStartupPolicies,
		HostDevSubSysTypes:     hostDevSubSysTypes,
		FeaturesGicSupport:     gicSupport,
		FeatureGenIdSupport:    genIdSupport,
		FeatureVMCoreInfo:      vmCoreInfoSupport,
		FeatureSevSupport:      sevSupport,
	}
}

// Running hypervisor: QEMU 4.2.1
func GetVersion(conn *libvirt.Connect) string {

	ver, verErr := conn.GetVersion()
	check(verErr)
	major := ver / 1000000
	ver %= 1000000
	minor := ver / 1000
	ver %= 1000
	release := ver
	return fmt.Sprintf("%d.%d.%d", major, minor, release)
}

// Using library: libvirt 6.0.0
func GetLibVersion(conn *libvirt.Connect) string {
	ver, verErr := conn.GetLibVersion()
	check(verErr)
	major := ver / 1000000
	ver %= 1000000
	minor := ver / 1000
	ver %= 1000
	release := ver
	return fmt.Sprintf("%d.%d.%d", major, minor, release)
}

func GetHypervisorsDomainType(conn *libvirt.Connect) []HostDomainType {
	GetCapXML(conn)
	arch, archErr := GetXPathsAttr("xml/capabilities.xml", "./capabilities/guest/arch", "name")
	check(archErr)
	var domainTypeList []HostDomainType
	for i := range arch {
		domainType, domainTypeErr := GetXPathsAttr("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/domain", arch[i]), "type")
		check(domainTypeErr)
		for j := range domainType {
			domainTypeList = append(domainTypeList, HostDomainType{Arch: arch[i], DomainType: domainType[j]})
		}
	}
	return domainTypeList
}

func GetHypervisorsMachines(conn *libvirt.Connect) []HostMachine {
	GetCapXML(conn)
	arch, archErr := GetXPathsAttr("xml/capabilities.xml", "./capabilities/guest/arch", "name")
	check(archErr)
	var hostMachines []HostMachine
	for i := range arch {
		machineType := GetMachineTypes(conn, arch[i])
		for j := range machineType {
			hostMachines = append(hostMachines, HostMachine{Arch: arch[i], Machine: machineType[j]})
		}
	}
	return hostMachines
}

func GetEmulator(conn *libvirt.Connect, arch string) string {
	GetCapXML(conn)
	emu, err := GetXPath("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/emulator", arch))
	check(err)
	return emu
}

func GetEmulators(conn *libvirt.Connect) []HostEmulator {
	GetCapXML(conn)
	arch, archErr := GetXPathsAttr("xml/capabilities.xml", "./capabilities/guest/arch", "name")
	check(archErr)
	emu, emuErr := GetXPaths("xml/capabilities.xml", "./capabilities/guest/arch[@name]/emulator")
	check(emuErr)
	result := make([]HostEmulator, len(emu))
	for i := range emu {
		result[i] = HostEmulator{Arch: arch[i], Emulator: emu[i]}
	}
	return result
}

func GetMachineTypes(conn *libvirt.Connect, arch string) []string {
	GetCapXML(conn)
	var blank []string
	canonicalName, canonicalNameErr := GetXPathsAttr("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine[@canonical]", arch), "canonical")
	check(canonicalNameErr)
	if len(canonicalName) > 0 {
		machineName, machineNameErr := GetXPaths("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine", arch))
		check(machineNameErr)
		return machineName
	}
	return blank
}

func GetOsLoaders(conn *libvirt.Connect, arch string, machine string) []string {
	if arch == "" {
		arch = "x86_64"
	}
	if machine == "" {
		machine = "pc"
	}
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/os/loader[@supported='yes']/value")
	check(valErr)
	return val
}

func GetOsLoaderEnums(conn *libvirt.Connect, arch string, machine string) []OsLoaderEnum {
	machine = GetDomCapXML(conn, arch, machine)
	enums, enumsErr := GetXPathsAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/os/loader[@supported='yes']/enum", "name")
	check(enumsErr)
	var (
		val    []string
		valErr error
		result []OsLoaderEnum
	)
	for i := range enums {
		val, valErr = GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), fmt.Sprintf("/domainCapabilities/os/loader[@supported='yes']/enum[@name='%s']/value", enums[i]))
		check(valErr)
		for j := range val {
			result = append(result, OsLoaderEnum{Enum: enums[i], Value: val[j]})
		}
	}
	return result
}

func GetDiskBusTypes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/disk/enum[@name='bus']/value")
	check(valErr)
	return val
}

func GetDiskDeviceTypes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/disk/enum[@name='diskDevice']/value")
	check(valErr)
	return val
}

func GetGraphicTypes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/graphics/enum[@name='type']/value")
	check(valErr)
	return val
}

func GetCPUModes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPathsAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/cpu/mode[@supported='yes']", "name")
	check(valErr)
	return val
}

func GetCPUCustomTypes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	var result []string
	usableYes, usableYesErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/cpu/mode[@name='custom'][@supported='yes']/model[@usable='yes']")
	check(usableYesErr)
	usableUnknown, usableUnknownErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/cpu/mode[@name='custom'][@supported='yes']/model[@usable='unknown']")
	check(usableUnknownErr)
	for i := range usableYes {
		result = append(result, usableYes[i])
	}
	for i := range usableUnknown {
		result = append(result, usableUnknown[i])
	}
	return result
}

func GetHostDevModes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/hostdev/enum[@name='mode']/value")
	check(valErr)
	return val
}

func GetHostDevStartupPolicies(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/hostdev/enum[@name='startupPolicy']/value")
	check(valErr)
	return val
}

func GetHostDevSubSysTypes(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/hostdev/enum[@name='subsysType']/value")
	check(valErr)
	return val
}

func GetVideoModels(conn *libvirt.Connect, arch string, machine string) []string {
	machine = GetDomCapXML(conn, arch, machine)
	videoEnumName, videoEnumNameErr := GetXPathsAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/video/enum", "name")
	check(videoEnumNameErr)
	videoEnum, videoEnumErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/video/enum[@name='modelType']/value")
	check(videoEnumErr)
	var result []string
	if videoEnumName[0] == "modelType" {
		result = append(result, videoEnum...)
	}
	return result
}

func GetInterface(conn *libvirt.Connect, name string) *libvirt.Interface {
	iface, err := conn.LookupInterfaceByName(name)
	check(err)
	return iface
}

func GetInterfaces(conn *libvirt.Connect) []string {
	var interfaces []string
	interfaceList, interfaceListErr := conn.ListInterfaces()
	check(interfaceListErr)
	definedInterfaceList, definedInterfaceListErr := conn.ListDefinedInterfaces()
	check(definedInterfaceListErr)
	for iface := range interfaceList {
		interfaces = append(interfaces, interfaceList[iface])
	}
	for iface := range definedInterfaceList {
		interfaces = append(interfaces, definedInterfaceList[iface])
	}
	return interfaces
}

func GetStorage(conn *libvirt.Connect, name string) *libvirt.StoragePool {
	storage, err := conn.LookupStoragePoolByName(name)
	check(err)
	return storage
}

func GetStorages(conn *libvirt.Connect, onlyActive bool) []string {
	var storages []string
	storagePool, storagePoolErr := conn.ListStoragePools()
	check(storagePoolErr)
	for pool := range storagePool {
		storages = append(storages, storagePool[pool])
	}
	if !onlyActive {
		definedStoragePool, definedStoragePoolErr := conn.ListDefinedStoragePools()
		check(definedStoragePoolErr)
		for pool := range definedStoragePool {
			storages = append(storages, definedStoragePool[pool])
		}
	}
	return storages
}

func GetVolumeByPath(conn *libvirt.Connect, path string) *libvirt.StorageVol {
	volume, err := conn.LookupStorageVolByPath(path)
	check(err)
	return volume
}

func GetNetwork(conn *libvirt.Connect, net string) *libvirt.Network {
	network, err := conn.LookupNetworkByName(net)
	check(err)
	return network
}

func GetNetworks(conn *libvirt.Connect) []string {
	var hostNets []string
	netList, netListErr := conn.ListNetworks()
	check(netListErr)
	definedNetList, definedNetListErr := conn.ListDefinedNetworks()
	check(definedNetListErr)
	for net := range netList {
		hostNets = append(hostNets, netList[net])
	}
	for net := range definedNetList {
		hostNets = append(hostNets, definedNetList[net])
	}
	return hostNets
}

func GetNetworkForward(conn *libvirt.Connect, net_name string) (string, error) {
	net_forward, err := GetNetwork(conn, net_name).GetXMLDesc(0)
	check(err)
	return net_forward, nil
}

func GetInstance(conn *libvirt.Connect, name string) *libvirt.Domain {
	instance, err := conn.LookupDomainByName(name)
	check(err)
	return instance
}

func GetInstances(conn *libvirt.Connect) []string {
	var instances []string
	domList, domListErr := conn.ListDomains()
	check(domListErr)
	for instID := range domList {
		dom, domErr := conn.LookupDomainById(uint32(domList[instID]))
		check(domErr)
		domName, domNameErr := dom.GetName()
		check(domNameErr)
		instances = append(instances, domName)
	}
	definedDomList, definedDomListErr := conn.ListDefinedDomains()
	check(definedDomListErr)
	for index := range definedDomList {
		instances = append(instances, definedDomList[index])
	}
	return instances
}

func GetSnapshots(conn *libvirt.Connect) []string {
	var instance []string
	domList, domListErr := conn.ListDomains()
	check(domListErr)
	for snapID := range domList {
		dom, domErr := conn.LookupDomainById(uint32(domList[snapID]))
		check(domErr)
		domSnapshotNum, domSnapshotNumErr := dom.SnapshotNum(0)
		check(domSnapshotNumErr)
		if domSnapshotNum != 0 {
			domName, domNameErr := dom.GetName()
			check(domNameErr)
			instance = append(instance, domName)
		}
	}
	definedDomList, definedDomListErr := conn.ListDefinedDomains()
	check(definedDomListErr)
	for index := range definedDomList {
		dom, domErr := conn.LookupDomainByName(definedDomList[index])
		check(domErr)
		domSnapshotNum, domSnapshotNumErr := dom.SnapshotNum(0)
		check(domSnapshotNumErr)
		if domSnapshotNum != 0 {
			domName, domNameErr := dom.GetName()
			check(domNameErr)
			instance = append(instance, domName)
		}
	}
	return instance
}

func GetNetDevices(conn *libvirt.Connect) []string {
	var netDevice []string
	deviceList, deviceListErr := conn.ListAllNodeDevices(0)
	for device_index := range deviceList {
		xml, xmlErr := deviceList[device_index].GetXMLDesc(0)
		check(xmlErr)
		WriteStringtoFile(xml, fmt.Sprintf("xml/device_list/device_%d.xml", device_index))
	}
	for index := range deviceList {
		devType, devTypeErr := GetXPath(fmt.Sprintf("xml/device_list/device_%d.xml", index), "./device/capability[@type='net']/interface")
		check(devTypeErr)
		iFace, iFaceErr := GetXPath(fmt.Sprintf("xml/device_list/device_%d.xml", index), "./device/capability/interface")
		check(iFaceErr)
		if devType != "" {
			netDevice = append(netDevice, iFace)
		}
	}
	check(deviceListErr)
	return netDevice
}

func GetHostInstances(conn *libvirt.Connect) []HostInstance {
	var vcpu string
	var vcpuErr error
	var mem int
	rawMemSize := false
	instances := GetInstances(conn)
	vname := make([]HostInstance, len(instances))
	for instIndex := range instances {
		dom := GetInstance(conn, instances[instIndex])
		domInfo, domInfoErr := dom.GetInfo()
		check(domInfoErr)
		domUUID, domUUIDErr := dom.GetUUIDString()
		check(domUUIDErr)
		xml, xmlErr := dom.GetXMLDesc(0)
		check(xmlErr)
		WriteStringtoFile(xml, "xml/host_instances.xml")
		memRaw, memRawErr := GetXPath("xml/host_instances.xml", "/domain/currentMemory")
		check(memRawErr)
		if rawMemSize {
			memInt, memIntErr := strconv.Atoi(memRaw) // 1024
			check(memIntErr)
			mem = memInt * (1024 * 1024)
		}
		currentVCpu, currentVCpuErr := GetXPath("xml/host_instances.xml", "/domain/vcpu/@current")
		check(currentVCpuErr)

		if currentVCpu != "" {
			vcpu = currentVCpu
		} else {
			vcpu, vcpuErr = GetXPath("xml/host_instances.xml", "/domain/vcpu")
			check(vcpuErr)
		}
		title, titleErr := GetXPath("xml/host_instances.xml", "/domain/title")
		check(titleErr)
		if title == "" {
			title = ""
		}
		description, descriptionErr := GetXPath("xml/host_instances.xml", "/domain/description")
		check(descriptionErr)
		if description == "" {
			description = ""
		}
		vname[instIndex] = HostInstance{
			Status:      domInfo.State,
			UUID:        domUUID,
			VCPU:        vcpu,
			Memory:      mem,
			Title:       title,
			Description: description,
		}
	}
	return vname
}

func GetUserInstances(conn *libvirt.Connect, name string) UserInstance {
	var vcpu string
	var vcpuErr error
	dom := GetInstance(conn, name)
	xml, xmlErr := dom.GetXMLDesc(0)
	check(xmlErr)
	WriteStringtoFile(xml, "xml/user_instances.xml")
	domName, domNameErr := dom.GetName()
	check(domNameErr)
	domInfo, domInfoErr := dom.GetInfo()
	check(domInfoErr)
	domUUID, domUUIDErr := dom.GetUUIDString()
	check(domUUIDErr)
	memRaw, memRawErr := GetXPath("xml/user_instances.xml", "/domain/currentMemory")
	check(memRawErr)
	mem, memErr := strconv.Atoi(memRaw) // 1024
	check(memErr)
	currentVCpu, currentVCpuErr := GetXPath("xml/user_instances.xml", "/domain/vcpu/@current")
	check(currentVCpuErr)
	if currentVCpu != "" {
		vcpu = currentVCpu
	} else {
		vcpu, vcpuErr = GetXPath("xml/user_instances.xml", "/domain/vcpu")
		check(vcpuErr)
	}
	title, titleErr := GetXPath("xml/user_instances.xml", "/domain/title")
	check(titleErr)
	if title == "" {
		title = ""
	}
	description, descriptionErr := GetXPath("xml/user_instances.xml", "/domain/description")
	check(descriptionErr)
	if description == "" {
		description = ""
	}
	return UserInstance{
		Name:        domName,
		Status:      domInfo.State,
		UUID:        domUUID,
		VCPU:        vcpu,
		Memory:      mem,
		Title:       title,
		Description: description,
	}
}

func ArchCanUEFI(arch string) bool {
	supportedArch := map[string]bool{
		"i686":    true,
		"x86_64":  true,
		"aarch64": true,
		"armv7l":  true,
	}
	return supportedArch[arch]
}

func IsQEMU(conn *libvirt.Connect) bool {
	uri, err := conn.GetURI()
	check(err)
	return strings.HasPrefix(uri, "qemu")
}

// Search the loader paths for one that matches the passed arch
func FindUEFIPathForArch(conn *libvirt.Connect, arch string, machine string) string {
	if !ArchCanUEFI(arch) {
		return ""
	}
	loaders := GetOsLoaders(conn, arch, machine)
	archPatterns := GetUEFIArchPatterns()
	var patterns []string
	for i := range archPatterns {
		if arch == archPatterns[i].Arch {
			patterns = append(patterns, archPatterns[i].UEFI...)
		}
	}
	for i := range patterns {
		for j := range loaders {
			match, _ := regexp.MatchString(patterns[i], loaders[j])
			if match {
				return loaders[j]
			}
		}
	}
	return ""
}

// Return a pretty label for passed path, based on if we know about it or no
func LabelForFirmwarePath(conn *libvirt.Connect, arch string, path string) string {
	archs := []string{"i686", "x86_64"}
	if path == "" {
		if contains(archs, arch) {
			return "BIOS"
		}
		return ""
	}
	archPatterns := GetUEFIArchPatterns()
	for i := range archPatterns {
		for j := range archPatterns[i].UEFI {
			match, _ := regexp.MatchString(archPatterns[i].UEFI[j], path)
			if match {
				return fmt.Sprintf("UEFI %s: %s", arch, path)
			}
		}
	}
	return fmt.Sprintf("Custom: %s", path)
}

// Return True if libvirt advertises support for proper UEFI setup
func SupportsUEFIXml(conn *libvirt.Connect, loader_enums []OsLoaderEnum) bool {
	hasReadonly := false
	hasYes := false
	for i := range loader_enums {
		if loader_enums[i].Enum == "readonly" {
			hasReadonly = true
			if loader_enums[i].Value == "yes" {
				hasYes = true
			}
		}
	}
	return (hasReadonly && hasYes)
}

func IsSupportsVirtio(conn *libvirt.Connect, arch string, machine string) bool {
	if !IsQEMU(conn) {
		return false
	}
	// These _only_ support virtio so don't check the OS
	archs := []string{"aarch64", "armv7l", "ppc64", "ppc64le", "s390x", "riscv64", "riscv32"}
	machines := []string{"virt", "pseries"}
	if contains(archs, arch) && contains(machines, machine) {
		return true
	}
	if contains([]string{"x86_64", "i686"}, arch) {
		return true
	}
	return false
}

func GetUEFIArchPatterns() []ArchUEFI {
	return []ArchUEFI{
		{
			Arch: "i686",
			UEFI: []string{
				`.*ovmf-ia32.*`, // fedora, gerd's firmware repo
			},
		},
		{
			Arch: "x86_64",
			UEFI: []string{
				`.*OVMF_CODE\.fd`,       // RHEL
				`.*ovmf-x64/OVMF.*\.fd`, // gerd's firmware repo
				`*ovmf-x86_64-.*`,       // SUSE
				`.*ovmf.*`,
				`.*OVMF.*`, // generic attempt at a catchall
			},
		},
		{
			Arch: "aarch64",
			UEFI: []string{
				`.*AAVMF_CODE\.fd`,     // RHEL
				`.*aarch64/QEMU_EFI.*`, // gerd's firmware repo
				`.*aarch64.*`,          // generic attempt at a catchall
			},
		},
		{
			Arch: "armv7l",
			UEFI: []string{
				`.*arm/QEMU_EFI.*`, // fedora, gerd's firmware repo
			},
		},
	}
}

// Get cache available modes
func GetCacheModes() CacheMode {
	return CacheMode{
		Default:      "Default",
		None:         "Disabled",
		WriteThrough: "Write through",
		WriteBack:    "Write back",
		DirectSync:   "Direct sync", //since libvirt 0.9.5
		Unsafe:       "Unsafe",      //since libvirt 0.9.7
	}
}

// return available io modes
func GetIOModes() IOMode {
	return IOMode{
		Default: "Default",
		Native:  "Native",
		Threads: "Threads",
	}
}

// return: available discard modes
func GetDiscardModes() DiscardMode {
	return DiscardMode{
		Default: "Default",
		Ignore:  "Ignore",
		Unmap:   "Unmap",
	}
}

// return: available detect zeroes modes
func GetDetectZeroModes() DetectZeroMode {
	return DetectZeroMode{
		Default: "Default",
		On:      "On",
		Off:     "Off",
		Unmap:   "Unmap",
	}
}

// return: network card models
func GetNetworkModels() NetworkModel {
	return NetworkModel{
		Default: "Default",
		E1000:   "e1000",
		Virtio:  "virtio",
	}
}

// return: available image formats
func GetImageFormats() ImageFormat {
	return ImageFormat{
		Raw:   "raw",
		Qcow:  "qcow",
		Qcow2: "qcow2",
	}
}

// return: available image filename extensions
func GetFileExtensions() FileExtension {
	return FileExtension{
		Img:   "img",
		Qcow:  "qcow",
		Qcow2: "qcow2",
	}
}

// * skip this feature
// func GetSecrets(conn *libvirt.Connect) []string {
// 	secrets, err := conn.ListSecrets()
// 	check(err)
// 	return secrets
// }

// * skip this feature
// func GetSecret(conn *libvirt.Connect, uuid string) *libvirt.Secret {
// 	secret, err := conn.LookupSecretByUUIDString(uuid)
// 	check(err)
// 	return secret
// }

// * skip this feature
// func GetNWFilter(conn *libvirt.Connect, name string) *libvirt.NWFilter {
// 	nwfilter, err := conn.LookupNWFilterByName(name)
// 	check(err)
// 	return nwfilter
// }
