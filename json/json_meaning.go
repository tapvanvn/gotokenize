package json

import (
	"github.com/tapvanvn/gotokenize"
)

type JSONRawMeaning struct {
	gotokenize.IMeaning
}

func CreateJSONMeaning() gotokenize.PatternMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{
		".{}[]-\\,\":": {TokenType: TokenJSONOperator, Separate: true},
		"0123456789":   {TokenType: TokenJSONNumber, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	jsonRawMeaning := JSONRawMeaning{
		IMeaning: &meaning,
	}
	return gotokenize.CreatePatternMeaning(&jsonRawMeaning, JSONPatterns, gotokenize.NoTokens, JSONGlobalNested)
}

func (meaning *JSONRawMeaning) Next() *gotokenize.Token {
	iter := meaning.GetIter()

	return meaning.getNextMeaningToken(iter)
}

func (meaning *JSONRawMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()

		if token.Content == "{" {

			tmpToken := &gotokenize.Token{
				Content: "{",
				Type:    TokenJSONBlock,
			}

			meaning.continueUntil(iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "[" {

			tmpToken := &gotokenize.Token{
				Content: "[",
				Type:    TokenJSONSquare,
			}

			meaning.continueUntil(iter, tmpToken, "]")

			return tmpToken

		} else if token.Content == "\"" {

			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenJSONString,
			}
			meaning.continueReadString(iter, tmpToken, token.Content)

			return tmpToken

		} else if token.Content == "." || token.Content == "-" || token.Type == TokenJSONNumber {
			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenJSONNumberString,
			}
			tmpToken.Children.AddToken(*token)
			meaning.continueNumber(iter, tmpToken)
			return tmpToken
		}

		return token
	}
	return nil
}

func (meaning *JSONRawMeaning) continueUntil(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

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

func (meaning *JSONRawMeaning) continueReadString(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

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

func (meaning *JSONRawMeaning) continueNumber(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	var token = iter.Get()

	for {
		if token != nil && (token.Type == TokenJSONNumber || token.Content == ".") {

			currentToken.Children.AddToken(*token)
			_ = iter.Read()
			token = iter.Get()

		} else {

			break
		}
	}
}
