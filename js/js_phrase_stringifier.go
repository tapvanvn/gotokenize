package js

import (
	"github.com/tapvanvn/gotokenize/v2"
)

func NewDefaultPhraseStringifier() *Stringifier {

	stringifier := NewDefaultRawStringifier()
	stringifier.SetProcessor(TokenJSPhrase, ProcessPhrasePhrase)
	stringifier.SetProcessor(TokenJSPhraseAssign, ProcessPhraseAssign)
	stringifier.SetProcessor(TokenJSPhraseInlineIf, ProcessPhraseInlineIf)
	stringifier.SetProcessor(TokenJSPhraseIfTrail, ProcessPhraseIfTrail)
	stringifier.SetProcessor(TokenJSPhraseSwitch, ProcessPhraseSwitch)
	stringifier.SetProcessor(TokenJSBreak, ProcessPhraseBreakStatement)
	stringifier.SetProcessor(TokenJSPhraseDo, ProcessPhraseDo)
	stringifier.SetProcessor(TokenJSPhraseFor, ProcessPhraseFor)
	stringifier.SetProcessor(TokenJSPhraseWhile, ProcessPhraseWhile)
	stringifier.SetProcessor(TokenJSPhraseTry, ProcessPhraseTry)
	stringifier.SetProcessor(TokenJSPhraseLambda, ProcessPhraseLambda)
	stringifier.SetProcessor(TokenJSPhraseFunction, ProcessPhraseFunction)
	stringifier.SetProcessor(TokenJSPhraseClass, ProcessPhraseClass)
	stringifier.SetProcessor(TokenJSPhraseClassFunction, ProcessPhraseClassFunction)
	stringifier.SetProcessor(TokenJSBracket, ProcessPhraseBracket)
	stringifier.SetProcessor(TokenJSBracketSquare, ProcessPhraseBracketSquare)
	stringifier.SetProcessor(TokenJSBlock, ProcessPhraseBlock)
	return stringifier
}

func ProcessPhraseBracket(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("(", &BreakAfterStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put(")", &DefaultStroke)
}
func ProcessPhraseBracketSquare(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("[", &BreakAfterStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("]", &DefaultStroke)
}
func ProcessPhraseBlock(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("{", &BreakAfterStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("}", &DefaultStroke)
}
func ProcessPhraseStream(stringifier *Stringifier, iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.PutToken(token)
	}
}
func ProcessPhraseStreamSpace(stringifier *Stringifier, iter *gotokenize.Iterator) {
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put(" ", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
}

func ProcessPhrasePhrase(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseSpacing(stringifier, token)
	stringifier.put(" ", &DefaultStroke)
}
func ProcessPhrasePhraseDefault(stringifier *Stringifier, parentToken *gotokenize.Token) {
	iter := parentToken.Children.Iterator()
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.PutToken(token)
	}
	stringifier.put("", &DefaultStroke)
}
func ProcessPhrasePhraseSpacing(stringifier *Stringifier, parentToken *gotokenize.Token) {
	iter := parentToken.Children.Iterator()
	for {
		token := iter.Read()
		if token == nil {
			break
		}
		stringifier.put(" ", &BreakAfterStroke)
		stringifier.PutToken(token)
	}
	stringifier.put("", &DefaultStroke)
}

func ProcessPhraseAssign(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("=", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
	stringifier.put("", &DefaultStroke)

}

func ProcessPhraseInlineIf(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &DefaultStroke)
}

func ProcessPhraseIfTrail(stringifier *Stringifier, token *gotokenize.Token) {
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
			stringifier.put("", &DefaultStroke)
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
				stringifier.put("", &DefaultStroke)
			}
		} else {
			break
		}
	}
}

func ProcessPhraseSwitch(stringifier *Stringifier, token *gotokenize.Token) {
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
					ProcessPhraseStream(stringifier, identity.Children.Iterator())
					stringifier.put(":", &BreakAfterStroke)
				} else if identity.Type == TokenJSDefault {
					stringifier.put("default:", &NeedAndHasBreakStroke)
				}
				ProcessPhraseStream(stringifier, iter)
			}
		}
	}
	stringifier.put("}", &BreakAfterStroke)
}
func ProcessPhraseBreakStatement(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("break ", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseSpacing(stringifier, token)
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseFor(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("for(", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseSpacing(stringifier, token.Children.GetTokenAt(1))

	stringifier.put(")", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseDo(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("do ", &NeedAndHasBreakStroke)
	body := token.Children.GetTokenAt(1)
	stringifier.PutToken(body)
	if body.Type == TokenJSBlock {
		stringifier.put("", &BreakAfterStroke)
	}
	stringifier.put("while ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(3))
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseWhile(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("while ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseTry(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &DefaultStroke)

}
func ProcessPhraseLambda(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put("=>", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
}
func ProcessPhraseFunction(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("function ", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(1))
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(2))
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(token.Children.GetTokenAt(3))
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseClass(stringifier *Stringifier, token *gotokenize.Token) {
	iter := token.Children.Iterator()
	iter.Read() //class
	stringifier.put("class ", &NeedAndHasBreakStroke)
	if next := iter.Get(); next != nil && (next.Type == TokenJSWord || (next.Type == TokenJSKeyWord && next.Content == "extends")) {
		iter.Read() //class name
		stringifier.PutToken(next)
		stringifier.put(" ", &BreakAfterStroke)
		if next.Type == TokenJSWord {
			if next2 := iter.Get(); next2 != nil && next2.Content == "extends" {
				iter.Read() //class name
				stringifier.PutToken(next)
				stringifier.put(" ", &BreakAfterStroke)

				stringifier.PutToken(iter.Read())
				stringifier.put(" ", &BreakAfterStroke)
			}
		} else {
			stringifier.PutToken(iter.Read())
			stringifier.put(" ", &BreakAfterStroke)
		}
	}
	stringifier.put("", &BreakAfterStroke)
	stringifier.PutToken(iter.Read())
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseClassFunction(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &DefaultStroke)
}
