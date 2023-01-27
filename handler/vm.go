// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetVM - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func GetVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	// username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// Getting VM's info
	info, err := internal.GetVM(urlStr, cookies)
	if err != nil {
		log.Println("Error: from getting VM's info :", err)
		return err
	}
	log.Printf("Got info from vmid : %s", vmid)
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// GetVMList - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu
/*
	using Query Params
	@username : account's username
	@node : node's name
*/
func GetVMList(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	// username := c.Query("username")
	node := c.Query("node")

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", node)
	urlStr := u.String()

	// Getting VM's info
	info, err := internal.GetVMList(urlStr, cookies)
	if err != nil {
		log.Println("Error: from getting VM's list :", err)
		return err
	}
	log.Printf("Got VM list from node : %s", node)
	return c.Status(200).JSON(fiber.Map{"status": "Success", "message": info})
}

// CreateVM - Create VM on specific node
// POST /api2/json/nodes/{node}/qemu
/*
	using Request's Body
	@vmid : VM's ID
	@name : VM's name
	@memory : e.g. 1024 (MB)
	@cores : e.g. 2 (cores)
	@sockets : e.g. 2 (sockets of cpu)
	@onboot : {0, 1}
	@scsi0 : "ceph-vm:32"
	@cdrom : "cephfs:iso/ubuntu-20.04.4-live-server-amd64.iso"
	@net0 : "virtio,bridge=vmbr0,firewall=1"
	@scsihw : "virtio-scsi-single"
*/
// TODO : Specific resource pool by add pool in request's body
func CreateVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's body
	createBody := new(model.CreateBody)
	if err := c.BodyParser(createBody); err != nil {
		return err
	}
	vmid := fmt.Sprint(createBody.VMID)

	// TODO: Another work around here?
	// Construct payload
	data := url.Values{}
	data.Set("vmid", vmid)
	data.Set("name", createBody.Name)
	data.Set("memory", fmt.Sprint(createBody.Memory))
	data.Set("cores", fmt.Sprint(createBody.Cores))
	data.Set("sockets", fmt.Sprint(createBody.Sockets))
	data.Set("onboot", fmt.Sprint(createBody.Onboot))
	data.Set("scsi0", createBody.SCSI0)
	data.Set("cdrom", createBody.CDROM)
	data.Set("net0", createBody.Net0)
	data.Set("scsihw", createBody.SCSIHW)

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Getting target node from node allocation
	workerNodes, target, nodeErr := internal.AllocateNode(cookies)
	if nodeErr != nil {
		return nodeErr
	}
	log.Printf("Create body : %s, target node : %s", data, target)

	// Check duplicate vmid
	for _, workerNode := range workerNodes {
		// Construct VM List URL
		vmListURL, _ := url.ParseRequestURI(hostURL)
		vmListURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", workerNode.Node) // all node
		vmListURLStr := vmListURL.String()

		vmList, vmListErr := internal.GetVMList(vmListURLStr, cookies)
		if vmListErr != nil {
			return vmListErr
		}

		// Checking duplicate VM by VMID
		var list []string
		for _, v := range vmList.Info {
			list = append(list, fmt.Sprintf("%d", v.VMID))
		}
		log.Printf("VMs in node : %s : %s", workerNode.Node, list)

		if config.Contains(list, string(vmid)) {
			log.Printf("Error: found duplicate VMID : %s in node : %s", vmid, workerNode.Node)
			return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Found duplicate VMID: %s in %s", vmid, workerNode.Node)})
		}
	}

	// Construct URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", target)
	urlStr := u.String()

	// Getting info
	info, err := internal.CreateVM(urlStr, data, cookies)
	if err != nil {
		log.Println("Error: from creating VM :", err)
		return err
	}
	log.Println(info)

	// Waiting until creating process has been complete
	created := internal.StatusVM(target, vmid, []string{"created", "starting", "running"}, cookies)
	if created {
		log.Printf("Finished creating VMID : %s in %s", vmid, target)
		return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Creating new VMID: %s in %s successfully", vmid, target)})
	}
	log.Printf("Error: Could not create VMID : %s in %s", vmid, target)
	return c.Status(500).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Creating new VMID: %s in %s has failed", vmid, target)})
}

// DeleteVM - Deleting specific VM
// DELETE /api2/json/nodes/{node}/qemu/{vmid}
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func DeleteVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's query params
	// username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Construct Getting info URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// First check that target VM has been stopped
	vm, err := internal.GetVM(urlStr, cookies)
	if err != nil {
		return err
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: deleting VMID : %s in %s due to VM has not been stopped", vmid, node)
		return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, node)})
	}

	// Construct Deleting API URL
	deleteURL, _ := url.ParseRequestURI(hostURL)
	deleteURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", node, vmid)
	deleteURLStr := deleteURL.String()

	// Deleting target VM
	_, deleteErr := internal.DeleteVM(deleteURLStr, cookies)
	if deleteErr != nil {
		log.Printf("Error: deleting VMID : %s in %s due to %s", vmid, node, deleteErr)
		return deleteErr
	}

	// Check that target VM has been deleted completely yet
	deleted := internal.DeleteCompletely(node, vmid, cookies)
	if deleted {
		log.Printf("Finished deleting VMID : %s in %s", vmid, node)
		return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been deleted", vmid, node)})
	}
	log.Printf("Error: Could not delete VMID : %s in %s", vmid, node)
	return c.Status(500).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been deleted", vmid, node)})
}

