package js

import (
	"fmt"
	"strings"

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
	/*if len(process.Context.AncestorTokens) == 0 && process.Iter.Offset == 0 {
		fmt.Print("\033[s") //save cursor the position
	}*/
	token := meaning.getNextMeaningToken(&process.Context, process.Iter)

	if token != nil {

		if token.Children.Length() > 0 && gotokenize.IsContainToken(JSGlobalNested, token.Type) {

			childProcess := gotokenize.NewMeaningProcessFromStream(append(process.Context.AncestorTokens, token.Type), &token.Children)

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
		process.Context.PreviousToken = token.Type
		process.Context.PreviousTokenContent = token.Content

	} else {
		process.Context.PreviousToken = gotokenize.TokenNoType
		process.Context.PreviousTokenContent = ""
	}
	/*if len(process.Context.AncestorTokens) == 0 {
		fmt.Print("\033[u\033[K") //restore
		fmt.Printf("%s percent: %f%%\n", meaning.GetName(), process.GetPercent())
		fmt.Print("\033[A")
	}*/
	return token
}

func (meaning *JSRawMeaning) getNextMeaningToken(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {

			break
		}
		token := iter.Read()

		if token.Content == "{" && token.Type != TokenJSBlock {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBlock, "{")

			meaning.continueUntil(context, iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "[" && token.Type != TokenJSBracketSquare {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBracketSquare, "[")

			meaning.continueUntil(context, iter, tmpToken, "]")

			return tmpToken

		} else if token.Content == "(" && token.Type != TokenJSBracket {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBracket, "(")

			meaning.continueUntil(context, iter, tmpToken, ")")

			return tmpToken

		} else if token.Type != TokenJSString && (token.Content == "\"" || token.Content == "'" || token.Content == "`") {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSString, token.Content)

			meaning.continueReadString(context, iter, tmpToken, token.Content)

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

					meaning.continueReadLineComment(context, iter, tmpToken)

					return tmpToken

				} else if nextToken.Content == "*" {

					tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBlockComment, "/*")

					_ = iter.Read()

					meaning.continueReadBlockComment(context, iter, tmpToken)

					return tmpToken

				} else if nextToken.Content == "=" {
					_ = iter.Read()
					return gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSAssign, "/=")

				} else {

					if meaning.testRegex(context, iter) {

						tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSRegex, "")

						tmpToken.Children.AddToken(gotokenize.Token{Type: TokenJSWord, Content: "/"})

						meaning.continueReadRegex(context, iter, tmpToken)

						return tmpToken
					} else {

						tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "/")
						return tmpToken
					}
				}
			}

		} else if token.Content == " " || token.Content == "\t" {

			return meaning.getNextMeaningToken(context, iter)

		} else if token.Content == ";" {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhraseBreak, ";")
			meaning.continueMergePhraseBreak(context, iter, tmpToken)
			return tmpToken

		} else if token.Content == "\n" || token.Content == "\r" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhraseBreak, "")
			meaning.continueMergePhraseBreak(context, iter, tmpToken)
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
			token.Type = TokenJSBinaryOperator
			if nextToken.Content == "&" {
				token.Content += "&"
				_ = iter.Read()
			} else if nextToken.Content == "=" {
				token.Content += "="
				token.Type = TokenJSAssign
				_ = iter.Read()
			}

			return token

		} else if token.Content == "|" { // |, ||
			nextToken := iter.Get()
			token.Type = TokenJSBinaryOperator
			if nextToken.Content == "|" {
				token.Content += "|"
				_ = iter.Read()
			} else if nextToken.Content == "=" {
				token.Content += "="
				token.Type = TokenJSAssign
				_ = iter.Read()
			}

			return token

		} else if token.Content == "+" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "+")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == "+" {
					tmpToken.Content += "+"
					tmpToken.Type = TokenJSUnaryOperator
					_ = iter.Read()
				} else if nextToken.Content == "=" {
					tmpToken.Content += "="
					tmpToken.Type = TokenJSAssign
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
				} else if nextToken.Content == "=" {
					tmpToken.Content += "="
					tmpToken.Type = TokenJSAssign
					_ = iter.Read()
				}
			}
			return tmpToken
		} else if token.Content == "%" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "%")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == "=" {
					tmpToken.Content += "="
					tmpToken.Type = TokenJSAssign
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
					if nextToken2 := iter.Get(); nextToken2 != nil && nextToken2.Content == "=" {
						tmpToken.Content += "="
						tmpToken.Type = TokenJSAssign
						_ = iter.Read()
					}
				} else if nextToken.Content == "=" {
					tmpToken.Content += "="
					tmpToken.Type = TokenJSAssign
					_ = iter.Read()
				}
			}
			return tmpToken

		} else if token.Content == ">" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, ">")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == ">" {
					tmpToken.Content += ">"
					_ = iter.Read()
					if nextToken2 := iter.Get(); nextToken2 != nil {
						if nextToken2.Content == ">" {
							tmpToken.Content += ">"
							tmpToken.Type = TokenJSBinaryOperator
							_ = iter.Read()
						} else if nextToken2.Content == "=" {
							tmpToken.Content += "="
							tmpToken.Type = TokenJSBinaryOperator
							_ = iter.Read()
						}
					}
				} else if nextToken.Content == "=" {
					tmpToken.Content += "="
					tmpToken.Type = TokenJSBinaryOperator
					_ = iter.Read()
				}
			}
			return tmpToken

		} else if token.Content == "<" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSBinaryOperator, "<")
			if nextToken := iter.Get(); nextToken != nil {
				if nextToken.Content == "<" {
					tmpToken.Content += "<"
					_ = iter.Read()
					if nextToken2 := iter.Get(); nextToken2 != nil {
						if nextToken2.Content == "<" {
							tmpToken.Content += "<"
							tmpToken.Type = TokenJSBinaryOperator
							_ = iter.Read()
						} else if nextToken2.Content == "=" {
							tmpToken.Content += "="
							tmpToken.Type = TokenJSBinaryOperator
							_ = iter.Read()
						}
					}
				} else if nextToken.Content == "=" {
					tmpToken.Content += "="
					tmpToken.Type = TokenJSBinaryOperator
					_ = iter.Read()
				}
			}
			return tmpToken

		} else if token.Content == "." {

			token.Type = TokenJSBinaryOperator
			nextToken := iter.Get()
			nextToken2 := iter.GetBy(1)
			if nextToken != nil && nextToken2 != nil && nextToken.Content == "." && nextToken2.Content == "." {
				token.Type = TokenJSTreeDotOperator
				token.Content = "..."
				iter.Read()
				iter.Read()
			}
		} else if token.Content == "?" {
			token.Type = TokenJSQuestionOperator
		} else if token.Content == "," {
			token.Type = TokenJSSoftBreak
		} else if token.Content == "case" {
			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSCase, "")
			meaning.continueCase(context, iter, tmpToken)

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
		} else if token.Content == "continue" {
			token.Type = TokenJSContinue
			token.Content = ""
		} else if token.Content == ":" {
			token.Type = TokenJSColonOperator
		} else if (token.Type == TokenJSWord || token.Type == 0) && strings.Index(JSKeyWords, fmt.Sprintf(",%s,", token.Content)) > 0 {
			token.Type = TokenJSKeyWord
		}

		if token.Type == 0 {

			token.Type = TokenJSWord
		}
		return token
	}
	return nil
}
func (meaning *JSRawMeaning) continueCase(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	for {
		nextToken := meaning.getNextMeaningToken(context, iter)
		if nextToken == nil || nextToken.Content == ":" {
			break
		}
		if nextToken.Type == TokenJSOperator && len(strings.TrimSpace(nextToken.Content)) == 0 {
			continue
		}
		currentToken.Children.AddToken(*nextToken)
	}
}

