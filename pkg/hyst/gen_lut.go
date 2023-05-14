//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	var name string
	var max int
	var start int
	var length int
	flag.IntVar(&max, "max", 0xFFFF, "maximum value")
	flag.IntVar(&start, "start", 0, "start value")
	flag.IntVar(&length, "length", 32, "length of the table")
	flag.StringVar(&name, "name", "", "name of the table")

	flag.Parse()

	if name == "" {
		name = "lut"
	}

	stride := (max - start) / length

	pkg := os.Getenv("GOPACKAGE")
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("package %s\n", pkg))
	buf.WriteString(fmt.Sprintf("var %s = []uint16{\n", name))

	v := start
	for i := 0; i < length; i++ {
		buf.WriteString(fmt.Sprintf("\t%d,\n", v))
		v += stride
	}

	buf.WriteString("}\n")

	err := os.WriteFile(name+"_gen.go", []byte(buf.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
