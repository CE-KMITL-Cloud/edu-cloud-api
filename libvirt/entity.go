package libvirt

import (
	"fmt"
	"strconv"

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
}

func GetCapXML(conn *libvirt.Connect) string {
	cap, err := conn.GetCapabilities()
	check(err)
	WriteStringtoFile(cap, "capabilities.xml")
	return cap
}

func GetEmulator(arch string) string {
	emu, err := GetXPath("capabilities.xml", fmt.Sprintf("./capabilities/guest/arch[@name='%s']/emulator", arch))
	check(err)
	return emu
}

func GetIFace(conn *libvirt.Connect, name string) *libvirt.Interface {
	iface, err := conn.LookupInterfaceByName(name)
	check(err)
	return iface
}

func GetSecrets(conn *libvirt.Connect) []string {
	secrets, err := conn.ListSecrets()
	check(err)
	return secrets
}

func GetSecret(conn *libvirt.Connect, uuid string) *libvirt.Secret {
	secret, err := conn.LookupSecretByUUIDString(uuid)
	check(err)
	return secret
}

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

func GetNWFilter(conn *libvirt.Connect, name string) *libvirt.NWFilter {
	nwfilter, err := conn.LookupNWFilterByName(name)
	check(err)
	return nwfilter
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
		WriteStringtoFile(xml, "device_list.xml")
		devType, devTypeErr := GetXPath("device_list.xml", "/device/capability/@type")
		check(devTypeErr)
		iFace, iFaceErr := GetXPath("device_list.xml", "/device/capability/interface")
		check(iFaceErr)
		if devType == "net" {
			netDevice = append(netDevice, iFace)
		}
	}
	check(deviceListErr)
	return netDevice
}

func GetHostInstances(conn *libvirt.Connect) []HostInstance {
	var vname []HostInstance
	var vcpu string
	var vcpuErr error
	var mem int
	rawMemSize := false
	instances := GetInstances(conn)
	for instIndex := range instances {
		dom := GetInstance(conn, instances[instIndex])
		domID, domIDErr := dom.GetID()
		check(domIDErr)
		domName, domNameErr := dom.GetName()
		check(domNameErr)
		domInfo, domInfoErr := dom.GetInfo()
		check(domInfoErr)
		domUUID, domUUIDErr := dom.GetUUIDString()
		check(domUUIDErr)
		xml, xmlErr := dom.GetXMLDesc(0)
		check(xmlErr)
		WriteStringtoFile(xml, "host_instances.xml")
		memRaw, memRawErr := GetXPath("host_instances.xml", "/domain/currentMemory")
		check(memRawErr)
		if rawMemSize {
			memInt, memIntErr := strconv.Atoi(memRaw)
			check(memIntErr)
			mem = memInt * (1024 * 1024)
		}
		currentVCpu, currentVCpuErr := GetXPath("host_instances.xml", "/domain/vcpu/@current")
		check(currentVCpuErr)

		if currentVCpu != "" {
			vcpu = currentVCpu
		} else {
			vcpu, vcpuErr = GetXPath("host_instances.xml", "/domain/vcpu")
			check(vcpuErr)
		}
		title, titleErr := GetXPath("host_instances.xml", "/domain/title")
		check(titleErr)
		if title == "" {
			title = ""
		}
		description, descriptionErr := GetXPath("host_instances.xml", "/domain/description")
		check(descriptionErr)
		if description == "" {
			description = ""
		}
		vname[domID] = HostInstance{
			Name:        domName,
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

func ArchCanUEFI(arch string) bool {
	supportedArch := map[string]bool{
		"i686":    true,
		"x86_64":  true,
		"aarch64": true,
		"armv7l":  true,
	}
	return supportedArch[arch]
}
