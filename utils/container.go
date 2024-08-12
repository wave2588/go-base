package utils

import (
	"github.com/wave2588/go-base/utils/internal/cgroup"
	"runtime"
)

var totalCPU, totalMemory = cgroup.TotalCPU(), cgroup.TotalMemory()

func init() {
	runtime.GOMAXPROCS(totalCPU)
}

func TotalCPU() int {
	return totalCPU
}

func TotalMemory() int {
	return totalMemory
}
