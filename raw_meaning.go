package gotokenize

import (
	"strings"
)

type RawMeaning struct {
	Meaning
	tokenMap map[string]RawTokenDefine
	separate bool
}

var NoTokens []int = []int{}

func CreateRawMeaning(tokenMap map[string]RawTokenDefine, outputSeparate bool) RawMeaning {
	return RawMeaning{
		Meaning:  CreateMeaning(nil),
		tokenMap: tokenMap,
		separate: outputSeparate,
	}
}

func (meaning *RawMeaning) Prepare(stream *TokenStream) {

	meaning.Meaning.Stream = CreateStream()

	curType := 0

	iter := stream.Iterator()

	var curContent = ""

	for {

		if iter.EOS() {
			break
		}

		token := iter.Read()

		if len(token.Content) > 0 {

			found := false

			char := token.Content[0:1]

			for key, value := range meaning.tokenMap {

				if strings.Index(key, char) >= 0 {

					if len(curContent) > 0 && (curType != value.TokenType || value.Separate) {

						meaning.Meaning.Stream.AddToken(Token{
							Content: curContent,
							Type:    curType,
						})

						curContent = ""
					}

					curContent += char

					curType = value.TokenType

					found = true

					break
				}
			}

			if !found {

				if curType != 0 || meaning.separate {

					meaning.Meaning.Stream.AddToken(Token{Content: curContent, Type: curType})
					curContent = ""
				}
				curContent += char
				curType = 0
			}
		}
	}

	if len(curContent) > 0 {

		meaning.Meaning.Stream.AddToken(Token{Content: curContent, Type: curType})
	}

	meaning.Meaning.Iter = meaning.Meaning.Stream.Iterator()
}
