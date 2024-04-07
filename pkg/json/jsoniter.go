package json

import (
	jsoniter "github.com/json-iterator/go"
)

var std = jsoniter.ConfigCompatibleWithStandardLibrary

// ToJSON transform struct to json
func ToJSON(v interface{}) string {
	b, err := std.Marshal(v)
	if err != nil {
		return ""
	}

	return string(b)
}

func ToJSONb(v interface{}) ([]byte, error) {
	b, err := std.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ToJSONf transform struct to json and text format
func ToJSONf(v interface{}) string {
	b, err := std.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func ToJSONs(v interface{}) string {
	b, err := std.MarshalIndent(v, "", "")
	if err != nil {
		return ""
	}
	return string(b)
}

// ToStruct ...json string to struct
func StringToStruct(s string, v interface{}) error {
	return std.Unmarshal([]byte(s), v)
}

func BytesToStruct(d []byte, v any) error {
	return std.Unmarshal(d, v)
}
