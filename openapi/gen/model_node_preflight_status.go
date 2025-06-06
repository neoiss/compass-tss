
// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// NodePreflightStatus struct for NodePreflightStatus
type NodePreflightStatus struct {
	// the next status of the node
	Status string `json:"status"`
	// the reason for the transition to the next status
	Reason string `json:"reason"`
	Code int64 `json:"code"`
}

// NewNodePreflightStatus instantiates a new NodePreflightStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNodePreflightStatus(status string, reason string, code int64) *NodePreflightStatus {
	this := NodePreflightStatus{}
	this.Status = status
	this.Reason = reason
	this.Code = code
	return &this
}

// NewNodePreflightStatusWithDefaults instantiates a new NodePreflightStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNodePreflightStatusWithDefaults() *NodePreflightStatus {
	this := NodePreflightStatus{}
	return &this
}

// GetStatus returns the Status field value
func (o *NodePreflightStatus) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *NodePreflightStatus) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *NodePreflightStatus) SetStatus(v string) {
	o.Status = v
}

// GetReason returns the Reason field value
func (o *NodePreflightStatus) GetReason() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Reason
}

// GetReasonOk returns a tuple with the Reason field value
// and a boolean to check if the value has been set.
func (o *NodePreflightStatus) GetReasonOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Reason, true
}

// SetReason sets field value
func (o *NodePreflightStatus) SetReason(v string) {
	o.Reason = v
}

// GetCode returns the Code field value
func (o *NodePreflightStatus) GetCode() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Code
}

// GetCodeOk returns a tuple with the Code field value
// and a boolean to check if the value has been set.
func (o *NodePreflightStatus) GetCodeOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Code, true
}

// SetCode sets field value
func (o *NodePreflightStatus) SetCode(v int64) {
	o.Code = v
}

func (o NodePreflightStatus) MarshalJSON_deprecated() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["status"] = o.Status
	}
	if true {
		toSerialize["reason"] = o.Reason
	}
	if true {
		toSerialize["code"] = o.Code
	}
	return json.Marshal(toSerialize)
}

type NullableNodePreflightStatus struct {
	value *NodePreflightStatus
	isSet bool
}

func (v NullableNodePreflightStatus) Get() *NodePreflightStatus {
	return v.value
}

func (v *NullableNodePreflightStatus) Set(val *NodePreflightStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableNodePreflightStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableNodePreflightStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNodePreflightStatus(val *NodePreflightStatus) *NullableNodePreflightStatus {
	return &NullableNodePreflightStatus{value: val, isSet: true}
}

func (v NullableNodePreflightStatus) MarshalJSON_deprecated() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNodePreflightStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


