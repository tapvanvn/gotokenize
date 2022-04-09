package gotokenize

import (
	"strings"
	"unicode/utf8"
)

type RawMeaning struct {
	*AbstractMeaning
	tokenMap map[string]RawTokenDefine
	separate bool
}

var NoTokens []int = []int{}

func CreateRawMeaning(tokenMap map[string]RawTokenDefine, outputSeparate bool) *RawMeaning {
	return &RawMeaning{
		AbstractMeaning: NewAbtractMeaning(nil),
		tokenMap:        tokenMap,
		separate:        outputSeparate,
	}
}

func (meaning *RawMeaning) Prepare(proc *MeaningProcess) {

	newStream := CreateStream(meaning.GetMeaningLevel())
	//fmt.Printf("raw meaning prepare:%d\n", proc.Stream.Length())
	curType := 0

	meaningLevel := meaning.GetMeaningLevel()

	iter := proc.Stream.Iterator()

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

						newStream.AddToken(Token{
							Content: curContent,
							Type:    curType,
							Children: TokenStream{
								MeaningLevel: meaningLevel,
							},
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

					newStream.AddToken(Token{
						Content: curContent,
						Type:    curType,
						Children: TokenStream{
							MeaningLevel: meaningLevel,
						},
					})
					curContent = ""
				}
				curContent += char
				curType = 0
			}
		}
	}

	if len(curContent) > 0 {

		newStream.AddToken(Token{
			Content: curContent,
			Type:    curType,
			Children: TokenStream{
				MeaningLevel: meaningLevel,
			},
		})
	}
	//fmt.Printf("after raw prepare:%d\n", newStream.Length())
	proc.SetStream(proc.ParentTokens, &newStream)
}

func (meaning *RawMeaning) Clone() IMeaning {

	return CreateRawMeaning(meaning.tokenMap, meaning.separate)
}

func (meaning *RawMeaning) GetName() string {
	return "RawMeaning"
}
