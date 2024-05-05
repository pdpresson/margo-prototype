/*
Copyright 2024.

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

package v1

import (
	modelsv1 "ghcr.io/pdpresson/symphony/orchestration-operator/api/v1/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Properties         map[string]string           `json:"properties,omitempty"`
	ProvisioningStatus modelsv1.ProvisioningStatus `json:"provisioningStatus"`
	LastModified       metav1.Time                 `json:"lastModified,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.properties.status`
// +kubebuilder:printcolumn:name="Targets",type=string,JSONPath=`.status.properties.targets`
// +kubebuilder:printcolumn:name="Deployed",type=string,JSONPath=`.status.properties.deployed`

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   modelsv1.InstanceSpec `json:"spec,omitempty"`
	Status InstanceStatus        `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
