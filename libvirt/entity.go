package libvirt

import (
	"libvirt.org/go/libvirt"
)

func Get_iface(conn *libvirt.Connect, name string) *libvirt.Interface {
	iface, err := conn.LookupInterfaceByName(name)
	check(err)
	return iface
}

func Get_secrets(conn *libvirt.Connect) []string {
	secrets, err := conn.ListSecrets()
	check(err)
	return secrets
}

func Get_secret(conn *libvirt.Connect, uuid string) *libvirt.Secret {
	secret, err := conn.LookupSecretByUUIDString(uuid)
	check(err)
	return secret
}

func Get_storage(conn *libvirt.Connect, name string) *libvirt.StoragePool {
	storage, err := conn.LookupStoragePoolByName(name)
	check(err)
	return storage
}

func Get_volume_by_path(conn *libvirt.Connect, path string) *libvirt.StorageVol {
	volume, err := conn.LookupStorageVolByPath(path)
	check(err)
	return volume
}

func Get_network(conn *libvirt.Connect, net string) *libvirt.Network {
	network, err := conn.LookupNetworkByName(net)
	check(err)
	return network
}

// TODO : have a look on util.py # https://github.com/retspen/webvirtcloud/blob/master/vrtManager/util.py
// func get_network_forward(conn *libvirt.Connect, net_name string) {}

func Get_nwfilter(conn *libvirt.Connect, name string) *libvirt.NWFilter {
	nwfilter, err := conn.LookupNWFilterByName(name)
	check(err)
	return nwfilter
}

func Get_instance(conn *libvirt.Connect, name string) *libvirt.Domain {
	instance, err := conn.LookupDomainByName(name)
	check(err)
	return instance
}

// TODO : what is lookById func?
// func get_instances(conn *libvirt.Connect) {}

func Get_cap_xml(conn *libvirt.Connect) string {
	cap, err := conn.GetCapabilities()
	check(err)
	return cap
}

func get_emulator(conn *libvirt.Connect, arch string) {
	// return util.get_xml_path(self.get_cap_xml(), "/capabilities/guest/arch[@name='{}']/emulator".format(arch))
}

func Get_dom_cap_xml(conn *libvirt.Connect, arch string, machine string) {
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
