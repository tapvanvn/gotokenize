package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewDefaultOperatorStringifier() *Stringifier {

	stringifier := NewDefaultPhraseStringifier()
	stringifier.SetProcessor(TokenJSOperatorTrail, ProcessOperatorOperatorTrailToken)
	return stringifier
}

func ProcessOperatorOperatorTrailToken(stringifier *Stringifier, token *gotokenize.Token) {

	stringifier.put("", &NeedAndHasBreakStroke)
	/*iter := token.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken != nil {
			break
		}
		if childToken.Type == TokenJSWord || childToken.Type == TokenJSKeyWord {
			stringifier.PutToken(childToken, )
		}
	}*/
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &DefaultStroke)

}
