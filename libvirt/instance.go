package libvirt

import (
	"fmt"
	"strconv"

	"libvirt.org/go/libvirt"
)

// CreateInstance - For only admin role
/*
TODO:
 - volumes
*/
func CreateInstance(conn *libvirt.Connect, name, memory, vcpu, vcpuMode, uuid, arch, machine, volumes, networks, nwfilter, graphics, virtio, listenerAddr, video, consolePass, mac, qemuGa, addCDRom, addInput string, firmware OsLoaderEnum) string {
	video = defaultValue(video, "vga")
	consolePass = defaultValue(consolePass, "random")
	mac = defaultValue(mac, "None")
	qemuGa = defaultValue(qemuGa, "True")
	addCDRom = defaultValue(addCDRom, "sata")
	addInput = defaultValue(addInput, "default")

	caps := GetCapabilities(conn, arch)
	domCaps := GetDomainCapabilities(conn, arch, machine)
	mem, memErr := strconv.Atoi(memory) // 1024
	check(memErr)
	mem = mem * 1024

	xml := fmt.Sprintf(
		`<domain type='%s'>
		<name>%s</name>
		<description>None</description>
		<uuid>%s</uuid>
		<memory unit='KiB'>%d</memory>
		<vcpu>%s</vcpu>`, domCaps.Domain, name, uuid, mem, vcpu)

	if domCaps.OsSupport == "yes" {
		xml += fmt.Sprintf(`
		<os>
		<type arch='%s' machine='%s'>%s</type>`, arch, machine, caps.OsType)
		xml += `<boot dev='hd'/>
		<boot dev='cdrom'/>
		<bootmenu enable='yes'/>`
		// if firmware != "" {
		// 	if firmware.
		// }
	}
	return xml
}
