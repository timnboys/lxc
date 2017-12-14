package lxc

import "github.com/virtmonitor/driver"

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

//Close Close driver
func (l *LXC) Close() {
	return
}
