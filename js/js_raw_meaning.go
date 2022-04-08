package js

import (
	"github.com/tapvanvn/gotokenize/v2"
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

	if token != nil && token.Children.Length() > 0 && gotokenize.IsContainToken(JSGlobalNested, token.Type) {

		childProcess := gotokenize.NewMeaningProcessFromStream(&token.Children)

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

		if token.Content == "{" && token.Type != TokenJSBlock {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBlock, "{")

			meaning.continueUntil(iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "[" && token.Type != TokenJSBracketSquare {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBracketSquare, "[")

			meaning.continueUntil(iter, tmpToken, "]")

			return tmpToken

		} else if token.Content == "(" && token.Type != TokenJSBracket {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBracket, "(")

			meaning.continueUntil(iter, tmpToken, ")")

			return tmpToken

		} else if token.Type != TokenJSString && (token.Content == "\"" || token.Content == "'" || token.Content == "`") {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSString, token.Content)

			meaning.continueReadString(iter, tmpToken, token.Content)

			return tmpToken

		} else if token.Content == "=" && token.Type != TokenJSAssign {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSAssign, "=")
			nextToken := iter.Get()

			if nextToken != nil {

				if nextToken.Content == ">" {
					tmpToken.Content += ">"
					tmpToken.Type = TokenJSRightArrow
					_ = iter.Read()
					return tmpToken

				} else if nextToken.Content == "=" {
					tmpToken.Content += "="
					_ = iter.Read()
					tmpToken.Type = TokenJSBinaryOperator

					if nextToken2 := iter.Get(); nextToken2 != nil && nextToken2.Content == "=" {
						tmpToken.Content += "="
						_ = iter.Read()
					}
					return tmpToken
				}
			}

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

						tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSRegex, "")

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

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhraseBreak, "")
			meaning.continueMergePhraseBreak(iter, tmpToken)
			return tmpToken

		} else if token.Content == "!" { //!, !=, !==

			token.Type = TokenJSUnaryOperator
			nextToken := iter.Get()
			if nextToken.Content == "=" {
				token.Content += "="
				_ = iter.Read()
				nextToken2 := iter.Get()
				if nextToken2.Content == "=" {
					token.Content += "="
					_ = iter.Read()
				}
				token.Type = TokenJSBinaryOperator
			}

			return token

		} else if token.Content == "&" { // &, &&

			nextToken := iter.Get()
			if nextToken.Content == "&" {
				token.Content += "&"
				_ = iter.Read()
			}
			token.Type = TokenJSBinaryOperator
			return token

		} else if token.Content == "|" { // |, ||
			nextToken := iter.Get()
			if nextToken.Content == "|" {
				token.Content += "|"
				_ = iter.Read()
			}
			token.Type = TokenJSBinaryOperator
			return token
		} else if token.Content == "+" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "+")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == "+" {
					tmpToken.Content += "+"
					tmpToken.Type = TokenJSUnaryOperator
					_ = iter.Read()
				}
			}
			return tmpToken
		} else if token.Content == "-" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "-")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == "-" {
					tmpToken.Content += "-"
					tmpToken.Type = TokenJSUnaryOperator
					_ = iter.Read()
				}
			}
			return tmpToken
		} else if token.Content == "~" {
			token.Type = TokenJSUnaryOperator
			return token
		} else if token.Content == "^" || token.Content == "%" || token.Content == "/" {
			token.Type = TokenJSBinaryOperator
			return token
		} else if token.Content == "*" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "*")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == "*" {
					tmpToken.Content += "*"
					_ = iter.Read()
				}
			}
			return tmpToken
		} else if token.Content == "case" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSCase, "")
			for {
				nextToken := iter.Read()
				if nextToken == nil || nextToken.Content == ":" {
					break
				}
				tmpToken.Children.AddToken(*nextToken)
			}
			return tmpToken
		} else if token.Content == "default" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSDefault, "")
			for {
				nextToken := iter.Read()
				if nextToken == nil || nextToken.Content == ":" {
					break
				}
			}
			return tmpToken
		} else if token.Content == "break" {
			token.Type = TokenJSBreak
			token.Content = ""
		} else if token.Content == "." {
			nextToken := iter.Get()
			nextToken2 := iter.GetBy(1)
			if nextToken != nil && nextToken2 != nil && nextToken.Content == "." && nextToken2.Content == "." {
				token.Type = TokenJSTreeDotOperator
				token.Content = "..."
				iter.Read()
				iter.Read()
			}
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

		token := iter.Read()

		if token == nil {
			return
		}
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

//testRegex test if reach regex
func (meaning *JSRawMeaning) testRegex(iter *gotokenize.Iterator) bool {

	var specialCharacter bool = false

	var i = iter.Offset
	for {
		tmpToken := iter.GetAt(i)

		if tmpToken == nil {
			return false
		}
		tmpContent := tmpToken.GetContent()

		if tmpContent == "\\" {

			specialCharacter = !specialCharacter

		} else if !specialCharacter && tmpContent == "/" {

			return true

		} else {

			specialCharacter = false
		}

		i += 1
	}
	return false
}

func (meaing *JSRawMeaning) continueMergePhraseBreak(iter *gotokenize.Iterator, currToken *gotokenize.Token) {
	for {
		token := iter.Get()
		if token == nil || !(token.Content == ";" || token.Content == "\r" || token.Content == "\n") {

			break
		}
		iter.Read()
	}
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
