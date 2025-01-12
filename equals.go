package fastjson

import (
	"fmt"
)

func Equals(a, b string, ignorePaths ...string) error {
	return EqualsBytes([]byte(a), []byte(b), ignorePaths...)
}

func EqualsBytes(a, b []byte, ignorePaths ...string) (err error) {
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

	path := "$"
	m := map[string]struct{}{}
	for _, v := range ignorePaths {
		m[v] = struct{}{}
	}

	if va.Type() != TypeArray && va.Type() != TypeObject {
		err = equals(path, va, vb, m)
		if err != nil {
			return
		}
	} else if va.Type() == TypeArray {
		err = equalsArray(path, va, vb, m)
		if err != nil {
			return
		}
	} else {
		err = equalsObject(path, va, vb, m)
		if err != nil {
			return
		}
	}
	return
}

func equalsObject(path string, va, vb *Value, m map[string]struct{}) (err error) {

	if _, ok := m[path]; ok {
		return nil
	}

	ob, err := vb.Object()
	if err != nil {
		err = compareError(path, fmt.Errorf("get b object error: %s", err))
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
		err = compareError(path, fmt.Errorf("object length: %d != %d", oa.Len(), ob.Len()))
		return
	}

	oa.Visit(func(key []byte, v *Value) {
		p := fmt.Sprintf("%s.%s", path, string(key))

		obv := ob.Get(string(key))
		if obv == nil {
			err = compareError(p, fmt.Errorf("b.Value is nil"))
			return
		}
		if v.Type() == TypeObject {
			err = equalsObject(p, v, obv, m)
			if err != nil {
				panic(err)
			}
		} else if v.Type() == TypeArray {
			err = equalsArray(p, v, obv, m)
			if err != nil {
				panic(err)
			}
		} else {
			err = equals(p, v, obv, m)
			if err != nil {
				panic(err)
			}
		}
	})

	return
}

func equals(path string, va, vb *Value, m map[string]struct{}) error {
	if _, ok := m[path]; ok {
		return nil
	}
	if va.Type() == TypeArray || va.Type() == TypeObject {
		return compareError(path, fmt.Errorf("internal error: sholud not use equals"))
	}

	if va.Raw() != vb.Raw() {
		return compareError(path, fmt.Errorf("%s != %s", va.Raw(), vb.Raw()))
	}

	return nil
}

func equalsArray(path string, va, vb *Value, m map[string]struct{}) (err error) {
	if _, ok := m[path]; ok {
		return nil
	}

	ab, err := vb.Array()
	if err != nil {
		err = compareError(path, fmt.Errorf("get b object error: %s", err))
		return
	}
	aa, _ := va.Array()

	if len(aa) != len(ab) {
		err = compareError(path, fmt.Errorf("array length not match, aa: %d, ab: %d", len(aa), len(ab)))
		return
	}

	for i, v := range aa {
		p := fmt.Sprintf("%s.%d", path, i)
		if v.Type() == TypeObject {
			err = equalsObject(p, v, ab[i], m)
		} else if v.Type() == TypeArray {
			err = equalsArray(p, v, ab[i], m)
		} else {
			err = equals(p, v, ab[i], m)
		}
		if err != nil {
			return
		}
	}

	return
}

func compareError(path string, err error) error {
	return fmt.Errorf("path: %s compare error: %s", path, err)
}
