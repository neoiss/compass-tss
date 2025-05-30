
// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// Coin struct for Coin
type Coin struct {
	Asset string `json:"asset"`
	Amount string `json:"amount"`
	Decimals *int64 `json:"decimals,omitempty"`
}

// NewCoin instantiates a new Coin object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCoin(asset string, amount string) *Coin {
	this := Coin{}
	this.Asset = asset
	this.Amount = amount
	return &this
}

// NewCoinWithDefaults instantiates a new Coin object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCoinWithDefaults() *Coin {
	this := Coin{}
	return &this
}

// GetAsset returns the Asset field value
func (o *Coin) GetAsset() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Asset
}

// GetAssetOk returns a tuple with the Asset field value
// and a boolean to check if the value has been set.
func (o *Coin) GetAssetOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Asset, true
}

// SetAsset sets field value
func (o *Coin) SetAsset(v string) {
	o.Asset = v
}

// GetAmount returns the Amount field value
func (o *Coin) GetAmount() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Amount
}

// GetAmountOk returns a tuple with the Amount field value
// and a boolean to check if the value has been set.
func (o *Coin) GetAmountOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Amount, true
}

// SetAmount sets field value
func (o *Coin) SetAmount(v string) {
	o.Amount = v
}

// GetDecimals returns the Decimals field value if set, zero value otherwise.
func (o *Coin) GetDecimals() int64 {
	if o == nil || o.Decimals == nil {
		var ret int64
		return ret
	}
	return *o.Decimals
}

// GetDecimalsOk returns a tuple with the Decimals field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Coin) GetDecimalsOk() (*int64, bool) {
	if o == nil || o.Decimals == nil {
		return nil, false
	}
	return o.Decimals, true
}

// HasDecimals returns a boolean if a field has been set.
func (o *Coin) HasDecimals() bool {
	if o != nil && o.Decimals != nil {
		return true
	}

	return false
}

// SetDecimals gets a reference to the given int64 and assigns it to the Decimals field.
func (o *Coin) SetDecimals(v int64) {
	o.Decimals = &v
}

func (o Coin) MarshalJSON_deprecated() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["asset"] = o.Asset
	}
	if true {
		toSerialize["amount"] = o.Amount
	}
	if o.Decimals != nil {
		toSerialize["decimals"] = o.Decimals
	}
	return json.Marshal(toSerialize)
}

type NullableCoin struct {
	value *Coin
	isSet bool
}

func (v NullableCoin) Get() *Coin {
	return v.value
}

func (v *NullableCoin) Set(val *Coin) {
	v.value = val
	v.isSet = true
}

func (v NullableCoin) IsSet() bool {
	return v.isSet
}

func (v *NullableCoin) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCoin(val *Coin) *NullableCoin {
	return &NullableCoin{value: val, isSet: true}
}

func (v NullableCoin) MarshalJSON_deprecated() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCoin) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


