package js

import (
	"github.com/tapvanvn/gotokenize"
)

type JSRawMeaning struct {
	gotokenize.IMeaning
}

func (meaning *JSRawMeaning) Prepare(stream *gotokenize.TokenStream) {
	//fmt.Println("jsrawmeaning prepare")
	meaning.IMeaning.Prepare(stream)
	newStream := gotokenize.CreateStream()
	for {
		token := meaning.IMeaning.Next()
		if token == nil {
			break
		}
		newStream.AddToken(*token)
	}
	meaning.IMeaning.SetStream(newStream)
	//fmt.Println("end jsrawmeaning prepare")
}

func CreateJSMeaning() gotokenize.PatternMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{

		"#%^&*-+/!<>=?:@\"'` \\;\r\n\t{}[](),.|": {TokenType: TokenJSOperator, Separate: true},
		//"0123456789":   {TokenType: TokenJSNumber, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	jsRawMeaning := JSRawMeaning{

		IMeaning: &meaning,
	}
	jsPhraseMeaning := JSPhraseMeaning{

		IMeaning: &jsRawMeaning,
	}
	return gotokenize.CreatePatternMeaning(&jsPhraseMeaning, JSPatterns, gotokenize.NoTokens, JSGlobalNested)
}

func (meaning *JSRawMeaning) Next() *gotokenize.Token {

	iter := meaning.GetIter()

	return meaning.getNextMeaningToken(iter)
}

func (meaning *JSRawMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {

			break
		}
		token := iter.Read()

		if token.Content == "{" {

			tmpToken := &gotokenize.Token{
				Content: "{",
				Type:    TokenJSBlock,
			}

			meaning.continueUntil(iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "[" {

			tmpToken := &gotokenize.Token{
				Content: "[",
				Type:    TokenJSBracketSquare,
			}

			meaning.continueUntil(iter, tmpToken, "]")

			return tmpToken
		} else if token.Content == "(" {

			tmpToken := &gotokenize.Token{
				Content: "(",
				Type:    TokenJSBracket,
			}

			meaning.continueUntil(iter, tmpToken, ")")

			return tmpToken

		} else if token.Content == "\"" || token.Content == "'" || token.Content == "`" {

			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenJSString,
			}
			meaning.continueReadString(iter, tmpToken, token.Content)

			return tmpToken

		} else if token.Content == "=" {

			nextToken := iter.Get()
			if nextToken != nil {
				if nextToken.Content == ">" {
					tmpToken := gotokenize.Token{Content: "=>", Type: TokenJSRightArrow}
					_ = iter.Read()
					return &tmpToken
				}
			}
			tmpToken := gotokenize.Token{Content: "=", Type: TokenJSAssign}
			return &tmpToken

		} else if token.Content == "/" {

			nextToken := iter.Get()

			if nextToken != nil {

				if nextToken.Content == "/" {

					tmpToken := gotokenize.Token{Content: "//", Type: TokenJSLineComment}

					_ = iter.Read()

					meaning.continueReadLineComment(iter, &tmpToken)

					return &tmpToken

				} else if nextToken.Content == "*" {

					tmpToken := gotokenize.Token{Content: "/*", Type: TokenJSBlockComment}

					_ = iter.Read()

					meaning.continueReadBlockComment(iter, &tmpToken)

					return &tmpToken
				} else {
					if meaning.testRegex(iter) {
						tmpToken := gotokenize.Token{Content: "/", Type: TokenJSRegex}
						tmpToken.Children.AddToken(gotokenize.Token{Type: TokenJSWord, Content: "/"})

						meaning.continueReadRegex(iter, &tmpToken)

						return &tmpToken
					} else {
						return &gotokenize.Token{Content: "/", Type: TokenJSOperator}
					}
				}
			}
		} else if token.Content == " " || token.Content == "\t" {

			return meaning.getNextMeaningToken(iter)

		} else if token.Content == ";" || token.Content == "\n" || token.Content == "\r" {

			token.Type = TokenJSPhraseBreak

			return token
		}

		token.Type = TokenJSWord
		return token
	}
	return nil
}

func (meaning *JSRawMeaning) continueUntil(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

	var specialCharacter bool = false
	var lastSpecialToken *gotokenize.Token = nil

	for {
		if iter.EOS() {
			break
		}

		token := meaning.getNextMeaningToken(iter)

		if token.Content == "\\" {

			specialCharacter = !specialCharacter
			lastSpecialToken = token

		} else if token.Content == reach {

			if specialCharacter {

				specialCharacter = false
				currentToken.Children.AddToken(*token)

			} else {

				break
			}
		} else {

			if specialCharacter {

				currentToken.Children.AddToken(*lastSpecialToken)
			}
			specialCharacter = false
			currentToken.Children.AddToken(*token)
		}

	}
}

