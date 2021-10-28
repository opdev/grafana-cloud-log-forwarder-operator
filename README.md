# OpenShift Grafana Cloud Log Forwarder

The OpenShift Grafana Cloud Log Forwarder Operator forwards cluster logs from an OpenShift 4x cluster cluster to a [Grafana Cloud](https://grafana.com/products/cloud/) Loki datasource. The operator was built with [operator-sdk](https://sdk.operatorframework.io).

## Overview

The OpenShift Grafana Cloud Log Forwarder Operator will automatically handle the installation of the [Cluster Logging Operator](https://github.com/openshift/cluster-logging-operator) and configure it to use a Grafana Cloud Loki datasource. The user only needs to provide the Grafana Cloud Loki datasource url, username, and api-key/passsword.


### GrafanaCloudLogForwarder Custom Resource

Example for GrafanaCloudLogForwarder spec:

```
apiVersion: grafana.example.com/v1alpha1
kind: GrafanaCloudLogForwarder
metadata:
  name: grafanacloudlogforwarder-sample
  namespace: openshift-logging
spec:
  url: "******"
  username: "******"
  apipassword: "******"
```

## Prerequisites


Since the GrafanaCloudLogForwarder Operator is designed to run inside an OpenShift cluster, hence set it up first. dFor local tests we recommend to use one of the following solutions:

* [OpenShift 4.8 or greater](try.openshift.com), which can be deployed via bare-metal, AWS, GCP, Azure, etc.

To run a single-node OpenShift cluster on your laptop, you can try CRC](https://github.com/code-ready/crc)

## Deployment

There are three ways to run/install the Operator:

* [Operator Lifecycle Manager](https://olm.operatorframework.io) CatalogSource and install via the OpenShift UI
* The `operator-sdk run bundle` command
* As a Go program outside the OpenShift cluster

### Managed by the Operator Lifecycle Manager (OLM) in bundle format

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










### As a Go program outside a cluster: 

First, clone the repository and change to the directory:

```sh
git clone https://github.com/yashoza19/grafana-cloud-log-forwarder-operator
cd grafanacloud-operator
```

Create the `openshift-logging` namespace and make sure that `Red Hat OpenShift Logging Operator` is also installed in the same namespace. 

```sh
oc new-project openshift-logging
```

To run the operator as a go program outside the cluster we will use the following command:

```sh
WATCH_NAMESPACE="openshift-logging" make run
```
