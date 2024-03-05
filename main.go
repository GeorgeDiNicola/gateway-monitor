package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/georgedinicola/gateway-monitor/network"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api/write"
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
	for {
		log.Out = os.Stdout
		log.Formatter = &logrus.JSONFormatter{}

		// TODO: add config website name or default gateway option
		gatewayIPAddr, err := getGatewayIpAddress()
		if err != nil {
			log.Fatalf("Could not get gateway IP address: %v", err)
			return
		}
		fmt.Println("Gateway IP:", gatewayIPAddr)
		fmt.Println("Measurements are classified into groups: Excellent, Good, Fair, Poor/Weak")

		var wg sync.WaitGroup

		type Result struct {
			Result interface{}
			Error  error
			ID     string
		}
		resultsChannel := make(chan Result, 3)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if result, err := network.GetGatewaySignalStrength(); err != nil {
				resultsChannel <- Result{Error: err, ID: "GetGatewaySignalStrength"}
			} else {
				resultsChannel <- Result{Result: result, ID: "GetGatewaySignalStrength"}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if result, err := network.CollectSpeedMetrics(); err != nil {
				resultsChannel <- Result{Error: err, ID: "CollectSpeedMetrics"}
			} else {
				resultsChannel <- Result{Result: result, ID: "CollectSpeedMetrics"}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if result, err := network.PingGatewayForStats(gatewayIPAddr, 25, 1, 60); err != nil {
				resultsChannel <- Result{Error: err, ID: "PingGatewayForStats"}
			} else {
				resultsChannel <- Result{Result: result, ID: "PingGatewayForStats"}
			}
		}()

		go func() {
			wg.Wait()
			close(resultsChannel)
		}()

		token := os.Getenv("INFLUXDB_TOKEN")
		url := os.Getenv("INFLUXDB_URL")
		client := influxdb2.NewClient(url, token)
		org := os.Getenv("INFLUXDB_ORG")
		bucket := os.Getenv("INFLUXDB_BUCKET")
		measurement := "network_metrics"

		writeAPI := client.WriteAPI(org, bucket)

		// process results as they arrive
		for result := range resultsChannel {
			if result.Error != nil {
				log.Errorf("%s returned an error: %v\n", result.ID, result.Error)
				continue
			}

			switch result.ID {
			case "GetGatewaySignalStrength":
				signalData := result.Result.(network.SignalData)
				fmt.Printf("Gateway Signal Strength: %v (%v)\n", signalData.SignalStrength, signalData.SignalStrengthClassification)

				// write to InfluxDB
				tags := map[string]string{
					"collection_type": "signal",
				}
				fields := map[string]interface{}{
					"signal_strength": signalData.SignalStrength,
				}
				point := write.NewPoint(measurement, tags, fields, time.Now())
				writeAPI.WritePoint(point)
			case "CollectSpeedMetrics":
				speedTestData := result.Result.(network.SpeedTestData)
				fmt.Printf("Download Speed: %v\n", speedTestData.DownloadSpeed)
				fmt.Printf("Upload Speed: %v\n", speedTestData.UploadSpeed)

				// write to InfluxDB
				tags := map[string]string{
					"collection_type": "ookla",
				}
				fields := map[string]interface{}{
					"download_speed": speedTestData.DownloadSpeed,
					"upload_speed":   speedTestData.UploadSpeed,
				}
				point := write.NewPoint(measurement, tags, fields, time.Now())
				writeAPI.WritePoint(point)
			case "PingGatewayForStats":
				pingData := result.Result.(network.PingData)
				fmt.Printf("Packet Loss percentage: %v\n", pingData.PacketLossPercentage)
				fmt.Printf("Round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
					pingData.MinRtt, pingData.AvgRtt, pingData.MaxRtt, pingData.StdDevRtt)

				jitter := float64(pingData.Jitter) / float64(time.Millisecond)

				// write to InfluxDB
				tags := map[string]string{
					"collection_type": "ping",
				}

				fields := map[string]interface{}{
					"packet_loss_percentage": pingData.PacketLossPercentage,
					"jitter":                 jitter,
					"min_packet_latency":     pingData.MinRtt,
					"avg_packet_latency":     pingData.AvgRtt,
					"max_packet_latency":     pingData.AvgRtt,
					"stddev_packet_latency":  pingData.StdDevRtt,
				}
				point := write.NewPoint(measurement, tags, fields, time.Now())
				writeAPI.WritePoint(point)
			}
			time.Sleep(60 * time.Second)
		}

	}
}
