package libvirt

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

func GetDomCapXML(conn *libvirt.Connect, arch string, machine string) {
	/*
	   """ Return domain capabilities xml"""

	   	emulatorbin = self.get_emulator(arch)
	   	virttype = "kvm" if "kvm" in self.get_hypervisors_domain_types()[arch] else "qemu"

	   	machine_types = self.get_machine_types(arch)
	   	if not machine or machine not in machine_types:
	   	    machine = "pc" if "pc" in machine_types else machine_types[0]
	   	return self.wvm.getDomainCapabilities(emulatorbin, arch, machine, virttype)
	*/

	// emulatorbin := GetEmulator(arch)
	// var virtType string
}

func GetCapXML(conn *libvirt.Connect) string {
	cap, err := conn.GetCapabilities()
	check(err)
	WriteStringtoFile(cap, "xml/capabilities.xml")
	return cap
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
	canonicalName, canonicalNameErr := GetXPaths("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine[@canonical]", arch))
	check(canonicalNameErr)
	log.Println(canonicalName)
	if len(canonicalName) > 0 {
		machineName, machineNameErr := GetXPaths("xml/capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/machine", arch))
		check(machineNameErr)
		return machineName
	}
	return blank
}

func GetIFace(conn *libvirt.Connect, name string) *libvirt.Interface {
	iface, err := conn.LookupInterfaceByName(name)
	check(err)
	return iface
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

func GetStorage(conn *libvirt.Connect, name string) *libvirt.StoragePool {
	storage, err := conn.LookupStoragePoolByName(name)
	check(err)
	return storage
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

func GetNetworkForward(conn *libvirt.Connect, net_name string) (string, error) {
	net_forward, err := GetNetwork(conn, net_name).GetXMLDesc(0)
	check(err)
	return net_forward, nil
}

// * skip this feature
// func GetNWFilter(conn *libvirt.Connect, name string) *libvirt.NWFilter {
// 	nwfilter, err := conn.LookupNWFilterByName(name)
// 	check(err)
// 	return nwfilter
// }

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
