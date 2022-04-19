package js

import "github.com/tapvanvn/gotokenize/v2"

func NewDefaultPatternStringifier() *Stringifier {

	stringifier := NewDefaultOperatorStringifier()
	stringifier.SetProcessor(TokenJSLabel, ProcessPatternLabel)

	return stringifier
}

func ProcessPatternLabel(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &BreakAfterStroke)
}
