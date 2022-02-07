// +build !windows

package lxc

import (
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strings"
	"time"

	"github.com/virtmonitor/driver"
	"github.com/virtmonitor/virNetTap"
	"github.com/lxc/go-lxc"
)

// Collect Collect domain statistics
func (l *LXC) Collect(cpu bool, block bool, network bool) (domains map[driver.DomainID]*driver.Domain, err error) {
	domains = make(map[driver.DomainID]*driver.Domain)

	containers := lxc.ActiveContainers(lxc.DefaultConfigPath())

	if len(containers) <= 0 {
		return
	}

	var procnet virNetTap.VirNetTap
	nstat := make(map[string]virNetTap.InterfaceStats)

	if nstat, err = procnet.GetAllVifStats(); err != nil {
		return
	}

	for _, container := range containers {
		domain := &driver.Domain{}

		domain.ID = driver.DomainID(container.InitPid())
		domain.Name = container.Name()

		if cpu {

			var cpus map[int]time.Duration
			if cpus, err = container.CPUTimePerCPU(); err != nil {
				log.Printf("Error collecting CPU stats for LXC container %s: %v", domain.Name, err)
			} else {
				for id, ti := range cpus {
					var cpu driver.CPU
					cpu.ID = uint64(id)
					cpu.Time = float64(int(ti)) / float64(1000000000)

					domain.Cpus = append(domain.Cpus, cpu)
				}
			}

		}

		if block {
			//TODO: parse out per-device block stats via cgroups blkio
		}

		if network {

			var pfx = "lxc.net"
			var iftype, ifname []string
			if !lxc.VersionAtLeast(2, 1, 0) {
				pfx = "lxc.network"
			}

			for i := 0; i < len(container.ConfigItem(pfx)); i++ {
				if iftype = container.RunningConfigItem(fmt.Sprintf("%s.%d.type", pfx, i)); iftype == nil {
					continue
				}

				if strings.ToUpper(iftype[0]) == "VETH" {
					ifname = container.RunningConfigItem(fmt.Sprintf("%s.%d.veth.pair", pfx, i))
				} else {
					ifname = container.RunningConfigItem(fmt.Sprintf("%s.%d.link", pfx, i))
				}

				if ifname == nil {
					continue
				}

				var dnetwork driver.NetworkInterface
				dnetwork.Name = ifname[0]

				var bridges []string
				if bridges, err = filepath.Glob(fmt.Sprintf("/sys/class/net/%s/upper_*", ifname[0])); err == nil && len(bridges) > 0 {
					for _, bridge := range bridges {
						dnetwork.Bridges = append(dnetwork.Bridges, strings.ToLower(strings.TrimPrefix(filepath.Base(bridge), "upper_")))
					}
				}

				var mac []string
				if mac = container.RunningConfigItem(fmt.Sprintf("%s.%d.hwaddr", pfx, i)); mac != nil {
					dnetwork.Mac, _ = net.ParseMAC(mac[0])
				}

				if vifstat, ok := nstat[strings.ToLower(ifname[0])]; ok {
					//we found statistics for interface
					dnetwork.RX = driver.NetworkIO{Bytes: vifstat.IN.Bytes, Packets: vifstat.IN.Pkts, Errors: vifstat.IN.Errs, Drops: vifstat.IN.Drops}
					dnetwork.TX = driver.NetworkIO{Bytes: vifstat.OUT.Bytes, Packets: vifstat.OUT.Pkts, Errors: vifstat.OUT.Errs, Drops: vifstat.OUT.Drops}
				}

				domain.Interfaces = append(domain.Interfaces, dnetwork)
			}
		}

		domains[domain.ID] = domain
	}
	return
}
