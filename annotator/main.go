//Original file obtained from https://github.com/kelseyhightower/scheduler/blob/master/annotator/main.go

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var listOnly bool

type NodeList struct {
	Items []Node `json:"items"`
}

type Node struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Name        string            `json:"name,omitempty"`
	Annotations map[string]string `json:"annotations"`
}

func main() {
	flag.BoolVar(&listOnly, "l", false, "List current annotations and exist")
	flag.Parse()

	resp_metric, err_metric := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemTotal")
    fmt.Println(resp_metric)
	if err_metric != nil {
		fmt.Println(err_metric)
		os.Exit(1)
	}

	prices := []string{"0.05", "0.10", "0.20", "0.40", "0.80", "1.60"}
	resp, err := http.Get("http://127.0.0.1:8001/api/v1/nodes")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Invalid status code", resp.Status)
		os.Exit(1)
	}

	var nodes NodeList
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&nodes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if listOnly {
		for _, node := range nodes.Items {
			price := node.Metadata.Annotations["hightower.com/cost"]
			fmt.Printf("%s %s\n", node.Metadata.Name, price)
		}
		os.Exit(0)
	}

	rand.Seed(time.Now().Unix())
	for _, node := range nodes.Items {
		price := prices[rand.Intn(len(prices))]
		annotations := map[string]string{
			"hightower.com/cost": price,
		}
		patch := Node{
			Metadata{
				Annotations: annotations,
			},
		}

		var b []byte
		body := bytes.NewBuffer(b)
		err := json.NewEncoder(body).Encode(patch)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		url := "http://127.0.0.1:8001/api/v1/nodes/" + node.Metadata.Name
		request, err := http.NewRequest("PATCH", url, body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		request.Header.Set("Content-Type", "application/strategic-merge-patch+json")
		request.Header.Set("Accept", "application/json, */*")

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if resp.StatusCode != 200 {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s %s\n", node.Metadata.Name, price)
	}
}
