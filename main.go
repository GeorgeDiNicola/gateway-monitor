package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/georgedinicola/gateway-monitor/network"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func getGatewayIpAddress() (string, error) {
	cmd := exec.Command("sh", "-c", "netstat -nr | grep default | awk '{print $2}' | head -n 1")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return "", err
	}

	gatewayIPAddr := out.String()
	gatewayIPAddr = string(bytes.TrimSpace([]byte(gatewayIPAddr)))

	return gatewayIPAddr, nil
}

func main() {
	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{}

	fmt.Println("Measurements are classified into groups: Excellent, Good, Fair, Poor/Weak")

	gatewayIPAddr, err := getGatewayIpAddress()
	if err != nil {
		log.Fatalf("Could not get gateway IP address: %v", err)
		return
	}
	fmt.Println("Gateway IP:", gatewayIPAddr)

	// TODO: use goroutines for better performance
	signalStrength, err := network.GetGatewaySignalStrength()
	if err != nil {
		log.Errorf("error performing signal test: %v", err)
	}
	category := network.ClassifySignalStrength(signalStrength)
	fmt.Printf("Gateway Signal Strength: %v (%v)\n", signalStrength, category)

	downloadSpeed, uploadSpeed, err := network.CollectSpeedMetrics()
	if err != nil {
		log.Errorf("error performing speed test: %v", err)
	}
	fmt.Println("Download Speed:", downloadSpeed)
	fmt.Println("Upload Speed:", uploadSpeed)

	err = network.PingGatewayForStats(gatewayIPAddr, 25, 1, 60)
	if err != nil {
		log.Errorf("error performing ping test: %v", err)
	}
}
