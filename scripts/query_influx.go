package main

import (
	"context"
	"fmt"
	"log"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

func main() {

	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	org := "georgedinicola"

	queryAPI := client.QueryAPI(org)
	query := `from(bucket: "gateway_metrics")
            |> range(start: -20m)
            |> filter(fn: (r) => r._measurement == "network_metrics")`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
}
