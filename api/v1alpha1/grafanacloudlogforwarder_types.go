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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaCloudLogForwarderSpec defines the desired state of GrafanaCloudLogForwarder
type GrafanaCloudLogForwarderSpec struct {

	// // Secret name that will be used for storing the loki username and API Key
	// SecretName string `json:"secretname"`

	// The username from the loki endpoint
	Username string `json:"username"`

	// Enter API key to authenticate clusterLogForwarder to loki datasource
	APIPassword string `json:"apipassword"`

	// URL to loki datasource
	URL string `json:"url"`
}

// GrafanaCloudLogForwarderStatus defines the observed state of GrafanaCloudLogForwarder
type GrafanaCloudLogForwarderStatus struct {
	// Sec
	SecretName      string `json:"secretName"`
	SecretNamespace string `json:"secretNamespace"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// GrafanaCloudLogForwarder is the Schema for the grafanacloudlogforwarders API
type GrafanaCloudLogForwarder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaCloudLogForwarderSpec   `json:"spec,omitempty"`
	Status GrafanaCloudLogForwarderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GrafanaCloudLogForwarderList contains a list of GrafanaCloudLogForwarder
type GrafanaCloudLogForwarderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaCloudLogForwarder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GrafanaCloudLogForwarder{}, &GrafanaCloudLogForwarderList{})
}
