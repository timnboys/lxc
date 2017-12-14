// +build windows

package lxc

import (
	"fmt"

	"github.com/virtmonitor/driver"
)

// Detect Detect dependencies
func (l *LXC) Detect() bool {
	return false
}

// Collect Collect domain statistics
func (l *LXC) Collect(cpu bool, block bool, network bool) (domains map[driver.DomainID]*driver.Domain, err error) {
	domains = make(map[driver.DomainID]*driver.Domain)
	err = fmt.Errorf("LXC not supported on this platform")
	return
}
