package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewJSPhraseMeaning(baseMeaning gotokenize.IMeaning) *JSPhraseMeaning {
	return &JSPhraseMeaning{
		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSPhraseMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSPhraseMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	iter := process.Iter

	token := meaning.getNextMeaningToken(iter)

	if token != nil && token.Type == TokenJSBracket {

		childProcess := gotokenize.NewMeaningProcessFromStream(&token.Children)

		meaning.Prepare(childProcess)

		subStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

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

func (meaning *JSPhraseMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	for {
		if iter.EOS() {

			break
		}
		token := iter.Read()

		if gotokenize.IsContainToken(JSPhraseAllow, token.Type) {

			tmpToken := gotokenize.NewToken(meaning.GetMeaningLevel(), TokenJSPhrase, token.Content)

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
func (meaning *JSPhraseMeaning) GetName() string {
	return "JSPhraseMeaning"
}
