# Custom-Kubernetes-Scheduler

## Setting up a multi-node cluster on a single machine by deploying nodes as Docker containers
1. cd DIND_cluster/
2. chmod +x dind-cluster-v1.13.sh
3. sudo ./dind-cluster-v1.13.sh up

<img src="https://github.com/meeramurali/Custom-Kubernetes-Scheduler/blob/master/images/1.png" height="60"/>

See [here](https://www.mirantis.com/blog/multi-kubernetes-kdc-quick-and-dirty-guide/) for detailed instructions

## Setting up Prometheus with node-exporter On Kubernetes Cluster: 
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

## Connecting To Prometheus (Web UI)
See "Using Kubectl Port Forwarding" instructions [here](https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/).

kubectl port-forward [your-prometheus-pod-name] 8080:9090 -n monitoring

## Running and testing the Custom Scheduler
1. Open 3 terminals
2. Terminal 1: kubectl proxy
3. Terminal 2: 
   - go build . (from within the scheduler/ folder)
   - ./scheduler.
4. Terminal 3: 
   - kubectl create -f deployments/testcustom.yaml 
   - kubectl get pods -o wide (to see which node its been scheduled on; should be the 'best node' identified by the custom scheduler. See logs in Terminal 2 to verify)
   
##
This project was guided by Kelsey Hightower’s [demo](https://www.youtube.com/watch?v=IYcL0Un1io0) for building a custom Kubernetes scheduler. Hightower demonstrates building a toy scheduler that schedules based on some manually added random annotations (each node is annotated with some random cost; and the scheduler picks the node with the least annotation value). We have used parts of his [source code](https://github.com/kelseyhightower/scheduler) as the base for our scheduler, specifically, the components for monitoring and finding unscheduled pods, running the default predicate checks to find nodes that satisfy the pods’ requested CPU and memory requirements, and binding the pods to the selected nodes. 

