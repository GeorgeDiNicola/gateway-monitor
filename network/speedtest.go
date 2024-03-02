package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

type SpeedTest interface {
	CollectSpeedMetrics() (downloadSpeed, uploadSpeed string, err error)
}

// capture and report upload and download speed
func CollectSpeedMetrics() (downloadSpeed, uploadSpeed string, err error) {
	cmd := exec.Command("speedtest-cli", "--simple")
	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return "", "", err
	}

	output := out.String()
	downloadSpeed = parseSpeed(output, "Download")
	uploadSpeed = parseSpeed(output, "Upload")

	return downloadSpeed, uploadSpeed, nil
}

func parseSpeed(output, speedType string) string {
	re := regexp.MustCompile(fmt.Sprintf("%s: ([0-9\\.]+) Mbit/s", speedType))
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1] + " Mbps"
	}
	return "Not found"
}