// CloneVM - Cloning specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/clone
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID

	using Request's Body
	@newid : new VM's ID
	@name : VM's name
*/
func CloneVM(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting request's body
	cloneBody := new(model.CloneBody)
	if err := c.BodyParser(cloneBody); err != nil {
		return err
	}
	newid := fmt.Sprint(cloneBody.NewID)

	// Getting data from query & Mapping values
	node := c.Query("node")
	vmid := c.Query("vmid")

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Check VM Template from vmid
	isTemplate := internal.IsTemplate(node, vmid, cookies)
	if isTemplate {
		// Getting target node from node allocation
		workerNodes, target, nodeErr := internal.AllocateNode(cookies)
		if nodeErr != nil {
			log.Println("Error: allocate node :", nodeErr)
			return nodeErr
		}

		for _, workerNode := range workerNodes {
			// Construct VM List URL
			vmListURL, _ := url.ParseRequestURI(hostURL)
			vmListURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu", workerNode.Node) // all node
			vmListURLStr := vmListURL.String()

			vmList, vmListErr := internal.GetVMList(vmListURLStr, cookies)
			if vmListErr != nil {
				return vmListErr
			}

			// Checking duplicate VM by VMID
			var list []string
			for _, v := range vmList.Info {
				list = append(list, fmt.Sprintf("%d", v.VMID))
			}
			log.Printf("VMs in node : %s : %s", workerNode.Node, list)

			if config.Contains(list, string(newid)) {
				log.Printf("Error: found duplicate VMID : %s in node : %s", newid, workerNode.Node)
				return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Found duplicate VMID: %s in %s", newid, workerNode.Node)})
			}
		}

		// Construct payload
		data := url.Values{}
		data.Set("newid", newid)
		data.Set("name", cloneBody.Name)
		data.Set("target", target)
		log.Println("clone body :", data)

		// Construct URL
		u, _ := url.ParseRequestURI(hostURL)
		u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/clone", node, vmid)
		urlStr := u.String()

		// Getting info
		info, err := internal.CloneVM(urlStr, data, cookies)
		if err != nil {
			log.Printf("Error: cloning VMID : %s in %s : %s", newid, target, err)
			return err
		}
		log.Println(info)

		// Waiting until cloning process has been completed
		cloned := internal.StatusVM(target, newid, []string{"created", "starting", "running"}, cookies)
		if cloned {
			log.Printf("Finished cloning VMID : %s in %s", newid, target)
			return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Cloning new VMID: %s to %s successfully", newid, target)})
		}
		log.Printf("Error: cloning VMID : %s in %s has failed", newid, target)
		return c.Status(500).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Cloning new VMID: %s to %s has failed", newid, target)})
	}
	log.Printf("Error: cloning VMID : %s from VMID : %s", newid, vmid)
	return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": "Could not clone VM"})
}

// CreateTemplate - Templating specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/template
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID
*/
func CreateTemplate(c *fiber.Ctx) error {
	// Get host's URL
	hostURL := config.GetFromENV("PROXMOX_HOST")

	// Getting data from query & Mapping values
	node := c.Query("node")
	vmid := c.Query("vmid")

	// Getting Cookie, CSRF Token
	cookies := model.Cookies{
		Cookie: http.Cookie{
			Name:  "PVEAuthCookie",
			Value: c.Cookies("PVEAuthCookie"),
		},
		CSRFPreventionToken: fiber.Cookie{
			Name:  "CSRFPreventionToken",
			Value: c.Cookies("CSRFPreventionToken"),
		},
	}

	// Construct Getting info URL
	u, _ := url.ParseRequestURI(hostURL)
	u.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid)
	urlStr := u.String()

	// First check that target VM has been stopped
	vm, err := internal.GetVM(urlStr, cookies)
	if err != nil {
		return err
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: Could not template VMID : %s in %s due to VM hasn't been stopped", vmid, node)
		return c.Status(400).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, node)})
	}

	// Construct URL
	templateURL, _ := url.ParseRequestURI(hostURL)
	templateURL.Path = fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/template", node, vmid)
	templateURLStr := templateURL.String()

	// Templating VM
	_, templateErr := internal.CreateTemplate(templateURLStr, cookies)
	if templateErr != nil {
		log.Printf("Error: Could not template VMID : %s in %s : %s", vmid, node, templateErr)
		return templateErr
	}

	// Waiting until templating process has been completed
	templated := internal.TemplateCompletely(node, vmid, []string{"created", "existing"}, cookies)
	if templated {
		log.Printf("Finished templating VMID : %s in %s", vmid, node)
		return c.Status(200).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been templated", vmid, node)})
	}
	log.Printf("Error: Could not template VMID : %s in %s", vmid, node)
	return c.Status(500).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been templated correctly", vmid, node)})
}
