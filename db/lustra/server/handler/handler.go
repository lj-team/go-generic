package handler

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/lj-team/go-generic/db/lustra/global"
	st "github.com/lj-team/go-generic/db/lustra/storage"
	"github.com/lj-team/go-generic/text/args"
	"github.com/lj-team/go-generic/time/strftime"
)

type Storage = *st.Storage
type hfunc func(Storage, []string) (string, string)

var (
	handlers map[string]hfunc
)

const (
	resPrefix         = "+"
	errPrefix         = "-"
	defAnswer         = "ok"
	defPong           = "pong"
	errUnknownCommand = "unknown command"
	errInvalidParams  = "invalid params"
)

func init() {

	handlers = map[string]hfunc{}

	handlers["cbadd"] = fCBAdd
	handlers["dec"] = fDec
	handlers["decby"] = fDecBy
	handlers["del"] = fDel
	handlers["get"] = fGet
	handlers["hdel"] = fHDel
	handlers["hget"] = fHGet
	handlers["hset"] = fHSet
	handlers["hsetifless"] = fHSetIfLess
	handlers["hsetifmore"] = fHSetIfMore
	handlers["hsetnx"] = fHSetNX
	handlers["hdec"] = fHDec
	handlers["hdecby"] = fHDecBy
	handlers["hinc"] = fHInc
	handlers["hincby"] = fHIncBy
	handlers["inc"] = fInc
	handlers["incby"] = fIncBy
	handlers["ping"] = fPing
	handlers["set"] = fSet
	handlers["setifless"] = fSetIfLess
	handlers["setifmore"] = fSetIfMore
	handlers["setnx"] = fSetNX
	handlers["time"] = fTime
	handlers["uheap"] = fUHeap
	handlers["version"] = fVersion
}

func Handler(storage Storage, in []byte) []byte {

	list := args.Read(bytes.NewReader(in))
	if len(list) < 2 {
		return nil
	}

	id := list[0]
	isExec := false

	switch id {
	case "exec":
		isExec = true
		id = ""

	case "async":
		id = ""

	default:

		if strings.Index(id, "exec") == 0 {
			id = id[4:]
			isExec = true
		} else {
			return nil
		}
	}

	cmd := list[1]
	list = list[2:]

	var answer string
	var err string

	if fn, h := handlers[cmd]; h {
		answer, err = fn(storage, list)
		if isExec {
			if err != "" {
				answer = errPrefix + id + err
			} else {
				answer = resPrefix + id + answer
			}
		}

	} else {
		if isExec {
			answer = errPrefix + id + errUnknownCommand
		}
	}

	if isExec {
		return []byte(answer)
	}

	return nil
}

func fVersion(storage Storage, list []string) (string, string) {
	if len(list) != 0 {
		return "", errInvalidParams
	}
	return global.Version, ""
}

func fPing(storage Storage, list []string) (string, string) {
	if len(list) != 0 {
		return "", errInvalidParams
	}
	return defPong, ""
}

func fTime(storage Storage, list []string) (string, string) {
	if len(list) != 0 {
		return "", errInvalidParams
	}
	return strftime.Format("%Y-%m-%d %H:%M:%S", time.Now()), ""
}

func fGet(storage Storage, list []string) (string, string) {
	if len(list) != 1 {
		return "", errInvalidParams
	}

	return storage.Get(list[0]), ""
}

func fSet(storage Storage, list []string) (string, string) {
	if len(list)%2 != 0 {
		return "", errInvalidParams
	}

	for len(list) >= 2 {
		storage.Set(list[0], list[1])
		list = list[2:]
	}

	return defAnswer, ""
}

func fSetNX(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	return storage.SetNX(list[0], list[1]), ""
}

func fSetIfMore(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	val, err := strconv.ParseUint(list[1], 10, 64)
	if err != nil {
		return "", errInvalidParams
	}

	return storage.SetIfMore(list[0], val), ""
}

