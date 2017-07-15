package Background

import (
	"fmt"
	"math"
	"time"

	libvirt "github.com/libvirt/libvirt-go"
)

type vT struct {
	name string
	time int64
	cpu  uint64
	Load float64
	Mem  uint64
}

var VTS = make(map[string]vT)

func getCPUUsage() {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		fmt.Println(err.Error())
	}else{
		doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
		if err != nil {
			fmt.Println(err.Error())
		}else{
			for _, dom := range doms {
				name, _ := dom.GetName()
				cpustats, cErr := dom.GetCPUStats(-1, 1, 0)
				memstats, mErr := dom.MemoryStats(1, 0)
				if (mErr == nil) && (cErr == nil) {
					tm := time.Now().Unix()
					cputime := cpustats[0].CpuTime - VTS[name].cpu
					load := float64(cputime*100.0) / float64(((tm - VTS[name].time) * 1000.0 * 1000.0 * 1000.0))
					per := math.Max(0.0, math.Min(100.0, load))
					if len(memstats) > 1 {
						VTS[name] = vT{name: name, time: tm, cpu: (cpustats[0].CpuTime), Load: per, Mem: memstats[1].Val}
					} else {
						VTS[name] = vT{name: name, time: tm, cpu: (cpustats[0].CpuTime), Load: per, Mem: 0}
					}
				}
			}
		}
		conn.Close()
	}
		//fmt.Println(VTS)
}

func LibConnect() (conn *libvirt.Connect) {
	//conx, err := libvirt.NewConnect("qemu:///system")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	return CONN
}

var CONN *libvirt.Connect

func StartResourceCollector() {
	CONN, _ = libvirt.NewConnect("qemu:///system")
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				getCPUUsage()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
