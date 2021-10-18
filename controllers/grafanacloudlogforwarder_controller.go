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
	err = r.Get(ctx, types.NamespacedName{Name: grafanaCloudLogForwarder.Spec.SecretName, Namespace: grafanaCloudLogForwarder.Namespace}, foundSecret)
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

	return ctrl.Result{RequeueAfter: defaulRetryPeriod}, nil
}

func (r *GrafanaCloudLogForwarderReconciler) newSecretForCR(gclf *grafanav1alpha1.GrafanaCloudLogForwarder) *corev1.Secret {
	labels := map[string]string{
		"app": gclf.Spec.SecretName,
	}
	data := make(map[string][]byte)
	data["username"] = []byte(gclf.Spec.Username)
	data["password"] = []byte(gclf.Spec.APIPassword)

	secretObject := &corev1.Secret{
		TypeMeta:   metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: gclf.Spec.SecretName, Namespace: gclf.Namespace, Labels: labels},
		Immutable:  new(bool),
		Data:       data,
	}

	return secretObject
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaCloudLogForwarderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&grafanav1alpha1.GrafanaCloudLogForwarder{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
