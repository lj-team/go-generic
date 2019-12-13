package mapcnt

import (
	"fmt"
	"math"
	"reflect"
	"sort"
)

type tagRec struct {
	Tag  string
	Orig reflect.Value
	Cnt  float64
}

type tagList []*tagRec

func (a tagList) Len() int {
	return len(a)
}

func (a tagList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a tagList) Less(i, j int) bool {
	return a[i].Cnt > a[j].Cnt || (a[i].Cnt == a[j].Cnt && a[i].Tag < a[j].Tag)
}

func val2float(v reflect.Value) float64 {
	val := float64(0)

	switch v.Type().Kind() {
	case reflect.Bool:
		if v.Bool() {
			val = 1
		}
	case reflect.Float32:
		val = v.Float()
	case reflect.Float64:
		val = v.Float()
	case reflect.Int:
		val = float64(v.Int())
	case reflect.Int8:
		val = float64(v.Int())
	case reflect.Int16:
		val = float64(v.Int())
	case reflect.Int32:
		val = float64(v.Int())
	case reflect.Int64:
		val = float64(v.Int())
	case reflect.Uint:
		val = float64(v.Uint())
	case reflect.Uint8:
		val = float64(v.Uint())
	case reflect.Uint16:
		val = float64(v.Uint())
	case reflect.Uint32:
		val = float64(v.Uint())
	case reflect.Uint64:
		val = float64(v.Uint())
	}

	return val
}

type ListOpts struct {
	Limit  int
	MinVal float64
}

var defListOpts *ListOpts = &ListOpts{
	Limit:  math.MaxInt32,
	MinVal: 0,
}

func List(src interface{}, opt *ListOpts) interface{} {

	if opt == nil {
		opt = defListOpts
	}

	mp := reflect.ValueOf(src)

	list := make([]*tagRec, 0, mp.Len())

	for _, k := range mp.MapKeys() {

		v := mp.MapIndex(k)
		val := val2float(v)

		if val < opt.MinVal {
			continue
		}

		list = append(list, &tagRec{
			Tag:  fmt.Sprintf("%v", k),
			Orig: k,
			Cnt:  val,
		})
	}

	sort.Sort(tagList(list))

	if opt.Limit > -1 && len(list) > opt.Limit {
		list = list[:opt.Limit]
	}

	res := reflect.MakeSlice(reflect.SliceOf(mp.Type().Key()), 0, len(list))

	for _, v := range list {
		res = reflect.Append(res, v.Orig)
	}

	return res.Interface()
}
