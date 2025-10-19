package server

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type Resources struct {
	Timestamp   time.Time
	CPUPercent  float64
	MemoryUsed  uint64
	MemoryTotal uint64
	DiskUsed    uint64
	DiskTotal   uint64
}

func GetStats() (*Resources, error) {

	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	stats := &Resources{
		Timestamp:   time.Now(),
		CPUPercent:  cpuPercents[0],
		MemoryUsed:  vmStat.Used,
		MemoryTotal: vmStat.Total,
		DiskUsed:    diskStat.Used,
		DiskTotal:   diskStat.Total,
	}

	return stats, nil
}

func ShowStats(interval int) {
	for {
		cpuPercents, _ := cpu.Percent(0, false)
		vmStat, _ := mem.VirtualMemory()
		diskStat, _ := disk.Usage("/")

		fmt.Printf("\r[MONITOR] CPU: %.1f%% | Mem: %.1f%% | Disk: %.1f%%\n",
			cpuPercents[0],
			vmStat.UsedPercent,
			diskStat.UsedPercent)

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
