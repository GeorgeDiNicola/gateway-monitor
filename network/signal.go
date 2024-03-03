package network

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Signal interface {
	GetGatewaySignalStrength() (SignalData, error)
	ClassifySignalStrength(signalStrength string) string
}

// measures the gateway' signal strength in dBm
func GetGatewaySignalStrength() (SignalData, error) {
	var data SignalData

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-I")
	case "linux":
		cmd = exec.Command("sh", "-c", "iwconfig 2>/dev/null | grep -i --color=never 'signal level'")
	default:
		return data, fmt.Errorf("unsupported platform")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return data, err
	}

	parsedSignalStrength, err := parseSignalStrength(out.String(), runtime.GOOS)
	if err != nil {
		return data, err
	}

	signalStrength, err := strconv.Atoi(strings.TrimSuffix(parsedSignalStrength, " dBm"))
	if err != nil {
		return data, err
	}
	classifiedSignalStrength := ClassifySignalStrength(signalStrength)

	data.SignalStrength = float64(signalStrength)
	data.SignalStrengthClassification = classifiedSignalStrength

	return data, nil
}

func parseSignalStrength(output, osType string) (string, error) {
	var re *regexp.Regexp
	var signalStrength string

	switch osType {
	case "darwin":
		re = regexp.MustCompile(`agrCtlRSSI:\s+(-?\d+)`)
		matches := re.FindStringSubmatch(output)
		if len(matches) >= 2 {
			signalStrength = matches[1] + " dBm"
		}
	case "linux":
		re = regexp.MustCompile(`Signal level=(-?\d+) dBm`)
		matches := re.FindStringSubmatch(output)
		if len(matches) >= 2 {
			signalStrength = matches[1] + " dBm"
		}
	default:
		return "", fmt.Errorf("unsupported platform for parsing")
	}

	if signalStrength == "" {
		return "", fmt.Errorf("signal strength not found")
	}

	return signalStrength, nil
}

func ClassifySignalStrength(signalStrength int) string {
	switch {
	case signalStrength >= -30:
		return "Excellent"
	case signalStrength >= -60:
		return "Good"
	case signalStrength >= -70:
		return "Fair"
	case signalStrength >= -80:
		return "Weak"
	default:
		return "Poor"
	}
}
