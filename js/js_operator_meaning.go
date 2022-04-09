package js

import (
	"strings"

	"github.com/tapvanvn/gotokenize/v2"
)

var notOperatorTrail = ",:,"

func NewJSOperatorMeaning(baseMeaning gotokenize.IMeaning) *JSOperatorMeaning {

	return &JSOperatorMeaning{

		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSOperatorMeaning struct {
	*gotokenize.AbstractMeaning
}

func (meaning *JSOperatorMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := meaning.getNextMeaningToken(process.Iter)

	if token != nil && gotokenize.IsContainToken(JSGlobalNested, token.Type) {

		proc := gotokenize.NewMeaningProcessFromStream(append(process.ParentTokens, token.Type), &token.Children)

		newStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

		for {
			token := meaning.Next(proc)
			if token == nil {
				break
			}
			newStream.AddToken(*token)
		}
		token.Children = newStream
	}

	return token
}

func (meaning *JSOperatorMeaning) getNextMeaningToken(iter *gotokenize.Iterator) *gotokenize.Token {

	if token := iter.Read(); token != nil {

		if token.Type == TokenJSWord || token.Type == TokenJSPhrase {

			meaning.continueOperatorPhrase(iter, token)
		}
		return token
	}

	return nil
}
func (meaning *JSOperatorMeaning) continueOperatorPhrase(iter *gotokenize.Iterator, parentToken *gotokenize.Token) {
	trail := gotokenize.CreateStream(meaning.GetMeaningLevel())
	trail.AddToken(*parentToken)
	for {
		token := iter.Get()
		if token == nil || token.Type != TokenJSBinaryOperator || strings.Contains(notOperatorTrail, ","+token.Content+",") {
			break
		}
		_ = iter.Read()
		trail.AddToken(*token)
		trail.AddToken(*iter.Read())
	}
	if trail.Length() > 2 {
		parentToken.Type = TokenJSOperatorTrail
		parentToken.Children = trail
	}
}
