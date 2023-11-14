package status

import "time"

type HostStatus struct {
	PacketLoss float64
	Rtts       []time.Duration
	AvgRtt     time.Duration
	MinRtt     time.Duration
	MaxRtt     time.Duration
}
