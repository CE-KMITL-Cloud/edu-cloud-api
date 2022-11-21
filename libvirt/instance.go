package libvirt

import (
	"fmt"
	"strconv"

	"libvirt.org/go/libvirt"
)

// CreateInstance - For only admin role
// func CreateInstance(conn *libvirt.Connect, name, memory, vcpu, vcpuMode, uuid, arch, machine, volumes, networks, nwfilter, graphics, virtio, listenerAddr, video, consolePass, mac, qemuGa, addCDRom, addInput string, firmware OsLoaderEnum) string {
// 	video = defaultValue(video, "vga")
// 	consolePass = defaultValue(consolePass, "random")
// 	mac = defaultValue(mac, "None")
// 	qemuGa = defaultValue(qemuGa, "True")
// 	addCDRom = defaultValue(addCDRom, "sata")
// 	addInput = defaultValue(addInput, "default")

// 	caps := GetCapabilities(conn, arch)
// 	domCaps := GetDomainCapabilities(conn, arch, machine)
// 	mem, memErr := strconv.Atoi(memory) // 1024
// 	check(memErr)
// 	mem = mem * 1024

// 	xml := fmt.Sprintf(
// 		`<domain type='%s'>
// 		<name>%s</name>
// 		<description>None</description>
// 		<uuid>%s</uuid>
// 		<memory unit='KiB'>%d</memory>
// 		<vcpu>%s</vcpu>`, domCaps.Domain, name, uuid, mem, vcpu)

// 	loaderEnum := GetOsLoaderEnums(conn, arch, machine)

// 	if domCaps.OsSupport == "yes" {
// 		xml += fmt.Sprintf(`
// 		<os>
// 		<type arch='%s' machine='%s'>%s</type>`, arch, machine, caps.OsType)
// 		xml += `<boot dev='hd'/>
// 		<boot dev='cdrom'/>
// 		<bootmenu enable='yes'/>`
// 		if len(loaderEnum) > 0 {
// 			for index := range loaderEnum {
// 				if loaderEnum[index].Enum == "secure" && loaderEnum[index].Value == "yes" {
// 				}
// 			}
// 		}
// 		xml += `
// 		</os>`
// 	}
// 	return xml
// }

// GetInstanceXML - Getting XML from specific instance's name
func GetInstanceXML(conn *libvirt.Connect, name string) string {
	inst := GetInstance(conn, name)
	xml, err := inst.GetXMLDesc(0)
	check(err)
	WriteStringtoFile(xml, fmt.Sprintf("xml/instance/%s_inst.xml", name))
	return xml
}

// GetInstanceStatus - Getting Instance status
func GetInstanceStatus(conn *libvirt.Connect, name string) libvirt.DomainState {
	inst := GetInstance(conn, name)
	info, err := inst.GetInfo()
	check(err)
	return info.State
}

// GetInstanceMemory - Getting Instance memory
func GetInstanceMemory(conn *libvirt.Connect, name string) int {
	GetInstanceXML(conn, name)
	memory, err := GetXPath(fmt.Sprintf("xml/instance/%s_inst.xml", name), "/domain/currentMemory")
	check(err)
	memInt, memIntErr := strconv.Atoi(memory) // 1024
	check(memIntErr)
	return memInt // 1024
}

// GetInstanceVCPU - Getting Instance VCPU
func GetInstanceVCPU(conn *libvirt.Connect, name string) string {
	GetInstanceXML(conn, name)
	vcpu, err := GetXPath(fmt.Sprintf("xml/instance/%s_inst.xml", name), "/domain/vcpu/@current")
	check(err)
	if vcpu == "" {
		staticVCPU, e := GetXPath(fmt.Sprintf("xml/instance/%s_inst.xml", name), "/domain/vcpu")
		check(e)
		vcpu = staticVCPU
	}
	return vcpu
}

// GetInstanceUUID - Getting UUID as string type
func GetInstanceUUID(conn *libvirt.Connect, name string) string {
	inst := GetInstance(conn, name)
	uuid, err := inst.GetUUIDString()
	check(err)
	return uuid
}

// StartInstance - If the call succeeds the domain moves from the defined to the running domains pools.
func StartInstance(conn *libvirt.Connect, name string) {
	dom := GetInstance(conn, name)
	err := dom.Create()
	check(err)
}

// ShutdownInstance - Shutdown a domain, the domain object is still usable thereafter, but the domain OS is being stopped
func ShutdownInstance(conn *libvirt.Connect, name string) {
	dom := GetInstance(conn, name)
	err := dom.Shutdown()
	check(err)
}
