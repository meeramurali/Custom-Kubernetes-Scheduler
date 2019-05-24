# Custom-Kubernetes-Scheduler

## Setting up Prometheus with node-exporter On Kubernetes Cluster: 

1. Create ‘monitoring’ namespace:
   kubectl create namespace monitoring
2. Set up node-exporter on all nodes using a DaemonSet: 
   kubectl create -f node-exporter/node-exporter-daemonset.yml
   Verify with: kubectl get pods -n monitoring (Must show 1 node-exporter pod per node)
3. Set up prometheus:
   (https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/)

