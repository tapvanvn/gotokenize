package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewJSInstructionMeaning(baseMeaning gotokenize.IMeaning) *JSInstructionMeaning {

	return &JSInstructionMeaning{

		AbstractMeaning: gotokenize.NewAbtractMeaning(baseMeaning),
	}
}

type JSInstructionMeaning struct {
	*gotokenize.AbstractMeaning
}

var jsBeautifyTokens = []int{
	TokenJSIf,
	TokenJSElseIf,
	TokenJSFor,
	TokenJSFunction,
	TokenJSWhile,
	TokenJSClass,
	TokenJSClassFunction,
	TokenJSAssignVariable,
}

func (meaning *JSInstructionMeaning) Next(process *gotokenize.MeaningProcess) *gotokenize.Token {

	token := process.Iter.Read()
	if token != nil {
		if gotokenize.IsContainToken(jsBeautifyTokens, token.Type) {
			if token.Children.Length() == 1 && token.Children.GetTokenAt(0).Type == TokenJSPhrase {

				firstToken := token.Children.GetTokenAt(0)

				childProcess := gotokenize.NewMeaningProcessFromStream(append(process.Context.AncestorTokens, firstToken.Type), &firstToken.Children)

				subStream := gotokenize.CreateStream(meaning.GetMeaningLevel())

				for {

					nestedToken := meaning.Next(childProcess)

					if nestedToken == nil {
						break
					}
					subStream.AddToken(*nestedToken)
				}

				token.Children = subStream
				return token
			}
		} else if token.Type == TokenJSPhrase && token.Children.Length() == 1 {

			*token = *token.Children.GetTokenAt(0)
		}
		if gotokenize.IsContainToken(JSInstructionGlobalNested, token.Type) {

			childProcess := gotokenize.NewMeaningProcessFromStream(append(process.Context.AncestorTokens, token.Type), &token.Children)

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
		process.Context.PreviousToken = token.Type
		process.Context.PreviousTokenContent = token.Content
	} else {
		process.Context.PreviousToken = gotokenize.TokenNoType
		process.Context.PreviousTokenContent = ""
	}

	return token
}

func (meaning *JSInstructionMeaning) GetName() string {

	return "JSInstructionMeaning"
}
