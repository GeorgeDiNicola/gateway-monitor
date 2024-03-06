# get the docker-compose.yml from github
curl -o docker-compose.yml https://github.com/GeorgeDiNicola/gateway-monitor/blob/main/docker-compose.yml

# empty current env file
echo "" > test_influxdb.env

# set up ENV
echo "INFLUXDB_USERNAME=$INFLUXDB_USERNAME" >> test_influxdb.env
echo "INFLUXDB_PASSWORD=$INFLUXDB_PASSWORD" >> test_influxdb.env
echo "INFLUXDB_ORG=$INFLUXDB_ORG" >> test_influxdb.env
echo "INFLUXDB_ORG=gateway_metrics" >> test_influxdb.env
echo "INFLUXDB_ORG=$INFLUXDB_TOKEN" >> test_influxdb.env
echo "INFLUXDB_ORG=http://localhost:8086" >> test_influxdb.env

# startup the app
docker-compose --env-file test_influxdb.env up