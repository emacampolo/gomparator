package json

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func Equal(b1, b2 []byte) bool {
	vx, vy, err := Unmarshal(b1, b2)
	if err != nil {
		return false
	}

	equals := bytes.Equal(b1, b2)
	if equals {
		return true
	}

	return deepEqual(vx, vy)
}

func deepEqual(vx, vy interface{}) bool {
	if reflect.TypeOf(vx) != reflect.TypeOf(vy) {
		return false
	}

	switch x := vx.(type) {
	case map[string]interface{}:
		y := vy.(map[string]interface{})

		if len(x) != len(y) {
			return false
		}

		for k, v := range x {
			val2 := y[k]

			if (v == nil) != (val2 == nil) {
				return false
			}

			if !deepEqual(v, val2) {
				return false
			}
		}

		return true
	case []interface{}:
		y := vy.([]interface{})

		if len(x) != len(y) {
			return false
		}

		var matches int
		flagged := make([]bool, len(y))
		for _, v := range x {
			for i, v2 := range y {
				if deepEqual(v, v2) && !flagged[i] {
					matches++
					flagged[i] = true
					break
				}
			}
		}
		return matches == len(x)
	default:
		return vx == vy
	}
}

func Unmarshal(b1 []byte, b2 []byte) (interface{}, interface{}, error) {
	var j1 interface{}
	var j2 interface{}

	err := json.Unmarshal(b1, &j1)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(b2, &j2)
	if err != nil {
		return nil, nil, err
	}

	return j1, j2, nil
}
