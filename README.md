# Custom-Kubernetes-Scheduler

#### Note: kubectl commands will not work from git repo root folder. Must run from outside or cd into one of the sub-folders.

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

## Run annotator
kubectl proxy (in a separate terminal)

kubectl port-forward [your-prometheus-pod-name] 8080:9090 -n monitoring (in a separate terminal)

go run annotator/main.go


Original version from https://github.com/kelseyhightower/scheduler/tree/master/annotator. Adds random cost values as annotations to each node. 

TO DO: Annotator needs to be modified to use metric data from prometheus instead of random cost values.
