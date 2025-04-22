package main

import (
	"fmt"
	"log"

	libvirt "github.com/libvirt/libvirt-go"
)

func main() {
	// Connect to the local libvirt daemon
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	// Get list of active (running) domains (VMs)
	domainIDs, err := conn.ListDomains()
	if err != nil {
		log.Fatalf("Failed to list domains: %v", err)
	}

	fmt.Printf("Running VMs: %d\n", len(domainIDs))

	for _, id := range domainIDs {
		dom, err := conn.LookupDomainById(id)
		if err != nil {
			log.Printf("Failed to lookup domain by ID %d: %v", id, err)
			continue
		}
		defer dom.Free()

		name, _ := dom.GetName()
		info, err := dom.GetInfo()
		if err != nil {
			log.Printf("Failed to get info for domain %s: %v", name, err)
			continue
		}

		fmt.Printf("\nVM Name: %s\n", name)
		fmt.Printf("  State: %v\n", info.State)
		fmt.Printf("  Max Memory: %d KB\n", info.MaxMem)
		fmt.Printf("  Used Memory: %d KB\n", info.Memory)
		fmt.Printf("  Number of VCPUs: %d\n", info.NrVirtCpu)
		fmt.Printf("  CPU Time: %d ns\n", info.CpuTime)


    memStats, err := dom.MemoryStats(11, 0)
    if err != nil {
    	log.Printf("Failed to get memory stats for domain %s: %v", name, err)
    } else {
    	for _, stat := range memStats {
    		switch libvirt.DomainMemoryStatTags(stat.Tag) {
    		case libvirt.DOMAIN_MEMORY_STAT_ACTUAL_BALLOON:
    			fmt.Printf("  Actual Memory (balloon): %d KB\n", stat.Val)
    		case libvirt.DOMAIN_MEMORY_STAT_RSS:
    			fmt.Printf("  Resident Set Size (RSS): %d KB\n", stat.Val)
    		case libvirt.DOMAIN_MEMORY_STAT_UNUSED:
    			fmt.Printf("  Unused Memory: %d KB\n", stat.Val)
    		}
    	}
    }
	}
}
