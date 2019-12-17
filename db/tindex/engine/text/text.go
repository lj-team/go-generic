package text

import (
	"bufio"
	"io"
	"os"
	"sort"
	"strings"

	en "github.com/lj-team/go-generic/db/tindex/engine"
	"github.com/lj-team/go-generic/text/args"
	"github.com/lj-team/go-generic/text/trgm"
)

type engine struct {
	Names map[string]int
	Data  map[int]*record
	Trgms map[uint32][]int
}

type record struct {
	Name        string
	Value       string
	TrgmCounter int
}

var (
	nextId chan int
)

func init() {

	nextId = make(chan int, 64)

	go func() {

		i := 1

		for {
			nextId <- i
			i++
		}

	}()

	en.Register("text", maker)
}

func maker() en.Engine {
	return New()
}

func New() *engine {
	return &engine{
		Names: make(map[string]int),
		Data:  make(map[int]*record),
		Trgms: make(map[uint32][]int),
	}
}

func (e *engine) Add(k, v string) error {

	kl := strings.ToLower(k)

	if id, h := e.Names[kl]; h {
		e.Data[id].Value = v
		return nil
	}

	nid := <-nextId

	e.Names[kl] = nid

	r := &record{
		Name:  k,
		Value: v,
	}

	e.Data[nid] = r

	for t := range trgm.MakeTrgms(strings.NewReader(kl)) {
		r.TrgmCounter++
		e.Trgms[t] = append(e.Trgms[t], nid)
	}

	return nil
}

func (e *engine) Del(k string) error {

	kl := strings.ToLower(k)

	kid, h := e.Names[kl]
	if !h {
		return nil
	}

	for t := range trgm.MakeTrgms(strings.NewReader(kl)) {

		list := e.Trgms[t]
		size := len(list)

		for i, v := range list {
			if v == kid {
				list[i] = list[size-1]
				break
			}
		}

		list = list[:size-1]
		e.Trgms[t] = list
	}

	delete(e.Names, kl)
	delete(e.Data, kid)

	return nil
}

func (e *engine) Open(dsn string) error {

	rh, err := os.Open(dsn)
	if err != nil {
		return err
	}
	defer rh.Close()

	e.Names = make(map[string]int)
	e.Data = make(map[int]*record)
	e.Trgms = make(map[uint32][]int)

	e.load(rh)

	return nil
}

func (e *engine) Close() {
	e.Names = make(map[string]int)
	e.Data = make(map[int]*record)
	e.Trgms = make(map[uint32][]int)
}

func (e *engine) Search(k string, sim float64) []*en.SearchRecord {

	k = strings.ToLower(k)

	ids := make(map[int]int)
	total := 0

	for t := range trgm.MakeTrgms(strings.NewReader(k)) {

		if list, h := e.Trgms[t]; h {
			for _, v := range list {
				ids[v]++
			}
		}

		total++
	}

	list := make([]*en.SearchRecord, 0, len(ids))

	for kid, v := range ids {

		rec := e.Data[kid]

		simCur := float64(v) / float64(total+rec.TrgmCounter-v)
		if simCur < sim {
			continue
		}

		list = append(list, &en.SearchRecord{
			Name:  rec.Name,
			Value: rec.Value,
			Sim:   simCur,
		})

	}

	sort.Sort(en.SearchRecordList(list))

	return list
}

func (e *engine) load(rh io.Reader) {

	br := bufio.NewReader(rh)

	for {

		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			return
		}

		toks := args.Parse(str)
		if len(toks) != 2 {
			continue
		}

		e.Add(toks[0], toks[1])
	}

}

func (e *engine) LoadString(txt string) {

	e.Names = make(map[string]int)
	e.Data = make(map[int]*record)
	e.Trgms = make(map[uint32][]int)

	e.load(strings.NewReader(txt))
}
