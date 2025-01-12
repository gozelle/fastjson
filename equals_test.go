package fastjson

import "testing"

type EqualsTest struct {
	a     string
	b     string
	error bool
}

func TestEquals(t *testing.T) {

	tests := []EqualsTest{
		{
			a:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d":[1,2,3]}`,
			b:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d": [1,2,3]}`,
			error: false,
		},
		{
			a:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d":[1,2,3]}`,
			b:     `{"a": 3, "b":{"b1":2, "b2": 3}, "d": [1,2,3]}`,
			error: true,
		},
		{
			a:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d":[1,2,3]}`,
			b:     `{"a": 1, "b":{"b1":4, "b2": 3}, "d":[1,2,3]}`,
			error: true,
		},
		{
			a:     `{"a": 1, "b":{"b1": {"c": "cc"}, "b2": 3}, "d":[1,2,3]}`,
			b:     `{"a": 1, "b":{"b1": {"c": "cx"}, "b2": 3}, "d":[1,2,3]}`,
			error: true,
		},
		{
			a:     `{"a": 1, "b":{"b2": 3}, "d":[1,2,3]}`,
			b:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d": [1,2,3]}`,
			error: true,
		},
		{
			a:     `{"a": 1, "b":{"b1": 1, "b2": 3}, "d":[2,3]}`,
			b:     `{"a": 1, "b":{"b2": 3}, "d": [1,2,3]}`,
			error: true,
		},
		{
			a:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d":[1,2,3]}`,
			b:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d": [1,3,2]}`,
			error: true,
		},
		{
			a:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d":[{"a":1}]}`,
			b:     `{"a": 1, "b":{"b1":2, "b2": 3}, "d": [{"a":2}]}`,
			error: true,
		},
	}

	for i, v := range tests {
		err := Equals(v.a, v.b)
		if v.error {
			if err == nil {
				t.Fatalf("test %d expect error, got nil", i)
			}
			t.Logf("test %d error: %s", i, err)
		} else {
			if err != nil {
				t.Fatalf("test %d unexpect error, got error: %s", i, err)
			}
		}
	}

	err := Equals(
		`{"a": 1, "b":{"b1":1, "b2": 3}, "d":[{"a":1}]}`,
		`{"a": 2, "b":{"b1":2, "b2": "3"}, "d": [{"a":2}]}`,
		"$.a",
		"$.b.b1",
		"$.b.b2",
		"$.d.0.a",
	)
	if err != nil {
		t.Fatalf("unexpect error, got error: %s", err)
	}
}