func fHSetIfMore(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	val, err := strconv.ParseUint(list[2], 10, 64)
	if err != nil {
		return "", errInvalidParams
	}

	return storage.HSetIfMore(list[0], list[1], val), ""
}

func fHSetIfLess(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	val, err := strconv.ParseUint(list[2], 10, 64)
	if err != nil {
		return "", errInvalidParams
	}

	return storage.HSetIfLess(list[0], list[1], val), ""
}

func fSetIfLess(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	val, err := strconv.ParseUint(list[1], 10, 64)
	if err != nil {
		return "", errInvalidParams
	}

	return storage.SetIfLess(list[0], val), ""
}

func fDel(storage Storage, list []string) (string, string) {

	for _, key := range list {
		storage.Set(key, "")
	}

	return defAnswer, ""
}

func fInc(storage Storage, list []string) (string, string) {
	if len(list) != 1 {
		return "", errInvalidParams
	}

	return storage.IncBy(list[0], 1), ""
}

func fIncBy(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	if val, err := strconv.ParseUint(list[1], 10, 64); err == nil {
		return storage.IncBy(list[0], val), ""
	}

	return "", errInvalidParams
}

func fDec(storage Storage, list []string) (string, string) {
	if len(list) != 1 {
		return "", errInvalidParams
	}

	return storage.DecBy(list[0], 1), ""
}

func fDecBy(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	if val, err := strconv.ParseUint(list[1], 10, 64); err == nil {
		return storage.DecBy(list[0], val), ""
	}

	return "", errInvalidParams
}

func fCBAdd(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	if val, err := strconv.ParseInt(list[2], 10, 64); err == nil {
		storage.CBAdd(list[0], list[1], int(val))
		return defAnswer, ""
	}

	return "", errInvalidParams
}

func fUHeap(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	if val, err := strconv.ParseInt(list[2], 10, 64); err == nil {
		storage.UHeap(list[0], list[1], int(val))
		return defAnswer, ""
	}

	return "", errInvalidParams
}

func fHGet(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	return storage.HGet(list[0], list[1]), ""
}

func fHSet(storage Storage, list []string) (string, string) {
	if len(list) < 3 || len(list)%2 != 1 {
		return "", errInvalidParams
	}

	hash := make(map[string]string)

	for i := 1; i+1 < len(list); i += 2 {
		hash[list[i]] = list[i+1]
	}

	storage.HSet(list[0], hash)

	return defAnswer, ""
}

func fHSetNX(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	return storage.HSetNX(list[0], list[1], list[2]), ""
}

func fHDel(storage Storage, list []string) (string, string) {
	if len(list) < 2 {
		return "", errInvalidParams
	}

	hash := make(map[string]string)

	for _, k := range list[1:] {
		hash[k] = ""
	}

	storage.HSet(list[0], hash)

	return defAnswer, ""
}

func fHInc(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	return storage.HIncBy(list[0], list[1], 1), ""
}

func fHIncBy(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	if val, err := strconv.ParseUint(list[2], 10, 64); err == nil {
		return storage.HIncBy(list[0], list[1], val), ""
	}

	return "", errInvalidParams
}

func fHDec(storage Storage, list []string) (string, string) {
	if len(list) != 2 {
		return "", errInvalidParams
	}

	return storage.HDecBy(list[0], list[1], 1), ""
}

func fHDecBy(storage Storage, list []string) (string, string) {
	if len(list) != 3 {
		return "", errInvalidParams
	}

	if val, err := strconv.ParseUint(list[2], 10, 64); err == nil {
		return storage.HDecBy(list[0], list[1], val), ""
	}

	return "", errInvalidParams
}

func Batch(storage Storage, in io.Reader) {
	br := bufio.NewReader(in)

	for {
		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			break
		}

		list := args.Read(strings.NewReader(str))
		if len(list) > 0 {
			cmd := list[0]
			list = list[1:]
			if fn, h := handlers[cmd]; h {
				fn(storage, list)
			}
		}
	}
}
