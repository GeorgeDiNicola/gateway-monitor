package network

type PingData struct {
	PacketLossPercentage float64
	Jitter               float64 // Ms
	JitterClassification string
	MinRtt               float64 // Ms
	AvgRtt               float64 // Ms
	MaxRtt               float64 // Ms
	StdDevRtt            float64 // Ms
}

type SignalData struct {
	SignalStrength               float64 // dBm
	SignalStrengthClassification string
}

type SpeedTestData struct {
	DownloadSpeed float64 // Mbps
	UploadSpeed   float64 // Mbps
}
