# ksops-poc, sops for Kubernetes

## Overview
sops for Kubernetes decrypts Kubernetes sops files that can be securely stored along side your code. Extends the Kubernetes API by declaring special _customer resource definition_ that extend their Kubernetes counterparts `kind`; 
`Deployment=ConfigDeploymentSops`, `Ingress=ConfigIngressSops`, `Service=ConfigServiceSops` ...

### Getting started
Directions are intended for Mac OSX users

### Prerequisites
- install brew and golang 1.8+
```
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
$ brew install go
```
- set your $GOPATH
```
$ export GOPATH="$HOME/go"
```

### Requirements
- Kubernetes 1.9+
- Kubebuilder 1.0.5+
~~- Minikube~~

## Build
Don't forget to login into your docker hub account
### Publishing
```
$ export IMG=jechocnct/ksops:alpine
$ make docker-build
$ make docker-push
$ make deploy
```

### Local (Minikube)
Clones, vendors and installs crds onto the cluster:
```
$ git clone git@github.com:jecho/ksops-test.git
$ cd ksops-test
$ dep ensure
$ make 
$ make install
```
Run a local copy
```
$ make run
```

### Verify
Verify that our _custom resource definitions_ are installed properly
```
$ kubectl get crd
NAME                                  CREATED AT
configdeploymentsops.mygroup.k8s.io   2018-12-17T21:50:44Z
configingresssops.mygroup.k8s.io      2018-12-17T22:47:57Z
configservicesops.mygroup.k8s.io      2018-12-17T21:50:44Z
```

## Testing

Files will be encrypted as such, snippet _ghost_deployment.yaml_

```
apiVersion: mygroup.k8s.io/v1beta1
kind: ConfigDeploymentSops
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: configdeploymentsops-sample
spec:
  manifest: |
    apiVersion: ENC[AES256_GCM,data:XEZLS/OKVA==,iv:N6o/g2EMb4oQsFN981uyq1wuXiG9cHM2D7KWLpf70bk=,tag:VVEMifJscE6y+GbIJsHpyA==,type:str]
    kind: ENC[AES256_GCM,data:yuEsTnj7DajaOQ==,iv:fxzwgGp57iEIMywYIxLDajdj1G5VcdDryQRrIjPKztQ=,tag:8A1GlO37Pj+Sm3MMxZuGuA==,type:str]
    metadata:
        name: ENC[AES256_GCM,data:z8fCmAJylx86d8YZ,iv:3yWOPJoUJrRAEe4L+5NXbs7USpzyGsgixu+UdmNcGUk=,tag:mdT79lema3gf5UvXnECcig==,type:str]
    spec:
        replicas: ENC[AES256_GCM,data:gQ==,iv:rDoSdFgE2UuSBHxyHrbU+FiCMCGjoJ8xyb/DBMz+Ojk=,tag:cMigwItqjaDCy0jNmvyklg==,type:int]
        selector:
            matchLabels:
                name: ENC[AES256_GCM,data:RPzgigA=,iv:bJzUKoFliPiw07GsyJUaspb+BMV/vGTMKHC3CpwRPnU=,tag:VSDzMGFDtOv/MP0Pz/c2GQ==,type:str]
                env: ENC[AES256_GCM,data:08kWxwa0oQ==,iv:MQgVcjpug4oiqWpwmuFXDBcYnYr82uhTZlE7YcS4+gQ=,tag:dyqGPdlpyxbcHjwl7vNUKQ==,type:str]
        template:
  ...
```
Running files as their respective `kind` CRDs will decrypt the resources and that Kubernetes can consume

To do this, run an instance of our `ghost` demo

```
$ kubectl create -f ghost_deployment.yaml
$ kubectl create -f ghost_svc.yaml
```

Verify that the deployment is healthy and running
```
$ kubectl get configdeploymentsops.mygroup.k8s.io
NAME                          CREATED AT
configdeploymentsops-sample   1h

$ kubectl get pods
NAME                            READY   STATUS    RESTARTS   AGE
ghost-deploy-5fc8f79f75-rcr65   1/1     Running   0          1h
```

Retrieve the `minikube ip` and the assigned `node port` and reach through your browser
```
$ NODE_PORT=$(kubectl get svc ghost-svc --output=jsonpath='{range .spec.ports[0]}{.nodePort}')
$ echo http://$(minikube ip):${NODE_PORT}
```

## Usage
stuff
