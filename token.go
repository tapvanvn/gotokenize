package gotokenize

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
