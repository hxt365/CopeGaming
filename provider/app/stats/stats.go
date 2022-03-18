package stats

import (
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SysInfo struct {
	HostName string
	Platform string
	CpuName  string
	CpuNum   int
	MemSize  float64
}

func GetSysInfo() (*SysInfo, error) {
	var s SysInfo

	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	s.HostName = hostInfo.Hostname
	s.Platform = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	s.CpuName = cpuInfo[0].ModelName
	s.CpuNum, err = cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	s.MemSize = math.Ceil(float64(memInfo.Total) / (1024 * 1024 * 1024))

	return &s, nil
}

type SysStats struct {
	CpuPercent float64
	MemPercent float64
}

func GetSysStats(interval time.Duration) (*SysStats, error) {
	var s SysStats

	cpuUsage, err := cpu.Percent(interval, false)
	if err != nil {
		return nil, err
	}
	s.CpuPercent = cpuUsage[0]

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	s.MemPercent = memInfo.UsedPercent

	return &s, nil
}
