package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
	"errors"
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
func getBestNodeName(nodes []Node) (string, error) {
	var nodeNames []string
	for _, n := range nodes {
		nodeNames = append(nodeNames, n.Metadata.Name)
	}

	// Execute a query over the HTTP API to get the metric node_memory_MemAvailable
	resp, err := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemAvailable")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Decode the JSON body of the HTTP response into a struct
	var metrics MetricResponse
	decodeJsonDataToStruct(&metrics, resp)


        // Iterate through the metric results to find the node with the best value
	max := 0
	bestNode := ""
	for _, m := range metrics.Data.Results {
		// Print metric value for the node
		fmt.Printf("Node name: %s\n", m.MetricInfo["instance"])
		fmt.Printf("Value: %s\n\n", m.MetricValue[1])

		// Convert string in metric results to an integer
		metricValue, err := strconv.Atoi(m.MetricValue[1].(string))
		if err != nil {
			return "", err
		}

		if metricValue > max {
			// Check if the node is in the list passed in (nodes the pod will fit on)
			available := nodeAvailable(nodeNames, m.MetricInfo["instance"])
			if available == true {
				max = metricValue
				bestNode = m.MetricInfo["instance"]
			}
		}
	}
	if bestNode == "" {
		return "", errors.New("No node found")
	} else {
		return bestNode, nil
	}
}


func nodeAvailable(nodeNames []string, name string) (result bool) {
	for _, n := range nodeNames {
		if name == n {
			return true
		}
	}
	return false
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

