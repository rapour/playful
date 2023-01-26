

## Installation

The project is meant to run over a Kubernetes cluster. We will use Minikube for the sake of simplicity. The instructions for installing minikube can be found [here](https://minikube.sigs.k8s.io/docs/start/).

After installing minikube, start it with your preferred driver. If you already have docker on your host machine, you can use docker as the driver:
```bash
minikube start --driver=docker --memory=4096
```

The next step is to install [Strimzi](https://strimzi.io/), an open source tool to manage Kafka on Kubernetes. Strimzi takes advantage of [operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) in k8s. Create a namespace on your k8s cluster using the following command:
```bash
kubectl create namespace kafka
```
We have named our new namespace `kafka`. Now, we apply Strimzi install files on k8s using the following command:
```bash
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```

More details can be found [here](https://strimzi.io/quickstarts/). Note that the kafka namespace we have created and the `namespace` query parameter inside the installation url of Strimzi must match. 

Check for the status of the cluster you are creating with one of the below commands:
```bash
kubectl get kafka -n kafka -w
kubectl get pods -n kafka -w
``

When the output of the first commands indicates `READY` is `TRUE`, proceed with making a kafka topic:
```bash

```