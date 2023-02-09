// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/internal/cluster"
	"github.com/edu-cloud-api/internal/qemu"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetVM - Getting specific VM's info
// GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
/*
	using Query Params
	? @username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func GetVM(c *fiber.Ctx) error {
	// ? username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")

	cookies := config.GetCookies(c)

	// Getting VM's info
	log.Printf("Getting detail from VMID : %s in %s", vmid, node)
	url := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
	info, err := qemu.GetVM(url, cookies)
	if err != nil {
		log.Println("Error: from getting VM's info :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, node, err)})
	}
	log.Printf("Got info from vmid : %s", vmid)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": info})
}

// GetVMList - Getting VM list from given node
// GET /api2/json/nodes/{node}/qemu
/*
	using Query Params
	? @username : account's username
	@node : node's name
*/
func GetVMList(c *fiber.Ctx) error {
	// ? username := c.Query("username")
	node := c.Query("node")

	cookies := config.GetCookies(c)

	// Getting VM list
	url := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu", node))
	log.Printf("Getting VM list from %s", node)
	vmList, err := qemu.GetVMList(url, cookies)
	if err != nil {
		log.Println("Error: from getting VM's list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list from %s due to %s", node, err)})
	}
	log.Printf("Got VM list from node : %s", node)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": vmList})
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
	createBody := new(model.CreateBody)
	if err := c.BodyParser(createBody); err != nil {
		log.Println("Error: Could not parse body parser to create VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to create VM's body"})
	}
	vmid := fmt.Sprint(createBody.VMID)

	// Construct payload
	data := url.Values{}
	data.Set("vmid", vmid)
	data.Set("name", createBody.Name)
	data.Set("memory", fmt.Sprint(createBody.Memory)) // memory (MB)
	data.Set("cores", fmt.Sprint(createBody.Cores))   // cpu (core)
	data.Set("sockets", fmt.Sprint(createBody.Sockets))
	data.Set("onboot", fmt.Sprint(createBody.Onboot))
	data.Set("scsi0", createBody.SCSI0) // "ceph-vm:32" have to use 32 (GB) as disk
	data.Set("cdrom", createBody.CDROM)
	data.Set("net0", createBody.Net0)
	data.Set("scsihw", createBody.SCSIHW)

	r := regexp.MustCompile(config.MatchNumber)
	maxDiskStr := r.FindStringSubmatch(createBody.SCSI0)[1]
	maxDisk, parseErr := strconv.ParseUint(maxDiskStr, 10, 64)
	if parseErr != nil {
		log.Println("Error: extract max disk :", parseErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to extract max disk field for creating VM due to %s", parseErr)})
	}

	// Parse mem, cpu, disk for checking free space
	vmSpec := model.VMSpec{
		Memory: config.MBtoByte(createBody.Memory),
		CPU:    createBody.Cores,
		Disk:   config.GBtoByte(maxDisk),
	}

	cookies := config.GetCookies(c)

	// Getting target node from node allocation
	workerNodes, target, nodeErr := cluster.AllocateNode(vmSpec, cookies)
	if nodeErr != nil {
		log.Println("Error: allocate node :", nodeErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to allocate node for creating VM due to %s", nodeErr)})
	}
	log.Printf("Create body : %s, target node : %s", data, target)

	// Check duplicate vmid
	for _, workerNode := range workerNodes {
		vmListURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu", workerNode.Node))
		vmList, vmListErr := qemu.GetVMList(vmListURL, cookies)
		if vmListErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list from %s due to %s", workerNode.Node, vmListErr)})
		}
		var list []string
		for _, v := range vmList.Info {
			list = append(list, fmt.Sprintf("%d", v.VMID))
		}
		log.Printf("VMs in node : %s : %s", workerNode.Node, list)

		if config.Contains(list, string(vmid)) {
			log.Printf("Error: found duplicate VMID : %s in node : %s", vmid, workerNode.Node)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Found duplicate VMID: %s in %s", vmid, workerNode.Node)})
		}
	}

	// Creating VM
	vmCreateURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu", target))
	log.Printf("Creating VMID : %s in %s", vmid, target)
	info, err := qemu.CreateVM(vmCreateURL, data, cookies)
	if err != nil {
		log.Println("Error: from creating VM :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed creating VMID : %s in %s due to %s", vmid, target, err)})
	}
	log.Println(info)

	// Waiting until creating process has been complete
	created := qemu.CheckStatus(target, vmid, []string{"created", "starting", "running"}, true, (10 * time.Minute), time.Second, cookies)
	if created {
		log.Printf("Finished creating VMID : %s in %s", vmid, target)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Creating new VMID: %s in %s successfully", vmid, target)})
	}
	log.Printf("Error: Could not create VMID : %s in %s", vmid, target)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Creating new VMID: %s in %s has failed", vmid, target)})
}

