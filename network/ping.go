package network

import (
	"fmt"
	"math"
	"time"

	"github.com/go-ping/ping"
)

type Ping interface {
	PingGatewayForStats(gatewayIpAddr string, pingCount int, pingInterval time.Duration, pingTimeout time.Duration) error
}

func PingGatewayForStats(gatewayIpAddr string, pingCount int, pingInterval time.Duration, pingTimeout time.Duration) error {
	pinger, err := ping.NewPinger(gatewayIpAddr)
	if err != nil {
		return err
	}

	pinger.Count = pingCount
	interval := time.Second * pingInterval
	timeout := time.Second * pingTimeout
	if interval > timeout {
		return fmt.Errorf("set timeout higher than interval")
	}

	pinger.Interval = interval
	pinger.Timeout = timeout

	var rtts []time.Duration
	pinger.OnRecv = func(pkt *ping.Packet) {
		rtts = append(rtts, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("Round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		jitter, jitterClassification := calculateJitter(rtts)
		fmt.Printf("Jitter: %vms (%v)\n", jitter, jitterClassification)
	}

	if err = pinger.Run(); err != nil {
		fmt.Printf("Failed to ping target host: %s\n", err)
	}

	return nil
}

func calculateJitter(rtts []time.Duration) (float64, string) {
	if len(rtts) < 2 {
		return 0, ""
	}

	var jitterSum time.Duration
	for i := 1; i < len(rtts); i++ {
		diff := rtts[i] - rtts[i-1]
		jitterSum += time.Duration(math.Abs(float64(diff)))
	}

	averageJitter := jitterSum / time.Duration(len(rtts)-1)
	averageJitterMs := float64(averageJitter.Nanoseconds()) / 1e6 // Convert jitter to ms

	var jitterCategory string
	switch {
	case averageJitterMs < 0:
		jitterCategory = "Invalid"
	case averageJitterMs == 0:
		jitterCategory = "Perfect"
	case averageJitterMs > 0 && averageJitterMs <= 30:
		jitterCategory = "Excellent"
	case averageJitterMs > 30 && averageJitterMs <= 50:
		jitterCategory = "Good"
	case averageJitterMs > 50 && averageJitterMs <= 70:
		jitterCategory = "Fair"
	case averageJitterMs > 70 && averageJitterMs <= 100:
		jitterCategory = "Poor"
	case averageJitterMs > 100:
		jitterCategory = "Bad"
	default:
		jitterCategory = "Unknown"
	}

	return averageJitterMs, jitterCategory
}
