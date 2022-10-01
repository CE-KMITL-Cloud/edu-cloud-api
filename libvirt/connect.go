package libvirt

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
		conn = ssh_connect(compute.Username, compute.Host)
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
/*
# command-line-arguments
ld: warning: -no_pie is deprecated when targeting new OS versions
(ce@10.20.20.100) Password:
2022/10/01 23:24:13 virError(Code=38, Domain=7, Message='End of file while reading data: nc: unix connect failed: No such file or directory
nc: /usr/local/var/run/libvirt/libvirt-sock: No such file or directory: Input/output error')
exit status 1
*/
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
