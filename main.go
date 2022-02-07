package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ivancorrales/hcl-by-example/dsl"
)

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "input", "example.hcl", "path to the input file")
	flag.Parse()
	pipeline, err := loadInputfile(inputPath)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	pipeline.Run()
}

func loadInputfile(path string) (*dsl.Pipeline, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file, diagnostics := hclsyntax.ParseConfig(content, path,
		hcl.Pos{Line: 1, Column: 1, Byte: 0})
	if diagnostics != nil && diagnostics.HasErrors() {
		return nil, diagnostics.Errs()[0]
	}

	out, decodeErr := dsl.Decode(file.Body)
	if decodeErr != nil {
		return nil, decodeErr
	}

	return out, nil
}
