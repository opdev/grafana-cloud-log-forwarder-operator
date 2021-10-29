/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	loggingv1 "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	log "sigs.k8s.io/controller-runtime/pkg/log"

	grafanav1alpha1 "github.com/example/grafana-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaulRetryPeriod = time.Second * 30
)

// GrafanaCloudLogForwarderReconciler reconciles a GrafanaCloudLogForwarder object
type GrafanaCloudLogForwarderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=grafana.example.com,resources=grafanacloudlogforwarders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=grafana.example.com,resources=grafanacloudlogforwarders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=grafana.example.com,resources=grafanacloudlogforwarders/finalizers,verbs=update
//+kubebuilder:rbac:groups=logging.openshift.io,resources=clusterlogforwarders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=logging.openshift.io,resources=clusterloggings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;create;watch;patch;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GrafanaCloudLogForwarder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *GrafanaCloudLogForwarderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// your logic here
	grafanaCloudLogForwarder := &grafanav1alpha1.GrafanaCloudLogForwarder{}
	err := r.Get(ctx, req.NamespacedName, grafanaCloudLogForwarder)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("GrafanaLogForwarderObject not found.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get GrafanaLogForwarderObject")
		return ctrl.Result{}, err
	}

	secretSet := grafanaCloudLogForwarder
	foundSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: "loki1", Namespace: grafanaCloudLogForwarder.Namespace}, foundSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define a new Secret object
			secret := r.newSecretForCR(grafanaCloudLogForwarder)
			if err := controllerutil.SetControllerReference(secretSet, secret, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			err = r.Create(ctx, secret)
			if err != nil {
				log.Error(err, "Failed to create Secret", "secret", secret)
				return ctrl.Result{}, err
			}

			// Secret created successfully - return and requeue after refreshInterval
			return ctrl.Result{RequeueAfter: defaulRetryPeriod}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}

	loggingSet := grafanaCloudLogForwarder
	logging := &loggingv1.ClusterLogging{}
	err = r.Get(ctx, types.NamespacedName{Name: "instance", Namespace: grafanaCloudLogForwarder.Namespace}, logging)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define a new ClusterLogging object
			loggingInstance := r.clusterLoggingForGrafanaCloud(grafanaCloudLogForwarder)
			if err := controllerutil.SetControllerReference(loggingSet, loggingInstance, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			log.Info("Creating a new loggingInstance.", "loggingInstance.Namespace", loggingInstance.Namespace, "loggingInstance.Name", loggingInstance.Name)
			err = r.Create(ctx, loggingInstance)
			if err != nil {
				log.Error(err, "Failed to create new loggingInstance.", loggingInstance)
				return ctrl.Result{}, err
			}
			// ClusterLogging created successfully - return and requeue after refreshInterval
			return ctrl.Result{RequeueAfter: defaulRetryPeriod}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ClusterLogging")
		return ctrl.Result{}, err
	}

	logForward := &loggingv1.ClusterLogForwarder{}
	err = r.Get(ctx, types.NamespacedName{Name: "instance", Namespace: grafanaCloudLogForwarder.Namespace}, logForward)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define a new ClusterLogForwarder object
			logForwardingInstance := r.clusterLogForwarderForGrafanaCloud(grafanaCloudLogForwarder)
			if err := controllerutil.SetControllerReference(loggingSet, logForwardingInstance, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			log.Info("Creating a new logForwarderInstance.", "loggingInstance.Namespace", logForwardingInstance.Namespace, "loggingInstance.Name", logForwardingInstance.Name)
			err = r.Create(ctx, logForwardingInstance)
			if err != nil {
				log.Error(err, "Failed to create new logForwarderInstance.", logForwardingInstance)
				return ctrl.Result{}, err
			}
			// ClusterLogForwarder created successfully - return and requeue after refreshInterval
			return ctrl.Result{RequeueAfter: defaulRetryPeriod}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ClusterLogForwarder")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: defaulRetryPeriod}, nil
}

func (r *GrafanaCloudLogForwarderReconciler) newSecretForCR(gclf *grafanav1alpha1.GrafanaCloudLogForwarder) *corev1.Secret {
	labels := map[string]string{
		"app": "loki1",
	}
	data := make(map[string][]byte)
	data["username"] = []byte(gclf.Spec.Username)
	data["password"] = []byte(gclf.Spec.APIPassword)

	secretObject := &corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "loki1", Namespace: gclf.Namespace, Labels: labels},
		Immutable:  new(bool),
		Data:       data,
	}

	return secretObject
}

func (r *GrafanaCloudLogForwarderReconciler) clusterLoggingForGrafanaCloud(gclf *grafanav1alpha1.GrafanaCloudLogForwarder) *loggingv1.ClusterLogging {

	clusterlogging := &loggingv1.ClusterLogging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "instance",
			Namespace: gclf.Namespace,
		},
		Spec: loggingv1.ClusterLoggingSpec{
			ManagementState: "Managed",
			Collection: &loggingv1.CollectionSpec{
				Logs: loggingv1.LogCollectionSpec{
					Type: loggingv1.LogCollectionTypeFluentd,
				},
			},
		},
	}
	return clusterlogging
}

func (r *GrafanaCloudLogForwarderReconciler) clusterLogForwarderForGrafanaCloud(gclf *grafanav1alpha1.GrafanaCloudLogForwarder) *loggingv1.ClusterLogForwarder {

	clusterLogForwarder := &loggingv1.ClusterLogForwarder{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "instance",
			Namespace: gclf.Namespace,
		},
		Spec: loggingv1.ClusterLogForwarderSpec{
			Outputs: []loggingv1.OutputSpec{{
				Name: "loki-secure",
				URL:  gclf.Spec.URL,
				Type: "loki",
				OutputTypeSpec: loggingv1.OutputTypeSpec{
					Loki: &loggingv1.Loki{},
				},
				Secret: &loggingv1.OutputSecretSpec{
					Name: "loki1",
				},
			},
			},
			Pipelines: []loggingv1.PipelineSpec{{
				Name:       "application-logs",
				InputRefs:  []string{loggingv1.InputNameApplication, loggingv1.InputNameAudit, loggingv1.InputNameInfrastructure},
				OutputRefs: []string{"loki-secure"},
			},
			},
		},
	}
	return clusterLogForwarder
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaCloudLogForwarderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaCloudLogForwarder{}).
		Owns(&corev1.Secret{}).
		Owns(&loggingv1.ClusterLogging{}).
		Owns(&loggingv1.ClusterLogForwarder{}).
		Complete(r)
}
