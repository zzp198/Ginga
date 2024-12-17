package stateyemp

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"time"
)

var state State

var lastStatus *State

type State struct {
	T   time.Time
	Cpu float64
	Mem struct {
		Current uint64
		Total   uint64
	}
	Swap struct {
		Current uint64
		Total   uint64
	}
	Disk struct {
		Current uint64
		Total   uint64
	}
	Uptime   uint64
	Loads    []float64
	TcpCount int
	UdpCount int
	NetIO    struct {
		Up   uint64
		Down uint64
	}
	NetTraffic struct {
		Sent uint64
		Recv uint64
	}
}

func (*State) Get() State {

	state := State{
		T: time.Now(),
	}

	percents, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("get cpu percent failed:", err)
	} else {
		state.Cpu = percents[0]
	}
	upTime, err := host.Uptime()
	if err != nil {
		log.Println("get uptime failed:", err)
	} else {
		state.Uptime = upTime
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Println("get virtual memory failed:", err)
	} else {
		state.Mem.Current = memInfo.Used
		state.Mem.Total = memInfo.Total
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		log.Println("get swap memory failed:", err)
	} else {
		state.Swap.Current = swapInfo.Used
		state.Swap.Total = swapInfo.Total
	}

	distInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("get dist usage failed:", err)
	} else {
		state.Disk.Current = distInfo.Used
		state.Disk.Total = distInfo.Total
	}

	avgState, err := load.Avg()
	if err != nil {
		log.Println("get load avg failed:", err)
	} else {
		state.Loads = []float64{avgState.Load1, avgState.Load5, avgState.Load15}
	}

	ioStats, err := net.IOCounters(false)
	if err != nil {
		log.Println("get io counters failed:", err)
	} else if len(ioStats) > 0 {
		ioStat := ioStats[0]
		state.NetTraffic.Sent = ioStat.BytesSent
		state.NetTraffic.Recv = ioStat.BytesRecv

		if lastStatus != nil {
			duration := state.T.Sub(lastStatus.T)
			seconds := float64(duration) / float64(time.Second)
			up := uint64(float64(state.NetTraffic.Sent-lastStatus.NetTraffic.Sent) / seconds)
			down := uint64(float64(state.NetTraffic.Recv-lastStatus.NetTraffic.Recv) / seconds)
			state.NetIO.Up = up
			state.NetIO.Down = down
		}
	} else {
		log.Println("can not find io counters")
	}

	tcpStat, err := net.Connections("tcp")
	if err != nil {
		log.Println("get tcp io counters failed:", err)
	} else {
		state.TcpCount = len(tcpStat)
	}

	udpStat, err := net.Connections("udp")
	if err != nil {
		log.Println("get udp io counters failed:", err)
	} else {
		state.UdpCount = len(udpStat)
	}

	lastStatus = &state

	return state
}
