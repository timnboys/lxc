// +build lxc all

package lxc

import "github.com/virtmonitor/driver"

const (
	//LXC_PATH Path to LXC
	LXC_PATH = "/usr/sbin/vzlist"
)

//LXC LXC struct
type LXC struct {
	driver.Driver
}

func init() {
	driver.RegisterDriver(&LXC{})
}

//Name Return driver name
func (l *LXC) Name() string {
	return "LXC"
}

func (l *LXC) Stop() {
	return
}
