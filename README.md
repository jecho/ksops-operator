# ksops-poc

## Overview

Extends the Kubernetes API by declaring custom CRDs ConfigServiceSops,ConfigIngressSopsConfigDeploymentSops that allow mainfest to be encrypted and stored along side code. (git repo)

### Getting started

### Prerequisites
make sure ur gopath is set

### Requirements
- Kubernetes 1.9+
- Kubebuilder 1.0.5
- Minikube

##Build
mac osx
```
git clone git@github.com:jecho/ksops-test.git
cd ksops-test
dep ensure
make 
make install
```

###Deploy
```
$ make run
```
## Testing with Minikube

```
kubectl create -f ghost_deployment.yaml
kubectl create -f ghost_svc.yaml
```