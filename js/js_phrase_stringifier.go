package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewDefaultPhraseStringifier() *JSPhraseStringifier {

	return &JSPhraseStringifier{

		JSRawStringifier: NewDefaultRawStringifier(),
	}
}

type JSPhraseStringifier struct {
	*JSRawStringifier
}

func (stringifier *JSPhraseStringifier) PutToken(token *gotokenize.Token) {
	if token == nil {
		return
	}
	switch token.Type {
	//phrase
	case TokenJSPhrase:
		stringifier.PutPhrase(token)
	case TokenJSPhraseAssign:
		stringifier.PutPhraseAssign(token)
	case TokenJSPhraseBreak:
		if token.Content == ";" {
			stringifier.put(";", &BreakAfterStroke)
		}
		break
	case TokenJSPhraseInlineIf:

		stringifier.PutPhraseInlineIf(token)
	case TokenJSPhraseIfTrail:
		stringifier.PutPhraseIfTrail(token)
	case TokenJSPhraseSwitch:
		stringifier.PutPhraseSwitch(token)
	case TokenJSBreak:
		stringifier.PutBreakStatement(token)
	case TokenJSPhraseDo:
		stringifier.PutPhraseDo(token)
	case TokenJSPhraseWhile:
		stringifier.PutPhraseWhile(token)
	case TokenJSPhraseFor:
		stringifier.PutPhraseFor(token)
	case TokenJSPhraseTry:
		stringifier.PutPhraseTry(token)
	case TokenJSPhraseLambda:
		stringifier.PutPhraseLambda(token)
	case TokenJSPhraseFunction:
	case TokenJSPhraseClass:
		stringifier.PutPhraseClass(token)
	case TokenJSPhraseClassFunction:
		stringifier.PutPhraseClassFunction(token)
		//block
	case TokenJSBracket:
		stringifier.put("(", &BreakAfterStroke)
		stringifier.PutPhraseDefault(token)
		stringifier.put(")", &DefaultStroke)
	case TokenJSBracketSquare:
		stringifier.put("[", &BreakAfterStroke)
		stringifier.PutPhraseDefault(token)
		stringifier.put("]", &DefaultStroke)
	case TokenJSBlock:
		stringifier.put("{", &BreakAfterStroke)
		stringifier.PutPhraseDefault(token)
		stringifier.put("}", &DefaultStroke)

	default:
		stringifier.JSRawStringifier.PutToken(token)
	}
}
func (stringifier *JSPhraseStringifier) PutStream(iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.PutToken(token)
	}
}
func (stringifier *JSPhraseStringifier) PutStreamSpace(iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put(" ", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
}
func (stringifier *JSPhraseStringifier) PutPhrase(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutPhraseDefault(token)
	stringifier.put(" ", &DefaultStroke)
}
func (stringifier *JSPhraseStringifier) PutPhraseDefault(parentToken *gotokenize.Token) {
	iter := parentToken.Children.Iterator()
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put("", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
}
func (stringifier *JSPhraseStringifier) PutPhraseSpacing(parentToken *gotokenize.Token) {
	iter := parentToken.Children.Iterator()
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put(" ", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
}

func (stringifier *JSPhraseStringifier) PutPhraseAssign(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutPhraseDefault(token)
	stringifier.put("", &DefaultStroke)

}

func (stringifier *JSPhraseStringifier) PutPhraseInlineIf(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutPhraseDefault(token)
	stringifier.put("", &DefaultStroke)
}

func (stringifier *JSPhraseStringifier) PutPhraseIfTrail(token *gotokenize.Token) {
	iter := token.Children.Iterator()
	stringifier.put("if ", &NeedAndHasBreakStroke)
	_ = iter.Read() //if
	conditionBracket := iter.Read()
	bodyPhrase := iter.Read()

	stringifier.PutToken(conditionBracket)
	stringifier.put("", &BreakAfterStroke)

	if bodyPhrase != nil {

		stringifier.PutToken(bodyPhrase)
		if bodyPhrase.Type != TokenJSBlock {

			stringifier.put(";", &BreakAfterStroke)
		} else {
			stringifier.put("", &BreakAfterStroke)
		}
	}

	for {
		if next := iter.Get(); next != nil && next.Content == "else" {
			stringifier.put("else ", &NeedAndHasBreakStroke)
			_ = iter.Read()
			if next := iter.Get(); next != nil && next.Content == "if" {
				stringifier.put("if ", &NeedAndHasBreakStroke)
				_ = iter.Read()
				stringifier.PutToken(iter.Read())
				stringifier.put("", &BreakAfterStroke)
			}
			bodyPhrase := iter.Read()
			stringifier.PutToken(bodyPhrase) //body block
			if bodyPhrase.Type != TokenJSBlock {

				stringifier.put(";", &BreakAfterStroke)
			} else {
				stringifier.put("", &BreakAfterStroke)
			}
		} else {
			break
		}
	}

}

func (stringifier *JSPhraseStringifier) PutPhraseSwitch(token *gotokenize.Token) {
	stringifier.put("switch ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("{", &BreakAfterStroke)
	body := token.Children.GetTokenAt(2)
	if body != nil {
		iter := body.Children.Iterator()
		for {
			childToken := iter.Read()
			if childToken == nil {
				break
			}
			if childToken.Type == TokenJSPhrase {
				iter := childToken.Children.Iterator()
				identity := iter.Read()
				if identity.Type == TokenJSCase {
					stringifier.put("case ", &NeedAndHasBreakStroke)
					stringifier.PutStream(identity.Children.Iterator())
					stringifier.put(":", &BreakAfterStroke)
				} else if identity.Type == TokenJSDefault {
					stringifier.put("default:", &NeedAndHasBreakStroke)
				}
				stringifier.PutStream(iter)
			}
		}
	}
	stringifier.put("}", &BreakAfterStroke)
}
func (stringifier *JSPhraseStringifier) PutBreakStatement(token *gotokenize.Token) {
	stringifier.put("break ", &NeedAndHasBreakStroke)
	stringifier.PutPhraseSpacing(token)
	stringifier.put("", &DefaultStroke)
}
func (stringifier *JSPhraseStringifier) PutPhraseFor(token *gotokenize.Token) {
	stringifier.put("for(", &NeedAndHasBreakStroke)
	stringifier.PutPhraseSpacing(token.Children.GetTokenAt(1))

	stringifier.put(")", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
	stringifier.put("", &DefaultStroke)
}
func (stringifier *JSPhraseStringifier) PutPhraseDo(token *gotokenize.Token) {
	stringifier.put("do", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("while", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(3))
	stringifier.put("", &DefaultStroke)
}
func (stringifier *JSPhraseStringifier) PutPhraseWhile(token *gotokenize.Token) {
	stringifier.put("while", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
}
func (stringifier *JSPhraseStringifier) PutPhraseTry(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutPhraseDefault(token)
	stringifier.put("", &DefaultStroke)

}
func (stringifier *JSPhraseStringifier) PutPhraseLambda(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutPhraseDefault(token)
	stringifier.put("", &DefaultStroke)
}
func (stringifier *JSPhraseStringifier) PutPhraseClass(token *gotokenize.Token) {
	iter := token.Children.Iterator()
	iter.Read() //class
	stringifier.put("class ", &NeedAndHasBreakStroke)
	if next := iter.Get(); next != nil && next.Type == TokenJSWord {
		iter.Read() //class name
		stringifier.PutToken(next)
		stringifier.put(" ", &BreakAfterStroke)
	}
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(iter.Read())
	stringifier.put("", &DefaultStroke)
}
func (stringifier *JSPhraseStringifier) PutPhraseClassFunction(token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutPhraseDefault(token)
	stringifier.put("", &DefaultStroke)
}
