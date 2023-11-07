package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/troublete/go-csi"
)

func main() {
	file := flag.String("file", "", "the file to convert")
	port := flag.String("port", "4321", "the target port to listen")
	flag.Parse()

	createHash := func(s []byte) string {
		return fmt.Sprintf("%x", sha512.Sum512(s))
	}

	var s string

	update := make(chan *[]byte)
	go func() {
		for {
			select {
			case content := <-update:
				toks, err := csi.Lexer(string(*content))
				if err != nil {
					panic(err)
				}

				res, err := csi.Interpret(toks)
				if err != nil {
					panic(err)
				}

				s = res
			}
		}
	}()

	var hash string
	go func() {
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
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(s))
		if err != nil {
			log.Printf("serve error: %v", err)
		}
	})

	log.Printf("serving file on :" + (*port))
	err := http.ListenAndServe(":"+(*port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
