build:
	docker build --tag gateway-monitor .

start-system:
	docker-compose --env-file test_influxdb.env up

stop-system:
	docker-compose down -v

start-influxdb-container:
	docker run -p 8086:8086 -v myInfluxVolume:/var/lib/influxdb2 influxdb:latest

start-local-influxdb:
	influxd

remove-local-influxdb:
	rm -rf ~/.influxdbv2/

start-local-influxdb-with-brew:
	brew services stop influxd

remove-local-homebrew-influxdb:
	rm -rf /opt/homebrew/etc/influxdb2