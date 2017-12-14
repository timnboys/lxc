// +build !windows

package lxc

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Detect Detect dependencies
func (l *LXC) Detect() bool {
	var err error

	var f *os.File
	if f, err = os.Open("/proc/self/cgroup"); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("Error detecting LXC driver: %v", err)
		return false
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "lxc") {
			return true
		}
	}

	return false
}
