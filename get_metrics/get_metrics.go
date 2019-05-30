package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"io/ioutil"
)

// Struct for decoded JSON from HTTP response
type MetricResponse struct {
	//Status string `json:"status"`
	Data Data `json:"data"`
}

type Data struct {
	//ResultType string `json:"resultType"`
	Results []Result `json:"result"`
}

// Idea to use interface for metric values (which have different types) from 
// https://stackoverflow.com/questions/38861295/how-to-parse-json-arrays-with-two-different-data-types-into-a-struct-in-go-lang
type Result struct {
	MetricInfo map[string]string  `json:"metric"`
	MetricValue []interface{} `json:"value"` //Index 0 is unix_time, index 1 is sample_value
}

func main() {
	// Execute a query over the HTTP API to get the metric node_memory_MemTotal
	resp, err := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemTotal")
        //fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Decode the JSON body of the HTTP response into a struct
	var metrics MetricResponse
	decodeJsonDataToStruct(&metrics, resp)
	//fmt.Println(metrics)

	// Print the metric values
	for i:=0; i<len(metrics.Data.Results); i++ {
		fmt.Printf("Node IP: %s\n", metrics.Data.Results[i].MetricInfo["instance"])
		fmt.Printf("Pod name: %s\n", metrics.Data.Results[i].MetricInfo["kubernetes_pod_name"])
		fmt.Printf("Time: %f\n", metrics.Data.Results[i].MetricValue[0])
		fmt.Printf("Value: %s\n\n", metrics.Data.Results[i].MetricValue[1])
	}
}

func decodeArbitraryJsonData(resp *http.Response) {
	// Decode JSON data into an interface (don't need to know the structure of the JSON)

	// Used https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go to
	// figured out how to read the bytes from the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Used https://blog.golang.org/json-and-go to figure out how to read arbitrary json data
	var metrics interface{}
	err = json.Unmarshal(bodyBytes, &metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//fmt.Println(metrics)
}

func decodeJsonDataToStruct(metrics *MetricResponse, resp *http.Response) {
	// Decode JSON data into a struct to get the metric values
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(metrics)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

