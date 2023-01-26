package fastjson

import (
	"fmt"
)

func Equals(a, b string) error {
	return EqualsBytes([]byte(a), []byte(b))
}

func EqualsBytes(a, b []byte) (err error) {
	
	va, err := ParseBytes(a)
	if err != nil {
		err = fmt.Errorf("init a error: %s", err)
		return
	}
	
	vb, err := ParseBytes(b)
	if err != nil {
		err = fmt.Errorf("init b error: %s", err)
		return
	}
	if va.Type() != TypeArray && va.Type() != TypeObject {
		return equals(va, vb)
	} else if va.Type() == TypeArray {
		return equalsArray(va, vb)
	} else {
		return equalsObject(va, vb)
	}
}

func equalsObject(va, vb *Value) (err error) {
	
	ob, err := vb.Object()
	if err != nil {
		err = fmt.Errorf("get b object error: %s", err)
		return
	}
	oa, _ := va.Object()
	
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	
	if oa.Len() != ob.Len() {
		err = fmt.Errorf("oa len: %d != ob len: %d", oa.Len(), ob.Len())
		return
	}
	
	oa.Visit(func(key []byte, v *Value) {
		obv := ob.Get(string(key))
		if obv == nil {
			panic(fmt.Errorf("b(%s) is nil", string(key)))
		}
		if v.Type() == TypeObject {
			err = equalsObject(v, obv)
			if err != nil {
				panic(err)
			}
		} else if v.Type() == TypeArray {
			err = equalsArray(v, obv)
			if err != nil {
				panic(err)
			}
		} else {
			err = equals(v, obv)
			if err != nil {
				panic(err)
			}
		}
	})
	
	return
}

func equals(va, vb *Value) error {
	
	if va.Type() == TypeArray || va.Type() == TypeObject {
		return fmt.Errorf("internal error: sholud not use equals")
	}
	
	if va.String() != vb.String() {
		return fmt.Errorf("%s != %s", va, vb)
	}
	
	return nil
}

func equalsArray(va, vb *Value) (err error) {
	
	ab, err := vb.Array()
	if err != nil {
		err = fmt.Errorf("get b object error: %s", err)
		return
	}
	aa, _ := va.Array()
	
	if len(aa) != len(ab) {
		err = fmt.Errorf("array length not match, aa: %d, ab: %d", len(aa), len(ab))
		return
	}
	for i, v := range aa {
		if v.Type() == TypeObject {
			err = equalsObject(v, ab[i])
		} else if v.Type() == TypeArray {
			err = equalsArray(v, ab[i])
		} else {
			err = equals(v, ab[i])
		}
		if err != nil {
			return
		}
	}
	
	return
}
