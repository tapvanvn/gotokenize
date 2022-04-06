package gotokenize

import (
	"strings"
	"unicode/utf8"
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
	//fmt.Println("rawmeaning prepare")
	meaning.Meaning.Stream = CreateStream()

	curType := 0

	iter := stream.Iterator()

	var curContent = ""

	for {

		if iter.EOS() {
			break
		}

		token := iter.Read()
		tmpContent := []byte(token.Content)

		if len(tmpContent) > 0 {

			found := false
			utf8Rune, size := utf8.DecodeRune(tmpContent)
			tmpContent = tmpContent[size:]
			char := string(utf8Rune)
			//char := token.Content[0:1]

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
