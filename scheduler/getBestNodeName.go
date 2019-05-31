package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
)

// Struct for decoded JSON from HTTP response
type MetricResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Results []Result `json:"result"`
}

// Idea to use interface for metric values (which have different types) from 
// https://stackoverflow.com/questions/38861295/how-to-parse-json-arrays-with-two-different-data-types-into-a-struct-in-go-lang
type Result struct {
	MetricInfo map[string]string  `json:"metric"`
	MetricValue []interface{} `json:"value"` //Index 0 is unix_time, index 1 is sample_value (metric value)
}

// Returns the name of the node with the best metric value
func getBestNodeName() (string, error) {
	// Execute a query over the HTTP API to get the metric node_memory_MemAvailable
	resp, err := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemAvailable")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Decode the JSON body of the HTTP response into a struct
	var metrics MetricResponse
	decodeJsonDataToStruct(&metrics, resp)

        // Iterate through the nodes to find the best value
	max, err := strconv.Atoi(metrics.Data.Results[0].MetricValue[1].(string))
	bestNode := metrics.Data.Results[0].MetricInfo["instance"]
	for i:=0; i<len(metrics.Data.Results); i++ {
		// Print metric value for each node
		fmt.Printf("Node name: %s\n", metrics.Data.Results[i].MetricInfo["instance"])
		fmt.Printf("Value: %s\n\n", metrics.Data.Results[i].MetricValue[1])

		metricValue, err := strconv.Atoi(metrics.Data.Results[i].MetricValue[1].(string))
		// To-do: figure out if more error handling is needed
		if err != nil {
			break
		}

		if metricValue > max {
			max = metricValue
			bestNode = metrics.Data.Results[i].MetricInfo["instance"]
		}
	}
	return bestNode, nil
}

// Decode JSON data into a struct to get the metric values
func decodeJsonDataToStruct(metrics *MetricResponse, resp *http.Response) {
        decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
