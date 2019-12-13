package resource

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileData struct {
	Name        string
	ContentType string
	Data        []byte
}

var (
	data map[string]*FileData
)

func init() {
	data = map[string]*FileData{}
}

func Add(name, contentType string, arch bool, content string) {

	cont, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		panic(err)
	}

	if arch {
		gz, _ := gzip.NewReader(bytes.NewReader(cont))

		data, err := ioutil.ReadAll(gz)
		if err != nil {
			panic(err)
		}

		cont = data
	}

	data[name] = &FileData{
		Name:        name,
		ContentType: contentType,
		Data:        cont,
	}
}

func Delete(name string) {
	delete(data, name)
}

func Get(name string) *FileData {

	val, _ := data[name]
	return val
}

func Convert(src string, dest string, pkg string, name string, noarch bool) error {
	if pkg == "" {
		return errors.New("package name not entered")
	}

	if src == "" {
		return errors.New("source file not entered")
	}

	if name == "" {
		return errors.New("internal resource name not entered")
	}

	if dest == "" {
		return errors.New("output file name not enetered")
	}

	data, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}

	rw, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer rw.Close()

	if !noarch {

		var b bytes.Buffer
		gz := gzip.NewWriter(&b)

		if _, err := gz.Write(data); err != nil {
			panic(err)
		}
		if err := gz.Close(); err != nil {
			panic(err)
		}
		data = b.Bytes()
	}

	result := base64.StdEncoding.EncodeToString(data)

	ext := filepath.Ext(src)
	ct := "application/octet-stream"

	if ext != "" {
		ct = mime.TypeByExtension(ext)
		if ct == "" {
			ct = "application/octet-stream"
		}
	}

	builder := strings.Builder{}

	fmt.Fprintf(rw, "package %s\n\n", pkg)

	fmt.Fprintf(rw, "import (\n\t\"github.com/lj-team/go-generic/resource\"\n)\n\n")
	fmt.Fprintf(rw, "func init() {\n\n")

	fmt.Fprintf(rw, "\tresource.Add(%s, %s, %t, `\n", strconv.Quote(name), strconv.Quote(ct), !noarch)

	i := 0

	for _, r := range result {
		i++
		builder.WriteRune(r)
		if i >= 64 {
			builder.WriteRune('\n')
			fmt.Fprint(rw, builder.String())
			i = 0
			builder.Reset()
		}
	}

	if i > 0 {
		fmt.Fprint(rw, builder.String())
	}

	fmt.Fprintf(rw, "`)\n\n}\n")

	return nil
}
