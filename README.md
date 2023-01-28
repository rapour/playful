

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

More details can be found [here](https://strimzi.io/quickstarts/). Note that the kafka namespace we have created and the `namespace` query parameter inside the installation url of Strimzi must match. And then create a kafka cluster using:

```bash
kubectl apply -f .\kafka.yaml -n kafka
```

Check for the status of the cluster you are creating with one of the below commands:
```bash
kubectl get kafka -n kafka -w
kubectl get pods -n kafka -w
``

When the output of the first commands indicates `READY` is `TRUE`, proceed with making a kafka topic:
```bash
kubectl create -f ./kafka-topic.yaml -n kafka
```

Now let's deploy our components. Create another namespace named `dev`:
```bash
kubectl create namespace dev
```
Create a config map in that namespace to hold the information that is needed to connect components:
```bash
kubectl apply -f .\dev-config.yaml -n dev
```
Deploy the `particle` and `app` components on k8s:
```bash
kubectl apply -f .\particles.yaml -n dev
kubectl apply -f .\app.yaml -n dev
```

Now let's create a cluster of [Cassandra](https://cassandra.apache.org/_/index.html) nodes as our data persistence medium using k8s statefulset:
```bash
kubectl apply -f .\cassandra.yml
```

Note that we maintain the Cassandra service within the default namespace of k8s. 