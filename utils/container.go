package utils

import (
	"github.com/wave2588/go-base/utils/internal/cgroup"
	"runtime"
)

var totalCPU = cgroup.TotalCPU()

func init() {
	runtime.GOMAXPROCS(totalCPU)
}

func TotalCPU() int {
	return totalCPU
}
