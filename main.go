package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	libvirt "github.com/libvirt/libvirt-go"
)

const nanosecondsPerSecond = 1e9 // 1 billion nanoseconds in a second

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

	// Configure sleep duration (default to 500ms)
	sleepDuration := 500 * time.Millisecond

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Libvirt VM Statistics</h1>")
		fmt.Fprintf(w, "<table border='1'><tr><th>VM Name</th><th>State</th><th>Used Memory</th><th>VCPUs</th><th>CPU Time</th><th>Unused Memory (GB)</th><th>Max Memory (GB)</th></tr>")

		// Capture initial CPU times
		initialCPUTimes := make(map[uint32]float64)
		for _, id := range domainIDs {
			dom, err := conn.LookupDomainById(id)
			if err != nil {
				log.Printf("Failed to lookup domain by ID %d: %v", id, err)
				continue
			}
			defer dom.Free()

			info, err := dom.GetInfo()
			if err != nil {
				name, err := dom.GetName()
				if err != nil {
					log.Printf("Failed to get info or name for domain ID %d: %v", id, err)
				} else {
					log.Printf("Failed to get info for domain %s: %v", name, err)
				}
				continue
			}

			initialCPUTimes[id] = float64(info.CpuTime)
		}

		// Wait for the specified duration
		time.Sleep(sleepDuration)

		// Capture final CPU times
		finalCPUTimes := make(map[uint32]float64)
		for _, id := range domainIDs {
			dom, err := conn.LookupDomainById(id)
			if err != nil {
				log.Printf("Failed to lookup domain by ID %d: %v", id, err)
				continue
			}
			defer dom.Free()

			info, err := dom.GetInfo()
			if err != nil {
				name, err := dom.GetName()
				if err != nil {
					log.Printf("Failed to get info or name for domain ID %d: %v", id, err)
				} else {
					log.Printf("Failed to get info for domain %s: %v", name, err)
				}
				continue
			}

			finalCPUTimes[id] = float64(info.CpuTime)
		}

		// Calculate CPU usage percentages
		for _, id := range domainIDs {
			dom, err := conn.LookupDomainById(id)
			if err != nil {
				log.Printf("Failed to lookup domain by ID %d: %v", id, err)
				continue
			}

			info, err := dom.GetInfo()
			if err != nil {
				name, err := dom.GetName()
				if err != nil {
					log.Printf("Failed to get info or name for domain ID %d: %v", id, err)
				} else {
					log.Printf("Failed to get info for domain %s: %v", name, err)
				}
				continue
			}

			cpuTimeDifference := finalCPUTimes[id] - initialCPUTimes[id]
			numVCPUs := float64(info.NrVirtCpu)
			cpuUsagePercentage := (cpuTimeDifference / (numVCPUs * nanosecondsPerSecond)) * 100

			// Get memory stats (for display)
			var unusedMemory int
			memStats, err := dom.MemoryStats(11, 0)
			if err != nil {
				name, err := dom.GetName()
				if err != nil {
					log.Printf("Failed to get memory stats for domain ID %d: %v", id, err)
				} else {
					log.Printf("Failed to get memory stats for domain %s: %v", name, err)
				}
			} else {
				for _, stat := range memStats {
					switch libvirt.DomainMemoryStatTags(stat.Tag) {
					case libvirt.DOMAIN_MEMORY_STAT_UNUSED:
						unusedMemory = int(stat.Val)
					}
				}
			}

			name, err := dom.GetName()
			if err != nil {
				log.Printf("Failed to get name for domain ID %d: %v", id, err)
				name = "Unknown"
			}

			memoryUsedPercentage := 0.0
			if info.MaxMem > 0 {
				memoryUsedPercentage = float64(info.MaxMem-uint64(unusedMemory)) / float64(info.MaxMem) * 100
			}

			maxMemoryGB := float64(info.MaxMem) / (1024 * 1024)
			unusedMemoryGB := float64(unusedMemory) / (1024 * 1024)
			fmt.Fprintf(w, "<tr><td>%s</td><td>%v</td><td>%.2f%%</td><td>%d</td><td>%.2f%%</td><td>%.2f</td><td>%.2f</td></tr>", name, info.State, memoryUsedPercentage, info.NrVirtCpu, cpuUsagePercentage, unusedMemoryGB, maxMemoryGB)

		}
		fmt.Fprintf(w, "</table>")
	})

	port := 8080
	log.Printf("Server listening on port %d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
