package xml

import "github.com/tapvanvn/gotokenize"

func NewXMLHighMeaning() *XMLHighMeaning {
	return &XMLHighMeaning{
		PatternMeaning: CreateXMLRawMeaning(),
	}
}

type XMLHighMeaning struct {
	*gotokenize.PatternMeaning
}

func (meaning *XMLHighMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	return meaning.getNextMeaningToken(process.Iter)
}

func (meaning *XMLHighMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	if iter.EOS() {
		return nil
	}

	token := iter.Read()
	if token.Type == TokenXMLTagBegin {
		tmpToken := &gotokenize.Token{
			Type:    TokenXMLElement,
			Content: token.Content,
		}
		tagName := token.Content
		token.Content = ""
		token.Type = TokenXMLElementAttributes

		tmpToken.Children.AddToken(*token)
		meaning.continueTag(tagName, iter, tmpToken)
		return tmpToken
	} else if token.Type == TokenXMLTagUnknown {

		tmpToken := &gotokenize.Token{
			Type:    TokenXMLEndElement,
			Content: token.Content,
		}
		token.Content = ""
		token.Type = TokenXMLElementAttributes
		tmpToken.Children.AddToken(*token)
		return tmpToken
	}
	return token

}

func (meaning *XMLHighMeaning) continueTag(name string, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	for {
		if iter.EOS() {
			break
		}
		token := meaning.getNextMeaningToken(iter)
		if token.Type == TokenXMLTagEnd && token.Content == name {
			break
		} else if token.Type == TokenXMLTagEnd {
			println("end ", token.Content, "want", name)
		}
		currentToken.Children.AddToken(*token)
	}
}
