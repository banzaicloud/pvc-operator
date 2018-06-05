package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ObjectStoreList struct for lists
type ObjectStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ObjectStore `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ObjectStore struct bounds together the related struct
type ObjectStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ObjectStoreSpec   `json:"spec"`
	Status            ObjectStoreStatus `json:"status,omitempty"`
}

// ObjectStoreSpec struct holds specs
type ObjectStoreSpec struct {
	Name string
}

// ObjectStoreStatus struct holds status related things
type ObjectStoreStatus struct {
	// Fill me
}
