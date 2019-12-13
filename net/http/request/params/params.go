package params

import (
	"net/http"
	"strconv"
)

type Params map[string][]string

func New(req *http.Request) Params {
	param := req.URL.Query()
	return Params(param)
}

func (p Params) GetInt(name string) int64 {
	data, has := p[name]

	if !has || data == nil || len(data) == 0 {
		return 0
	}

	v, err := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return 0
	}

	return v
}

func (p Params) GetBool(name string) bool {

	data, has := p[name]
	if !has || data == nil || len(data) == 0 {
		return false
	}

	v, err := strconv.ParseBool(data[0])
	if err != nil {
		v = false
	}

	return v
}

func (p Params) GetString(name string) string {
	data, has := p[name]
	if !has || data == nil || len(data) == 0 {
		return ""
	}

	return data[0]
}

func (p Params) GetFloat(name string) float64 {
	data, has := p[name]
	if !has || data == nil || len(data) == 0 {
		return 0
	}

	v, err := strconv.ParseFloat(data[0], 64)
	if err != nil {
		return 0
	}

	return v
}
