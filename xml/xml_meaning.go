package xml

import (
	"strings"

	"github.com/tapvanvn/gotokenize"
)

type XMLRawMeaning struct {
	gotokenize.IMeaning
}

func CreateXMLRawMeaning() gotokenize.PatternMeaning {

	tokenMap := map[string]gotokenize.RawTokenDefine{
		"=<>/\\\"'!-": {TokenType: TokenXMLOperator, Separate: true},
		" \r\n":       {TokenType: TokenXMLSpace, Separate: false},
	}
	meaning := gotokenize.CreateRawMeaning(tokenMap, false)

	xmlRawMeaning := XMLRawMeaning{
		IMeaning: &meaning,
	}
	return gotokenize.CreatePatternMeaning(&xmlRawMeaning, XMLPatterns, XMLIgnores, XMLGlobalNested)
}

func (meaning *XMLRawMeaning) Next() *gotokenize.Token {

	iter := meaning.GetIter()

	return meaning.getNextMeaningToken(iter)
}

func (meaning *XMLRawMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	if iter.EOS() {
		return nil
	}
	token := iter.Read()
	if token.Content == "<" {
		nextToken := iter.Get()
		third := iter.GetBy(1)
		forth := iter.GetBy(2)
		check := nextToken != nil && nextToken.Content == "!"
		check = check && third != nil && forth != nil
		check = check && third.Content == forth.Content && third.Content == "-"
		if check {

			tmpToken := &gotokenize.Token{
				Type: TokenXMLComment,
			}
			meaning.continueComment(iter, tmpToken)
			return tmpToken

		} else {
			tagToken := &gotokenize.Token{
				Type: TokenXMLTagUnknown,
			}
			meaning.continueTag(iter, tagToken)
			return tagToken
		}
	} else if token.Content == "\"" || token.Content == "'" {

		tmpToken := &gotokenize.Token{
			Content: token.Content,
			Type:    TokenXMLString,
		}
		meaning.continueReadString(iter, tmpToken, token.Content)
		return tmpToken
	}
	return token

}

func (meaning *XMLRawMeaning) continueTag(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	var reach string = ">"
	var closeTag = false
	var posEndTag = false
	var stackContent = ""
	for {
		if iter.EOS() {
			break
		}

		token := meaning.getNextMeaningToken(iter)

		if token.Content == "/" {

			if currentToken.Content == "" {
				closeTag = true
			} else {
				posEndTag = true
			}
		} else if token.Content == reach {

			break
		} else if token.Type == 0 && token.Content != " " {
			if currentToken.Content == "" {

				currentToken.Content = stackContent + token.Content
			} else {
				currentToken.Children.AddToken(*token)
			}
			posEndTag = false
		} else if token.Type != TokenXMLSpace {
			if currentToken.Content == "" {
				stackContent += token.Content
			} else {
				currentToken.Children.AddToken(*token)
			}
			posEndTag = false
		}
	}
	if closeTag {
		currentToken.Type = TokenXMLTagEnd
	} else if !posEndTag {
		currentToken.Type = TokenXMLTagBegin
	}
}

func (meaning *XMLRawMeaning) continueReadString(iter *gotokenize.Iterator, currentToken *gotokenize.Token, reach string) {

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

			if specialCharacter {

				currentToken.Children.AddToken(*lastSpecialToken)
			}
			specialCharacter = false
			currentToken.Children.AddToken(*token)
		}

	}
}

func (meaning *XMLRawMeaning) continueComment(iter *gotokenize.Iterator, currentToken *gotokenize.Token) {

	lastContent := ""
	for {
		if iter.EOS() {
			break
		}
		token := iter.Read()
		if token.Content == ">" && strings.LastIndex(lastContent, "--") > -1 {
			break
		} else if token.Content == "-" {
			lastContent += token.Content
		} else {
			lastContent = ""
		}
		currentToken.Children.AddToken(*token)
	}
}
