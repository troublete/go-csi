package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/troublete/go-csi"
)

func main() {
	file := flag.String("file", "", "the file to convert")
	flag.Parse()

	content, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	toks, err := csi.Lexer(string(content))
	if err != nil {
		panic(err)
	}

	res, err := csi.Interpret(toks)
	if err != nil {
		panic(err)
	}

	fmt.Print(res)
}
