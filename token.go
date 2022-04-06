package gotokenize

import (
	"fmt"
	"strings"
)

const (
	TokenUnknown = iota
)

type Token struct {
	Type     int
	Content  string
	Children TokenStream
}

func (token *Token) GetType() int {
	return token.Type
}

func (token *Token) GetContent() string {
	return token.Content
}

func (token *Token) GetChildren() *TokenStream {
	return &token.Children
}

func IndexOf(runes []rune, ch rune) int {
	tmpOffset := 0
	for {
		if tmpOffset == len(runes) {
			break
		}
		tmpRune := runes[tmpOffset]
		if tmpRune == ch {
			return tmpOffset
		}
		tmpOffset++
	}
	return -1
}
func (token *Token) Debug(level int, fnName func(int) string) {

	trimContent := strings.TrimSpace(token.Content)

	for i := 0; i <= level; i++ {

		if i == 0 {

			fmt.Printf("|%s ", ColorType(token.Type))

		} else {

			fmt.Print("| ")
		}
	}

	if fnName != nil {

		fmt.Printf("%s", ColorContent(trimContent))

		fmt.Printf("-%s\n", ColorName(fnName(token.Type)))

	} else {

		fmt.Println(trimContent)

	}

	token.Children.Debug(level+1, fnName)
}
