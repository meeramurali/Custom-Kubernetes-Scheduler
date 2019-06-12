package main

import (
    "fmt"
)

func getBestNode(nodes []Node) (Node, error) {
    // Get name of best node per metrics
    bestNodeName, err := getBestNodeName(nodes)
    fmt.Printf("Best Node Name: %s\n", bestNodeName)
    if err != nil {
		return Node{}, err
    }

    // Initialize best node to nil
    var bestNode *Node
    bestNode = &Node{}

    // Find actual node with that node
	for _, n := range nodes {
        // if name of node n matches bestNodeName
        // save that as bestNode
        nodeName := n.Metadata.Name
        if nodeName == bestNodeName {
            fmt.Printf("Found Best Node\n")
            bestNode = &n
            break
        }
	}

    // return best node and error if any
    if bestNode != nil {
        return *bestNode, nil
    } else {
        return Node{}, fmt.Errorf("Unable to match best node name")
    }
}