// DeleteVM - Deleting specific VM
// DELETE /api2/json/nodes/{node}/qemu/{vmid}
/*
	using Request's Body
	@username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func DeleteVM(c *fiber.Ctx) error {
	// Getting request's body
	deleteBody := new(model.DeleteBody)
	if err := c.BodyParser(deleteBody); err != nil {
		log.Println("Error: Could not parse body parser to delete VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to delete VM's body"})
	}
	vmid := fmt.Sprint(deleteBody.VMID)

	cookies := config.GetCookies(c)

	// First check that target VM has been stopped
	vmGetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", deleteBody.Node, vmid))
	vm, err := qemu.GetVM(vmGetURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, deleteBody.Node, err)})
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: deleting VMID : %s in %s due to VM has not been stopped", vmid, deleteBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, deleteBody.Node)})
	}

	// Deleting target VM
	vmDeleteURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", deleteBody.Node, vmid))
	log.Printf("Deleting VMID : %s in %s", vmid, deleteBody.Node)
	_, deleteErr := qemu.DeleteVM(vmDeleteURL, cookies)
	if deleteErr != nil {
		log.Printf("Error: deleting VMID : %s in %s due to %s", vmid, deleteBody.Node, deleteErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting VMID : %s in %s due to %s", vmid, deleteBody.Node, deleteErr)})
	}

	// Check that target VM has been deleted completely yet
	deleted := qemu.DeleteCompletely(deleteBody.Node, vmid, cookies)
	if deleted {
		log.Printf("Finished deleting VMID : %s in %s", vmid, deleteBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been deleted", vmid, deleteBody.Node)})
	}
	log.Printf("Error: Could not delete VMID : %s in %s", vmid, deleteBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been deleted", vmid, deleteBody.Node)})
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
	// Getting request's body
	cloneBody := new(model.CloneBody)
	if err := c.BodyParser(cloneBody); err != nil {
		log.Println("Error: Could not parse body parser to clone VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to clone VM's body"})
	}
	newid := fmt.Sprint(cloneBody.NewID)

	// Getting data from query & Mapping values
	node := c.Query("node")
	vmid := c.Query("vmid")

	cookies := config.GetCookies(c)

	// Check VM Template from vmid
	isTemplate := qemu.IsTemplate(node, vmid, cookies)
	if isTemplate {
		// Check spec of the VM before allocate node
		vmGetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
		vm, vmInfoErr := qemu.GetVM(vmGetURL, cookies)
		if vmInfoErr != nil {
			log.Println(vm)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to get VM Template info for creating VM due to %s", vmInfoErr)})
		}

		// Parse mem, cpu, disk for checking free space
		vmSpec := model.VMSpec{
			Memory: vm.Info.MaxMem,
			CPU:    vm.Info.CPUs,
			Disk:   vm.Info.MaxDisk,
		}

		// Getting target node from node allocation
		workerNodes, target, nodeErr := cluster.AllocateNode(vmSpec, cookies)
		if nodeErr != nil {
			log.Println("Error: allocate node :", nodeErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to allocate node for creating VM due to %s", nodeErr)})
		}

		for _, workerNode := range workerNodes {
			vmListURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu", workerNode.Node))
			vmList, vmListErr := qemu.GetVMList(vmListURL, cookies)
			if vmListErr != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list from %s due to %s", workerNode.Node, vmListErr)})
			}

			// Checking duplicate VM by VMID
			var list []string
			for _, v := range vmList.Info {
				list = append(list, fmt.Sprintf("%d", v.VMID))
			}
			// log.Printf("VMs in node : %s : %s", workerNode.Node, list)

			if config.Contains(list, string(newid)) {
				log.Printf("Error: found duplicate VMID : %s in node : %s", newid, workerNode.Node)
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Found duplicate VMID: %s in %s", newid, workerNode.Node)})
			}
		}

		// Construct payload
		data := url.Values{}
		data.Set("newid", newid)
		data.Set("name", cloneBody.Name)
		data.Set("target", target)
		log.Println("clone body :", data)

		// Cloning VM
		log.Printf("Cloning VMID : %s in %s", newid, target)
		vmCloneURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/clone", node, vmid))
		info, cloneErr := qemu.CloneVM(vmCloneURL, data, cookies)
		if cloneErr != nil {
			log.Printf("Error: cloning VMID : %s in %s : %s", newid, target, cloneErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed cloning VMID : %s in %s due to %s", vmid, target, cloneErr)})
		}
		log.Println(info)

		// Waiting until cloning process has been completed
		cloned := qemu.CheckStatus(target, newid, []string{"created", "starting", "running"}, true, (10 * time.Minute), time.Second, cookies)
		if cloned {
			log.Printf("Finished cloning VMID : %s in %s", newid, target)
			return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Cloning new VMID: %s to %s successfully", newid, target)})
		}
		log.Printf("Error: cloning VMID : %s in %s has failed", newid, target)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Cloning new VMID: %s to %s has failed", newid, target)})
	}
	log.Printf("Error: cloning VMID : %s from VMID : %s", newid, vmid)
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Could not clone VM"})
}

// CreateTemplate - Templating specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/template
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID
*/
func CreateTemplate(c *fiber.Ctx) error {
	// Getting request's body
	templateBody := new(model.TemplateBody)
	if err := c.BodyParser(templateBody); err != nil {
		log.Println("Error: Could not parse body parser to create template VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to create template VM's body"})
	}
	vmid := fmt.Sprint(templateBody.VMID)

	cookies := config.GetCookies(c)

	// First check that target VM has been stopped
	vmGetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", templateBody.Node, vmid))
	vm, err := qemu.GetVM(vmGetURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s in %s due to %s", vmid, templateBody.Node, err)})
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: Could not template VMID : %s in %s due to VM hasn't been stopped", vmid, templateBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, templateBody.Node)})
	}

	// Templating VM
	log.Printf("Creating template from VMID : %s in %s", vmid, templateBody.Node)
	vmTemplateURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/template", templateBody.Node, vmid))
	_, templateErr := qemu.CreateTemplate(vmTemplateURL, cookies)
	if templateErr != nil {
		log.Printf("Error: Could not template VMID : %s in %s : %s", vmid, templateBody.Node, templateErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed templating VMID : %s in %s due to %s", vmid, templateBody.Node, templateErr)})
	}

	// Waiting until templating process has been completed
	templated := qemu.TemplateCompletely(templateBody.Node, vmid, []string{"created", "existing"}, cookies)
	if templated {
		log.Printf("Finished templating VMID : %s in %s", vmid, templateBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been templated", vmid, templateBody.Node)})
	}
	log.Printf("Error: Could not template VMID : %s in %s", vmid, templateBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been templated correctly", vmid, templateBody.Node)})
}
