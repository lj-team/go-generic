package tmpl

import (
	"bytes"
	"errors"
	"io"
	"net/url"

	"github.com/CloudyKit/jet"
	"github.com/lj-team/go-generic/log"
)

type VarMap = jet.VarMap

type Tmpl struct {
	set *jet.Set
}

var defTmpl *Tmpl = nil

func New(dirs []string, def bool) *Tmpl {

	tmpl := &Tmpl{
		set: jet.NewHTMLSet(dirs...),
	}

	// add noescape filter
	tmpl.set.AddGlobal("noescape", jet.SafeWriter(func(w io.Writer, b []byte) {
		w.Write(b)
	}))

	tmpl.set.AddGlobal("pathescape", jet.SafeWriter(func(w io.Writer, b []byte) {
		w.Write([]byte(url.PathEscape(string(b))))
	}))

	tmpl.set.AddGlobal("urlescape", jet.SafeWriter(func(w io.Writer, b []byte) {
		w.Write([]byte(url.QueryEscape(string(b))))
	}))

	if def {
		defTmpl = tmpl
	}

	return tmpl
}

func Set(t *Tmpl) {
	defTmpl = t
}

func (t *Tmpl) Vars() VarMap {
	return Vars()
}

func Vars() VarMap {
	return make(VarMap)
}

func (tmpl *Tmpl) Render(name string, vars VarMap) (string, error) {
	if tmpl == nil {
		err := errors.New("nit template object")
		log.Error(err)
		return "", err
	}

	t, err := tmpl.set.GetTemplate(name)
	if err != nil {
		log.Error(err)
		return "", err
	}

	var w bytes.Buffer

	if err = t.Execute(&w, vars, nil); err != nil {
		log.Error(err)
		return "", err
	}

	return w.String(), nil
}

func Render(name string, vars VarMap) (string, error) {
	return defTmpl.Render(name, vars)
}

func (tmpl *Tmpl) LoadTmplString(name string, content string) error {
	if tmpl == nil {
		err := errors.New("nil template pointer")
		return err
	}
	_, err := tmpl.set.LoadTemplate(name, content)
	return err
}

func LoadTmplString(name string, content string) error {
	return defTmpl.LoadTmplString(name, content)
}
