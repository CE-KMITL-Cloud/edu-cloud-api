package libvirt_connection

import (
	"fmt"
	"log"

	"libvirt.org/go/libvirt"
)

type Connection struct {
	Host      string // e.g. "qemu+tls://domain/system, qemu+libssh2://user@host/system?known_hosts=/home/user/.ssh/known_hosts"
	Username  string
	Passwd    string
	Conn_type string // 'ssh' or 'tls'
}

// TODO : implement singleton pattern & and Mutex lock
// func (c *Connection) Init(host string, username string, passwd string, conn_type string) {
// 	c.Host = host
// 	c.Username = username
// 	c.Passwd = passwd
// 	c.Conn_type = conn_type
// }

func CreateCompute(compute Connection) *libvirt.Connect {
	var conn *libvirt.Connect

	if compute.Conn_type == "ssh" {
		conn = ssh_connect(compute.Host, compute.Passwd)
	} else if compute.Conn_type == "tls" {
		conn = tls_connect(compute.Username, compute.Passwd, compute.Host)
	} else if compute.Conn_type == "socket" {
		conn = socket_connect()
	} else {
		log.Fatal("Invalid connection type")
	}
	return conn
}

func socket_connect() *libvirt.Connect {
	uri := "qemu:///system"

	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

// TODO : need connection testing
func ssh_connect(username string, host string) *libvirt.Connect {
	// uri := "qemu+libssh2://user@host/system?known_hosts=/home/user/.ssh/known_hosts"
	uri := fmt.Sprintf("qemu+ssh://%s@%s/system", username, host)

	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func tls_connect(auth_name string, passphase string, host string) *libvirt.Connect {
	/*
		Reference link to see how function use
		https://github.com/libvirt/libvirt-go/blob/master/integration_test.go
	*/

	callback := func(creds []*libvirt.ConnectCredential) {
		for _, cred := range creds {
			if cred.Type == libvirt.CRED_AUTHNAME {
				cred.Result = auth_name
				cred.ResultLen = len(cred.Result)
			} else if cred.Type == libvirt.CRED_PASSPHRASE {
				cred.Result = passphase
				cred.ResultLen = len(cred.Result)
			}
		}
	}
	auth := &libvirt.ConnectAuth{
		CredType: []libvirt.ConnectCredentialType{
			libvirt.CRED_AUTHNAME, libvirt.CRED_PASSPHRASE,
		},
		Callback: callback,
	}

	// uri = "qemu+tls://captain-2.ce.kmitl.cloud/system"
	uri := fmt.Sprintf("qemu+tls://%s/system", host)

	conn, err := libvirt.NewConnectWithAuth(uri, auth, 0)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func get_iface(conn *libvirt.Connect, name string) *libvirt.Interface {
	iface, err := conn.LookupInterfaceByName(name)
	if err != nil {
		log.Fatal(err)
	}
	return iface
}

func get_secrets(conn *libvirt.Connect) []string {
	secrets, err := conn.ListSecrets()
	if err != nil {
		log.Fatal(err)
	}
	return secrets
}

func get_secret(conn *libvirt.Connect, uuid string) *libvirt.Secret {
	secret, err := conn.LookupSecretByUUIDString(uuid)
	if err != nil {
		log.Fatal(err)
	}
	return secret
}

func get_storage(conn *libvirt.Connect, name string) *libvirt.StoragePool {
	storage, err := conn.LookupStoragePoolByName(name)
	if err != nil {
		log.Fatal(err)
	}
	return storage
}

func get_volume_by_path(conn *libvirt.Connect, path string) *libvirt.StorageVol {
	volume, err := conn.LookupStorageVolByPath(path)
	if err != nil {
		log.Fatal(err)
	}
	return volume
}

func get_network(conn *libvirt.Connect, net string) *libvirt.Network {
	network, err := conn.LookupNetworkByName(net)
	if err != nil {
		log.Fatal(err)
	}
	return network
}

// TODO : have a look on util.py # https://github.com/retspen/webvirtcloud/blob/master/vrtManager/util.py
// func get_network_forward(conn *libvirt.Connect, net_name string) {}

func get_nwfilter(conn *libvirt.Connect, name string) *libvirt.NWFilter {
	nwfilter, err := conn.LookupNWFilterByName(name)
	if err != nil {
		log.Fatal(err)
	}
	return nwfilter
}

func get_instance(conn *libvirt.Connect, name string) *libvirt.Domain {
	instance, err := conn.LookupDomainByName(name)
	if err != nil {
		log.Fatal(err)
	}
	return instance
}

// TODO : what is lookById func?
// func get_instances(conn *libvirt.Connect) {}
