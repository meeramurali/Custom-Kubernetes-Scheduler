# Custom-Kubernetes-Scheduler

## Setting up a multi-node cluster on a single machine by deploying nodes as Docker containers
1. chmod +x DIND_cluster/dind-cluster-v1.13.sh
2. sudo DIND_cluster/./dind-cluster-v1.13.sh up

See [here](https://www.mirantis.com/blog/multi-kubernetes-kdc-quick-and-dirty-guide/) for detailed instructions

## Setting up Prometheus with node-exporter On Kubernetes Cluster: 
1. Create ‘monitoring’ namespace:
   - kubectl create namespace monitoring
2. Set up node-exporter on all nodes using a DaemonSet: 
   - kubectl create -f node-exporter/node-exporter-daemonset.yml
   - Verify with: kubectl get pods -n monitoring (Must show 1 node-exporter pod per node)
3. Set up prometheus (relevant files in prometheus/):
   - https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/ (skip creating 'monitoring' namespace since already done in step 1)

