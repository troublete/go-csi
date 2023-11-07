package csi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Lexer(t *testing.T) {
	for c, ex := range []struct {
		in, want string
	}{
		{
			in:   "<head></head><body><main><!--{GET www.hello-world.de}--></main></body>",
			want: TokenSource + TokenStart + TokenVerbGet + TokenSource + TokenEnd + TokenSource,
		},
	} {
		t.Run(fmt.Sprintf("example #%v", c), func(t *testing.T) {
			tokens, err := Lexer(ex.in)
			if err != nil {
				t.Error(err)
			}

			var names []string
			for _, t := range tokens {
				names = append(names, t.Name)
			}

			got := strings.Join(names, "")
			if ex.want != got {
				t.Errorf("want %v, got %v", ex.want, got)
			}

		})
	}
}

func Test_Interpret(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer s.Close()

	for c, ex := range []struct {
		in, want string
	}{
		{
			in:   fmt.Sprintf("<head></head><body><main><!--{GET %s}--></main></body>", s.URL),
			want: "<head></head><body><main>hello world</main></body>",
		},
		{
			in:   "<head></head><body><main><!--{GET https://127.0.0.1:9999}--></main></body>",
			want: "<head></head><body><main></main></body>",
		},
	} {
		t.Run(fmt.Sprintf("example #%v", c), func(t *testing.T) {
			toks, err := Lexer(ex.in)
			if err != nil {
				t.Error(err)
			}

			result, err := Interpret(toks)
			if err != nil {
				t.Error(err)
			}

			if result != ex.want {
				t.Errorf("want %v, got %v", ex.want, result)
			}
		})
	}
}