func (meaning *JSRawMeaning) continueReadString(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

	var specialCharacter = false
	var lastSpecialToken *gotokenize.Token = nil

	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()

		if token.Content == "\\" {

			specialCharacter = !specialCharacter
			lastSpecialToken = token

		} else if token.Content == reach {

			if specialCharacter {

				specialCharacter = false
				currentToken.Children.AddToken(*token)

			} else {

				break
			}
		} else {

			if specialCharacter && token.Content != "{" && token.Content != "}" && token.Content != "\"" && token.Content != "'" {

				currentToken.Children.AddToken(*lastSpecialToken)
			}
			specialCharacter = false
			currentToken.Children.AddToken(*token)
		}

	}
}

//testRegex test if reach regex
func (meaning *JSRawMeaning) testRegex(iter *gotokenize.Iterator) bool {

	var i = iter.Offset + 1
	for {
		tmpToken := iter.GetAt(i)

		if tmpToken == nil {
			return false
		}

		if tmpToken.Content == "/" {

			testToken := iter.GetAt(i + 1)
			if testToken.Content == "i" || testToken.Content == "m" || testToken.Content == "g" {
				return true

			} else {
				return false
			}
		}
		i += 1
	}
	return false
}

func (meaning *JSRawMeaning) continueReadRegex(iter *gotokenize.Iterator, currToken *gotokenize.Token) {

	//todo: check syntax violence
	var specialCharacter bool = false
	var gotClose bool = false

	for {
		if iter.EOS() {

			break
		}

		tmpToken := iter.Get()

		tmpContent := tmpToken.GetContent()

		if tmpContent == "\\" {

			specialCharacter = !specialCharacter

			currToken.Children.AddToken(*tmpToken)

			_ = iter.Read()

		} else if tmpContent == "/" {

			if specialCharacter {

				specialCharacter = false

			} else {

				gotClose = true
			}

			currToken.Children.AddToken(*tmpToken)

			_ = iter.Read()

		} else {

			if gotClose && tmpContent != "i" && tmpContent != "m" && tmpContent != "g" {

				break

			} else {

				_ = iter.Read()

				specialCharacter = false

				currToken.Children.AddToken(*tmpToken)
			}
		}
	}
}

func (meaning *JSRawMeaning) continueReadLineComment(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	for {
		if iter.EOS() {

			break
		}

		tmpToken := iter.Get()

		if tmpToken.Content == "\n" || tmpToken.Content == "\r" {

			break

		} else {
			_ = iter.Read()
			currentToken.Children.AddToken(*tmpToken)
		}
	}
}

func (meaning *JSRawMeaning) continueReadBlockComment(iter *gotokenize.Iterator, currToken *gotokenize.Token) {

	for {
		if iter.EOS() {

			break
		}
		tmpToken := iter.Read()

		if tmpToken.Content == "*" {

			nextToken := iter.Get()

			if nextToken != nil && nextToken.Content == "/" {

				_ = iter.Read()

				return
			}
		} else {

			currToken.Children.AddToken(*tmpToken)
		}
	}
}

type JSPhraseMeaning struct {
	gotokenize.IMeaning
}

func (meaning *JSPhraseMeaning) Prepare(stream *gotokenize.TokenStream) {
	//fmt.Println("jsphrase prepare")
	meaning.IMeaning.Prepare(stream)
	newStream := gotokenize.CreateStream()
	for {
		token := meaning.IMeaning.Next()
		if token == nil {
			break
		}
		newStream.AddToken(*token)
	}
	meaning.IMeaning.SetStream(newStream)
	//fmt.Println("end jsphrase prepare")
}

func (meaning *JSPhraseMeaning) Next() *gotokenize.Token {

	iter := meaning.GetIter()

	return meaning.getNextMeaningToken(iter)
}
func (meaning *JSPhraseMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {

			break
		}
		token := iter.Read()

		if gotokenize.IsContainToken(JSPhraseAllow, token.Type) {
			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenJSPhrase,
			}
			tmpToken.Children.AddToken(*token)
			meaning.continuePhrase(iter, tmpToken)
			return tmpToken
		} else if token.Type == TokenJSPhraseBreak {
			continue
		}
		return token
	}
	return nil
}

func (meaning *JSPhraseMeaning) continuePhrase(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	for {

		if iter.EOS() {

			break
		}

		tmpToken := iter.Get()

		if gotokenize.IsContainToken(JSPhraseAllow, tmpToken.Type) {
			_ = iter.Read()
			currentToken.Children.AddToken(*tmpToken)
			continue
		}
		break
	}
}
