/*
Copyright 2019 The Pharmer Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CloudProviderSpec defines the desired state of CloudProvider
type CloudProviderSpec struct {
	Regions            []Region            `json:"regions,omitempty"`
	MachineTypes       []MachineType       `json:"machineTypes,omitempty"`
	CredentialFormat   CredentialFormat    `json:"credentialFormat,omitempty"`
	KubernetesVersions []KubernetesVersion `json:"kubernetesVersions,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus,watch
// +kubebuilder:object:root=true

// CloudProvider is the Schema for the cloudproviders API
type CloudProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CloudProviderSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// CloudProviderList contains a list of CloudProvider
type CloudProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CloudProvider `json:"items"`
}

// Region defines the desired state of Region
type Region struct {
	Region   string   `json:"region"`
	Zones    []string `json:"zones,omitempty"`
	Location string   `json:"location,omitempty"`
}
