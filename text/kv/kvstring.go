package kv

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func New(str string) (KV, error) {

	res := make(map[string]string)

	k := strings.Builder{}
	v := strings.Builder{}

	mode := 0

	for _, r := range str {

		switch mode {

		case 0:

			if unicode.IsSpace(r) {

			} else if r == '\\' {
				mode = 1
				k.Reset()
			} else if r == '=' {
				k.Reset()
				mode = 3
			} else {
				k.Reset()
				k.WriteRune(r)
				mode = 2
			}

		case 1:

			if c, h := mapper[r]; h {
				k.WriteRune(c)
				mode = 2
			} else {
				return nil, ErrParse
			}

		case 2:

			if r == '=' {
				mode = 3
			} else if r == '\\' {
				mode = 1
			} else if unicode.IsSpace(r) {
				return nil, ErrParse
			} else {
				k.WriteRune(r)
			}

		case 3:

			v.Reset()

			if unicode.IsSpace(r) {
				key := k.String()
				if _, h := res[key]; h {
					return nil, ErrKeyAlreadyExists
				}
				res[key] = ""
				mode = 0
			} else if r == '=' {
				return nil, ErrParse
			} else if r == '\\' {
				mode = 4
			} else {
				v.WriteRune(r)
				mode = 5
			}

		case 4:

			if c, h := mapper[r]; h {
				v.WriteRune(c)
				mode = 5
			} else {
				return nil, ErrParse
			}

		case 5:

			if unicode.IsSpace(r) {
				key := k.String()
				if _, h := res[key]; h {
					return nil, ErrKeyAlreadyExists
				}
				res[key] = v.String()
				mode = 0
			} else if r == '=' {
				return nil, ErrParse
			} else if r == '\\' {
				mode = 4
			} else {
				v.WriteRune(r)

			}

		}

	}

	switch mode {

	case 1, 2, 4:
		return nil, ErrParse

	case 3, 5:

		if mode == 3 {
			v.Reset()
		}

		key := k.String()
		if _, h := res[key]; h {
			return nil, ErrKeyAlreadyExists
		}
		res[key] = v.String()

	}

	return res, nil
}

func (kv KV) GetString(key string, defval string) string {

	if v, h := kv[key]; h {
		return v
	}

	return defval
}

func (kv KV) SetString(key string, val string) {
	kv[key] = val
}

func (kv KV) GetInt64(key string, defval int64) int64 {

	if v, h := kv[key]; h {

		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			return val
		}
	}

	return defval
}

func (kv KV) SetInt64(key string, val int64) {
	kv[key] = strconv.FormatInt(val, 10)
}

func (kv KV) GetInt(key string, defval int) int {

	if v, h := kv[key]; h {

		if val, err := strconv.Atoi(v); err == nil {
			return val
		}
	}

	return defval
}

func (kv KV) SetInt(key string, val int) {
	kv[key] = strconv.Itoa(val)
}

func (kv KV) GetBool(key string, defval bool) bool {

	if v, h := kv[key]; h {

		if val, err := strconv.ParseBool(v); err == nil {
			return val
		}
	}

	return defval
}

func (kv KV) SetBool(key string, val bool) {
	kv[key] = strconv.FormatBool(val)
}

func (kv KV) GetFloat(key string, defval float64) float64 {

	if v, h := kv[key]; h {

		if val, err := strconv.ParseFloat(v, 64); err == nil {
			return val
		}
	}

	return defval
}

func (kv KV) SetFloat(key string, val float64) {
	kv[key] = fmt.Sprintf("%.6f", val)
}

func (kv KV) String() string {

	notFirst := false
	maker := strings.Builder{}

	for k, v := range kv {

		if notFirst {
			maker.WriteRune(' ')
		}

		maker.WriteString(encode(k))
		maker.WriteRune('=')
		maker.WriteString(encode(v))

		notFirst = true
	}

	return maker.String()
}
