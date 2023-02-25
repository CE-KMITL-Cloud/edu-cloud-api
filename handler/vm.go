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
	using Params
	? @username : account's username
	@node : node's name
	@vmid : VM's ID
*/
func GetVM(c *fiber.Ctx) error {
	// ? username := c.Query("username")
	node := c.Params("node")
	vmid := c.Params("vmid")

	cookies := config.GetCookies(c)
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

// GetVMListByNode - Getting VM list from given node
// GET /api2/json/nodes/{node}/qemu
/*
	using Query Params
	? @username : account's username

	using Params
	@node : node's name
*/
func GetVMListByNode(c *fiber.Ctx) error {
	// ? username := c.Query("username")
	node := c.Params("node")

	cookies := config.GetCookies(c)
	url := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu", node))
	log.Printf("Getting VM list from %s", node)
	vmList, err := qemu.GetVMListByNode(url, cookies)
	if err != nil {
		log.Println("Error: from getting VM's list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list from %s due to %s", node, err)})
	}
	log.Printf("Got VM list from node : %s", node)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": vmList})
}

// GetVMList - Getting VM list
// GET /api2/json/cluster/resources
func GetVMList(c *fiber.Ctx) error {
	cookies := config.GetCookies(c)
	vmList, err := qemu.GetVMList(cookies)
	if err != nil {
		log.Println("Error: from getting VM list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list due to %s", err)})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": vmList})
}

// CreateVM - Create VM on specific node
// POST /api2/json/nodes/{node}/qemu
/*
	using Request's Body
	@vmid : VM's ID
	@name : VM's name
	@memory : 1024 (MB)
	@cores : 2 (cores)
	@sockets : 1 (sockets of cpu fixed to 1)
	@onboot : {0, 1}
	@scsi0 : "ceph-vm:32"
	@cdrom : "cephfs:iso/ubuntu-20.04.4-live-server-amd64.iso"
	@net0 : "virtio,bridge=vmbr0,firewall=1"
	@scsihw : "virtio-scsi-single"
*/
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
		vmList, vmListErr := qemu.GetVMListByNode(vmListURL, cookies)
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
			vmList, vmListErr := qemu.GetVMListByNode(vmListURL, cookies)
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

// GetTemplateList - Getting VM Template list
// GET /api2/json/cluster/resources
func GetTemplateList(c *fiber.Ctx) error {
	cookies := config.GetCookies(c)
	log.Println("Getting VM Template list")
	templateList, err := qemu.GetTemplateList(cookies)
	if err != nil {
		log.Println("Error: from getting VM's list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM Template list due to %s", err)})
	}
	log.Println("Got VM Template list")
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": templateList})
}

