# OpenShift Grafana Cloud Log Forwarder

The OpenShift Grafana Cloud Log Forwarder Operator forwards cluster logs from an [OpenShift](https://try.openshift.com) cluster to a [Grafana Cloud](https://grafana.com/products/cloud/) Loki datasource. The Operator was built with [operator-sdk](https://sdk.operatorframework.io).

## Overview

The OpenShift Grafana Cloud Log Forwarder Operator will automatically handle the installation of the [Cluster Logging Operator](https://github.com/openshift/cluster-logging-operator) and configure it to use a Grafana Cloud Loki datasource. The user only needs to provide the Grafana Cloud Loki datasource url, username, and api-key/passsword.

## Prerequisites

* OpenShift 4.8 or greater

To run a single-node OpenShift cluster on your laptop, you can try [CRC](https://github.com/code-ready/crc).

## Deployment

There are three ways to run/install the Operator:

* Create a [Operator Lifecycle Manager CatalogSource](https://olm.operatorframework.io/docs/concepts/crds/catalogsource/) and install the Operator via the OpenShift Console/UI
* Run the `operator-sdk run bundle` command
* Run the Operator-SDK's `make run` target outside the OpenShift cluster.

### OLM CatalogSource and install via the OpenShift UI

To install the Operator through the Openshift UI, we first create a index:

```sh
opm index add --bundles quay.io/yoza/grafanacloud-operator-bundle:v1.0.0 --tag quay.io/yoza/grafanacloud-operator-index:1.0.0 -c docker
```

Once the index is created, we will push the index image to any repository

```sh
docker push quay.io/yoza/grafanacloud-operator-index:1.0.0
```

Next, we will be creating a CatalogSource and adding the newly created operator index image to the catalogSource. The CatalogSource file is present in this repository

```sh
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: my-test-operators
  namespace: openshift-marketplace
spec:
  sourceType: grpc
  image: quay.io/yoza/grafanacloud-operator-index:1.0.0
  displayName: Test Operators
  publisher: Red Hat Partner
```

```sh
oc create -f catalogsource.yaml
```

Once, the CatalogSource is created, we will go to the Openshift UI and install the Operator through the OperatorHub.

![OperatorHub Tile](https://user-images.githubusercontent.com/4207880/139361897-3210dcbf-3289-44ef-b3f0-e99c107ddf3e.png)

![Install Operator](https://user-images.githubusercontent.com/4207880/139362045-2d141ced-bf0b-4c3b-89ec-e12b7de2c968.png)

Installing the Operator using OLM CatalogSource will automatically create the `openshift-logging` namespace and will deploy the Cluster logging Operator to the same namespace.

### Operator-SDK Run Bundle

Bundle your operator, then build and push the bundle image. The bundle target generates a bundle in the bundle directory containing manifests and metadata defining your Operator. bundle-build and bundle-push build and push a bundle image defined by bundle.Dockerfile.

```sh
make bundle bundle-build bundle-push
```

Make sure that the bundle image is public. Also, make sure that you are using `openshift-logging` namespace before running the `run bundle` command

```sh
operator-sdk run bundle quay.io/yoza/grafanacloud-operator-bundle:v1.0.0
```

### Run outside the OpenShift cluster(`make run`)

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
## Creating the CR

Once the controller/operator is running, we want to create our custom CR in the `config/sample` directory. To create the CR, we first need to add the `URL`, `Username`, and `APIKey/Password` to the sample CR.

The Sample CR should look like this:

```
apiVersion: logs.grafana.com/v1alpha1
kind: GrafanaCloudLogForwarder
metadata:
  name: grafanacloudlogforwarder-sample
  namespace: openshift-logging
spec:
  url: "******"
  username: "******"
  apipassword: "******"
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

Alternatively, the Custom Resource can be created from the UI as well.

<img width="876" alt="Screen Shot 2021-10-28 at 4 16 50 PM" src="https://user-images.githubusercontent.com/29581754/139329181-e6013bdc-8502-498e-b3a6-1544b6222a45.png">

## Cluster logs in loki datasource on GrafanaCloud

After creating GrafanaCloudLogForwarder Custom Resource, we will head over to grafanacloud to check the cluster logs on Loki datasource

<img width="1786" alt="Screen Shot 2021-10-28 at 4 24 09 PM" src="https://user-images.githubusercontent.com/29581754/139330368-8de136db-59c5-4563-94c0-e65e247cc608.png">

## Cleanup

Clean up the GrafanaCloudLogForwarder CR first:
```
oc delete -f config/samples/grafana_v1alpha1_grafanacloudlogforwarder.yaml
```

**Note:** Make sure the above custom resource has been deleted before proceeding to stop the go program. Otherwise your cluster may have dangling custom resource objects that cannot be deleted.

## Useful Links
 
1. [Getting Started with building Go-based Operator](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)
2. [Adding 3rd party resources to your Operator](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#adding-3rd-party-resources-to-your-operator)
3. [Using 3rd party APIs in operator-sdk projects](https://developers.redhat.com/blog/2020/02/04/how-to-use-third-party-apis-in-operator-sdk-projects#step_4__use_the_api_in_the_controllers)
4. [Deploying on OpenShift using opm](https://redhat-connect.gitbook.io/certified-operator-guide/ocp-deployment/openshift-deployment)
