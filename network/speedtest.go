package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
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
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("speedtest-cli", "--simple")
	case "linux":
		cmd = exec.Command("speedtest", "--accept-license", "--accept-gdpr")
	default:
		return data, fmt.Errorf("unsupported platform")
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return data, err
	}

	output := out.String()

	var parsedDownloadSpeed string
	var parsedUploadSpeed string

	switch runtime.GOOS {
	case "darwin":
		parsedDownloadSpeed = strings.TrimSpace(strings.TrimSuffix(parseSpeedMacOs(output, "Download"), " Mbps"))
		parsedUploadSpeed = strings.TrimSpace(strings.TrimSuffix(parseSpeedMacOs(output, "Upload"), " Mbps"))
	case "linux":
		parsedDownloadSpeed = strings.TrimSpace(strings.TrimSuffix(parseSpeedLinux(output, "Download"), " Mbps"))
		parsedUploadSpeed = strings.TrimSpace(strings.TrimSuffix(parseSpeedLinux(output, "Upload"), " Mbps"))
	default:
		return data, fmt.Errorf("unsupported platform")
	}

	downloadSpeed, err := strconv.ParseFloat(parsedDownloadSpeed, 64)
	if err != nil {
		log.Error("could not parse download speed")
	}
	data.DownloadSpeed = downloadSpeed

	uploadSpeed, err := strconv.ParseFloat(parsedUploadSpeed, 64)
	if err != nil {
		log.Error("could not parse upload speed")
	}
	data.UploadSpeed = uploadSpeed

	return data, nil
}

func parseSpeedMacOs(output, speedType string) string {
	re := regexp.MustCompile(fmt.Sprintf("%s: ([0-9\\.]+) Mbit/s", speedType))
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1] + " Mbps"
	}
	return "Not found"
}

func parseSpeedLinux(output, speedType string) string {
	var re *regexp.Regexp
	if speedType == "Download" {
		re = regexp.MustCompile(`Download:\s+(\d+\.\d+)\s+Mbps`)
	} else {
		re = regexp.MustCompile(`Upload:\s+(\d+\.\d+)\s+Mbps`)
	}

	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1] + " Mbps"
	}
	return "Not found"
}
