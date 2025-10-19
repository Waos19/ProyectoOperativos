package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

func GenerateReport() (string, error) {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		return "", err
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}
	procs, _ := process.Processes()
	var procNames []string
	for _, p := range procs {
		name, _ := p.Name()
		if name != "" {
			procNames = append(procNames, name)
		}
		if len(procNames) >= 10 {
			break
		}
	}

	var b strings.Builder
	b.WriteString("\n╔════════════════════════════════════════════════════╗\n")
	b.WriteString("║                 REPORTE DE SISTEMA                ║\n")
	b.WriteString("╠════════════════════════════════════════════════════╣\n")
	b.WriteString(fmt.Sprintf("║ 	Fecha:     %-38s ║\n", time.Now().Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("║ 	CPU:       %-6.2f%%                                  ║\n", cpuPercents[0]))
	b.WriteString(fmt.Sprintf("║ 	Memoria:   %.1f / %.1f GB (%.1f%% usada)         ║\n",
		float64(vmStat.Used)/1e9,
		float64(vmStat.Total)/1e9,
		vmStat.UsedPercent))
	b.WriteString(fmt.Sprintf("║ 	Disco:     %.1f / %.1f GB (%.1f%% usada)         ║\n",
		float64(diskStat.Used)/1e9,
		float64(diskStat.Total)/1e9,
		diskStat.UsedPercent))
	b.WriteString("╠════════════════════════════════════════════════════╣\n")
	b.WriteString("║ 	Procesos activos:                               ║\n")
	for _, name := range procNames {
		b.WriteString(fmt.Sprintf("║   • %-44s ║\n", name))
	}
	b.WriteString("╚════════════════════════════════════════════════════╝\n")

	return b.String(), nil
}
