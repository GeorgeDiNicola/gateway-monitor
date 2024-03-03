package network

import (
	"fmt"
	"math"
	"time"

	"github.com/go-ping/ping"
	"github.com/labstack/gommon/log"
)

type Ping interface {
	PingGatewayForStats(gatewayIpAddr string, pingCount int, pingInterval time.Duration, pingTimeout time.Duration) (PingData, error)
	ClassifyJitter(avgJitterMs float64) string
}

func PingGatewayForStats(gatewayIpAddr string, pingCount int, pingInterval time.Duration, pingTimeout time.Duration) (PingData, error) {
	var data PingData

	pinger, err := ping.NewPinger(gatewayIpAddr)
	if err != nil {
		return data, err
	}

	pinger.Count = pingCount
	interval := time.Second * pingInterval
	timeout := time.Second * pingTimeout
	if interval > timeout {
		return data, fmt.Errorf("set timeout higher than interval")
	}

	pinger.Interval = interval
	pinger.Timeout = timeout

	var rtts []time.Duration
	pinger.OnRecv = func(pkt *ping.Packet) {
		rtts = append(rtts, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		data.PacketLossPercentage = float64(stats.PacketLoss)
		// Convert time.Duration to float64 representing milliseconds
		data.MinRtt = float64(stats.MinRtt) / float64(time.Millisecond)
		data.AvgRtt = float64(stats.AvgRtt) / float64(time.Millisecond)
		data.MaxRtt = float64(stats.MaxRtt) / float64(time.Millisecond)
		data.StdDevRtt = float64(stats.StdDevRtt) / float64(time.Millisecond)
		jitter, err := calculateJitter(rtts)
		if err != nil {
			log.Error("error calculating jitter")
		}
		data.Jitter = jitter / float64(time.Millisecond)
		data.JitterClassification = ClassifyJitter(jitter)
	}

	if err = pinger.Run(); err != nil {
		fmt.Printf("Failed to ping target host: %s\n", err)
	}

	return data, nil
}

func calculateJitter(rtts []time.Duration) (float64, error) {
	if len(rtts) < 2 {
		return 0.0, nil
	}

	var jitterSum time.Duration
	for i := 1; i < len(rtts); i++ {
		diff := rtts[i] - rtts[i-1]
		jitterSum += time.Duration(math.Abs(float64(diff)))
	}

	averageJitter := jitterSum / time.Duration(len(rtts)-1)
	averageJitterMs := float64(averageJitter.Nanoseconds()) / 1e6 // Convert jitter to ms

	return averageJitterMs, nil
}

func ClassifyJitter(avgJitterMs float64) string {
	var jitterCategory string
	switch {
	case avgJitterMs < 0:
		jitterCategory = "Invalid"
	case avgJitterMs == 0:
		jitterCategory = "Perfect"
	case avgJitterMs > 0 && avgJitterMs <= 30:
		jitterCategory = "Excellent"
	case avgJitterMs > 30 && avgJitterMs <= 50:
		jitterCategory = "Good"
	case avgJitterMs > 50 && avgJitterMs <= 70:
		jitterCategory = "Fair"
	case avgJitterMs > 70 && avgJitterMs <= 100:
		jitterCategory = "Poor"
	case avgJitterMs > 100:
		jitterCategory = "Bad"
	default:
		jitterCategory = "Unknown"
	}

	return jitterCategory
}
