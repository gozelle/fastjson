package fastjson

import "testing"

func TestEquals(t *testing.T) {
	
	a := `{"a": 1, "b":{"b1":2, "b2": 3}, "d":[1,2,3]}`
	b := `{"a": 1, "b":{"b1":2, "b2": 3}, "d": [1,3,2]}`
	
	t.Log(Equals(a, b))
}
