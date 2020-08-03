package main

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Equal checks equality between 2 Body-encoded data.
func Equal(vx, vy interface{}) bool {
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

			if !Equal(v, val2) {
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
				if Equal(v, v2) && !flagged[i] {
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

func Remove(i interface{}, path string) {
	if path == "" {
		return
	}

	var next, current string
	index := strings.IndexRune(path, '.')

	if index == -1 {
		current = path
	} else {
		current = path[:index]
		next = path[index+1:]
	}

	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			if k == current {
				// If there is no more nodes to traverse we can remove it and terminate the routine
				if next == "" {
					delete(t, current)
					return
				}
				Remove(v, next)
			}
		}
	case []interface{}:
		if current == "#" {
			for _, v := range t {
				Remove(v, next)
			}
		}
	}
}

// Unmarshal parses the Body-encoded data into an interface{}.
func Unmarshal(b []byte) (interface{}, error) {
	var j interface{}

	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return j, nil
}
