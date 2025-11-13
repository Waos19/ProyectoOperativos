package monitor

import (
	"fmt"
	"strings"
	"time"

	"proyoper/internal/shell"
)

func GetCPUUsage() (string, error) {
	out, err := shell.RunCommand("top -bn1", "")
	if err != nil {
		return "", fmt.Errorf("error CPU: %w", err)
	}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu(s):") {
			line = strings.ReplaceAll(line, ",", ".")
			fields := strings.Fields(line)
			us := toFloat(fields[1])
			sy := toFloat(fields[3])
			return fmt.Sprintf("%.1f", us+sy), nil
		}
	}
	return "", fmt.Errorf("no se encontró línea de CPU")
}

func GetMemoryUsage() (used string, total string, percent string, err error) {
	out, err := shell.RunCommand("free -m", "")
	if err != nil {
		return "", "", "", fmt.Errorf("error memoria: %w", err)
	}
	lines := strings.Split(out, "\n")
	fields := strings.Fields(lines[1])
	used = fields[2]
	total = fields[1]
	percent = fmt.Sprintf("%.1f", (toFloat(used)/toFloat(total))*100)
	return
}

func GetDiskUsage() (used string, total string, percent string, err error) {
	out, err := shell.RunCommand("df -h /", "")
	if err != nil {
		return "", "", "", fmt.Errorf("error disco: %w", err)
	}
	lines := strings.Split(out, "\n")
	fields := strings.Fields(lines[1])
	used = fields[2]
	total = fields[1]
	percent = fields[4]
	return
}

func GetProcessList() ([]string, error) {
	out, err := shell.RunCommand("ps -eo comm", "")
	if err != nil {
		return nil, fmt.Errorf("error procesos: %w", err)
	}
	lines := strings.Split(out, "\n")
	var names []string
	for _, line := range lines[1:] {
		name := strings.TrimSpace(line)
		if name != "" {
			names = append(names, name)
		}
		if len(names) >= 10 {
			break
		}
	}
	return names, nil
}

func toFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func GenerateReport() (string, error) {
	cpu, err := GetCPUUsage()
	if err != nil {
		return "", err
	}
	memUsed, memTotal, memPercent, err := GetMemoryUsage()
	if err != nil {
		return "", err
	}
	diskUsed, diskTotal, diskPercent, err := GetDiskUsage()
	if err != nil {
		return "", err
	}
	procs, err := GetProcessList()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("\n╔════════════════════════════════════════════════════╗\n")
	b.WriteString("║                 REPORTE DE SISTEMA                ║\n")
	b.WriteString("╠════════════════════════════════════════════════════╣\n")
	b.WriteString(fmt.Sprintf("║  Fecha:     %-38s ║\n", time.Now().Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf("║  CPU:       %-6s%%                                  ║\n", cpu))
	b.WriteString(fmt.Sprintf("║  Memoria:   %s / %s MB (%s%% usada)         ║\n", memUsed, memTotal, memPercent))
	b.WriteString(fmt.Sprintf("║  Disco:     %s / %s (%s usada)              ║\n", diskUsed, diskTotal, diskPercent))
	b.WriteString("╠════════════════════════════════════════════════════╣\n")
	b.WriteString("║  Procesos activos:                               ║\n")
	for _, name := range procs {
		b.WriteString(fmt.Sprintf("║   • %-44s ║\n", name))
	}
	b.WriteString("╚════════════════════════════════════════════════════╝\n")

	return b.String(), nil
}
