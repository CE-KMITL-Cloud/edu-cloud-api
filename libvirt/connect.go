package libvirt_connection

import (
	"log"

	"libvirt.org/go/libvirt"
)

func TCP_Connect(AUTHNAME string, PASSPHASE string) *libvirt.Connect {
	/*
		Reference link to see how function use
		https://github.com/libvirt/libvirt-go/blob/master/integration_test.go
	*/

	callback := func(creds []*libvirt.ConnectCredential) {
		for _, cred := range creds {
			if cred.Type == libvirt.CRED_AUTHNAME {
				cred.Result = AUTHNAME
				cred.ResultLen = len(cred.Result)
			} else if cred.Type == libvirt.CRED_PASSPHRASE {
				cred.Result = PASSPHASE
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

	// TODO: We have many URI, Store as array might be better
	uri := "qemu+tls://captain-2.ce.kmitl.cloud/system"

	conn, err := libvirt.NewConnectWithAuth(uri, auth, 0)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
