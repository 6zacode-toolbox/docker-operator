/*
Copyright 2023 6zacode-toolbox.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DockerComposeRunnerSpec defines the desired state of DockerComposeRunner
type DockerComposeRunnerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DockerComposeRunner. Edit dockercomposerunner_types.go to remove/update
	HostIP        string   `json:"hostip"`
	ComposeFile   string   `json:"composeFile"`
	ExecutionPath string   `json:"executionPath,omitempty"`
	RepoAddress   string   `json:"repoAddress"`
	MountVars     []string `json:"mountVars,omitempty"`
}

// DockerComposeRunnerStatus defines the observed state of DockerComposeRunner
type DockerComposeRunnerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Instanced     bool            `json:"instanced,omitempty"`
	ComposeStatus []ComposeStatus `json:"composeStatus,omitempty"`
}

type ComposeStatus struct {
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	ConfigFiles string `json:"configFiles,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DockerComposeRunner is the Schema for the dockercomposerunners API
type DockerComposeRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DockerComposeRunnerSpec   `json:"spec,omitempty"`
	Status DockerComposeRunnerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DockerComposeRunnerList contains a list of DockerComposeRunner
type DockerComposeRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DockerComposeRunner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DockerComposeRunner{}, &DockerComposeRunnerList{})
}
