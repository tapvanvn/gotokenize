package json

import "github.com/tapvanvn/gotokenize"

type JSONMeaning struct {
	gotokenize.IMeaning
}

func CreateJSONMeaning() JSONMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{
		".{}[]-\\,\":": {TokenType: TokenJSONOperator, Separate: true},
		"0123456789":   {TokenType: TokenJSONNumber, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	patternMeaning := gotokenize.CreatePatternMeaning(&meaning, JSONPatterns, gotokenize.NoTokens, gotokenize.NoTokens)
	return JSONMeaning{
		IMeaning: &patternMeaning,
	}
}

func (meaning *JSONMeaning) Next() *gotokenize.Token {
	iter := meaning.GetIter()

	return meaning.getNextMeaningToken(iter)
}

func (meaning *JSONMeaning) getNextMeaningToken(iter *gotokenize.TokenStreamIterator) *gotokenize.Token {

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

func (meaning *JSONMeaning) continueUntil(iter *gotokenize.TokenStreamIterator, currentToken *gotokenize.Token, reach string) {

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

func (meaning *JSONMeaning) continueReadString(iter *gotokenize.TokenStreamIterator, currentToken *gotokenize.Token, reach string) {

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

func (meaning *JSONMeaning) continueNumber(iter *gotokenize.TokenStreamIterator, currentToken *gotokenize.Token) {

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
