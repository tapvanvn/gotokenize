package gotokenize

import (
	"fmt"
	"strings"
)

const (
	TokenUnknown = iota
	TokenNoType  = -1
)

func NewToken(atLevel int, tokenType int, content string) *Token {
	return &Token{
		Type:    tokenType,
		Content: content,
		Children: TokenStream{
			MeaningLevel: atLevel,
		},
	}
}

type Token struct {
	Type     int
	Content  string
	Children TokenStream
}
type DebugOption struct {
	StringifyTokens []int
}

var EmptyDebugOption = &DebugOption{}

func (token *Token) GetType() int {
	return token.Type
}
func (token *Token) AddChild(childToken Token) {

	token.Children.AddToken(childToken)
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
func (token *Token) Debug(level int, fnName func(int) string, options *DebugOption) {

	printContent := ""
	printChildren := true
	if options != nil && IsContainToken(options.StringifyTokens, token.Type) {
		printContent = strings.TrimSpace(token.Children.ConcatStringContent())
		printChildren = false
	} else {
		printContent = strings.TrimSpace(token.Content)
	}
	if len(printContent) > 20 {

		printContent = string(printContent[0:20]) + "..."
	}

	for i := 0; i <= level; i++ {

		if i == 0 {

			fmt.Printf("|%s ", ColorType(token.Type))

		} else {

			fmt.Print("| ")
		}
	}

	if fnName != nil {

		fmt.Printf("%s", ColorContent(printContent))

		fmt.Printf("-%s\n", ColorName(fnName(token.Type)))

	} else {

		fmt.Println(printContent)

	}
	if !printChildren {
		return
	}
	token.Children.Debug(level+1, fnName, options)
}
