package css

import "github.com/tapvanvn/gotokenize"

func NewCSSRawMeaning(baseMeaning gotokenize.IMeaning) *CSSRawMeaning {
	return &CSSRawMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type CSSRawMeaning struct {
	*gotokenize.AbstractMeaning
}

func CreateCSSMeaning() *gotokenize.PatternMeaning {
	tokenMap := map[string]gotokenize.RawTokenDefine{
		"=<>+*\"'!-:{};,()[]": {TokenType: TokenCSSOperator, Separate: true},
		" \r\n":               {TokenType: TokenCSSSpace, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)
	cssRawMeaning := NewCSSRawMeaning(meaning)
	return gotokenize.CreatePatternMeaning(cssRawMeaning, CSSPatterns, CSSIgnores, CSSGlobalNested)
}

func (meaning *CSSRawMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	return meaning.getNextMeaningToken(process.Iter)
}

func (meaning *CSSRawMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()

		if token.Content == "{" {

			tmpToken := &gotokenize.Token{
				Content: "{",
				Type:    TokenCSSBlock,
			}

			meaning.continueUntil(iter, tmpToken, "}")

			return tmpToken

		} else if token.Content == "[" {

			tmpToken := &gotokenize.Token{
				Content: "[",
				Type:    TokenCSSSquare,
			}

			meaning.continueUntil(iter, tmpToken, "]")

			return tmpToken

		} else if token.Content == "\"" {

			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenCSSString,
			}
			meaning.continueReadString(iter, tmpToken, token.Content)

			return tmpToken

		} /*else if token.Content == "." || token.Content == "-" || token.Type == TokenJSONNumber {
			tmpToken := &gotokenize.Token{
				Content: token.Content,
				Type:    TokenJSONNumberString,
			}
			tmpToken.Children.AddToken(*token)
			meaning.continueNumber(iter, tmpToken)
			return tmpToken
		}*/

		return token
	}
	return nil

}

func (meaning *CSSRawMeaning) continueUntil(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

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

func (meaning *CSSRawMeaning) continueReadString(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

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
