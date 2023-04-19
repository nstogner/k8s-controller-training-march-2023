/*
Copyright 2023.

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

// LoadTestSpec defines the desired state of LoadTest
type LoadTestSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Pattern:=`^(http|https)://(\S+)`
	Address string `json:"address"`
	// +kubebuilder:validation:Enum:={GET,POST,DELETE}
	// +kubebuilder:default:="GET"
	Method string `json:"method"`

	// Duration uses Golang duration formatting from the standard library time package.
	// +kubebuilder:validation:Format:=duration
	Duration metav1.Duration `json:"duration"`
}

// LoadTestStatus defines the observed state of LoadTest
type LoadTestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Enum=Running;Completed
	Phase Phase `json:"phase"`

	StartTime metav1.Time `json:"startTime"`
	EndTime   metav1.Time `json:"endTime,omitempty"`

	RequestCount int `json:"requestCount"`
	SuccessCount int `json:"successCount"`
}

type Phase string

const (
	PhaseRunning   = "Running"
	PhaseCompleted = "Completed"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LoadTest is the Schema for the loadtests API
type LoadTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadTestSpec   `json:"spec,omitempty"`
	Status LoadTestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LoadTestList contains a list of LoadTest
type LoadTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoadTest{}, &LoadTestList{})
}
