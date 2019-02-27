package json

import (
	"encoding/json"
	"reflect"
)

func Equal(b1 []byte, b2 [] byte) (bool, error) {
	var j1 interface{}
	var j2 interface{}

	err := json.Unmarshal(b1, &j1)
	if err != nil {
		return false, nil
	}

	err = json.Unmarshal(b2, &j2)
	if err != nil {
		return false, nil
	}

	v1 := reflect.ValueOf(j1)
	v2 := reflect.ValueOf(j2)

	if v1.Type() != v2.Type() {
		return false, nil
	}

	return deepValueEqual(j1, j2), nil
}

func deepValueEqual(v1, v2 interface{}) bool {
	switch vv1 := v1.(type) {
	case map[string]interface{}:
		vv2 := v2.(map[string]interface{})
		if len(vv1) != len(vv2) {
			return false
		}

		for k, v := range vv1 {
			val2 := vv2[k]
			if !deepValueEqual(v, val2) {
				return false
			}
		}
		return true
	case []interface{}:
		var matches int
		vv2 := v2.([]interface{})
		if len(vv1) != len(vv2) {
			return false
		}
		for _, v := range vv1 {
			for _, v2 := range vv2 {
				if deepValueEqual(v, v2) {
					matches++
				}
			}
		}
		return matches == len(vv1)
	default:
		return v1 == v2
	}
}
