package status

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type ICMPStatus struct {
	DefaultInterval time.Duration
}

func NewICMPStatus() *ICMPStatus {
	return &ICMPStatus{
		DefaultInterval: time.Millisecond * 100,
	}
}

func (icmp *ICMPStatus) CheckStatus(host string, count int) (*HostStatus, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return nil, err
	}
	pinger.Count = count
	pinger.Interval = icmp.DefaultInterval

	if err = pinger.Run(); err != nil {
		return nil, err
	}
	stats := pinger.Statistics()

	return &HostStatus{
		PacketLoss: stats.PacketLoss,
		Rtts:       stats.Rtts,
		AvgRtt:     stats.AvgRtt,
		MinRtt:     stats.MinRtt,
		MaxRtt:     stats.MinRtt,
	}, nil
}
