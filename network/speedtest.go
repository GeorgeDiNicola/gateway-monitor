package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
)

type SpeedTest interface {
	CollectSpeedMetrics() (SpeedTestData, error)
}

// capture and report upload and download speed
func CollectSpeedMetrics() (SpeedTestData, error) {
	var data SpeedTestData

	cmd := exec.Command("speedtest-cli", "--simple")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return data, err
	}

	output := out.String()

	parsedDownloadSpeed := strings.TrimSpace(strings.TrimSuffix(parseSpeed(output, "Download"), " Mbps"))
	downloadSpeed, err := strconv.ParseFloat(parsedDownloadSpeed, 64)
	if err != nil {
		log.Error("could not parse download speed")
	}
	data.DownloadSpeed = downloadSpeed

	parsedUploadSpeed := strings.TrimSpace(strings.TrimSuffix(parseSpeed(output, "Upload"), " Mbps"))
	uploadSpeed, err := strconv.ParseFloat(parsedUploadSpeed, 64)
	if err != nil {
		log.Error("could not parse upload speed")
	}
	data.UploadSpeed = uploadSpeed

	return data, nil
}

func parseSpeed(output, speedType string) string {
	re := regexp.MustCompile(fmt.Sprintf("%s: ([0-9\\.]+) Mbit/s", speedType))
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1] + " Mbps"
	}
	return "Not found"
}
