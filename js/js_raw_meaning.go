package js

import (
	"fmt"

	"github.com/tapvanvn/gotokenize"
)

func NewJSRawMeaning() *JSRawMeaning {
	return &JSRawMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(nil),
	}
}

type JSRawMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSRawMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := meaning.getNextMeaningToken(process.Iter)

	if token != nil && gotokenize.IsContainToken(JSGlobalNested, token.Type) {

		childProcess := gotokenize.NewMeaningProcessFromStream(&token.Children)

		fmt.Println("childToken level", token.Children.MeaningLevel)

		subStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

		meaning.Prepare(childProcess)

		for {

			nestedToken := meaning.Next(childProcess)

			if nestedToken == nil {
				break
			}
			subStream.AddToken(*nestedToken)
		}

		token.Children = subStream
	}
	return token
}

func (meaning *JSRawMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {

			break
		}
		token := iter.Read()

		if token.Content == "{" {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBlock, "{")

			meaning.continueUntil(iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "[" {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBracketSquare, "[")

			meaning.continueUntil(iter, tmpToken, "]")

			return tmpToken

		} else if token.Content == "(" {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBracket, "(")

			meaning.continueUntil(iter, tmpToken, ")")

			return tmpToken

		} else if token.Content == "\"" || token.Content == "'" || token.Content == "`" {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSString, token.Content)

			meaning.continueReadString(iter, tmpToken, token.Content)

			return tmpToken

		} else if token.Content == "=" {

			nextToken := iter.Get()

			if nextToken != nil {

				if nextToken.Content == ">" {
					tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSRightArrow, "=>")

					_ = iter.Read()
					return tmpToken

				} else if nextToken.Content == "=" {
					tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "==")

					_ = iter.Read()
					return tmpToken
				}
			}
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSAssign, "=")

			return tmpToken

		} else if token.Content == "/" {

			nextToken := iter.Get()

			if nextToken != nil {

				if nextToken.Content == "/" {

					tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSLineComment, "//")

					_ = iter.Read()

					meaning.continueReadLineComment(iter, tmpToken)

					return tmpToken

				} else if nextToken.Content == "*" {

					tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBlockComment, "/*")

					_ = iter.Read()

					meaning.continueReadBlockComment(iter, tmpToken)

					return tmpToken

				} else {
					if meaning.testRegex(iter) {
						tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSRegex, "/")

						tmpToken.Children.AddToken(gotokenize.Token{Type: TokenJSWord, Content: "/"})

						meaning.continueReadRegex(iter, tmpToken)

						return tmpToken
					} else {
						tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSOperator, "/")
						return tmpToken
					}
				}
			}
		} else if token.Content == " " || token.Content == "\t" {

			return meaning.getNextMeaningToken(iter)

		} else if token.Content == ";" || token.Content == "\n" || token.Content == "\r" {

			token.Type = TokenJSPhraseBreak

			return token
		}

		if token.Type == 0 {

			token.Type = TokenJSWord
		}
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
func (meaning *JSRawMeaning) GetName() string {

	return "JSRawMeaning"
}
