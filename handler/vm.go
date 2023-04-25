// Package handler - handling context
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/edu-cloud-api/config"
	"github.com/edu-cloud-api/database"
	"github.com/edu-cloud-api/internal/cluster"
	"github.com/edu-cloud-api/internal/qemu"
	"github.com/edu-cloud-api/model"
	"github.com/gofiber/fiber/v2"
)

// GetVM - Getting specific VM's info from Proxmox
// GET /api2/json/nodes/{node}/qemu/{vmid}/status/current
/*
	using Params
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func GetVM(c *fiber.Ctx) error {
	node := c.Params("node")
	vmid := c.Params("vmid")
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)
	log.Printf("Getting detail from VMID : %s in %s", vmid, node)
	url := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", node, vmid))
	info, err := qemu.GetVM(url, cookies)
	if err != nil {
		log.Println("Error: from getting VM's info :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s due to %s", vmid, err)})
	}
	log.Printf("Got info from vmid : %s", vmid)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": info})
}

// GetVMListByNode - Getting VM list from given node 🚫
// GET /api2/json/nodes/{node}/qemu
/*
	using Params
	@node : node's name
*/
// ! this function might never be used. to be deprecated
func GetVMListByNode(c *fiber.Ctx) error {
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

// GetVMList - Getting VM list (VM Template not included)
// GET /api2/json/cluster/resources
/*
	using Query
	username : account's username
*/
func GetVMList(c *fiber.Ctx) error {
	var returnList []model.VMsInfo
	username := c.Query("username")
	cookies := config.GetCookies(c)
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	vmList, err := qemu.GetVMList(cookies)
	if err != nil {
		log.Println("Error: from getting VM list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list due to %s", err)})
	}
	if group == config.ADMIN {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": vmList})
	}
	list, _ := database.GetAllInstancesIDByOwner(username)
	for _, vm := range vmList {
		if config.Contains(list, fmt.Sprint(vm.VMID)) {
			returnList = append(returnList, vm)
		}
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": returnList})
}

// CreateVM - Create VM on specific node
// POST /api2/json/nodes/{node}/qemu
/*
	using Request's Body
	@name : VM's name
	@memory : 1024 (MB)
	@cores : 2 (cores)
	@storage : ceph-vm
	@disk : 32 (Amount of disk in GiB)
	@cdrom : "cephfs:iso/" + "ubuntu-20.04.4-live-server-amd64.iso"

	using Query
	@username : account's username
*/
func CreateVM(c *fiber.Ctx) error {
	createBody := new(model.CreateBody)
	if err := c.BodyParser(createBody); err != nil {
		log.Println("Error: Could not parse body parser to create VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to create VM's body"})
	}
	// check faculty, admin role
	username := c.Query("username")
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.STUDENT {
		log.Println("Error: user's group is not allowed to create VM")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to create VM due to user's group is not allowed"})
	}
	cookies := config.GetCookies(c)
	vmid, getVMIDErr := qemu.GetVMID(cookies)
	if getVMIDErr != nil {
		log.Println("Error: while getting vmid due to :", getVMIDErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting vmid due to %s", getVMIDErr)})
	}
	scsi0 := fmt.Sprintf("%s:%s", createBody.Storage, createBody.Disk)
	cdrom := config.ISO + createBody.CDROM

	// Construct payload
	data := url.Values{}
	data.Set("vmid", vmid)
	data.Set("name", createBody.Name)
	data.Set("memory", fmt.Sprint(createBody.Memory)) // memory (MB)
	data.Set("cores", fmt.Sprint(createBody.Cores))   // cpu (core)
	data.Set("sockets", fmt.Sprint(config.SOCKET))
	data.Set("onboot", fmt.Sprint(config.ONBOOT))
	data.Set("scsi0", scsi0) // "ceph-vm:32"
	data.Set("cdrom", cdrom)
	data.Set("net0", config.NET0)
	data.Set("scsihw", config.SCSIHW)

	maxDisk, parseErr := strconv.ParseUint(createBody.Disk, 10, 64)
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
	// Getting target node from node allocation
	workerNodes, target, nodeErr := cluster.AllocateNode(vmSpec, createBody.Storage, cookies)
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
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM list due to %s", vmListErr)})
		}
		var list []string
		for _, v := range vmList.Info {
			list = append(list, fmt.Sprintf("%d", v.VMID))
		}
		// log.Printf("VMs in node : %s : %s", workerNode.Node, list)

		if config.Contains(list, string(vmid)) {
			log.Printf("Error: found duplicate VMID : %s in node : %s", vmid, workerNode.Node)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Found duplicate VMID: %s ", vmid)})
		}
	}

	// Creating VM in Proxmox
	vmCreateURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu", target))
	log.Printf("Creating VMID : %s in %s", vmid, target)
	_, err := qemu.CreateVM(vmCreateURL, data, cookies)
	if err != nil {
		log.Println("Error: from creating VM :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed creating VMID : %s due to %s", vmid, err)})
	}

	// Waiting until creating process has been complete
	created := qemu.CheckStatus(target, vmid, []string{"created", "starting", "running"}, true, (time.Minute), time.Second)
	if created {
		// Creating VM in DB
		if _, createInstanceErr := database.CreateInstance(vmid, username, target, createBody.Name, vmSpec); createInstanceErr != nil {
			log.Printf("Error: Could not create VMID : %s in %s due to %s", vmid, target, createInstanceErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Creating new VMID: %s has failed due to %s", vmid, createInstanceErr)})
		}

		// todo : pull mac addr
		// todo : insert into proxy table
		log.Printf("Finished creating VMID : %s in %s", vmid, target)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Creating new VMID: %s successfully", vmid)})
	}
	log.Printf("Error: Could not create VMID : %s in %s", vmid, target)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Creating new VMID: %s has failed", vmid)})
}

