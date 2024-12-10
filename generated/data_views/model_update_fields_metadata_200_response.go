/*
Data views

OpenAPI schema for data view endpoints

API version: 0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package data_views

import (
	"encoding/json"
)

// checks if the UpdateFieldsMetadata200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdateFieldsMetadata200Response{}

// UpdateFieldsMetadata200Response struct for UpdateFieldsMetadata200Response
type UpdateFieldsMetadata200Response struct {
	Acknowledged *bool `json:"acknowledged,omitempty"`
}

// NewUpdateFieldsMetadata200Response instantiates a new UpdateFieldsMetadata200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateFieldsMetadata200Response() *UpdateFieldsMetadata200Response {
	this := UpdateFieldsMetadata200Response{}
	return &this
}

// NewUpdateFieldsMetadata200ResponseWithDefaults instantiates a new UpdateFieldsMetadata200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateFieldsMetadata200ResponseWithDefaults() *UpdateFieldsMetadata200Response {
	this := UpdateFieldsMetadata200Response{}
	return &this
}

// GetAcknowledged returns the Acknowledged field value if set, zero value otherwise.
func (o *UpdateFieldsMetadata200Response) GetAcknowledged() bool {
	if o == nil || IsNil(o.Acknowledged) {
		var ret bool
		return ret
	}
	return *o.Acknowledged
}

// GetAcknowledgedOk returns a tuple with the Acknowledged field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateFieldsMetadata200Response) GetAcknowledgedOk() (*bool, bool) {
	if o == nil || IsNil(o.Acknowledged) {
		return nil, false
	}
	return o.Acknowledged, true
}

// HasAcknowledged returns a boolean if a field has been set.
func (o *UpdateFieldsMetadata200Response) HasAcknowledged() bool {
	if o != nil && !IsNil(o.Acknowledged) {
		return true
	}

	return false
}

// SetAcknowledged gets a reference to the given bool and assigns it to the Acknowledged field.
func (o *UpdateFieldsMetadata200Response) SetAcknowledged(v bool) {
	o.Acknowledged = &v
}

func (o UpdateFieldsMetadata200Response) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdateFieldsMetadata200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Acknowledged) {
		toSerialize["acknowledged"] = o.Acknowledged
	}
	return toSerialize, nil
}

type NullableUpdateFieldsMetadata200Response struct {
	value *UpdateFieldsMetadata200Response
	isSet bool
}

func (v NullableUpdateFieldsMetadata200Response) Get() *UpdateFieldsMetadata200Response {
	return v.value
}

func (v *NullableUpdateFieldsMetadata200Response) Set(val *UpdateFieldsMetadata200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateFieldsMetadata200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateFieldsMetadata200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateFieldsMetadata200Response(val *UpdateFieldsMetadata200Response) *NullableUpdateFieldsMetadata200Response {
	return &NullableUpdateFieldsMetadata200Response{value: val, isSet: true}
}

func (v NullableUpdateFieldsMetadata200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateFieldsMetadata200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}