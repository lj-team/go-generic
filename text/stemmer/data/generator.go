// +build ignore

package main

import (
	"os"
	"os/exec"

	"github.com/lj-team/go-generic/resource"
)

func main() {

	exec.Command("sh", "-c", "sort stemmer.txt > stemmer.txt.sort").Run()
	exec.Command("sh", "-c", "uniq stemmer.txt.sort stemmer.txt").Run()

	os.Remove("stemmer.txt.sort")

	if err := resource.Convert("stemmer.txt", "stemmer.txt.go", "data",
		"github.com/lj-team/go-generic/text/stemmer/data/stemmer.txt", false); err != nil {
		panic(err)
	}
}
