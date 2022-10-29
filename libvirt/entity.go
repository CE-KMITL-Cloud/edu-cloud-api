package libvirt

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"libvirt.org/go/libvirt"
)

func GetDomCapXML(conn *libvirt.Connect, arch string, machine string) string {
	emulatorBin := GetEmulator(conn, arch)
	virtType := "qemu"
	hypervisorDomainTypes := GetHypervisorsDomainType(conn)
	for i := range hypervisorDomainTypes {
		if hypervisorDomainTypes[i].Arch == arch {
			if hypervisorDomainTypes[i].DomainType == "kvm" {
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
	return domCap
}

func GetCapXML(conn *libvirt.Connect) string {
	cap, err := conn.GetCapabilities()
	check(err)
	WriteStringtoFile(cap, "xml/capabilities.xml")
	return cap
}

// Host Capabilities for specified architecture
// func GetCapabilities(conn *libvirt.Connect, arch string) string {
// 	GetCapXML(conn)
// 	archElement, archElementErr := GetXPath("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']", arch))
// 	check(archElementErr)
// 	return archElement
// }

// func GetDomainCapabilities()

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
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "./domainCapabilities/os/loader[@supported='yes']/value")
	check(valErr)
	return val
}

func GetOsLoaderEnums(conn *libvirt.Connect, arch string, machine string) []OsLoaderEnum {
	GetDomCapXML(conn, arch, machine)
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
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/disk/enum[@name='bus']/value")
	check(valErr)
	return val
}

func GetDiskDeviceTypes(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/disk/enum[@name='diskDevice']/value")
	check(valErr)
	return val
}

func GetGraphicTypes(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/graphics/enum[@name='type']/value")
	check(valErr)
	return val
}

func GetCPUModes(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPathsAttr(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/cpu/mode[@supported='yes']", "name")
	check(valErr)
	return val
}

func GetCPUCustomTypes(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
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
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/hostdev/enum[@name='mode']/value")
	check(valErr)
	return val
}

func GetHostDevStartupPolicies(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/hostdev/enum[@name='startupPolicy']/value")
	check(valErr)
	return val
}

func GetHostDevSubSysTypes(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
	val, valErr := GetXPaths(fmt.Sprintf("xml/dom_cap_%s_%s.xml", arch, machine), "/domainCapabilities/devices/hostdev/enum[@name='subsysType']/value")
	check(valErr)
	return val
}

func GetVideoModels(conn *libvirt.Connect, arch string, machine string) []string {
	GetDomCapXML(conn, arch, machine)
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
	log.Println("loaders :", loaders)
	var patterns []string
	if arch == "i686" {
		patterns = UEFIArchPatterns().i686
	} else if arch == "x86_64" {
		patterns = UEFIArchPatterns().x86_64
	} else if arch == "aarch64" {
		patterns = UEFIArchPatterns().aarch64
	} else if arch == "armv7l" {
		patterns = UEFIArchPatterns().armv7l
	}
	log.Println("patterns :", patterns)
	for i := range patterns {
		for j := range loaders {
			match, _ := regexp.MatchString(patterns[i], loaders[j])
			if match {
				log.Println("loaders[j] :", loaders[j])
				return loaders[j]
			}
		}
	}
	return ""
}

// Return a pretty label for passed path, based on if we know about it or no
// func LabelForFirmwarePath(conn *libvirt.Connect, arch string, path string) string {
// 	var archs[]{"i686", "x86_64"}
// 	if path == "" {
// 		if contains(archs, arch) {
// 			return "BIOS"
// 		}
// 		return ""
// 	}

// }

// func SupportsUEFIXml()

// func IsSupportsVirtio()

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
