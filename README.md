# OpenShift Grafana Cloud Log Forwarder

The OpenShift Grafana Cloud Log Forwarder Operator forwards cluster logs from an OpenShift 4x cluster cluster to a [Grafana Cloud](https://grafana.com/products/cloud/) Loki datasource. The operator was built with [operator-sdk](https://sdk.operatorframework.io).

## Overview

OpenShift 4x features a variety of Operators available for download in the embedded [OperatorHub](https://docs.openshift.com/container-platform/4.9/operators/understanding/olm-understanding-operatorhub.html)

The OpenShift Grafana Cloud Log Forwarder Operator will automatically handle the installation of the [Cluster Logging Operator](https://github.com/openshift/cluster-logging-operator) and configure it to use a Grafana Cloud Loki datasource. The user only needs to provide the Grafana Cloud Loki datasource url, username, and api-key/passsword.
gs (from system pods and node logs), and audit logs (special node logs with legal/security implications)

### GrafanaCloudLogForwarder Custom Resource

Example for GrafanaCloudLogForwarder spec:

```
spec:
  username: "******"
  apipassword: "******"
  url: "******"
```

Once the operator is deployed on an OCP cluster, it would watch for creation/deletion of GrafanaCloudLogForwarder CR. 

## Prerequisites

Since the GrafanaCloudLogForwarder Operator is designed to run inside an OpenShift cluster, hence set it up first. For local tests we recommend to use one of the following solutions:
* [crc](https://github.com/code-ready/crc), which creates a single-node openshift cluster on your laptop.

The Red Hat OpenShift Logging Operator is responsible for handling ClusterLogging and ClusterLogForwarder CR. Red Hat OpenShift Logging Operator is only deployable to the `openshift-logging ` namespace. This namespace must be explicitly created by a cluster administrator (e.g. `oc new-project openshift-logging`). Red Hat OpenShift Logging Operator operator has been added to the list of dependency package in our bundle, to ensure that it is installed prior to the installation of GrafanaCloudLogForwarder Operator.

## Deployment

There are two ways to run the operator:

* As Go program outside a cluster
* Managed by the Operator Lifecycle Manager (OLM) in bundle format

### As Go program outside a cluster: 

The GrafanaCloudLogForwarder Operator can be installed simply by running it as a go program outside the cluster:

First, clone the repository and change to the directory:

```sh
git clone https://github.com/yashoza19/GrafanaLogForwarder-Operator.git
cd grafanacloud-operator
```

We would have to create the `openshift-logging` namespace and make sure that `Red Hat OpenShift Logging Operator` is also installed in the same namespace. 

```sh
oc new-project openshift-logging
```

Make sure you can currently access the namespace:

```sh
oc project
```

To run the operator as a go program outside the cluster we will use the following command:

```sh
WATCH_NAMESPACE="openshift-logging" make run
```

### Managed by the Operator Lifecycle Manager (OLM) in bundle format:

Firstly we would want to install [OLM](https://sdk.operatorframework.io/docs/olm-integration/tutorial-bundle/#enabling-olm).

```sh
operator-sdk olm install
```

Bundle your operator, then build and push the bundle image. The bundle target generates a bundle in the bundle directory containing manifests and metadata defining your operator. bundle-build and bundle-push build and push a bundle image defined by bundle.Dockerfile.

```sh
make bundle bundle-build bundle-push
```

Make sure that the bundle image is public, and then run your bundle.

```sh
operator-sdk run bundle quay.io/yoza/grafanacloud-operator-bundle:v0.0.1
```

### Creating the CR:

Once the controller/operator is running, we want to create our custom CR in the `config/sample` directory. To create the CR, we first need to add the `Username`, `APIPassword`, and `URL` to the sample CR.

The Sample CR should look like this:

```
apiVersion: grafana.example.com/v1alpha1
kind: GrafanaCloudLogForwarder
metadata:
  name: grafanacloudlogforwarder-sample
  namespace: openshift-logging
spec:
  username: "******"
  apipassword: "******"
  url: "******"
```

Create the GrafanaCloudLogForwarder CR that was modified:

```sh
oc apply -f config/samples/grafana_v1alpha1_grafanacloudlogforwarder.yaml
```

To verify the creation of the GrafanaCloudLogForwarder CR:

```sh
oc get grafanacloudlogforwarder
NAME                              AGE
grafanacloudlogforwarder-sample   1m
```

Verify the creation of the Operands (ClusterLogging, ClusterLogForwarder):

```sh
oc get secrets
NAME                                       TYPE                                  DATA   AGE
loki1                                      Opaque                                2      8s
```

```sh
oc get clusterlogging
NAME       MANAGEMENT STATE
instance   Managed
```

```sh
oc get clusterlogforwarder
NAME       AGE
instance   3m48s
```

### Cleanup:

Clean up the GrafanaCloudLogForwarder CR first:
```
oc delete -f config/samples/grafana_v1alpha1_grafanacloudlogforwarder.yaml
```

**Note:** Make sure the above custom resource has been deleted before proceeding to stop the go program. Otherwise your cluster may have dangling custom resource objects that cannot be deleted.