// DeleteVM - Deleting specific VM
// DELETE /api2/json/nodes/{node}/qemu/{vmid}
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func DeleteVM(c *fiber.Ctx) error {
	// Getting request's body
	deleteBody := new(model.DeleteBody)
	if err := c.BodyParser(deleteBody); err != nil {
		log.Println("Error: Could not parse body parser to delete VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to delete VM's body"})
	}
	vmid := fmt.Sprint(deleteBody.VMID)
	username := c.Query("username")
	cookies := config.GetCookies(c)

	// Check that user is owner of given VM
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		log.Printf("Error: Could not delete VMID : %s in %s", vmid, deleteBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to delete VMID: %s", vmid)})
	}

	// First check that target VM has been stopped
	vmGetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", deleteBody.Node, vmid))
	vm, err := qemu.GetVM(vmGetURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s due to %s", vmid, err)})
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: deleting VMID : %s in %s due to VM has not been stopped", vmid, deleteBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s hasn't been stopped", vmid)})
	}

	// Delete target VM
	vmDeleteURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s", deleteBody.Node, vmid))
	log.Printf("Deleting VMID : %s in %s", vmid, deleteBody.Node)
	_, deleteErr := qemu.DeleteVM(vmDeleteURL, cookies)
	if deleteErr != nil {
		log.Printf("Error: deleting VMID : %s in %s due to %s", vmid, deleteBody.Node, deleteErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting VMID : %s due to %s", vmid, deleteErr)})
	}

	// Check that target VM has been deleted completely yet
	deleted := qemu.DeleteCompletely(deleteBody.Node, vmid)
	if deleted {
		log.Printf("Finished deleting VMID : %s in %s", vmid, deleteBody.Node)

		// Delete VM in DB
		if deleteInstanceErr := database.DeleteInstance(vmid); deleteInstanceErr != nil {
			log.Printf("Error: Deleting instance ID : %s from DB due to %s", vmid, deleteInstanceErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed deleting instance ID : %s from DB due to %s", vmid, deleteInstanceErr)})
		}

		// Delete VM in Pool DB
		pools, getPoolsErr := database.GetPoolsByVMID(vmid)
		if getPoolsErr != nil {
			log.Printf("Error: Getting pools from given vmid : %s from DB due to %s", vmid, getPoolsErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed removing instance ID : %s from pools DB due to %s", vmid, getPoolsErr)})
		}
		for _, pool := range pools {
			pool.VMID = config.FilterString(pool.VMID, vmid)
			updateErr := database.AddPoolInstances(pool.Code, pool.Owner, pool.VMID)
			if updateErr != nil {
				log.Printf("Error: updating instances of pool code : %s, owner : %s due to %s", pool.Code, pool.Owner, updateErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Internal server error", "message": fmt.Sprintf("Failed updating instances of pool code : %s, owner : %s due to %s", pool.Code, pool.Owner, updateErr)})
			}
			log.Printf("Successfully removed template ID : %s to pool code : %s, owner : %s", vmid, pool.Code, pool.Owner)
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s has been deleted", vmid)})
	}
	log.Printf("Error: Could not delete VMID : %s in %s", vmid, deleteBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s hasn't been deleted", vmid)})
}

// CloneVM - Cloning specific VM
// POST /api2/json/nodes/{node}/qemu/{vmid}}/clone
/*
	using Query Params
	@username : account's username
	@node : node's name
	@vmid : VM's ID

	using Request's Body
	@name : VM's name
	@storage : Storage's name
	@full : 1
	@ciuser : cloudinit's username
	@cipassword : cloudinit's password
*/
func CloneVM(c *fiber.Ctx) error {
	// Getting request's body
	cloneBody := new(model.CloneBody)
	if err := c.BodyParser(cloneBody); err != nil {
		log.Println("Error: Could not parse body parser to clone VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to clone VM's body"})
	}
	cookies := config.GetCookies(c)

	// getting data from query & Mapping values
	username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")

	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}

	// able to clone only own template or sizing template except man who request is admin
	isSizingTemplate, _ := database.IsSizingTemplate(vmid)
	if !isSizingTemplate {
		// get template from every pools that username is member
		var poolInstances []string
		pools, getPoolsErr := database.GetAllPoolsByMember(username)
		if getPoolsErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting pool templates list from DB due to %s", getPoolsErr)})
		}
		for _, pool := range pools {
			for _, instance := range pool.VMID {
				if !config.Contains(poolInstances, instance) {
					poolInstances = append(poolInstances, instance)
				}
			}
		}
		log.Println(poolInstances)
		instanceTemplateOwner, _ := database.CheckInstanceTemplateOwner(username, vmid)
		if !instanceTemplateOwner && group != config.ADMIN && !config.Contains(poolInstances, vmid) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed cloning VMID : %s due to VM is not template or user is not owner", vmid)})
		}
	}

	// getting new vmid
	newid, getVMIDErr := qemu.GetVMID(cookies)
	if getVMIDErr != nil {
		log.Println("Error: while getting vmid due to :", getVMIDErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting vmid due to %s", getVMIDErr)})
	}

	// Check VM Template from vmid
	isTemplate := qemu.IsTemplate(node, vmid)
	if isTemplate || group == config.ADMIN {
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
		workerNodes, target, nodeErr := cluster.AllocateNode(vmSpec, cloneBody.Storage, cookies)
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
			var list []string
			for _, v := range vmList.Info {
				list = append(list, fmt.Sprintf("%d", v.VMID))
			}
			// log.Printf("VMs in node : %s : %s", workerNode.Node, list)
			if config.Contains(list, string(newid)) {
				log.Printf("Error: found duplicate VMID : %s in node : %s", newid, workerNode.Node)
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Found duplicate VMID: %s", newid)})
			}
		}

		// Construct payload
		data := url.Values{}
		data.Set("newid", newid)
		data.Set("name", cloneBody.Name)
		data.Set("target", target)
		data.Set("full", "1") // ! fixed to `1` for full clone
		log.Println("clone body :", data)

		// Cloning VM in Proxmox
		log.Printf("Cloning VMID : %s in %s", newid, target)
		vmCloneURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/clone", node, vmid))
		_, cloneErr := qemu.CloneVM(vmCloneURL, data, cookies)
		if cloneErr != nil {
			log.Printf("Error: cloning VMID : %s in %s : %s", newid, target, cloneErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed cloning VMID : %s due to %s", vmid, cloneErr)})
		}

		// Waiting until cloning process has been completed
		cloned := qemu.CheckStatus(target, newid, []string{"created", "stopped", "running"}, false, (10 * time.Minute), time.Second)
		if cloned {

			// Creating VM in DB
			if _, createInstanceErr := database.CreateInstance(newid, username, target, cloneBody.Name, vmSpec); createInstanceErr != nil {
				log.Printf("Error: Could not create VMID : %s in %s due to %s", newid, target, createInstanceErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Creating new VMID: %s has failed due to %s", newid, createInstanceErr)})
			}

			// resize disk to sizing template's disk in DB
			if isSizingTemplate {
				log.Printf("Resizing VMID : %s in %s", newid, target)
				sizing, _ := database.GetTemplate(vmid)
				sizingDiskByte := config.GBtoByteFloat(sizing.MaxDisk) - vm.Info.MaxDisk
				sizingDisk := fmt.Sprint(`+`, sizingDiskByte)

				resizeData := url.Values{}
				resizeData.Set("disk", "scsi0")
				resizeData.Set("size", sizingDisk)
				log.Println("resize data:", resizeData)

				// Resizing Disk in Proxmox
				resizeDiskURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/resize", target, newid))
				_, resizeInfoErr := qemu.ResizeDisk(resizeDiskURL, resizeData, cookies)
				if resizeInfoErr != nil {
					log.Printf("Error: resizing disk of VMID : %s in %s : %s", newid, target, resizeInfoErr)
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed resizing disk of VMID : %s due to %s", newid, resizeInfoErr)})
				}

				// Resizing Disk in DB
				resizeErr := database.ResizeDisk(newid, sizing.MaxDisk)
				if resizeErr != nil {
					log.Printf("Error: resizing disk of VMID : %s in DB : %s", newid, resizeErr)
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed resizing disk of VMID : %s in DB due to %s", newid, resizeErr)})
				}
			}
			// config ciuser, cipassword
			editData := url.Values{}
			editData.Set("ciuser", cloneBody.CIUser)
			editData.Set("cipassword", cloneBody.CIPass)
			log.Println("edit body :", editData)

			log.Printf("Editing VMID : %s in %s", newid, target)
			vmEditURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/config", target, newid))
			_, editErr := qemu.EditVM(vmEditURL, editData, cookies)
			if editErr != nil {
				log.Printf("Error: editing VMID : %s in %s : %s", newid, target, editErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed editing VMID : %s in %s due to %s", newid, target, editErr)})
			}
			// todo : pull mac addr
			// todo : insert into proxy table
			log.Printf("Finished cloning VMID : %s in %s", newid, target)
			return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Cloning new VMID: %s successfully", newid)})
		}
		log.Printf("Error: cloning VMID : %s in %s has failed", newid, target)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed cloning new VMID: %s", newid)})
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

	using Query
	@username : account's username
*/
func CreateTemplate(c *fiber.Ctx) error {
	// Getting request's body
	templateBody := new(model.TemplateBody)
	if err := c.BodyParser(templateBody); err != nil {
		log.Println("Error: Could not parse body parser to create template VM's body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to create template VM's body"})
	}
	vmid := fmt.Sprint(templateBody.VMID)

	// check faculty, admin role
	username := c.Query("username")
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.STUDENT {
		log.Println("Error: user's group is not allowed to create VM")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Failed to templating VM due to user's group is not allowed"})
	}
	// Check that user is owner of given VM
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		log.Printf("Error: templating VMID : %s in %s due to user is not owner of VM", vmid, templateBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed templating VMID : %s due to user is not owner of VM", vmid)})
	}

	// First check that target VM has been stopped
	cookies := config.GetCookies(c)
	vmGetURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/status/current", templateBody.Node, vmid))
	vm, err := qemu.GetVM(vmGetURL, cookies)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting detail from VMID: %s due to %s", vmid, err)})
	}

	// If target VM's status is not "stopped" then return
	if vm.Info.Status != "stopped" {
		log.Printf("Error: Could not template VMID : %s in %s due to VM hasn't been stopped", vmid, templateBody.Node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Target VMID: %s hasn't been stopped", vmid)})
	}

	// Templating VM
	log.Printf("Creating template from VMID : %s in %s", vmid, templateBody.Node)
	vmTemplateURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/template", templateBody.Node, vmid))
	_, templateErr := qemu.CreateTemplate(vmTemplateURL, cookies)
	if templateErr != nil {
		log.Printf("Error: Could not template VMID : %s in %s : %s", vmid, templateBody.Node, templateErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed templating VMID : %s due to %s", vmid, templateErr)})
	}

	// Waiting until templating process has been completed
	templated := qemu.TemplateCompletely(templateBody.Node, vmid, []string{"created", "existing"})
	if templated {
		updateErr := database.TemplateInstance(vmid)
		if updateErr != nil {
			log.Printf("Error: Could not update template status in DB VMID : %s in %s : %s", vmid, templateBody.Node, updateErr)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed updating template status in DB VMID : %s due to %s", vmid, updateErr)})
		}
		log.Printf("Finished templating VMID : %s in %s", vmid, templateBody.Node)
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Target VMID: %s in %s has been templated", vmid, templateBody.Node)})
	}
	log.Printf("Error: Could not template VMID : %s in %s", vmid, templateBody.Node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Target VMID: %s hasn't been templated correctly", vmid)})
}

// GetTemplateList - Getting VM Template list
// GET /api2/json/cluster/resources
/*
	using Query
	@username - account's username
*/
func GetTemplateList(c *fiber.Ctx) error {
	var returnList []model.VMsInfo
	cookies := config.GetCookies(c)
	username := c.Query("username")
	log.Println("Getting VM Template list")
	templateList, err := qemu.GetTemplateList(cookies)
	if err != nil {
		log.Println("Error: from getting VM's list :", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VM Template list due to %s", err)})
	}
	group, getGroupErr := database.GetUserGroup(username)
	if getGroupErr != nil {
		log.Println("Error: while getting user's group due to :", getGroupErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed to getting user's group due to %s", getGroupErr)})
	}
	if group == config.ADMIN {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": templateList})
	}

	// get sizing templates's id
	templates, getTemplatesErr := database.GetAllTemplatesID()
	if getTemplatesErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting sizing templates list from DB due to %s", getTemplatesErr)})
	}
	// get template from every pools that username is member
	var poolInstances []string
	pools, getPoolsErr := database.GetAllPoolsByMember(username)
	if getPoolsErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting pool templates list from DB due to %s", getPoolsErr)})
	}
	for _, pool := range pools {
		for _, instance := range pool.VMID {
			if !config.Contains(poolInstances, instance) {
				poolInstances = append(poolInstances, instance)
			}
		}
	}
	// get templates from given username
	list := database.GetAllInstanceTemplatesIDByOwner(username)
	for _, template := range templateList {
		if config.Contains(list, fmt.Sprint(template.VMID)) || config.Contains(templates, fmt.Sprint(template.VMID)) || config.Contains(poolInstances, fmt.Sprint(template.VMID)) {
			returnList = append(returnList, template)
		}
	}
	log.Println("Got VM Template list")
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": returnList})
}

// EditVM - Set virtual machine options (asynchrounous API).
// POST /api2/json/nodes/{node}/qemu/{vmid}/config
// PUT /api2/json/nodes/{node}/qemu/{vmid}/resize
/*
	using Query Params
	@username : account's username
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
	username := c.Query("username")
	node := c.Query("node")
	vmid := c.Query("vmid")
	cookies := config.GetCookies(c)

	// able to edit only own vm except requester is admin
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed deleting VMID : %s due to user is not owner of VM", vmid)})
	}
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
				log.Printf("Error: editing disk on VMID : %s in %s : %s", vmid, node, resizeInfoErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed editing disk on VMID : %s in %s due to %s", vmid, node, resizeInfoErr)})
			}

			// update vm spec in DB
			updateErr := database.EditInstance(model.Instance{
				MaxCPU:  editBody.Cores,
				MaxRAM:  config.MBtoGB(editBody.Memory),
				MaxDisk: float64(editBody.Disk),
			})
			if updateErr != nil {
				log.Printf("Error: updating VM spec on VMID : %s in %s due to %s", vmid, node, updateErr)
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Internal server error", "message": fmt.Sprintf("Failed updating VM spec on VMID : %s due to %s", vmid, updateErr)})
			}
			return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": fmt.Sprintf("Edited VM : %s in %s successfully", vmid, node)})
		}
		log.Printf("Error: editing VMID : %s in %s due to request spec is lower or equal to current spec", vmid, node)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": "Could not edit VM due to request spec is lower or equal to current spec"})
	}
	log.Printf("Error: editing VMID : %s in %s due to have no enough free space", vmid, node)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Internal server error", "message": "Node have no enough free space"})
}

// GetVncTicket - Get VNC Ticket from given VMID
// POST /api2/json/nodes/{node}/qemu/{vmid}/vncproxy
/*
	using Request's Body
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func GetVncTicket(c *fiber.Ctx) error {
	// Getting request's body
	vncProxyBody := new(model.VncProxyBody)
	if err := c.BodyParser(vncProxyBody); err != nil {
		log.Println("Error: Could not parse body parser to VNC Proxy body")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": "Failed parsing body parser to VNC Proxy body"})
	}
	vmid := fmt.Sprint(vncProxyBody.VMID)
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	cookies := config.GetCookies(c)
	getVncTicketURL := config.GetURL(fmt.Sprintf("/api2/json/nodes/%s/qemu/%s/vncproxy", vncProxyBody.Node, vmid))
	data := url.Values{}
	data.Set("websocket", "0")
	ticket, getTicketErr := qemu.VncProxy(getVncTicketURL, data, cookies)
	if getTicketErr != nil {
		log.Printf("Error: getting VNC Proxy ticket from VMID : %s in %s : %s", vmid, vncProxyBody.Node, getTicketErr)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "Failure", "message": fmt.Sprintf("Failed getting VNC Proxy ticket from VMID : %s in %s due to %s", vmid, vncProxyBody.Node, getTicketErr)})
	}
	log.Printf("Finished getting VNC Proxy ticket from VMID : %s in %s", vmid, vncProxyBody.Node)

	ticket.Detail.Url = fmt.Sprintf("wss://edu.ce.kmitl.cloud/api2/json/nodes/%s/qemu/%s/vncwebsocket?port=%s&vncticket=%s", vncProxyBody.Node, vmid, ticket.Detail.Port, ticket.Detail.Ticket)

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": ticket})
}

// GetVncConsole - Get VNC console URL from given VMID
/*
	using Params
	@node : node's name
	@vmid : VM's ID

	using Query
	@username : account's username
*/
func GetVncConsole(c *fiber.Ctx) error {
	vmid := c.Params("vmid")
	node := c.Params("node")
	username := c.Query("username")
	owner, checkOwnerErr := database.CheckInstanceOwner(username, vmid)
	if checkOwnerErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to %s", vmid, checkOwnerErr)})
	}
	if !owner {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "Bad request", "message": fmt.Sprintf("Failed getting VMID : %s due to user is not owner of VM", vmid)})
	}
	consoleURL := config.GetFromENV("PROXMOX_HOST") + fmt.Sprintf("/?console=kvm&novnc=1&vmid=%s&vmname=&node=%s&resize=off&cmd=", vmid, node)
	log.Printf("Finished getting VNC console URL from VMID : %s in %s", vmid, node)
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "Success", "message": consoleURL})
}
