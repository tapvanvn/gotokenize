package js

import (
	"fmt"

	"github.com/tapvanvn/gotokenize/v2"
)

func NewDefaultRawStringifier() *JSRawStringifier {
	return &JSRawStringifier{
		Stringifier: NewStringfier(),
	}
}

type JSRawStringifier struct {
	*Stringifier
}

func (stringifier *JSRawStringifier) PutToken(token *gotokenize.Token) {

	switch token.Type {
	case TokenJSBreak:
		stringifier.PutBracket(token)
	case TokenJSBracketSquare:
		stringifier.PutBracketSquare(token)
	case TokenJSBlock:
		stringifier.PutBlock(token)
	case TokenJSRegex:
		stringifier.PutRegexToken(token)
	case TokenJSString:
		stringifier.PutString(token)
	case TokenJSWord:
		stringifier.PutWord(token)
	case TokenJSPhraseBreak:
		if token.Content == ";" {
			stringifier.put(";", &BreakAfterStroke)
		}
		break
	case TokenJSStrongBreak:
		stringifier.put(";", &BreakAfterStroke)
	default:
		if token.Content == "," {
			stringifier.put(",", &BreakAfterStroke)
		} else if token.Type == TokenJSLineComment || token.Type == TokenJSBlockComment {
		} else {
			stringifier.put(token.Content, &DefaultStroke)
		}
	}
}

func (stringifier *JSRawStringifier) PutBlock(token *gotokenize.Token) {

	stringifier.put("{", &BreakAfterStroke)
	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		stringifier.PutToken(childToken)
	}
	stringifier.put("}", &BreakAfterStroke)
}
func (stringifier *JSRawStringifier) PutBracket(token *gotokenize.Token) {

	stringifier.put("(", &BreakAfterStroke)

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		fmt.Println("-----", childToken.Content)
		stringifier.PutToken(childToken)
	}
	stringifier.put(")", &DefaultStroke)
}
func (stringifier *JSRawStringifier) PutBracketSquare(token *gotokenize.Token) {

	stringifier.put("[", &BreakAfterStroke)

	iter := token.Children.Iterator()
	for {
		childToken := iter.Read()

		if childToken == nil {
			break
		}
		stringifier.PutToken(childToken)
	}
	stringifier.put("]", &DefaultStroke)
}
func (stringifier *JSRawStringifier) PutRegexToken(token *gotokenize.Token) {

	if token.Type == TokenJSBracket ||
		token.Type == TokenJSBlock ||
		token.Type == TokenJSBracketSquare {
		stringifier.PutRegexBlockStyle(token)
	} else if token.Type == TokenJSRegex {
		stringifier.PutRegexStream(token)
	} else {
		stringifier.put(token.Content, &DefaultStroke)
	}
}
func (stringifier *JSRawStringifier) PutRegexStream(parentToken *gotokenize.Token) {

	iter := parentToken.Children.Iterator()
	for {
		childToken := iter.Read()
		if childToken == nil {
			break
		}
		if childToken.Type == TokenJSBracket ||
			childToken.Type == TokenJSBlock ||
			childToken.Type == TokenJSBracketSquare {
			stringifier.PutRegexBlockStyle(childToken)
		} else {
			stringifier.PutRegexToken(childToken)
		}
	}
}
func (stringifier *JSRawStringifier) PutRegexBlockStyle(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.PutRegexStream(token)
	if token.Type == TokenJSBracket {
		stringifier.put(")", &DefaultStroke)
	} else if token.Type == TokenJSBlock {
		stringifier.put("}", &DefaultStroke)
	} else {
		stringifier.put("]", &DefaultStroke)
	}

}

func (stringifier *JSRawStringifier) PutString(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultStroke)
	stringifier.put(token.Children.ConcatStringContent(), &DefaultStroke)
	stringifier.put(token.Content, &DefaultStroke)
}
func (stringifier *JSRawStringifier) PutWord(token *gotokenize.Token) {
	stringifier.put(token.Content, &DefaultSpaceStroke)
}
