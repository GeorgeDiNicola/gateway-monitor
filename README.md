# gateway-monitor
A diagnostic monitoring tool that collects metrics about a home internet gateway's performance. Useful for measuring internet speed over time, detecting performance degradation, being notified of early signs of service outages, etc.

## Statistics Reported
- Gateway Signal Strenth (dBm - decibels relative to a milliwatt)
- Download Speed (Mbps)
- Upload Speed (Mbps)
- Packet Loss/Error Rate (% of packets transmitted & received)
- Round-trip packet traversal time (ms) min/avg/max/stddev
- Jitter (average ms) - the variability in the latency of packets over a gateway network

## Supported Operating Systems
macOS, Linux

## Dashboard of Metrics Streamed to the InfluxDB
![Dashboard](influx_dashboard.png)