package csi

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

const (
	TokenStart   = "start"
	TokenEnd     = "end"
	TokenSource  = "src"
	TokenVerbGet = "get"
)

var (
	Dict = map[*regexp.Regexp]Rule{
		regexp.MustCompile(`^<!--{$`): {
			TokenName: TokenStart,
			PeekAhead: 4,
		},
		regexp.MustCompile(`^}-->$`): {
			TokenName: TokenEnd,
			PeekAhead: 3,
		},
		regexp.MustCompile(`^GET$`): {
			TokenName: TokenVerbGet,
			PeekAhead: 2,
		},
	}
	Snippets = map[string]Snippet{
		TokenStart + TokenVerbGet + TokenSource + TokenEnd: func(ts []Token) string {
			t, err := url.Parse(ts[2].Value[1:])
			if err != nil {
				log.Printf("warn: %v", err)
				return ""
			}

			resp, err := http.Get(t.String())
			if err != nil {
				log.Printf("warn: %v", err)
				return ""
			}

			if resp.StatusCode == 200 {
				b, err := io.ReadAll(resp.Body)
				defer func() {
					err := resp.Body.Close()
					if err != nil {
						log.Printf("warn: %v", err)
					}
				}()
				if err != nil {
					log.Printf("warn: %v", err)
					return ""
				}

				return string(b)
			} else {
				log.Printf("warn: HTTP %v request to '%v'", resp.StatusCode, t.String())
				return ""
			}
		},
		TokenSource: func(ts []Token) string {
			return ts[0].Value
		},
	}
)

type Snippet func([]Token) string

type Rule struct {
	TokenName string
	PeekAhead int
}

type Token struct {
	Name  string
	Value string
}

func Lexer(input string) ([]Token, error) {
	var tokens []Token
	currentToken := Token{
		Name:  TokenSource,
		Value: "",
	}

	addToken := func(ts []Token, t Token) []Token {
		if t.Value != "" {
			ts = append(ts, t)
		}
		return ts
	}

	peekAhead := func(input string, start, length int) string {
		end := start + length + 1
		if end > len(input) {
			end = len(input)
		}
		return input[start+1 : end]
	}

	for cursor := 0; cursor < len(input); cursor++ {
		ruleMatch := false
		char := fmt.Sprintf("%c", input[cursor])

		for re, r := range Dict {
			if ruleMatch {
				break
			}

			if m := re.MatchString(char + peekAhead(input, cursor, r.PeekAhead)); m {
				ruleMatch = true
				tokens = addToken(tokens, currentToken)
				tokens = addToken(tokens, Token{
					Name:  r.TokenName,
					Value: char + peekAhead(input, cursor, r.PeekAhead),
				})
				cursor = cursor + r.PeekAhead
				currentToken = Token{
					Name:  TokenSource,
					Value: "",
				}
			}
		}

		if !ruleMatch {
			currentToken.Value += char
		}
	}

	tokens = addToken(tokens, currentToken)
	return tokens, nil
}

func Interpret(toks []Token) (string, error) {
	result := ""
	currentStream := ""
	var tokenList []Token
	for i := 0; i < len(toks); i++ {
		match := false
		currentStream += toks[i].Name
		tokenList = append(tokenList, toks[i])

		for p, s := range Snippets {
			if match {
				break
			}

			if currentStream == p {
				result += s(tokenList)
				match = true
			}
		}

		if match {
			currentStream = ""
			tokenList = []Token{}
			continue
		}
	}

	if currentStream != "" {
		return "", errors.New(fmt.Sprintf("syntax error: no matching expression for '%s'", currentStream))
	}

	return result, nil
}
