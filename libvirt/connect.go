package libvirt_connection

import (
	"fmt"
	"log"

	"libvirt.org/go/libvirt"
)

/*
TODO
* create struct for connection function
* 2 connection types -> SSH, TLS
* health check
*/

type Connection struct {
	Host      string // e.g. "qemu+tls://domain/system, qemu+libssh2://user@host/system?known_hosts=/home/user/.ssh/known_hosts"
	Username  string
	Passwd    string
	Conn_type string // 'ssh' or 'tls'
}

func CreateCompute(compute Connection) *libvirt.Connect {
	var conn *libvirt.Connect
	if compute.Conn_type == "ssh" {
		conn = ssh_connect(compute.Host, compute.Passwd)
	} else if compute.Conn_type == "tls" {
		conn = tls_connect(compute.Username, compute.Passwd, compute.Host)
	} else {
		log.Fatal("Invalid connection type")
	}
	return conn
}

// TODO : need connection testing
func ssh_connect(host string, passwd string) *libvirt.Connect {
	/*
		Reference link to see how function use
		https://github.com/libvirt/libvirt-go/blob/master/integration_test.go
	*/
	callback := func(creds []*libvirt.ConnectCredential) {
		for _, cred := range creds {
			if cred.Type == libvirt.CRED_AUTHNAME {
				cred.Result = host
				cred.ResultLen = len(cred.Result)
			} else if cred.Type == libvirt.CRED_PASSPHRASE {
				cred.Result = passwd
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

	// uri := "qemu+libssh2://user@host/system?known_hosts=/home/user/.ssh/known_hosts"
	uri := fmt.Sprintf("qemu+libssh2://%s@%s/system", host, passwd)

	conn, err := libvirt.NewConnectWithAuth(uri, auth, 0)
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