func (meaning *JSRawMeaning) continueUntil(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

	var specialCharacter bool = false
	var lastSpecialToken *gotokenize.Token = nil

	for {
		if iter.EOS() {
			break
		}

		token := meaning.getNextMeaningToken(context, iter)

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

func (meaning *JSRawMeaning) continueReadString(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

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
func (meaning *JSRawMeaning) isStartPhrase(context *gotokenize.MeaningContext) bool {
	if context.PreviousToken == gotokenize.TokenNoType ||
		context.PreviousToken == TokenJSAssign {
		return true
	}
	return false
}

//testRegex test if reach regex
func (meaning *JSRawMeaning) testRegex(context *gotokenize.MeaningContext, iter *gotokenize.Iterator) bool {
	if !meaning.isStartPhrase(context) {
		return false
	}
	var specialCharacter bool = false

	var i = iter.Offset
	for {
		tmpToken := iter.GetAt(i)

		if tmpToken == nil || tmpToken.Type == TokenJSPhraseBreak || tmpToken.Content == "\n" || tmpToken.Content == "\r" {

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

func (meaing *JSRawMeaning) continueMergePhraseBreak(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currToken *gotokenize.Token) {
	for {
		token := iter.Get()
		if token == nil || !(token.Content == "\r" || token.Content == "\n") {

			break
		}
		iter.Read()
	}
}

func (meaning *JSRawMeaning) continueReadRegex(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currToken *gotokenize.Token) {

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

func (meaning *JSRawMeaning) continueReadLineComment(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

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

func (meaning *JSRawMeaning) continueReadBlockComment(context *gotokenize.MeaningContext, iter *gotokenize.Iterator, currToken *gotokenize.Token) {

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
