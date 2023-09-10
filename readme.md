# Minikube testbench

## Overview

This repo is to demonstrate of deploying local Minikube cluster, local Docker registry and runs an Go application in it. 

## Prerequisite

### Skills

* Kubernetes

* Go

### In your local machine

* minikube CLI

* kubectl CLI  
  You may use built-in kubectl inside minikube, `minikube kubectl`.

* task CLI - [Link](https://taskfile.dev/installation/)

* Docker desktop or server.

* Go CLI

## Steps

### Minikube 

1. You may list the available scripts for the demostration:
   ```
   > task --list
   task: Available tasks for this project:
   * k8s:apply-env:                             Apply Kubernetes manifest with Kustomize, based on variable ENV.
   * k8s:build:                                 Builds docker image.
   ...
   ```

2. To start minikube cluster, run `task minikube:start`. Please allow the command to run for few minutes, on top of that, this commands starts few things:
   * In-cluster `localhost:5000` Docker image registry,
   * `csi-hostpath-sc` storage class driver,
   * `metrics-server` inside the minikube,
   * Minikube container runtime is `containerd`,
   * Minikube CNI is `calico`.

3. To use the local registry, run `task minikube:registry`, this puts a passthrough of ur machine localhost:5000 to in-cluster's Docker image registry.

4. Run `kubectl cluster-info`, to validate the cluster information:
   ```
   > kubectl cluster-info
   Kubernetes control plane is running at https://127.0.0.1:32769
   CoreDNS is running at https://127.0.0.1:32769/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
   
   To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
   ```

### Workload

Run `ENV=production task k8s:full-flow`, this will:

1. Run `docker build`, tags to latest and value evaluated in bash of `date +"%H%M%S-%d%m%y"`,

2. Push to local registry `localhost:5000`,

3. Emits `tag.txt` file with the date-time tag.

4. With `yq` processor, it edits the `kustomization.yaml` per `ENV`, in this case is `production`, this will update the `.images[].newTag` to the latest tag based on `tag.txt`

5. Deploy to the cluster based on current context (in this case is your Minikube cluster).

6. Run `kubectl -n production get pod`, you will see a `ubuntu` pod and 10 `api` deployment replicas.

### Horizontal Pod Autoscaling (HPA)

To demonstrate HPA in action, run `kubectl exec -it ubuntu -- ab -k -c 800 -n 1000000 api:8080/hello`. 

With elevated loads across the `api` pod replicas, the HPA will scale the pods from 10 to 20, as exemplified in `kubernetes/overlays/production/hpa.yaml`.