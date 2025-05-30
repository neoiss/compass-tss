
// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// KeygenBlock struct for KeygenBlock
type KeygenBlock struct {
	// the height of the keygen block
	Height *int64 `json:"height,omitempty"`
	Keygens []Keygen `json:"keygens"`
}

// NewKeygenBlock instantiates a new KeygenBlock object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKeygenBlock(keygens []Keygen) *KeygenBlock {
	this := KeygenBlock{}
	this.Keygens = keygens
	return &this
}

// NewKeygenBlockWithDefaults instantiates a new KeygenBlock object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKeygenBlockWithDefaults() *KeygenBlock {
	this := KeygenBlock{}
	return &this
}

// GetHeight returns the Height field value if set, zero value otherwise.
func (o *KeygenBlock) GetHeight() int64 {
	if o == nil || o.Height == nil {
		var ret int64
		return ret
	}
	return *o.Height
}

// GetHeightOk returns a tuple with the Height field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeygenBlock) GetHeightOk() (*int64, bool) {
	if o == nil || o.Height == nil {
		return nil, false
	}
	return o.Height, true
}

// HasHeight returns a boolean if a field has been set.
func (o *KeygenBlock) HasHeight() bool {
	if o != nil && o.Height != nil {
		return true
	}

	return false
}

// SetHeight gets a reference to the given int64 and assigns it to the Height field.
func (o *KeygenBlock) SetHeight(v int64) {
	o.Height = &v
}

// GetKeygens returns the Keygens field value
func (o *KeygenBlock) GetKeygens() []Keygen {
	if o == nil {
		var ret []Keygen
		return ret
	}

	return o.Keygens
}

// GetKeygensOk returns a tuple with the Keygens field value
// and a boolean to check if the value has been set.
func (o *KeygenBlock) GetKeygensOk() ([]Keygen, bool) {
	if o == nil {
		return nil, false
	}
	return o.Keygens, true
}

// SetKeygens sets field value
func (o *KeygenBlock) SetKeygens(v []Keygen) {
	o.Keygens = v
}

func (o KeygenBlock) MarshalJSON_deprecated() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Height != nil {
		toSerialize["height"] = o.Height
	}
	if true {
		toSerialize["keygens"] = o.Keygens
	}
	return json.Marshal(toSerialize)
}

type NullableKeygenBlock struct {
	value *KeygenBlock
	isSet bool
}

func (v NullableKeygenBlock) Get() *KeygenBlock {
	return v.value
}

func (v *NullableKeygenBlock) Set(val *KeygenBlock) {
	v.value = val
	v.isSet = true
}

func (v NullableKeygenBlock) IsSet() bool {
	return v.isSet
}

func (v *NullableKeygenBlock) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeygenBlock(val *KeygenBlock) *NullableKeygenBlock {
	return &NullableKeygenBlock{value: val, isSet: true}
}

func (v NullableKeygenBlock) MarshalJSON_deprecated() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeygenBlock) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


