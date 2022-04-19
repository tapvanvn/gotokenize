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
	stringifier.SetProcessor(TokenJSContinue, ProcessPhraseContinueStatement)
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
	stringifier.SetProcessor(TokenJSReturnStatement, ProcessPhraseReturnStatement)
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
	if token.Children.Length() == 0 {
		return
	}
	stringifier.put("", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	//TODO: improve this
	/*#if token.Children.Length() > 0 {
		if last := token.Children.GetTokenAt(token.Children.Length() - 1); last != nil &&
			(last.Type == TokenJSUnaryOperator || last.Type == TokenJSBracket || last.Type == TokenJSBracketSquare) {
			stringifier.put(" ", &DefaultStroke)
			fmt.Println("last", JSTokenName(last.Type))
		}
	}*/
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
	/*if stringifier.HasStrongBreak {
		stringifier.put(" ", &BreakAfterStroke)
	} else {
		stringifier.put(" ", &DefaultStroke)
	}*/
}

func ProcessPhraseAssign(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("", &NeedAndHasBreakStroke)
	stringifier.PutToken(token.Children.GetTokenAt(0))
	stringifier.put(token.Children.GetTokenAt(1).Content, &BreakAfterStroke)

	right := token.Children.GetTokenAt(2)
	stringifier.PutToken(right)
	//if right.Type == TokenJSBlock {
	//	stringifier.put("", &DefaultStroke)
	//}
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

			if childToken.Type == TokenJSCase {
				stringifier.put("case ", &NeedAndHasBreakStroke)
				ProcessPhraseStream(stringifier, childToken.Children.Iterator())
				stringifier.put(":", &BreakAfterStroke)
			} else if childToken.Type == TokenJSDefault {
				stringifier.put("default:", &NeedAndHasBreakStroke)
			}
			stringifier.PutToken(iter.Read())

		}
	}
	stringifier.put("}", &BreakAfterStroke)
}
func ProcessPhraseBreakStatement(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("break ", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseContinueStatement(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("continue ", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token)
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseReturnStatement(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("return ", &NeedAndHasBreakStroke)
	iter := token.Children.Iterator()
	iter.Read()
	ProcessPhraseStream(stringifier, iter)
	stringifier.put("", &DefaultStroke)
}
func ProcessPhraseFor(stringifier *Stringifier, token *gotokenize.Token) {
	stringifier.put("for(", &NeedAndHasBreakStroke)
	ProcessPhrasePhraseDefault(stringifier, token.Children.GetTokenAt(1))
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
	if token.Children.GetTokenAt(2).Type == TokenJSBlock {
		stringifier.put("", &BreakAfterStroke)
	} else {
		stringifier.put("", &DefaultStroke)
	}

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