// EditVM - Set virtual machine options (asynchrounous API).
// POST /api2/json/nodes/{node}/qemu/{vmid}/config
// PUT /api2/json/nodes/{node}/qemu/{vmid}/resize
/*
	using Query Params
	@node : node's name
	@vmid : VM's ID

	using Request's Body
	@cores : Amount of CPU core
	@memory : Amount of RAM in (MB)
	@disk : Amount of Disk (scsi0) to increase in Size_in_GiB format
*/
func EditVM(c *fiber.Ctx) error {
	// Getting request's body
	editBody := new(model.EditBody)
	if err := c.BodyParser(editBody); err != nil {
		log.Println("Error: Could not parse body parser to edit VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to edit VM's body"})
	}
	editCores := fmt.Sprint(editBody.Cores)
	editMemory := fmt.Sprint(editBody.Memory)
	editMaxMemory := config.MBtoByte(editBody.Memory)

	// Getting data from query & Mapping values
	node := c.Query("node")
	vmid := c.Query("vmid")

	cookies := config.GetCookies(c)
	nodeInfo, nodeInfoErr := cluster.GetNode(node, cookies)
	if nodeInfoErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to get node info for editing VM due to %s", nodeInfoErr)})
	}
	freeMemory, freeCPU, freeDisk := nodeInfo.MaxMem-nodeInfo.Mem, nodeInfo.MaxCPU-nodeInfo.CPU, nodeInfo.MaxDisk-nodeInfo.Disk

	// Check VM spec before edit configuration
	vmGetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
	vm, vmInfoErr := qemu.GetVM(vmGetURL, cookies)
	if vmInfoErr != nil {
		log.Println(vm)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to get VM info for editing VM due to %s", vmInfoErr)})
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: editing VMID : %s in %s due to VM has not been stopped", vmid, node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s in %s hasn't been stopped", vmid, node)})
	}

	// Parse mem, cpu, disk for checking free space
	vmSpec := model.VMSpec{
		Memory: vm.Info.MaxMem,
		CPU:    vm.Info.CPUs,
		Disk:   vm.Info.MaxDisk,
	}
	editBodyByteDisk := config.GBtoByte(editBody.Disk)
	editBodyDisk := fmt.Sprint(`+`, editBodyByteDisk)

	// Construct payload
	data := url.Values{}
	data.Set("cores", editCores)
	data.Set("memory", editMemory)

	resizeData := url.Values{}
	resizeData.Set("disk", "scsi0")
	resizeData.Set("size", editBodyDisk)
	log.Println("resize data:", resizeData)

	extendMemory, extendCPU, extendDisk := editMaxMemory-vmSpec.Memory, editBody.Cores-vmSpec.CPU, editBodyByteDisk+vmSpec.Disk
	if editMaxMemory < vmSpec.Memory {
		extendMemory = vmSpec.Memory - editMaxMemory
	}
	if editBody.Cores < vmSpec.CPU {
		extendCPU = vmSpec.CPU - editBody.Cores
	}

	log.Printf("VM spec : {cpu: %f, mem: %d, disk: %d}", vmSpec.CPU, vmSpec.Memory, vmSpec.Disk)
	log.Printf("Edit spec : {cpu: %f, mem: %d, disk: %d}", editBody.Cores, editMaxMemory, extendDisk)
	log.Printf("Free Node spec : {cpu: %f, mem: %d, disk: %d}", freeCPU, freeMemory, freeDisk)
	log.Printf("Extended spec : {cpu: %f, mem: %d, disk: %d}", extendCPU, extendMemory, editBodyByteDisk)
	log.Printf("Remain Node spec : {cpu: %f, mem: %d, disk: %d}", (freeCPU - extendCPU), (freeMemory - extendMemory), (freeDisk - editBodyByteDisk))

	// Check free space of node
	if freeCPU > extendCPU && freeMemory > extendMemory && freeDisk > editBodyByteDisk {
		// Approve if request is spec increasing only
		if config.GreaterOrEqual(editBody.Cores, vmSpec.CPU, editMaxMemory, vmSpec.Memory, editBodyByteDisk+vmSpec.Disk, vmSpec.Disk) {
			log.Println("Able to edit VM config")
			log.Printf("Editing VMID : %s in %s", vmid, node)
			vmEditURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/config", node, vmid))
			info, editErr := qemu.EditVM(vmEditURL, data, cookies)
			if editErr != nil {
				log.Printf("Error: editing VMID : %s in %s : %s", vmid, node, editErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed editing VMID : %s in %s due to %s", vmid, node, editErr)})
			}
			log.Println(info)
			resizeDiskURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/resize", node, vmid))
			_, resizeInfoErr := qemu.ResizeDisk(resizeDiskURL, resizeData, cookies)
			if resizeInfoErr != nil {
				log.Printf("Error: editing disk on VMID : %s in %s : %s", vmid, node, editErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed editing disk on VMID : %s in %s due to %s", vmid, node, resizeInfoErr)})
			}
			return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Edited VM : %s in %s successfully", vmid, node)})
		}
		log.Printf("Error: editing VMID : %s in %s due to request spec is lower or equal to current spec", vmid, node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Could not edit VM due to request spec is lower or equal to current spec"})
	}
	log.Printf("Error: editing VMID : %s in %s due to have no enough free space", vmid, node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Internal server error", "message": "Node have no enough free space"})
}
