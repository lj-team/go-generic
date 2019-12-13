package slice

import (
	"math/rand"
	"reflect"
	"time"
)

func Shuffle(slice interface{}) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()

	for i := 0; i < length; i++ {
		j := rnd.Intn(length)
		if i != j {
			swap(i, j)
		}
	}
}

func Reverse(slice interface{}) {

	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	j := length - 1

	for i := 0; i < length/2; i++ {
		swap(i, j)
		j--
	}
}

func Contain(list interface{}, obj interface{}) (bool, int) {

	if list == nil {
		return false, -1
	}

	if reflect.TypeOf(list).Kind() == reflect.Slice || reflect.TypeOf(list).Kind() == reflect.Array {

		listvalue := reflect.ValueOf(list)

		for i := 0; i < listvalue.Len(); i++ {
			if obj == listvalue.Index(i).Interface() {
				return true, i
			}
		}
	}

	return false, -1
}

func Equal(l1 interface{}, l2 interface{}) bool {

	if l1 == nil && l2 == nil {
		return true
	}

	if l1 == nil {

		if reflect.TypeOf(l2).Kind() == reflect.Slice || reflect.TypeOf(l2).Kind() == reflect.Array {
			return reflect.ValueOf(l2).Len() == 0
		}

		return false
	}

	if l2 == nil {

		if reflect.TypeOf(l1).Kind() == reflect.Slice || reflect.TypeOf(l1).Kind() == reflect.Array {
			return reflect.ValueOf(l1).Len() == 0
		}

		return false
	}

	v1 := reflect.ValueOf(l1)
	k1 := v1.Type().Kind()

	if k1 == reflect.Slice || k1 == reflect.Array {

		v2 := reflect.ValueOf(l2)
		k2 := v2.Type().Kind()

		if k2 == reflect.Slice || k2 == reflect.Array {

			if v1.Len() == v2.Len() {

				for i := 0; i < v1.Len(); i++ {
					if v1.Index(i).Interface() != v2.Index(i).Interface() {
						return false
					}

				}

				return true

			}
		}
	}

	return false
}

func AppendNewValue(listPtr interface{}, obj interface{}) {

	if listPtr == nil {
		return
	}

	vp := reflect.ValueOf(listPtr)
	if vp.Type().Kind() != reflect.Ptr {
		return
	}

	vlist := vp.Elem()
	if vlist.Type().Kind() != reflect.Slice {
		return
	}

	if h, _ := Contain(vlist.Interface(), obj); h {
		return
	}

	vobj := reflect.ValueOf(obj)

	vlist.Set(reflect.Append(vlist, vobj))
}
