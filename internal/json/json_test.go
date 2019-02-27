package json

import (
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name    string
		isEqual bool
		b1      []byte
		b2      []byte
	}{
		{
			name:    "empty array",
			isEqual: true,
			b1:      []byte(`[]`),
			b2:      []byte(`[]`),
		},
		{
			name:    "empty object",
			isEqual: true,
			b1:      []byte(`{}`),
			b2:      []byte(`{}`),
		},
		{
			name:    "different basic types",
			isEqual: false,
			b1:      []byte(`[]`),
			b2:      []byte(`{}`),
		},
		{
			name:    "different first name",
			isEqual: false,
			b1:      []byte(`{"FirstName":"Alan"}`),
			b2:      []byte(`{"FirstName":"Galileo"}`),
		},
		{
			name:    "equal first name",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan"}`),
			b2:      []byte(`{"FirstName":"Alan"}`),
		},
		{
			name:    "equal fields different order",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan", "LastName": "Turing"}`),
			b2:      []byte(`{"LastName": "Turing", "FirstName":"Alan"}`),
		},
		{
			name:    "equal fields different order",
			isEqual: true,
			b1:      []byte(`{"x": {"t": 1, "s": 2}, "z": 1}`),
			b2:      []byte(`{"z": 1, "x": {"s": 2, "t": 1}}`),
		},
		{
			name:    "equal array objects",
			isEqual: true,
			b1:      []byte(`[{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}]`),
			b2:      []byte(`[{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}]`),
		},
		{
			name:    "equal exact same objects",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
		},
		{
			name:    "equal flip one field at the root",
			isEqual: true,
			b1:      []byte(`{"LastName": "Turing", "FirstName":"Alan", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
		},
		{
			name:    "flip one field inside a node",
			isEqual: true,
			b1:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Martin Fowler", "Rob Pike"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
		},
		{
			name:    "Flip one field inside inner node",
			isEqual: true,
			b1:      []byte(`{"LastName": "Turing", "FirstName":"Alan", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Country":  "US", "Number": 111}]}`),
			b2:      []byte(`{"LastName": "Turing", "FirstName":"Alan", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"], "Info" : [{"Number": 111, "Country":  "US"}]}`),
		},
		{
			name:    "flip object in root array",
			isEqual: true,
			b1:      []byte(`[{"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}, {"FirstName":"Martin", "LastName": "Fowler", "Age" : 30, "Friends" : []}]`),
			b2:      []byte(`[{"FirstName":"Martin", "LastName": "Fowler", "Age" : 30, "Friends" : []}, {"FirstName":"Alan", "LastName": "Turing", "Age" : 20, "Friends" : ["Rob Pike", "Martin Fowler"]}]`),
		},
		/*		{
					name: "some complex example",
					isEqual: true,
					b1: []byte(`[{"a": 1, "b": [{"c": [1,5,2,4]}, {"d": [1]}]}]`),
					b2: []byte(`[{"b": [{"d": [1]}, {"c": [1,2,4,5]}], "a": 1}]`),
				},*/
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			equal, _ := Equal(test.b1, test.b2)
			assert.Equal(t, test.isEqual, equal)
		})
	}
}
