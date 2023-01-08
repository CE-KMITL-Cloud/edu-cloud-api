// Package model - structs
package model

// VM - struct for VM
type VM struct {
	Info []VMInfo `json:"data"`
}

// VMInfo - struct for VM's info
type VMInfo struct {
	CPU       float64 `json:"cpu"`
	NetOut    uint64  `json:"netout"`
	DiskWrite uint64  `json:"diskwrite"`
	UpTime    uint64  `json:"uptime"`
	MaxMem    uint64  `json:"maxmem"`
	PID       uint32  `json:"pid"`
	Name      string  `json:"name"`
	VMID      uint16  `json:"vmid"`
	MaxDisk   uint64  `json:"maxdisk"`
	Status    string  `json:"status"`
	Disk      uint16  `json:"disk"`
	DiskRead  uint64  `json:"diskread"`
	CPUs      uint16  `json:"cpus"`
	Mem       uint64  `json:"mem"`
	NetIn     uint64  `json:"netin"`
}

// VMResponse - struct for VM's response
type VMResponse struct {
	Info string `json:"data"`
}

// CloneBody - struct for request Cloning VM
type CloneBody struct {
	NewID uint32 `form:"newid"`
	Name  string `form:"name"`
}

/*
* Cloning status
{
	"data": {
		"mem": 0,
		"netin": 0,
		"ha": {
			"managed": 0
		},
		"maxmem": 536870912,
		"name": "VM 205",
		"disk": 0,
		"diskread": 0,
		"vmid": 205,
		"maxdisk": 0,
		"status": "stopped",
		"cpus": 1,
		"diskwrite": 0,
		"lock": "clone",
		"uptime": 0,
		"cpu": 0,
		"netout": 0,
		"qmpstatus": "stopped"
	}
}

* After Cloning finished - taking around ~10 mins
{
	"data": {
		"ha": {
			"managed": 0
		},
		"mem": 0,
		"netin": 0,
		"cpus": 2,
		"maxdisk": 34359738368,
		"status": "stopped",
		"vmid": 205,
		"disk": 0,
		"diskread": 0,
		"name": "vmclone1",
		"maxmem": 2147483648,
		"uptime": 0,
		"diskwrite": 0,
		"qmpstatus": "stopped",
		"netout": 0,
		"cpu": 0
	}
}

* Running
{
	"data": {
		"name": "vmtest3",
		"maxmem": 2147483648,
		"maxdisk": 34359738368,
		"diskread": 3937701920,
		"ha": {
			"managed": 0
		},
		"mem": 1345167360,
		"netin": 707164116,
		"nics": {
			"tap203i0": {
				"netin": 707164116,
				"netout": 8943566
			}
		},
		"blockstat": {
			"scsi0": {
				"account_invalid": true,
				"flush_operations": 31125,
				"unmap_merged": 0,
				"idle_time_ns": 35030651650,
				"failed_wr_operations": 0,
				"unmap_total_time_ns": 103503900,
				"timed_stats": [],
				"invalid_unmap_operations": 0,
				"rd_bytes": 2123895808,
				"wr_bytes": 12350885376,
				"flush_total_time_ns": 4686807599,
				"wr_operations": 178026,
				"failed_rd_operations": 0,
				"rd_operations": 41194,
				"unmap_bytes": 61916037120,
				"unmap_operations": 3978,
				"failed_flush_operations": 0,
				"invalid_flush_operations": 0,
				"failed_unmap_operations": 0,
				"rd_merged": 0,
				"account_failed": true,
				"wr_highest_offset": 34359738368,
				"wr_merged": 0,
				"wr_total_time_ns": 3034183589413,
				"invalid_wr_operations": 0,
				"invalid_rd_operations": 0,
				"rd_total_time_ns": 157531084777
			},
			"ide2": {
				"failed_unmap_operations": 0,
				"rd_merged": 0,
				"account_failed": true,
				"wr_merged": 0,
				"wr_highest_offset": 0,
				"wr_total_time_ns": 0,
				"invalid_wr_operations": 0,
				"invalid_rd_operations": 0,
				"rd_total_time_ns": 44455191852,
				"flush_total_time_ns": 0,
				"wr_operations": 0,
				"failed_rd_operations": 0,
				"rd_operations": 66093,
				"unmap_bytes": 0,
				"unmap_operations": 0,
				"failed_flush_operations": 0,
				"invalid_flush_operations": 0,
				"timed_stats": [],
				"invalid_unmap_operations": 0,
				"rd_bytes": 1813806112,
				"wr_bytes": 0,
				"account_invalid": true,
				"flush_operations": 0,
				"unmap_merged": 0,
				"idle_time_ns": 595596406762012,
				"failed_wr_operations": 0,
				"unmap_total_time_ns": 0
			}
		},
		"ballooninfo": {
			"last_update": 1673112186,
			"total_mem": 2079428608,
			"free_mem": 734261248,
			"minor_page_faults": 5552853,
			"actual": 2147483648,
			"max_mem": 2147483648,
			"major_page_faults": 1999,
			"mem_swapped_out": 0,
			"mem_swapped_in": 0
		},
		"diskwrite": 12350885376,
		"freemem": 734261248,
		"pid": 307519,
		"cpus": 2,
		"status": "running",
		"balloon": 2147483648,
		"vmid": 203,
		"disk": 0,
		"netout": 8943566,
		"cpu": 0.0406608122965322,
		"qmpstatus": "running",
		"running-qemu": "7.1.0",
		"running-machine": "pc-i440fx-7.1+pve0",
		"proxmox-support": {
			"pbs-library-version": "1.3.1 (4d450bb294cac5316d2f23bf087c4b02c0543d79)",
			"pbs-dirty-bitmap-savevm": true,
			"pbs-dirty-bitmap": true,
			"backup-max-workers": true,
			"pbs-dirty-bitmap-migration": true,
			"pbs-masterkey": true,
			"query-bitmap-info": true
		},
		"uptime": 630927
	}
}
*/
