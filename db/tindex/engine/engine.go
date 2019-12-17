package engine

import (
	"errors"
)

type Engine interface {
	Open(dsn string) error
	Close()
	Add(key, value string) error
	Del(key string) error
	Search(key string, sim float64) []*SearchRecord
}

type SearchRecord struct {
	Name  string
	Value string
	Sim   float64
}

type SearchRecordList []*SearchRecord

func (l SearchRecordList) Len() int {
	return len(l)
}

func (l SearchRecordList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l SearchRecordList) Less(i, j int) bool {
	return l[i].Sim > l[j].Sim || l[i].Sim == l[j].Sim && l[i].Name < l[j].Name
}

type pair struct {
	Name  string
	Maker func() Engine
}

var (
	drivers []*pair
)

func Register(name string, fn func() Engine) {

	for _, v := range drivers {
		if v.Name == name {
			v.Maker = fn
		}
	}

	drivers = append(drivers, &pair{
		Name:  name,
		Maker: fn,
	})
}

func Open(driver string, dsn string) (Engine, error) {

	for _, v := range drivers {
		if v.Name == driver {
			en := v.Maker()

			if driver == "text" && dsn != "" {

				if err := en.Open(dsn); err != nil {
					return nil, err
				}
			}

			return en, nil
		}
	}

	return nil, errors.New("driver not found")
}
