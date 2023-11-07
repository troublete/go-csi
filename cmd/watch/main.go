package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/troublete/go-csi"
)

func main() {
	file := flag.String("file", "", "the file to convert")
	output := flag.String("out", "./out.html", "the target to output")
	flag.Parse()

	createHash := func(s []byte) string {
		return fmt.Sprintf("%x", sha512.Sum512(s))
	}

	update := make(chan *[]byte)
	go func() {
		for {
			select {
			case content := <-update:
				log.Printf("info: detected update, writing file '%v'", *output)
				toks, err := csi.Lexer(string(*content))
				if err != nil {
					panic(err)
				}

				res, err := csi.Interpret(toks)
				if err != nil {
					panic(err)
				}

				err = os.WriteFile(*output, []byte(res), 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	var hash string
	for {
		content, err := os.ReadFile(*file)
		if err != nil {
			panic(err)
		}

		h := createHash(content)
		if h != hash {
			update <- &content
			hash = h
		}
	}
}
