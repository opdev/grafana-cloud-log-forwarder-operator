module github.com/example/grafana-operator

go 1.16

require (
	cloud.google.com/go v0.83.0 // indirect
	github.com/Azure/go-autorest/autorest v0.11.19 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.14 // indirect
	github.com/ViaQ/logerr v1.0.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.7.0 // indirect
	github.com/brancz/gojsontoyaml v0.0.0-20191212081931-bf2969bbd742 // indirect
	github.com/brancz/kube-rbac-proxy v0.5.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/go-openapi/validate v0.19.5 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/hashicorp/go-version v1.1.0 // indirect
	github.com/inhies/go-bytesize v0.0.0-20151001220322-5990f52c6ad6 // indirect
	github.com/jsonnet-bundler/jsonnet-bundler v0.3.1 // indirect
	github.com/kylelemons/godebug v0.0.0-20170820004349-d65d576e9348 // indirect
	github.com/mitchellh/hashstructure v0.0.0-20170609045927-2bca23e0e452 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/openshift/api v0.0.0-20210713130143-be21c6cb1bea // indirect
	github.com/openshift/cluster-logging-operator v0.0.0-20211019033137-b03379598eba
	github.com/openshift/prom-label-proxy v0.1.1-0.20191016113035-b8153a7f39f1 // indirect
	github.com/pavel-v-chernykh/keystore-go/v4 v4.1.0 // indirect
	github.com/thanos-io/thanos v0.11.0 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/net v0.0.0-20210610132358-84b48f89b13b // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.3 // indirect
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/controller-tools v0.2.4 // indirect
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.22.2 // Required by prometheus-operator
