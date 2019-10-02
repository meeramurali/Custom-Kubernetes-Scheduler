# Custom-Kubernetes-Scheduler
## Project goal
The goal was to build a custom Kubernetes scheduler that schedules containerized applications on nodes based on actual memory usage instead of requested memory (default). For each unscheduled pod, a compatible node (that meets the pod's CPU and memory requests) with the maximum available memory is selected. 

## The default scheduler
When a pod is created, it does not actually start running. The desired state is simply stored in the Kubernetes API server. The scheduler is responsible for assigning the pod to a compatible node. Then, kubelet, an agent that runs on every node in the cluster, executes the assigned pod, i.e. runs all the containers in the pod. 

The default scheduler determines the placement using some constraints and chooses the ones with highest score. It first applies predicates, which will remove nodes that don't meet the requirements for the pod or don't have enough resources available [10]. The available resources for a node are determined by adding the resource requests (including CPU and memory) of all pods running on the node and dividing by the capacity of the node [10].

Then, the scheduler applies priority functions to determine the best node to schedule the pod on.  The goal is to distribute the pods across different nodes and zones, and still favor nodes that are least loaded (based on the resources available). Additional factors, such as reducing the number of pods from the same service running on the same node, are considered as well [10]. The node with the highest priority (or a node tied for the highest priority) is chosen [9].

## Our custom scheduler
The objective of our custom scheduler is to optimize memory utilization across the nodes.  The default scheduler considers cpu and memory requests [10] (which are optional in the yaml files specifying pods or deployments of pods) [8], while our scheduler focuses on memory usage on the nodes.  It's possible for pods to use much more memory than requested, or to use much less memory than requested, or to not specify a request.  After the pods are scheduled, they may not use the estimated amount of memory, so focusing on memory requests may not accurately reflect how much memory will be used.

Our custom scheduler makes decisions based on how much memory is available on the nodes, to address the case where the actual usage is very different from the requested usage.  We use a metric (node_memory_memAvailable) to determine which node has the most memory available.  With the default Kubernetes scheduler, if one node is running a pod that requested a lot of memory but was not actually using very much memory, and another node was running a pod that didn't have a memory request but was actually using a lot of memory, Kubernetes would likely schedule a new pod on the node that was using a lot of memory.  Based on requests, it appears to be the least-loaded node.  Our custom scheduler could spread out the memory usage across the nodes, instead of ending up with memory-intensive applications running on the same node. 

## Methodology
### Approach:
The Kubernetes scheduler needs to do three main tasks: First, constantly monitor pods to find any unscheduled ones. Second, for each unscheduled pod, find the ‘best’ node. ‘Best’, in the context of our custom scheduler, means a compatible node with the most available memory. Third, bind the pod to the selected node. 
Node selection (the second task) is the one that we have customized for our purposes - to favor the one with most available memory. This involves the following steps:
* Find all compatible nodes that ‘fit’ the pod i.e. satisfy the pod’s requested CPU and memory requirements.
* Retrieve the relevant metric data (available memory) for each compatible node.
* Select the node with the optimum (maximum) metric value. 

### Tech stack:
We used the following tools to build and test our custom scheduler:
* Kubernetes + Docker
* Kubeadm-dind-cluster (KDC):  to create a local multi-node cluster made of Docker containers
* YAML files:  to define pods and deployments
* Prometheus Node Exporter:  to scrape node metrics
* GoLang:  custom scheduler development
* PromQL + Prometheus HTTP API:  to query metric data

### Metric Tool:
We have used Prometheus [5], an open source, systems monitoring toolkit, specifically the Node Exporter tool, to scrape metric data for all active nodes. Node exporter offers 43 metrics that expose memory information. “The amount of available memory on a Linux system is not just the reported “free” memory metric. Unix systems rely heavily on memory that is not in use by applications to share code (buffers) and to cache disk pages (cached). So one measure of available memory is: sum(node_memory_MemFree + node_memory_Cached + node_memory_Buffers)” [7]. Newer Linux kernels (after 3.14) expose a better free memory metric, node_memory_MemAvailable, which is the one we have chosen to use. 

## Implementation
Our implementation was guided by Kelsey Hightower’s demo for building a custom Kubernetes scheduler [12]. Hightower demonstrates building a toy custom scheduler that schedules based on some manually added random annotations (each node is annotated with some random cost; and the scheduler picks the least cost node). We have used parts of his source code [13] as the base for our scheduler, specifically, the components for monitoring and finding unscheduled pods, running the default predicate checks to find nodes that satisfy the pods’ requested CPU and memory requirements, and binding the pods to the selected nodes. 

Upon this base, we built components to factor in the node metric values from Prometheus, into the scheduling decision. Once the default predicate checks are run to find the list of compatible nodes that satisfy the pod’s requested CPU and memory requirements, a PromQL query is made via Prometheus’ HTTP API to get the node_memory_memAvailable metric values for the compatible nodes: 
```
resp, err := http.Get("http://localhost:8080/api/v1/query?query=node_memory_MemAvailable") 
```

Then, the node with the maximum value is identified and assigned to the pod. The following figure outlines the major components of the custom scheduler. 

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/impl.png" width="900"/>

## Testing
To test the custom scheduler against the default kubernetes scheduler, we deploy five pods: 
Sleep(2 replicas), Nginx(1 replica) on each scheduler (default and the custom scheduler) and sysbench(1 replica). 

The sleep.yaml is a pod that requests for 1800Mi of memory and sleeps infinitely. We deploy two pods of sleep, requesting for 3600Mi of memory in total. The testdefault.yaml runs the Nginx application on the default scheduler and testcustom.yaml  runs the Nginx application on the custom scheduler.   

The sysbench.yaml runs the sysbench workload for the memory test. The benchmark application allocates a buffer of size memory-block-size and then reads/writes from the buffer until a specific volume (memory-total-size) is reached. The number of threads and the type of operation(read or write, sequential or random) can be specified by the user[4]. 
In the current project, we have allocated a buffer size of 2GB and the total volume to be read is 5TB. The read operation is performed in random by 10 threads. 

In the beginning, the two pods running the sleep application and one pod running the sysbench workload are deployed. The three pods are deployed on the three different nodes. The fourth pod is deployed using custom scheduler (testcustom.yaml) and default scheduler (testdefault.yaml). The pods are deployed on nodes based on the metrics in the custom scheduler and the default scheduler. The results of the experiments performed to schedule the pods using the custom and the default scheduler are discussed in the next section. 

## Results
### Setting up the initial state:
Initial pods running on the nodes, with sleep pods running on two nodes, and the sysbench pod running on the third node:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r1.png" width="900"/>

### Initial memory requests for the nodes, shown by the output from the kubectl describe nodes command:
Node 42zw (running sysbench) has 24% of memory requested:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r2.png" width="900"/>

Node mcnx (running sleep) has 84% of memory requested:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r3.png" width="900"/>

Node xrg1 (running sleep) has 84% of memory requested:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r4.png" width="900"/>

### Initial memory utilization:
The initial node_memory_memAvailable metrics for each node:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r5.png" width="900"/>

The node 42zw, running the sysbench pod, has the least memory available.  The node mcnx, running a sleep pod, has the most memory available.

### Scheduling a pod with default scheduler:
When a new pod is scheduled with the default scheduler, it is scheduled on the node running sysbench, which is using the most memory:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r6.png" width="900"/>

### Scheduling a pod with custom scheduler:
When a new pod is scheduled with the custom scheduler, it is scheduled on the node using the least memory:

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r7.png" width="900"/>

### Custom scheduler log, showing the metric values captured:
<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/r8.png" width="900"/>


## Running the custom scheduler
### Setting up a multi-node cluster on a single machine by deploying nodes as Docker containers
1. cd DIND_cluster/
2. chmod +x dind-cluster-v1.13.sh
3. sudo ./dind-cluster-v1.13.sh up

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/1.png" height="60"/>

See [here](https://www.mirantis.com/blog/multi-kubernetes-kdc-quick-and-dirty-guide/) for detailed instructions

### Setting up Prometheus with node-exporter On Kubernetes Cluster: 
1. Create ‘monitoring’ namespace:
   - kubectl create namespace monitoring
2. Set up node-exporter on all nodes using a DaemonSet: 
   - cd node-exporter/
   - kubectl create -f node-exporter-daemonset.yml
   - Verify with: kubectl get pods -n monitoring (Must show 1 node-exporter pod per node)
3. Set up prometheus (relevant files in prometheus/):
   - cd prometheus/
   - kubectl create -f clusterRole.yaml
   - kubectl create -f config-map.yaml -n monitoring
   - kubectl create  -f prometheus-deployment.yaml --namespace=monitoring
   
See [here](https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/) for detailed instructions

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/2.png" height="75"/>

### Connecting To Prometheus (Web UI)
See "Using Kubectl Port Forwarding" instructions [here](https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/).

kubectl port-forward [your-prometheus-pod-name] 8080:9090 -n monitoring

### Running and testing the Custom Scheduler
1. Open 3 terminals
2. Terminal 1: kubectl proxy
3. Terminal 2: 
   - go build . (from within the scheduler/ folder)
   - ./scheduler.
4. Terminal 3: 
   - kubectl create -f deployments/testcustom.yaml 
   - kubectl get pods -o wide (to see which node its been scheduled on; should be the 'best node' identified by the custom scheduler. See logs in Terminal 2 to verify)
  
## References
[1]  Kubernetes tutorial. Url: https://www.tutorialspoint.com/kubernetes/index.htm

[2]  Kubernetes 101: Pods, Nodes, Containers, and Clusters. Url: https://medium.com/google-cloud/kubernetes-101-pods-nodes-containers-and-clusters-c1509e409e16

[3]  Kubernetes concepts. Url: https://kubernetes.io/docs/concepts/

[4]  Sysbench workload. Url: https://wiki.gentoo.org/wiki/Sysbench#Using_the_memory_workload

[5]  Prometheus. Url: https://prometheus.io/

[6]  Node Exporter. Url: https://prometheus.io/docs/guides/node-exporter/

[7]  A Deep Dive into Kubernetes Metrics. Url: https://blog.freshtracks.io/a-deep-dive-into-kubernetes-metrics-part-2-c869581e9f29

[8]  The Kubernetes Authors. Managing Compute Resources for Containers. url: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/ (accessed: 06.10.2019)

[9]  Eduar Tua. Scheduler Algorithm in Kubernetes. Url: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scheduling/scheduler_algorithm.md (accessed: 06.10.2019)

[10]  Eduar Tua. The Kubernetes Scheduler. Url: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-scheduling/scheduler.md. (accessed: 06.10.2019)

[11]  4 cool Kubernetes tools for mastering clusters. Url: https://www.infoworld.com/article/3196250/4-cool-kubernetes-tools-for-mastering-clusters.html

[12]  GopherCon 2016: Kelsey Hightower - Building a custom Kubernetes scheduler. Url: https://www.youtube.com/watch?v=IYcL0Un1io0

[13]  Hightower Toy Scheduler. Url: https://github.com/kelseyhightower/scheduler

[14] Google Kubernetes Engine. Url: https://cloud.google.com/kubernetes-engine/
