package xml

import "github.com/tapvanvn/gotokenize"

type XMLHightMeaning struct {
	gotokenize.PatternMeaning
}

func CreateXMLMeaning() XMLHightMeaning {
	return XMLHightMeaning{
		PatternMeaning: CreateXMLRawMeaning(),
	}
}

func (meaning *XMLHightMeaning) Prepare(stream *gotokenize.TokenStream) {
	meaning.PatternMeaning.Prepare(stream)
	tmpStream := gotokenize.CreateStream()
	for {
		token := meaning.PatternMeaning.Next()
		if token == nil {
			break
		}
		tmpStream.AddToken(*token)
	}
	meaning.SetStream(tmpStream)
}

func (meaning *XMLHightMeaning) Next() *gotokenize.Token {

	return meaning.getNextMeaningToken(meaning.GetIter())
}

func (meaning *XMLHightMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	if iter.EOS() {
		return nil
	}

	token := iter.Read()
	if token.Type == TokenXMLTagBegin {
		tmpToken := &gotokenize.Token{
			Type:    TokenXMLElement,
			Content: token.Content,
		}
		token.Content = ""
		token.Type = TokenXMLElementAttributes
		tmpToken.Children.AddToken(*token)
		meaning.continueTag(token.Content, iter, tmpToken)
		return tmpToken
	} else if token.Type == TokenXMLTagUnknown {
		tmpToken := &gotokenize.Token{
			Type:    TokenXMLElement,
			Content: token.Content,
		}
		token.Content = ""
		token.Type = TokenXMLElementAttributes
		tmpToken.Children.AddToken(*token)
		return tmpToken
	}
	return token

}
func (meaning *XMLHightMeaning) continueTag(name string, iter *gotokenize.Iterator, currentToken *gotokenize.Token) {
	for {
		if iter.EOS() {
			break
		}
		token := meaning.getNextMeaningToken(iter)
		if token.Type == TokenXMLTagEnd && token.Content == name {
			break
		}
		currentToken.Children.AddToken(*token)
	}
}
